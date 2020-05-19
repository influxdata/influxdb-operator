/*
Copyright 2020 InfluxData

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/domain"
	influxdbv1 "github.com/influxdata/influxdb-operator/api/v1"
)

// TokenReconciler reconciles a Token object
type TokenReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;patch;update;delete
// +kubebuilder:rbac:groups=influxdb.influxdata.com,resources=tokens,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=influxdb.influxdata.com,resources=tokens/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=influxdb.influxdata.com,resources=influxdbs,verbs=get
// +kubebuilder:rbac:groups=influxdb.influxdata.com,resources=influxdbs/status,verbs=get

// Reconcile shit
func (r *TokenReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("token", req.NamespacedName)

	var result = ctrl.Result{RequeueAfter: time.Second * 10}

	var token influxdbv1.Token
	if err := r.Get(ctx, req.NamespacedName, &token); err != nil {
		log.Error(err, "unable to fetch Token")
		return result, client.IgnoreNotFound(err)
	}

	// Check if our InfluxDB exists
	var influxDB influxdbv1.InfluxDB
	var influxDBName = client.ObjectKey{Name: token.Spec.InfluxDB, Namespace: req.Namespace}

	if err := r.Get(ctx, influxDBName, &influxDB); err != nil {
		r.Recorder.Event(&token, core.EventTypeWarning, "Token", "Couldn't find InfluxDB for token")
	}

	// Lets load our Secret for our InfluxDB Token
	var influxDBSecret core.Secret
	var influxDBSecretName = client.ObjectKey{Name: influxDB.Spec.Token.SecretName, Namespace: req.Namespace}

	if err := r.Get(ctx, influxDBSecretName, &influxDBSecret); err != nil {
		r.Recorder.Event(&token, core.EventTypeWarning, "InfluxDB", err.Error())
		return result, nil
	}

	influxDBTokenValue := string(influxDBSecret.Data[influxDB.Spec.Token.SecretKey])
	influxDBClient := influxdb2.NewClient(influxDB.Spec.URL, string(influxDBTokenValue))

	// Check if Secret exists
	// Lets load our Secret for our InfluxDB Token
	var tokenSecret core.Secret
	var tokenSecretName = client.ObjectKey{Name: token.Spec.SecretName, Namespace: req.Namespace}

	if err := r.Get(ctx, tokenSecretName, &tokenSecret); err != nil {
		r.Recorder.Event(&token, core.EventTypeWarning, "Token", "Couldn't find Secret for token")

		authAPI := influxDBClient.AuthorizationsApi()
		permission := &domain.Permission{
			Action: domain.PermissionActionWrite,
			Resource: domain.Resource{
				Type: domain.ResourceTypeBuckets,
			},
		}
		permissions := []domain.Permission{*permission}
		log.Info(fmt.Sprintf("Creating a token for org %v with token %v", influxDB.Spec.Organization, influxDBTokenValue))
		auth, err := authAPI.CreateAuthorizationWithOrgId(context.Background(), influxDB.Spec.Organization, permissions)
		if err != nil {
			log.Error(err, "Help me")
		}

		log.Info(fmt.Sprintf("Creating a secret with token %v", auth.Token))

		tokenSecret.Name = token.Spec.SecretName
		tokenSecret.Namespace = req.Namespace

		tokenSecret.Data = map[string][]byte{
			"token": []byte(*auth.Token),
		}

		err = r.Create(ctx, &tokenSecret)
		if err != nil {
			log.Error(err, "Failed to create secret")

			return result, err
		}
	}

	token.Status.Exists = true
	token.Status.Enabled = true
	err := r.Status().Update(ctx, &token)

	if err != nil {
		return result, err
	}

	return ctrl.Result{}, nil
}

func (r *TokenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&influxdbv1.Token{}).
		Complete(r)
}
