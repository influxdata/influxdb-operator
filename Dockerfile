FROM scratch
MAINTAINER Gianluca Arbezzano <gianarb92@gmail.com>

ADD bin/linux/influxdb-operator influxdb-operator

ENTRYPOINT ["/influxdb-operator"]
