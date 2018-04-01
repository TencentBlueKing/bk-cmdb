#!/bin/bash
set -e

ver="88ceed1-18.03.29"
image=""

while getopts :i:h opt
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

if [ -z "$image" ];then
    echo "please set the base image, eg: ./image.sh -i <base_image>"
    exit
fi

#mypwd=$(echo $PWD | sed 's_/_\\/_g')
FILES=$(find "$(pwd)/docker" -maxdepth 1 -type f | grep Dockerfile)

for tmp_file in $FILES;do
    echo "building image: ${tmp_file##*.}:${ver} ..."
    sed -e "s/image_placeholder/${image}/g" $tmp_file > "$(pwd)/Dockerfile"
    # build image
    docker build -t "${tmp_file##*.}:${ver}" .
    rm "$(pwd)/Dockerfile"
done
