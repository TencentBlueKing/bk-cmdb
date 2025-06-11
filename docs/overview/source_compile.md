# CMDB 编译指南

---

## 编译环境

- golang >= 1.20

- python >= 2.7.5

  注：尽量使用 python2 环境编译，python3 有可能导致脚本运行失败

- nodejs >= 4.0.0（编译过程中需要可以连公网下载依赖包）

- [铜锁](https://github.com/Tongsuo-Project/Tongsuo/) 8.3.2 版本

#### 安装编译依赖的铜锁环境

因为 cloud-server 使用了 crypto-golang-sdk ，所以需要安装铜锁环境并设置环境变量进行编译。参考文档：[crypto-golang-sdk](https://github.com/TencentBlueKing/crypto-golang-sdk/blob/master/readme.md)
Makefile 里内置了安装 8.3.2 版本铜锁环境的步骤，默认安装路径为`${cmdb编译路径}/tongsuo`，可以通过设置`TONGSUO_PATH`环境变量指定该安装路径

  注：最好使用nodejs14 LTS版本，如14.13.1、14.18.1、14.21.3

#### 将go mod设置为auto
```
go env -w GO111MODULE="auto"
```

## 源码下载

``` shell
git clone https://github.com/TencentBlueKing/bk-cmdb configcenter

# clone 速度较慢或超时建议配置代理：
git config --global http.proxy IP:PORT
git config --global https.proxy IP:PORT
```

## 下载项目所需依赖
``` shell
cd configcenter

go mod tidy

# go依赖下载失败或超时建议修改代理：
go env -w GOPROXY=https://goproxy.cn,direct
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

执行打包后进入归档包存放位置，解压cmdb.tar.gz，进入cmdb目录执行以下命令：

``` shell
./image.sh
```

## 编译问题及解决

### 源码下载问题

clone 速度较慢或超时建议配置代理：

```shell
git config --global http.proxy IP:PORT
git config --global https.proxy IP:PORT
git config --global --list//查看全局代理配置
```

### python版本问题

尽量使用 python2 环境编译，python3 有可能导致脚本运行失败；

当系统同时存在安装了 python2 和 python3 时，编译脚本可能会报错 `Failed: Command 'python' not found`，想在执行`python`命令时自动执行Python 2.x版本，可以按照以下步骤进行配置：（Ubuntu系统，其它系统可自行查找）

1. 确认Python 2.x已经安装：首先请确保已经在Ubuntu中安装了Python 2.x版本。
2. 使用`update-alternatives`设置优先级：Ubuntu提供了`update-alternatives`命令，可以用于管理系统中的可选命令。使用该命令来设置Python版本优先级：

```shell
sudo update-alternatives --install /usr/bin/python python /usr/bin/python2 1
```

上述命令将Python 2 设置为优先级为1的备选项，这将使`python`命令自动关联到Python 2.x版本。`/usr/bin/python2`为`python2`文件目录

### go mod tidy 失败、速度较慢或超时问题

`go env`命令检查go环境

```shell
//开启包管理工具
go env -w GO111MODULE="on"
//配置代理
go env -w GOPROXY=https://goproxy.cn,direct
```

### 前端编译失败问题

#### nodejs 版本问题：

最好使用nodejs14 LTS版本，如：14.13.1、14.18.1、14.21.3

#### 编译前端 `make ui`时可能报错如下：

```shell
npm ERR! code ELIFECYCLE
npm ERR! syscall spawn
npm ERR! file sh
npm ERR! errno ENOENT
npm ERR! fibers@5.0.3 install: 'node build.js || nodejs build.js'
npm ERR! spawn ENOENT
npm ERR!
npm ERR! Failed at the fibers@5.0.3 install script.
npm ERR! This is probably not a problem with npm. There is likely additional logging output above
npm ERR! A complete log of this run can be found in:/root/.npm/_logs/2023-07-31T12_06_25_972Z-debug.log
```
或
```shell
> fibers@5. 0.3 install /root/configcenter/src/ui/node_modules/fibers
> node build.js || nodejs build.js 

internal/modules/cjs/loader.js:883
  throw err;
  
Error: Cannot find module '/root/configcenter/src/ui/node_modules/fibers/build.js'
  at Function.Module.resolveFilename (internal/modules/cjs/loader.js:880:15)
  at Function.Module._load (internal/ modules/cjs/loader.js:725:27)
  at Function . executeUserEntryPoint [as runMain] ( internal/ modules/run_ main. js:72:12)
  at internal/main/run_main module.js:17:47 {
 code: 'MODULE_ NOT_ FOUND',
 requireStack: []
}
sh: 1: nodejs: Permission denied
```

**解决办法：**
```
1、检查主机内存是否充足，编译时内存不足可能会导致编译失败
2、安装 nvm 工具，通过 nvm 工具安装 nodejs
```
##### 通过 git 安装 nvm 步骤如下：

1. 终端执行 git

```shell
git clone https://github.com/creationix/nvm.git ~/.nvm && cd ~/.nvm && git checkout `git describe --abbrev=0 --tags`
```

2. 回到root目录执行，编辑环境变量配置文件

```shell
cd ~
vim .bashrc
```

将

```shell
source ~/.nvm/nvm.sh
```

写入环境变量配置文件并保存退出

3. 执行

```shell
source  .bashrc
和
nvm -v  //查看nvm版本号
```

显示版本号即安装nvm成功

##### 通过 nvm 工具安装 nodejs

1. 
```shell
   nvm ls-remote //查看能够使用的node版本号  
   ```

2. 这里选择了 v14.21.3，使用 nvm 命令来安装，而且将其设置为默认版本号。 分别执行：

```shell
nvm install 14.21.3
nvm alias default 14.21.3
```

3. 安装好的 nodejs 中是默认安装 npm 的，接着 `vim .bashrc` 打开环境变量配置文件查看是否有以下两句

```shell
export NVM_DIR="/Users/YOURUSERNAME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"  # This loads nvm
```

没有加上后保存退出执行 ` source .bashrc`

4. 重新打开一个终端分别执行：成功输出版本号即可

```shell
root@LAPTOP-0RIAHE03:~# nvm -v
0.39.4
root@LAPTOP-ORIAHE03:~# node -v
V14.21.3
root@LAPTOP-ORIAHE03:~# npm -V
6.14.18
```

5. 再次进入项目src目录执行`make ui`进行编译，编译成功查看相应输出目录


### 其他问题
- 查看 [cmdb项目issues地址](https://github.com/TencentBlueKing/bk-cmdb/issues) ,寻找相同或类似问题的解决办法
- 您也可以创建issue, 带上版本号+错误日志文件+配置文件等信息，我们看到后，会第一时间为您解答；
- 同时我们也鼓励，有问题互相解答，提PR参与开源贡献，共建开源社区。
