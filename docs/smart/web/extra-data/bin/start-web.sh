#! env bash
cat /app/configure/web.yaml.tmpl | /app/bin/envsutsb > /app/configure/web.yaml

web_server --logtostderr=false --v=3 --config=/app/configure/web.yaml --log-dir=/data/cmdb/cmdb_webserver/logs --addrport=0.0.0.0:80 --deployment-method=blueking
