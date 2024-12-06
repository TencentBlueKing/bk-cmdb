#!/bin/bash
set -e

# get local IP.
localIp=`python ip.py`
curl -X POST -H 'Content-Type:application/json' -H 'X-Bkcmdb-User:migrate' -H 'X-Bk-Tenant-Id:0' http://${localIp}:60004/migrate/v3/authcenter/init

echo ""
