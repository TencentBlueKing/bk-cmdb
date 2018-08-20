#!/bin/bash

set -e

# read the new version tarbar position
new_version_package=$1
if [ -z $new_version_package ]; then
    read -p "please input the new version package position:" new_version_package
fi

# check the file
if [ -f $new_version_package ];then
    old_pwd=$(pwd)
    # delete the old data
    if [ -d $(pwd)/upgrade_tmp ];then
        rm -rf $(pwd)/upgrade_tmp
    fi
    mkdir $(pwd)/upgrade_tmp
    # upgrade action
    echo "unpacking"
    tar -zxf $new_version_package -C $(pwd)/upgrade_tmp/
    echo "unpackcomplete"
    pushd $(pwd) > /dev/null
        cd $(pwd)/upgrade_tmp/cmdb
        tmp_dirs=$(find * -maxdepth 0 -type d | grep cmdb_)
        for tmp_item in $tmp_dirs;do
            if [ -d $old_pwd/$tmp_item ];then
                pushd $(pwd) > /dev/null
                    echo stop $tmp_item
                    cd $old_pwd/$tmp_item && bash ./stop.sh || true
                popd > /dev/null

                echo cp `realpath $old_pwd/$tmp_item/cmdb_*`
                cp -f $tmp_item/cmdb_* $old_pwd/$tmp_item/

                echo cp $tmp_item/conf/errors to `realpath $old_pwd/$tmp_item/conf/errors`
                cp -R -f $tmp_item/conf/errors $old_pwd/$tmp_item/conf

                echo cp $tmp_item/conf/language to `realpath $old_pwd/$tmp_item/conf/language`
                cp -R -f $tmp_item/conf/language $old_pwd/$tmp_item/conf
            else
                mkdir -p $old_pwd/$tmp_item
                cp -R -f $tmp_item $old_pwd/$tmp_item
            fi
        done
    popd > /dev/null

    if [ -d $(pwd)/upgrade_tmp/cmdb/web ]; then 
        echo cp `realpath $(pwd)/web`
        cp -R -f $(pwd)/upgrade_tmp/cmdb/web .
    fi

    if [ -f $(pwd)/upgrade_tmp/cmdb/init.py ]; then
        cp -R -f $(pwd)/upgrade_tmp/cmdb/init.py .
    fi

    if [ -f $(pwd)/upgrade_tmp/cmdb/upgrade.sh ]; then
        cp -R -f $(pwd)/upgrade_tmp/cmdb/upgrade.sh .
    fi

    # delete the template directory
    rm -rf $(pwd)/upgrade_tmp/

    echo -e "\033[35mall cmdb process stoped, please restart them manually \033[0m"
else
    echo "the $new_version_package is not a file"
    exit
fi
