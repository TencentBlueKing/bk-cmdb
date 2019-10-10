#!/bin/bash
set -e
pushd $(pwd) > /dev/null
    cd ${GOPATH}/src/configcenter/src
    DIRS=$(find * -maxdepth 0 -type d)
    for tmp in $DIRS;do
        FILES=$(find $tmp -name 'Makefile')
        for tmp_file in $FILES;do
            target_makefile_path=$(pwd)/$tmp_file
            if [ -f $target_makefile_path ];then
                pushd $(pwd) > /dev/null
                    cd $(dirname $target_makefile_path)
		    echo "enter directory: " $(pwd)
                    if [ "$1" = "debug" ];then
                        export ISDEBUG=true
                    fi
                    make -f Makefile
                    if [ $? -ne 0 ];then
                        exit
                    fi
                popd > /dev/null
            fi
        done
    done
popd > /dev/null

