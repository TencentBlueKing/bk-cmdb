#!/bin/bash

# image.sh用于一键打包cmdb镜像，在编译打包好的cmdb目录下执行

# copy dockerfile catalog
cp -r ../../../../docs/support-file/dockerfile/ .

# 获取版本信息
version=$(./cmdb_adminserver/cmdb_adminserver --version | grep "Version" | head -n 1 | awk '{print $3}')

# service list
services=(adminserver authserver coreservice eventserver operationserver toposerver apiserver cloudserver hostserver procserver taskserver webserver cacheservice datacollection synchronizeserver migrate)

# cp binary file and conf dir
for service in "${services[@]}"; do
    if [[ ${service} == "migrate" ]]; then
        continue
    fi
    mkdir -p "dockerfile/${service}/cmdb_${service}"
    cp -f "cmdb_${service}/cmdb_${service}" "dockerfile/${service}/cmdb_${service}/"

    mkdir -p "dockerfile/${service}/cmdb_${service}/conf"
    cp -r "cmdb_${service}/conf/errors" "dockerfile/${service}/cmdb_${service}/conf/"
    cp -r "cmdb_${service}/conf/language" "dockerfile/${service}/cmdb_${service}/conf/"
done

# 处理webserver
cp -dpr web "dockerfile/webserver/cmdb_webserver/"
cp -dpr changelog_user "dockerfile/webserver/cmdb_webserver/"

# 打包镜像
for service in "${services[@]}"; do
      cd dockerfile/${service}/
      cat dockerfile
      docker build -t "cmdb_${service}:${version}" -f dockerfile .
      cd ../../
done