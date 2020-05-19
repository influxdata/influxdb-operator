/*
Copyright 2020 InfluxData

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package controllers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	influxdbv1 "github.com/influxdata/influxdb-operator/api/v1"
)

// InfluxDBReconciler reconciles a InfluxDB object
type InfluxDBReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get
// +kubebuilder:rbac:groups=influxdb.influxdata.com,resources=influxdbs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=influxdb.influxdata.com,resources=influxdbs/status,verbs=get;update;patch

func (r *InfluxDBReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("influxdb", req.NamespacedName)

	var result = ctrl.Result{RequeueAfter: time.Second * 10}

	var influxdb influxdbv1.InfluxDB
	if err := r.Get(ctx, req.NamespacedName, &influxdb); err != nil {
		log.Error(err, "unable to fetch InfluxDB")
		return result, client.IgnoreNotFound(err)
	}

	// Assume false until we get what we need
	influxdb.Status.Available = false
	influxdb.Status.Authenticated = false

	// Lets load our Secret for our InfluxDB Token
	var tokenSecret core.Secret
	var tokenSecretName = client.ObjectKey{Name: influxdb.Spec.Token.SecretName, Namespace: req.Namespace}

	if err := r.Get(ctx, tokenSecretName, &tokenSecret); err != nil {
		r.Recorder.Event(&influxdb, core.EventTypeWarning, "Authentication", err.Error())

		err := r.Status().Update(ctx, &influxdb)
		if err != nil {
			return result, err
		}

		return result, nil
	}

	// Lets check if our InfluxDB URL responds to /health
	// We'll time this also
	var start, connect, dns, tlsHandshake time.Time
	var trace = &httptrace.ClientTrace{}

	// Refresh stats every 10 seconds
	if time.Now().Unix()-influxdb.Status.StatsRefresh > 9 {
		influxdb.Status.StatsRefresh = time.Now().Unix()

		trace = &httptrace.ClientTrace{
			DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
			DNSDone: func(ddi httptrace.DNSDoneInfo) {
				influxdb.Status.DNSResolution = time.Since(dns).Milliseconds()
			},
			TLSHandshakeStart: func() { tlsHandshake = time.Now() },
			TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
				influxdb.Status.TLSHandshake = time.Since(tlsHandshake).Milliseconds()
			},

			ConnectStart: func(network, addr string) { connect = time.Now() },
			ConnectDone: func(network, addr string, err error) {
				influxdb.Status.ConnectTime = time.Since(connect).Milliseconds()
			},

			GotFirstResponseByte: func() {
				influxdb.Status.FirstByte = time.Since(start).Milliseconds()
			},
		}
	}

	httpReq, _ := http.NewRequest("GET", influxdb.Spec.URL, nil)
	httpReq = httpReq.WithContext(httptrace.WithClientTrace(httpReq.Context(), trace))

	start = time.Now()
	httpResp, httpErr := http.DefaultTransport.RoundTrip(httpReq)

	if httpErr != nil {
		r.Recorder.Event(&influxdb, core.EventTypeWarning, "Connection", httpErr.Error())
		err := r.Status().Update(ctx, &influxdb)

		return result, err
	}

	// Currently this is almost impossible to hit ... but you never know ...
	// https://github.com/influxdata/influxdb/blob/master/http/health.go
	if httpResp.StatusCode != http.StatusOK {
		r.Recorder.Event(&influxdb, core.EventTypeWarning, "Connection", httpErr.Error())
		err := r.Status().Update(ctx, &influxdb)

		return result, err
	}

	// OK. We consider InfluxDB Available
	influxdb.Status.Available = true

	// Test Token with /authorizations
	httpReq, httpErr = http.NewRequest("GET", influxdb.Spec.URL, nil)
	if httpErr != nil {
		r.Recorder.Event(&influxdb, core.EventTypeWarning, "Connection", httpErr.Error())
		err := r.Status().Update(ctx, &influxdb)

		return result, err
	}

	token := string(tokenSecret.Data[influxdb.Spec.Token.SecretKey])

	httpReq.URL.Path = "/api/v2/authorizations"
	httpReq.Header.Add("Authorization", fmt.Sprintf("Token %v", token))
	log.Info(fmt.Sprintf("Using token %v", token))

	client := &http.Client{}
	httpResp, httpErr = client.Do(httpReq)

	if httpErr != nil {
		r.Recorder.Event(&influxdb, core.EventTypeWarning, "Connection", httpErr.Error())
		err := r.Status().Update(ctx, &influxdb)

		return result, err
	}

	if httpResp.StatusCode != http.StatusOK {
		r.Recorder.Event(&influxdb, core.EventTypeWarning, "Authorization", "Unauthorized")
		err := r.Status().Update(ctx, &influxdb)

		return result, err
	}

	// OK. Now we consider InfluxDB Authenticated
	influxdb.Status.Authenticated = true

	err := r.Status().Update(ctx, &influxdb)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *InfluxDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&influxdbv1.InfluxDB{}).
		Complete(r)
}
