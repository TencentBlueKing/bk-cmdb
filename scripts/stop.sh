#!/bin/bash
set -e

pushd ${BASH_SOURCE%/*} >/dev/null

# list the current directories
DIRS=$(find * -maxdepth 0 -type d | grep cmdb_)

# into the directory to stop the all cmdb process
for tmp_dir in $DIRS;do
    pushd $(pwd) > /dev/null
        num=`ps -efww | grep $tmp_dir | grep -v grep | wc -l`
        if [ "$num" -gt 0 ];then
            if [ -f "$tmp_dir/stop.sh" ];then
                cd $tmp_dir && bash stop.sh || true
            fi
        fi
    popd > /dev/null
done

for tmp_dir in $DIRS;do
    num=`ps -efww | grep $tmp_dir | grep -v grep | wc -l`
    if [ "$num" -ge 1 ];then
        echo "Stoped: $tmp_dir"
    fi
done

ps -ef| grep [c]mdb_  || true
cnt=$(pgrep cmdb_ | wc -l)
echo "Running process count: $cnt"

popd >/dev/null