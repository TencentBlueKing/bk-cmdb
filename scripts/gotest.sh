#!/bin/bash
set -e

function check_and_execute_gotest(){
    count=`ls -lst $1 | grep _test.go | wc -l`
    if [ ${count} -ne 0 ]
    then
        go test -v $1 -coverprofile=covprofile_$(echo $1 | sed 's/\//_/g').profile
    fi
}

function walk_dir(){
    for element in `ls  $1`
    do  
        dir_or_file=$1"/"${element}
        if [ -d ${dir_or_file} ]
       then 
            check_and_execute_gotest ${dir_or_file}
            walk_dir ${dir_or_file}
       fi  
    done
}


# main loop
start_dir=""
while getopts :d:h opt
do  
    case $opt in
        d)  
            start_dir=$OPTARG
            ;;
        :)
            echo "-$OPTARG needs an argument"
            exit
            ;;
        h)  
            echo "-d <start directory>"
            exit
            ;;
        *)  
            echo "-$opt not recognized"
            exit
            ;;
    esac
done


if [ -z "${start_dir}" ];then
    echo "please set the start directory, eg: ./gotest.sh -d <start directory>"
    exit
fi

echo "walk the directory:${start_dir}"
walk_dir ${start_dir}