###  根据条件查询主机
*  API: POST /api/{version}/hosts/search
* API名称： search_host
* 功能说明：
	* 中文：根据条件查询主机
	* English ：search host by condition
* input body：
```
{
    "page":{
        "start":0,
        "limit":10,
        "sort":"bk_host_id"
    },
    "pattern":"",
    "bk_biz_id":2,
    "ip":{
        "flag":"bk_host_innerip|bk_host_outerip",
        "exact":1,
        "data":[

        ]
    },
    "condition":[
        {
            "bk_obj_id":"host",
            "fields":[

            ],
            "condition":[

            ]
        },
        {
            "bk_obj_id":"module",
            "fields":[

            ],
            "condition":[

            ]
        },
        {
            "bk_obj_id":"set",
            "fields":[

            ],
            "condition":[

            ]
        },
        {
            "bk_obj_id":"biz",
            "fields":[

            ],
            "condition":[
                {
                    "field":"default",
                    "operator":"$ne",
                    "value":1
                }
            ]
        }
    ]
}
```

* input参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| ip| object| 否| 无|主机ip列表|ip condition|
| condition|object | 否| 无|组合条件|comb condition|
| page| object| 否| 无|查询条件|page condition for  search|
| pattern| string| 否| 无|按表达式搜索|search by pattern condition|


ip参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| data | ip 数组| 否| 无|ip list for search| the list for search |
| exact| int| 否| 无|是否根据ip精确搜索| is the exact query |
| flag| string| 否| 空|bk_host_innerip只匹配内网ip,bk_host_outerip只匹配外网ip, bk_host_innerip,bk_host_outerip同时匹配|bk_host_innerip match lan ip,bk_host_outerip match wan ip|

condition 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_obj_id| string| 否| 无|对象名,可以为biz,set,module,host,object|object name, it can be biz,set,module,host,object|
| fields| string数组| 否| 无|查询输出字段|fields output|
| condition| object array| 否| 无|查询条件|search condition|

二级condition 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| field| string| 否| 无|对象的字段|field of object|
| operator| string| 否| 无|操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于|$eq is equal,$in is belongs, $nin is not belong,$neq is not equal|
| value| string| 否| 无|字段对应的值|the value of field|

可以指定特定的提交查询，例如设置biz 中default =1 查资源池下主机， BK_SUPPLIER_ID_FIELD= 查询开发商下主机


page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|


* output
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"success",
    "data":{
        "count":1,
        "info":[
            {
                "biz":[
                    {
                        "bk_biz_developer":"",
                        "bk_biz_id":2,
                        "bk_biz_maintainer":"admin",
                        "bk_biz_name":"蓝鲸"
                    }
                ],
                "host":{
                    "bk_asset_id":"DKUXHBUH189",
                    "bk_bak_operator":"admin",
                    "bk_cloud_id":[
                        {
                            "id":"0",
                            "bk_obj_id":"plat",
                            "bk_obj_icon":"",
                            "bk_inst_id":0,
                            "bk_obj_name":"",
                            "bk_inst_name":"default area"
                        }
                    ],
                    "bk_comment":"",
                    "bk_cpu":8,
                    "bk_cpu_mhz":2609,
                    "bk_cpu_module":"E5-2620",
                    "bk_disk":300000,
                    "bk_host_id":17,
                    "bk_host_innerip":"192.168.1.1",
                    "bk_host_name":"nginx-1",
                    "bk_host_outerip":"",
                    "bk_isp_name":null,
                    "bk_mac":"",
                    "bk_mem":32000,
                    "bk_os_bit":""
                },
                "module":[
                    {
                        "TopModuleName":"蓝鲸##公共组件##consul",
                        "bk_bak_operator":"",
                        "bk_biz_id":2,
                        "bk_module_id":35,
                        "bk_module_name":"consul",
                        "bk_module_type":"1",
                        "bk_parent_id":8,
                        "bk_set_id":8,
                        "bk_supplier_account":"0",
                        "create_time":"2018-05-16T21:03:22.724+08:00",
                        "default":0,
                        "last_time":"2018-05-16T21:03:22.724+08:00",
                        "operator":""
                    }
                ],
                "set":[
                    {
                        "TopSetName":"蓝鲸##公共组件",
                        "bk_biz_id":2,
                        "bk_capacity":null,
                        "bk_parent_id":3,
                        "bk_service_status":"1",
                        "bk_set_desc":"111",
                        "bk_set_env":"3",
                        "bk_set_id":8,
                        "bk_set_name":"公共组件",
                        "bk_supplier_account":"0",
                        "create_time":"2018-05-16T21:03:22.692+08:00",
                        "default":0,
                        "description":"",
                        "last_time":"2018-05-18T11:50:53.947+08:00"
                    }
                ]
            }
        ]
    }
}
```

*  output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| count| int| 记录条数 |the num of record|
| info| object array | 主机实际数据 |host data|

info 字段说明:

| 名称  | 类型  | 说明 |Description|
|---|---|---|---| 
| biz | object array| 主机所属的业务信息 |host biz info|
| set| object array | 主机所属的集群信息 |host set info|
| module| object array| 主机所属的模块信息 |host module info|
| host| object | 主机自身属性|host attr info|

###  获取主机详情

* API: GET /api/{version}/hosts/{bk_supplier_account}/{bk_host_id}
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

*  API:   GET /api/{version}/hosts/snapshot/{bk_host_id}
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




