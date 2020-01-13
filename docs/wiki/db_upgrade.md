# CMDB版本升级说明

> CMDB数据升级主要两个部分：一是DB数据升级，与其它项目类似，主要是新特性带来的新增的表结构初始化及部分已有表结构的扩展，有时会包含一些不完整数据的检测和修复；二是权限中心数据升级（如果对接了蓝鲸权限中心），主要是做资源模板的初始化，权限中心通过资源模板识别注册过来的资源，以及用户升级权限时显示相应的资源，另外还包含一些默认用户组（比如admin用户组）的默认权限策略。



## DB数据升级
**注意** db数据升级前请做好db备份。

db数据升级通过调用`cmdb_adminserver`的如下接口实现:

```bash
curl -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${cmdb_adminserver_host}:${cmdb_adminserver_port}/migrate/v3/migrate/community/0
```

其中 `${cmdb_adminserver_host}`和`${cmdb_adminserver_port}`分别是 adminserver 的监听的地址和端口。

如果您使用的是开源版本，可以使用构建输出目录下的 `init_db.sh` 脚本实现DB升级。


### DB 升级程序原理及执行过程
cmdb 数据版本升级由当前版本号和标记了版本号的升级程序组成，当前版本可在mongodb中执行 `db.cc_System.find({type: "version"}).pretty()` 输出当前版本信息，升级程序为定义在 scene_server/admin_server/upgrader 目录下的一个个模块，比如 y3.6.201911261109 为一个升级程序，版本号为 y3.6.201911261109，执行完改目录下的升级程序后，cmdb数据版本会变更到 y3.6.201911261109。

```bash
(bk-cmdb) ➜  src git:(v3.6.x) ✗ ll scene_server/admin_server/upgrader
total 40
-rw-r--r--   1 hoffermei  staff   2.9K Jan 13 10:26 compare.go
-rw-r--r--   1 hoffermei  staff   994B Aug 12 16:08 doc.go
-rw-r--r--   1 hoffermei  staff   7.6K Jan 13 10:26 register.go
-rw-r--r--   1 hoffermei  staff   3.2K Jan 13 10:26 util.go
drwxr-xr-x  10 hoffermei  staff   320B Jan 13 10:26 v3.0.8
...
drwxr-xr-x   5 hoffermei  staff   160B Dec 19 20:04 y3.6.201911261109
drwxr-xr-x   4 hoffermei  staff   128B Jan  8 15:54 y3.6.201912241627
(bk-cmdb) ➜  src git:(v3.6.x) ✗
```

db 升级程序执行时，先将所有的升级程序按版本号升序排序，然后从小到大与当前版本号挨个比对，如果版本号大于当前版本号，则执行改版本的升级程序，然后开始下一轮循环，整个db升级可以用如下python脚本描述：

```python
all_versions = sorted(all_versions)
for version in all_versions:
    if version > current_version:
        print("run %s" % version)
        run(version)
        current_version = version
```



### DB 升级输出及故障排查

如下内容是从 v3.5.x 升级到 v3.6.x 的输出：

```bash
(bk-cmdb) ➜  v3.6.x git:(v3.6.x) ✗ ./init_db.sh
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "",
  "permission": null,
  "data": "migrate success",
  "pre_version": "x19_10_22_03",
  "current_version": "y3.6.201912241627",
  "finished_migrations": [
   "y3.6.201909062359",
   "y3.6.201909272359",
   "y3.6.201910091234",
   "y3.6.201911121930",
   "y3.6.201911122106",
   "y3.6.201911141015",
   "y3.6.201911141516",
   "y3.6.201911261109",
   "y3.6.201912241627"
  ]
 }
(bk-cmdb) ➜  v3.6.x git:(v3.6.x) ✗
```

#### 输出字段说明

通用字段

- result 字段表示升级是否出错
- bk_error_code 为出错时的具体原因
- bk_error_msg 为错误描述
- 正常执行完升级 data 字段内容为 "migrate success"
- permission为cmdb接口统一返回字段，这里直接忽略

db 升级相关字段

- pre_version 为db升级前数据库中cmdb数据处于的版本
- current_version 为执行完升级脚本后cmdb数据处于的版本
- finished_migrations 为本次升级脚本执行时完成的升级版本


#### db升级故障排查
如果在db升级过程中程序运行出错，上述输出字段 result，bk_error_code，bk_error_msg会提示错误原因，更具体的原因可通过查看 adminserver 的运行日志 `cmdb_adminserver/logs/cmdb_adminserver.ERROR` 获取

注意：如果输出信息中没有出现 result，bk_error_code，bk_error_msg 等字段，可能是没有成功访问到cmdb_adminserver导致的，如果怀疑是网络原因可以结合curl命令的 -vvv 参数重试，会有详细的debug信息输出。

```bash
curl -vvv -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${cmdb_adminserver_host}:${cmdb_adminserver_port}/migrate/v3/migrate/community/0
```

升级接口请求失败，常见的一类原因是配置了错误的ip和端口，其中cmdb_adminserver ip需要是监听的ip，可能你的主机有多个网卡，adminserver只监听了其中一个，如果还不是很确定，可在adminserver所在服务器上通过如下命令获取：

`ps aux | grep cmdb | grep adminserver | awk '{print $12}' | awk -F '=' '{print $2}'`


## 权限数据升级

> 权限数据升级指对接了蓝鲸权限中心的权限升级，如果没有对接蓝鲸权限中心(开源版本用户)，可跳过这部分说明信息

权限升级也是通过方案cmdb_adminserver的http接口进行，需要访问的接口如下：

```bash
curl -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${cmdb_adminserver_host}:${cmdb_adminserver_port}/migrate/v3/authcenter/init
```


如果在升级过程中出现问题，上述输出字段 result，bk_error_code，bk_error_msg会提示错误原因，更具体的原因可通过查看 adminserver 的运行日志 `cmdb_adminserver/logs/cmdb_adminserver.ERROR` 获取


