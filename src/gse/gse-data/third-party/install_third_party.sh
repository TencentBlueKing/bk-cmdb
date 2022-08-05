#!/usr/bin/env bash

set -eo pipefail

# libunwind
cd /tmp
wget -N https://github.com/libunwind/libunwind/releases/download/v1.5/libunwind-1.5.0.tar.gz
tar zxvf libunwind-1.5.0.tar.gz && rm -f libunwind-1.5.0.tar.gz && cd libunwind-1.5.0/
./configure  && make -j 16 && make install

# tcmalloc
cd /tmp
wget -N https://github.com/gperftools/gperftools/releases/download/gperftools-2.9.1/gperftools-2.9.1.tar.gz
tar zxvf gperftools-2.9.1.tar.gz  && rm -f cgperftools-2.9.1.tar.gz  && cd gperftools-2.9.1
./configure --disable-shared  && make -j 16 && make install

# kafka
cd /tmp
wget -N  https://github.com/edenhill/librdkafka/archive/refs/tags/v1.2.2.tar.gz
tar zxvf v1.2.2.tar.gz  && rm -f v1.2.2.tar.gz  && cd librdkafka-1.2.2/
./configure --enable-gssapi
make && make install

# pulsar
cd /tmp
wget -N https://archive.apache.org/dist/pulsar/pulsar-2.7.3/apache-pulsar-2.7.3-src.tar.gz
tar zxvf apache-pulsar-2.7.3-src.tar.gz  && rm -f apache-pulsar-2.7.3-src.tar.gz  && cd apache-pulsar-2.7.3/pulsar-client-cpp
mkdir -p build && cd build
cmake .. -DBUILD_PYTHON_WRAPPER=OFF -DLINK_STATIC=OFF -DUSE_LOG4CXX=OFF -DBUILD_TESTS=OFF -DBUILD_DYNAMIC_LIB=OFF
make install
