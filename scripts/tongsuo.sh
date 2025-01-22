#!/bin/bash
set -e

if [ "${DISABLE_CRYPTO}" = true ];then
  echo "crypto is disabled"
  exit 0
fi

if [ -d "$TONGSUO_PATH" ]; then
  echo "tongsuo already exists"
  exit 0
fi

echo -e "\033[34mpreparing tongsuo... \033[0m"

wget --no-check-certificate https://github.com/Tongsuo-Project/Tongsuo/archive/refs/tags/8.3.2.tar.gz
tar zxvf 8.3.2.tar.gz
cd Tongsuo-8.3.2/

if [ "$IS_STATIC" == true ]; then
  yum -y install glibc-static
  ./config --prefix=$TONGSUO_PATH -static
else
  ./config --prefix=$TONGSUO_PATH
fi

make -j
make install
