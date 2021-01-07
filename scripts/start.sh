#!/bin/bash

set -e

proc_num=12
withSynchronizeServer=false 


for opt in "$@"
do
    case $opt in 
    "--with-synchronize")
        withSynchronizeServer=true
        proc_num=$(($proc_num+1))
        ;;
    esac
done
pushd ${BASH_SOURCE%/*} >/dev/null
DIRS=$(find * -maxdepth 0 -type d | grep cmdb_)

# into the directory to start the all cmdb process
for tmp_dir in $DIRS;do
    if [  $tmp_dir == "cmdb_synchronizeserver" ] && [ $withSynchronizeServer == false  ] ;then 
        continue 
    fi
    pushd $(pwd) > /dev/null
    echo -e "starting: $tmp_dir"
    num=`ps -efww | grep $tmp_dir | grep -v grep | grep -v tail | wc -l`
    if [ "$num" -le 0 ];then
        if [ -f "$tmp_dir/start.sh" ];then
            cd "$tmp_dir" && bash start.sh
        fi
    fi
        
    popd > /dev/null
done

ps -ef| grep [c]mdb_ || true
cnt=$(pgrep cmdb_ | wc -l)
echo "process count should be: $proc_num , now: $cnt"

for tmp_dir in $DIRS;do
    if [  $tmp_dir == "cmdb_synchronizeserver" ] && [ $withSynchronizeServer == false  ] ;then 
        continue 
    fi
    num=`ps -efww | grep $tmp_dir | grep -v grep | grep -v tail | wc -l`
    if [ "$num" -le 0 ];then
        echo "Not Running: $tmp_dir"
    fi
done


popd >/dev/null

