#!/bin/bash

# read the new version tarbar position
read -p "please input the new version package position:" new_version_package

# check the file
if [ -f $new_version_package ];then
    old_pwd=$(pwd)
    # delete the old data
    if [ -d $(pwd)/upgrade_tmp ];then
        rm -rf $(pwd)/upgrade_tmp
    fi
    mkdir $(pwd)/upgrade_tmp
    # upgrade action
    tar -zvxf $new_version_package -C $(pwd)/upgrade_tmp/
    pushd $(pwd) > /dev/null
        cd $(pwd)/upgrade_tmp/cmdb
        tmp_dirs=$(find * -maxdepth 0 -type d | grep cmdb_)
        for tmp_item in $tmp_dirs;do
            if [ -d $old_pwd/$tmp_item ];then
                pushd $(pwd) > /dev/null
                    cd $old_pwd/$tmp_item && bash ./stop.sh || true
                popd > /dev/null
                cp -f $tmp_item/cmdb_* $old_pwd/$tmp_item/
            else
                mkdir -p $old_pwd/$tmp_item
                cp -R $tmp_item $old_pwd/$tmp_item
            fi
        done
    popd > /dev/null
    # delete the template directory
    rm -rf $(pwd)/upgrade_tmp/
else
    echo "the $new_version_package is not a file"
    exit
fi
