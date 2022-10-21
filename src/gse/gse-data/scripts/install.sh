#!/usr/bin/env bash

set -eo pipefail

source ./generate.sh -e ./gse_data.env -t ../etc/gse_data.conf.template > ../etc/gse_data.conf

mkdir -p ${BK_GSE_HOME_DIR}
mkdir -p ${BK_GSE_LOG_PATH}

# install
cp -rf ../etc/ ${BK_GSE_HOME_DIR}
cp -rf ../gse_data ${BK_GSE_HOME_DIR}
