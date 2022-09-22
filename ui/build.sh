#!/bin/bash
set -e

if [ -d "$2" ];then
    mkdir -p "$2"
fi

$1 install
$1 run build BUILD_OUTPUT="$2"