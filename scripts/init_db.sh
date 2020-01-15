
#!/bin/bash
set -e

# get local IP.
localIp=`python ip.py`

# 判断是否为IPV6，是则在地址两端加中括号
if [[ ${localIp} =~ ":" ]]
then
  localIp="[${localIp}]"
fi
echo "localIp:${localIp}"

curl -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${localIp}:60004/migrate/v3/migrate/community/0

echo ""
