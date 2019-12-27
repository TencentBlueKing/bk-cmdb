FROM golang:1.13.0
RUN apt-get update
RUN apt-get install -y git python jq curl

RUN curl -sL https://deb.nodesource.com/setup_13.x | bash -
RUN apt-get update && apt-get install -y nodejs
RUN npm install gulp -g
RUN npm install yarn -g

# usage: docker build --build-arg user=release-v3.6.3 .

ARG branch=release-v3.6.3
RUN mkdir -p /data/src/ && cd /data/src/ && git clone -b ${branch} --single-branch https://github.com/Tencent/bk-cmdb.git 
RUN apt-get install nodejs

# make server
RUN cd /data/src/ && ls -al ./ && mv bk-cmdb configcenter && cd /data/src/configcenter/src/ && export GOPATH=/data/ && make server

# generate configurations
RUN cd /data/src/configcenter/src/bin/build/*/ && ls -al . &&     python init.py  \
      --discovery          127.0.0.1:2181 \
      --database           cmdb \
      --redis_ip           127.0.0.1 \
      --redis_port         6379 \
      --redis_pass         1111 \
      --mongo_ip           127.0.0.1 \
      --mongo_port         27017 \
      --mongo_user         cc \
      --mongo_pass         cc \
      --blueking_cmdb_url  http://127.0.0.1:8080/ \
      --blueking_paas_url  http://paas.domain.com \
      --listen_port        8080 \
      --auth_scheme        internal \
      --auth_enabled       false \
      --auth_address       https://iam.domain.com/ \
      --auth_app_code      bk_cmdb \
      --auth_app_secret    xxxxxxx \
      --auth_sync_workers  1 \
      --auth_sync_interval_minutes  45 \
      --full_text_search   off \
      --es_url             http://127.0.0.1:9200 \
      --log_level          3


# make ui
RUN cd /data/src/configcenter/src/ && export GOPATH=/data/ && make ui

# copy to bin directory
RUN mkdir -p /data/bin/ && mv /data/src/configcenter/src/bin/build/* /data/bin/bk-cmdb && ls -al /data/bin/bk-cmdb/
