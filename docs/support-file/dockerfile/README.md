# BK-CMDB

蓝鲸配置平台（蓝鲸CMDB）是一个面向资产及应用的企业级配置管理平台。

本文档内容为如何根据提供的dockerfile制作cmdb镜像。

### 操作步骤
#### 对于cmdb各个服务。这里以adminserver为例：
（1）在adminserver目录里创建cmdb_adminserver目录
```
mkdir dockerfile/adminserver/cmdb_adminserver
```
（2）将adminserver的二进制拷贝到上述的cmdb_adminserver目录中

（3）在上述cmdb_adminserver目录创建conf目录，将errors,language
```
mkdir dockerfile/adminserver/cmdb_adminserver/conf
cp -r cmdb/{errors,language} dockerfile/adminserver/cmdb_adminserver/conf
```

（4）执行docker build构建镜像

注：其中webserver比较特殊，还需要将web目录拷贝到cmdb_webserver下，即与conf和二进制同级的目录下

