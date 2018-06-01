#!/bin/bash
set -e

# create configure
get_local_ip() {
    ip addr | \
        awk -F'[ /]+' '/inet/{
               split($3, N, ".")
               if ($3 ~ /^192.168/) {
                   print $3
               }
               if (($3 ~ /^172/) && (N[2] >= 16) && (N[2] <= 31)) {
                   print $3
               }
               if ($3 ~ /^10\./) {
                   print $3
               }
          }'

   return $?
}

localIP=`get_local_ip`

rd_server=''
redis_ip=''
redis_port=''
redis_user=''
redis_pass=''
redis
while getopts :d:r:p:x:s:m:P:X:S:u:U opt
do  
    case $opt in
        i)  
            image=$(echo $OPTARG | sed 's_/_\\/_g')
            ;;
        :)
            echo "-$OPTARG needs an argument"
            exit
            ;;
        h)  
            echo "-i <base_image>"
            exit
            ;;
        *)  
            echo "-$opt not recognized"
            exit
            ;;
    esac
done

echo ${localIP}

# list the current directories
#DIRS=$(find * -maxdepth 0 -type d | grep cmdb_)


# generate the right configure
#for tmp_dir in $DIRS;do
#    pushd $(pwd) >/dev/null
#        cd $tmp_dir
#        if [ -f $old_pwd/cmdb_adminserver/configures/${cmdbNameToConfigureMap[${tmp_dir}]}.conf ];then
#            cp -f  $old_pwd/cmdb_adminserver/configures/${cmdbNameToConfigureMap[${tmp_dir}]}.conf conf/${tmp_dir}.conf
#        fi
#        if [ $tmp_dir != "cmdb_adminserver" ];then
#            sed  -e "s/cmdb-name-placeholder/${tmp_dir}/g;s/\${localIp}/${cmdbArrayIP[${tmp_dir}]}/g;s/cmdb-port-placeholder/${cmdbArrayPort[${tmp_dir}]}/g;s/--regdiscv=127.0.0.1:2181/--regdiscv=${regdiscv}/g" template.start > start.sh
#        else
#            sed -e "s/cmdb-name-placeholder/${tmp_dir}/g;s/\${localIp}/${cmdbArrayIP[${tmp_dir}]}/g;s/cmdb-port-placeholder/${cmdbArrayPort[${tmp_dir}]}/g;s/--regdiscv=127.0.0.1:2181/--config=conf\/cmdb_adminserver.conf/g" template.start > start.sh
#        fi
#    popd >/dev/null
#done


echo "---Success---"



