# CMDB 编译指南

---

## 编译环境

- golang >= 1.17
- python >= 2.7.5
- nodejs >= 4.0.0（编译过程中需要可以连公网下载依赖包）

#### 将go mod设置为auto
```
go env -w GO111MODULE="auto"
```

## 源码下载

``` shell
git clone https://github.com/Tencent/bk-cmdb configcenter
```

## 下载项目所需依赖
``` shell
cd configcenter

go mod tidy
```

 go mod是Golang的包管理工具，若没有开启，可以进行下面操作:
 ``` shell
 go env -w GO111MODULE="auto"

或

 go env -w GO111MODULE="on"
 ```

## 编译

### 进入源码根目录：

``` shell
cd configcenter/src
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

大陆地区用户推荐使用npm镜像cnpm进行前端编译，cnpm安装参考[cnpmjs.org](https://cnpmjs.org/)，编译时需要采用以下命令：

``` shell
make NPM=cnpm
```

**注：使用其他npm镜像与此类似**


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

大陆地区用户推荐使用npm镜像cnpm进行前端编译，cnpm安装参考[cnpmjs.org](https://cnpmjs.org/)，编译时需要采用以下命令：

``` shell
make ui NPM=cnpm
```

**注：使用其他npm镜像与此类似**


此模式下仅会编译生成前端UI文件。

### 打包

``` shell
make package
```

归档包存放位置： configcenter/src/bin/pub/cmdb.tar.gz 


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
