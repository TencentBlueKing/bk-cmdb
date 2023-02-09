###  获取主机详情

* API: GET /api/v3/hosts/{bk_supplier_account}/{bk_host_id}
* API名称： get_host_base_info
* 功能说明：
	* 中文：获取主机基础信息详情
	* English ：get host base info
* input body：
无
* input参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商账号 |supplier account code |
| bk_host_id| int| 是|无|主机ID | host ID |

* output:
```
{
  "result": true, 
  "bk_error_code": 0, 
  "bk_error_msg": "", 
  "data": [
    {
      "bk_property_id": "bk_host_name", 
      "bk_property_name": "主机名", 
      "bk_property_value": "centos7"
    }, 
    {
      "bk_property_id": "bk_host_id", 
      "bk_property_name": "主机ID", 
      "bk_property_value": "1007"
    }
  ]
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data说明：

| 名称  | 类型  | 说明 | Description|
| ---  | ---  | --- | ---|
| bk_property_id| string| 属性id | property ID |
| bk_property_name| string| 属性名称 |property name |
| bk_property_value| string| 属性值 | property value |


### 根据主机id获取主机快照数据

*  API:   GET /api/v3/hosts/snapshot/{bk_host_id}
* API名称： get_host_snapshot
* 功能说明：
	* 中文：获取主机详情
	* English ：get host detail
* input body：
无
* input参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | --- |
| bk_host_id| int| 是|无|主机id | host ID |


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
        "Cpu":1,
        "Disk":49,
        "HostName":"VM_0_31_centos",
        "Mem":997,
        "OsName":"linux009",
        "bootTime":1505463112,
        "cpuUsage":30.2,
        "diskUsage":0,
        "hosts":[
            "127.0.0.1 localhost localhost.localdomain VM_0_31_centos",
            "::1 localhost localhost.localdomain localhost6 localhost6.localdomain6",
            ""
        ],
        "loadavg":"0 0 0",
        "memUsage":2287,
        "memUsed":228,
        "rcvRate":0,
        "route":[
            "Kernel IP routing table",
            "Destination Gateway Genmask Flags Metric Ref Use Iface",
            "10.0.0.0 0.0.0.0 255.255.255.0 U 0 0 0 eth0",
            "169.254.0.0 0.0.0.0 255.255.0.0 U 1002 0 0 eth0",
            "0.0.0.0 10.0.0.1 0.0.0.0 UG 0 0 0 eth0",
            ""
        ],
        "iptables":[
            "",
            ""
        ],
        "sendRate":0,
        "timezone":"Asia/Shanghai",
        "timezone_number":8,
        "upTime":"2017-09-19 16:57:07"
    }
}
```

* output字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| Cpu| int|  cpu个数 | cpu number|
| Mem| int| 内存大小单位M | memory size|
| bootTime| int| 系统启动时间时间戳 | boot time|
| cpuUsage| int| cpu利用率，这个是乘以100后的值，展示需要除以100 eg:101 =1.01% | cpu usage|
| diskUsage| int| 磁盘利用率，这个是乘以100后的值，展示需要除以100 eg:1100 = 11% | disk usage|
| hosts| 字符串数组| 系统hosts文件| server hosts info |
| loadavg| string| 系统负载 | load avg|
| memUsage| int| 内存使用率，这个是乘以100后的值，展示需要除以100 eg:101 =1.01%  | memory usage|
| memUsed| init| 已经用的内存大小，单位M | the mem used|
| rcvRate| int| 系统总入流量，这个是乘以100后的值，展示需要除以100 eg:101 =1.01 | receive rate|
| route| 字符串数组| 路由信息|route info|
| iptables| 字符串数组| iptable信息 | iptables array|
| sendRate| int| 系统总流出，这个是乘以100后的值，展示需要除以100 eg:111=1.11 |send rate|
| timezone_number| int | 数字时区 | time zone number|
| upTime| string | 最近更新时间 |data update time|




