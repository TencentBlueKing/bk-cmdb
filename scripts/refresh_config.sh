#!/bin/bash
# 用途：刷新配置文件内容到配置中心
set -e

# 配置名列表
config_names="all common mongodb redis extra error language"

# 参数错误
Param_Error="first param must be one of the config name in [${config_names}]
eg: ./refresh_config all"



# 刷新配置
function refresh_config(){
  # 获取local IP.
  localIp=`python ip.py`

  # 判断是否为IPV6，是则在地址两端加中括号
  if [[ ${localIp} =~ ":" ]]
  then
    localIp="[${localIp}]"
  fi
  echo "localIp:${localIp}"

  curl -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${localIp}:60004/migrate/v3/migrate/config/refresh --data "{\"config_name\":\"$1\"}"
}

# 判断入参
is_correct=false
for name in ${config_names};do
  if [ "$1" == ${name} ];then
    is_correct=true
    break
  fi
done

if [ "${is_correct}" == true ]; then
 refresh_config $1
else
  echo -e "$Param_Error"
  exit
fi