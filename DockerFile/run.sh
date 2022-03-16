#!/bin/bash
ip=127.0.0.1
port=8090

cd /data/cmdb
python init.py --discovery zookeeper:2181 --database cmdb --redis_ip redis-master --redis_port 6379 --redis_pass cc --mongo_ip mongo-mongodb --mongo_port 27017 --mongo_user cc --mongo_pass cc --blueking_cmdb_url http://${ip}:${port} --listen_port 8090 --user_info admin:admin --auth_enabled false  --full_text_search off --log_level 3

# invalid in skip-login mode

cd /data/cmdb/cmdb_adminserver/configures/

sed -i 's/opensource/skip-login/g' common.conf
sed -i 's/opensource/skip-login/g' common.yaml

# start cmdb
cd /data/cmdb

./start.sh
echo "" > /data/cmdb/cmdb_cacheservice/logs/std.log
echo "" > /data/cmdb/cmdb_toposerver/logs/std.log
echo "" > /data/cmdb/cmdb_taskserver/logs/std.log
sleep 60
./start.sh

# init data
cd cmdb_adminserver && ./init_db.sh
# hold on
tail -f /dev/null
