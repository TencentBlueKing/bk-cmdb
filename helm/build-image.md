# 制作docker镜像

进入helm/image目录构建docker镜像，如果顺利执行完成将会生成 bk-cmdb 和 bk-cmdb-dev 两个镜像，其中 bk-cmdb-dev 镜像包含编译环境

```bash
(bk-cmdb) ➜  configcenter git:(v3.6.x) ✗ cd helm/image
(bk-cmdb) ➜  image git:(v3.6.x) ✗ ll
total 24
-rw-r--r--  1 hoffer  hoffer   1.8K Jan  8 11:16 Dockerfile
-rw-r--r--  1 hoffer  hoffer    95B Jan  8 14:55 Dockerfile.product
-rwxr-xr-x  1 hoffer  hoffer   288B Jan  8 15:23 build.sh
(bk-cmdb) ➜  image git:(v3.6.x) ✗

./build v3.6.3
(bk-cmdb) ➜  image git:(v3.6.x) ✗ ./build.sh
7e285647cdc48cd94bd5f9de7b6a317e9bc994500d7500ddb6590530162249ca
Sending build context to Docker daemon  475.2MB
Step 1/3 : FROM centos:7
 ---> 5e35e350aded
Step 2/3 : RUN mkdir -p /data/bin/
 ---> Using cache
 ---> 58f588baaf16
Step 3/3 : COPY bk-cmdb /data/bin/bk-cmdb/
 ---> Using cache
 ---> 62ba968a5ffa
Successfully built 62ba968a5ffa
Successfully tagged bk-cmdb:v3.6.3

(bk-cmdb) ➜  image git:(v3.6.x) ✗ docker images | grep cmdb
bk-cmdb                                                                v3.6.3                62ba968a5ffa        23 minutes ago      677MB
bk-cmdb-dev                                                            v3.6.3                79599fb4cda7        11 days ago         3.21GB
```


也可以从docker上直接拉取笔者已经编译好的docker镜像

```bash
docker pull docker.io/hoffermei/bk-cmdb:v3.6.3  # 二进制
docker pull docker.io/hoffermei/bk-cmdb-dev:v3.6.3  # 二进制和编译环境
```
