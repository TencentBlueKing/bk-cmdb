#! /usr/bin/env bash
cat /app/configure/web.yaml.tmpl | /app/bin/envsubst > /app/configure/web.yaml

/app/bin/cmdb_webserver --logtostderr=false --v=3 --config=/app/configure/web.yaml --log-dir=/app/logs --addrport=0.0.0.0:80 --deployment-method=blueking
