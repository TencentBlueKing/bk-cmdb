# CMDB 编译指南

---

## 编译环境

- golang >= 1.8
- python >= 2.7.5
- nodejs >= 4.0.0（编译过程中需要可以连公网下载依赖包）

## 源码下载

``` shell
cd $GOPATH/src
git clone https://github.com/Tencent/bk-cmdb  configcenter
```

**GOPATH 是go项目的根目录，配置GOPATH的示例如下:**

``` shell
mkdir -p /data/abc
export GOPATH=/data/abc
```


## 编译



### 进入源码根目录：

``` shell
cd $GOPATH/src/configcenter/src
```

#### 编译共有三种模式


编译过程中如果需要特别指定版本号需要加入以下参数：


``` shell
make VERSION=xxxx
```

**注:xxx需要替换为需要需要指定的版本号**

##### 模式一：同时编译前端UI和后端服务

``` shell
make 
```

如因环境限制需要使用特殊的npm镜像地址，编译时需要采用以下命令：



``` shell
make NPM=xxx
```

**注:xxx需要替换为需要使用的npm的名字，如:cnpm。特殊的npm 镜像需要自行安装，此处不会自动配置。**

此模式编译后会同时生成前端UI文件和后端服务文件。


##### 模式二：仅编译后端服务

``` shell
make server
```

此模式下仅会编译生成后端服务文件。

##### 模式三：仅编译前端UI

``` shell
make ui
```

如因环境限制需要使用特殊的npm镜像地址，编译时需要采用以下命令：

``` shell
make ui NPM=xxx
```

**注:xxx需要替换为需要使用的npm的名字，如:cnpm。特殊的npm 镜像需要自行安装，此处不会自动配置。**


此模式下仅会编译生成前端UI文件。

### 打包

``` shell
make package
```

归档包存放位置： $GOPATH/src/configcenter/src/bin/pub/cmdb.tar.gz 


### Docker 镜像制作

解压cmdb.tar.gz，进入cmdb目录执行以下命令：

``` shell
./image.sh -i <base_image>
```

**示例：**

``` shell
./image.sh -i linux:latest
``` 

**注：-i 参数后面配置的参数是基础镜像，基础镜像可以自己制作，也可以使用公共镜像。**
