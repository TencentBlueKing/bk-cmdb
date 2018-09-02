#!/bin/bash
set -e

start_dir=""
output_dir=${PWD}/coverprofile

function check_and_execute_gotest(){
    echo "dir:"$2
    count=`ls -lst $1 | grep _test.go | wc -l`
    if [ ${count} -ne 0 ]
    then
        profile_name=${output_dir}"/"covprofile_$(echo $1 | sed 's/\//_/g')
        go test -v $1 -coverprofile=${profile_name}.profile
        go tool cover -html=${profile_name}.profile -o ${profile_name}.html
    fi
}

function walk_dir(){
    for element in `ls  $1`
    do  
        dir_or_file=$1"/"${element}
        if [ -d ${dir_or_file} ]
       then 
            check_and_execute_gotest ${dir_or_file} ${output_dir}
            walk_dir ${dir_or_file}
       fi  
    done
}


# main loop
while getopts :d:o:h opt
do  
    case $opt in
        d)  
            start_dir=$OPTARG
            ;;
        o)
            output_dir=$OPTARG
            ;;
        :)
            echo "-$OPTARG needs an argument"
            exit
            ;;
        h)  
            echo "-d <start directory> -o <output directory>"
            exit
            ;;
        *)  
            echo "-$opt not recognized"
            exit
            ;;
    esac
done


if [ -z "${start_dir}" ];then
    echo "please set the start directory, eg: ./gotest.sh -d <start directory> -o <output directory>"
    exit
fi

echo "walk the directory:${start_dir}"
echo "output directory:${output_dir}"
if [ ! -d ${output_dir} ];then
    mkdir -p ${output_dir}
fi
walk_dir ${start_dir} ${output_dir}