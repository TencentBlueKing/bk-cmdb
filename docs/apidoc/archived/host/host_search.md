###  根据条件查询主机
*  API: POST /api/v3/hosts/search
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

可以指定特定的提交查询，例如设置biz 中default =1 查资源池下主机


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
                            "bk_inst_name":"Default Area"
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


###  查询业务下的主机
*  API: POST /api/v3/hosts/app/{bk_biz_id}/list_hosts
* API名称： list_biz_hosts
* 功能说明：
  * 中文：查询业务下的主机
  * English ：list hosts in special business
* input body：
```
{
    "page":{
        "start":0,
        "limit":10,
        "sort":"bk_host_id"
    },
    "bk_biz_id":2,
    "set_ids": [1, 2],
    "module_ids": [23, 24],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "operator": "equal",
                "field": "bk_host_innerip",
                "value": "127.0.0.1"
            }
        ]
    },
}
```

* input参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| bk_biz_id| int| 是| 无|业务ID|biz condition|
| set_ids|array | 否| 无|集群ID列表|set condition|
| module_ids|array | 否| 无|模块列表|module condition|
| page| object| 否| 无|查询条件|page condition for  search|
| host_property_filter| object| 否| 无|组合查询条件||


host_property_filter 参数说明：
该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| field|string|是|无|字段名 ||
| operator|string|是|无|操作符 |可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between,begins_with,not_begins_with,contains,not_contains,ends_with,not_ends_with,is_empty,is_not_empty,is_null,is_not_null |
| value| - | 否| 无|操作数|不同的operator对应不同的value格式|

组装规则可参考: <https://querybuilder.js.org/index.html>

demo:

```json
{
  "condition": "AND",
  "rules": [
    {
      "field": "bk_host_outerip",
      "operator": "begins_with",
      "value": "127.0"
    },
    {
      "condition": "OR",
      "rules": [
        {
          "field": "bk_os_type",
          "operator": "not_in",
          "value": ["3"]
        },
        {
          "field": "bk_sla",
          "operator": "equal",
          "value": "1"
        }
      ]
    }
  ]
}
```
page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|


* output
```
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "data": {
    "count": 1,
    "info": [
      {
        "bk_asset_id": "DKUXHBUH189",
        "bk_bak_operator": "admin",
        "bk_cloud_id": "0",
        "bk_comment": "",
        "bk_cpu": 8,
        "bk_cpu_mhz": 2609,
        "bk_cpu_module": "E5-2620",
        "bk_disk": 300000,
        "bk_host_id": 17,
        "bk_host_innerip": "192.168.1.1",
        "bk_host_name": "nginx-1",
        "bk_host_outerip": "",
        "bk_isp_name": "1",
        "bk_mac": "",
        "bk_mem": 32000,
        "bk_os_bit": "",
        "create_time": "2019-07-22T01:52:21.737Z",
        "last_time": "2019-07-22T01:52:21.737Z",
        "bk_os_version": "",
        "bk_os_type": "1",
        "bk_service_term": 5,
        "bk_sla": "1",
        "import_from": "1",
        "bk_province_name": "广东",
        "bk_supplier_account": "0",
        "bk_state_name": "CN",
        "bk_outer_mac": "",
        "operator": "admin",
        "bk_sn": ""
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
| bk_isp_name| string | 所属运营商 | 0:其它；1:电信；2:联通；3:移动|
| bk_sn | string | 设备SN | |
| operator | string | 主要维护人 | |
| bk_outer_mac | string | 外网MAC | |
| bk_state_name | string | 所在国家 |CN:中国，详细值，请参考CMDB页面 |
| bk_province_name | string | 所在省份 |  |
| import_from | string | 录入方式 | 1:excel;2:agent;3:api |
| bk_sla | string | SLA级别 | 1:L1;2:L2;3:L3 |
| bk_service_term | int | 质保年限 | 1-10 |
| bk_os_type | string | 操作系统类型 | 1:Linux;2:Windows;3:AIX |
| bk_os_version | string | 操作系统版本 | |
| bk_os_bit | int | 操作系统位数 | |
| bk_mem | string | 内存容量| |
| bk_mac | string | 内网MAC地址 | |
| bk_host_outerip | string | 外网IP | |
| bk_host_name | string | 主机名称 |  |
| bk_host_innerip | string | 内网IP | |
| bk_host_id | int | 主机ID | |
| bk_disk | int | 磁盘容量 | |
| bk_cpu_module | string | CPU型号 | |
| bk_cpu_mhz | int | CPU频率 | |
| bk_cpu | int | CPU逻辑核心数 | 1-1000000 |
| bk_comment | string | 备注 | |
| bk_cloud_id | int | 云区域 | |
| bk_bak_operator | string | 备份维护人 | |
| bk_asset_id | string | 固资编号 | |



###  根据主机条件查询主机(无需指定具体的业务查询)
*  API: POST /api/v3/hosts/list_hosts_without_app
* API名称： list_hosts_without_app
* 功能说明：
  * 中文：根据主机条件查询主机
  * English ：query host based on host conditions
* input body：
```
{
    "page":{
        "start":0,
        "limit":10,
        "sort":"bk_host_id"
    },
    "bk_biz_id":2,
    "set_ids": [1, 2],
    "module_ids": [23, 24],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "operator": "equal",
                "field": "bk_host_innerip",
                "value": "127.0.0.1"
            }
        ]
    },
}
```

* input参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| bk_biz_id| int| 是| 无|业务ID|biz condition|
| set_ids|array | 否| 无|集群ID列表|set condition|
| module_ids|array | 否| 无|模块列表|module condition|
| page| object| 否| 无|查询条件|page condition for  search|
| host_property_filter| object| 否| 无|组合查询条件||


host_property_filter 参数说明：
该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| field|string|是|无|字段名 ||
| operator|string|是|无|操作符 |可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between,begins_with,not_begins_with,contains,not_contains,ends_with,not_ends_with,is_empty,is_not_empty,is_null,is_not_null |
| value| - | 否| 无|操作数|不同的operator对应不同的value格式|

组装规则可参考: <https://querybuilder.js.org/index.html>

demo:

```json
{
  "condition": "AND",
  "rules": [
    {
      "field": "bk_host_outerip",
      "operator": "begins_with",
      "value": "127.0"
    },
    {
      "condition": "OR",
      "rules": [
        {
          "field": "bk_os_type",
          "operator": "not_in",
          "value": ["3"]
        },
        {
          "field": "bk_sla",
          "operator": "equal",
          "value": "1"
        }
      ]
    }
  ]
}
```
page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|


* output
```
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "data": {
    "count": 1,
    "info": [
      {
        "bk_asset_id": "DKUXHBUH189",
        "bk_bak_operator": "admin",
        "bk_cloud_id": "0",
        "bk_comment": "",
        "bk_cpu": 8,
        "bk_cpu_mhz": 2609,
        "bk_cpu_module": "E5-2620",
        "bk_disk": 300000,
        "bk_host_id": 17,
        "bk_host_innerip": "192.168.1.1",
        "bk_host_name": "nginx-1",
        "bk_host_outerip": "",
        "bk_isp_name": "1",
        "bk_mac": "",
        "bk_mem": 32000,
        "bk_os_bit": "",
        "create_time": "2019-07-22T01:52:21.737Z",
        "last_time": "2019-07-22T01:52:21.737Z",
        "bk_os_version": "",
        "bk_os_type": "1",
        "bk_service_term": 5,
        "bk_sla": "1",
        "import_from": "1",
        "bk_province_name": "广东",
        "bk_supplier_account": "0",
        "bk_state_name": "CN",
        "bk_outer_mac": "",
        "operator": "admin",
        "bk_sn": ""
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
| bk_isp_name| string | 所属运营商 | 0:其它；1:电信；2:联通；3:移动|
| bk_sn | string | 设备SN | |
| operator | string | 主要维护人 | |
| bk_outer_mac | string | 外网MAC | |
| bk_state_name | string | 所在国家 |CN:中国，详细值，请参考CMDB页面 |
| bk_province_name | string | 所在省份 |  |
| import_from | string | 录入方式 | 1:excel;2:agent;3:api |
| bk_sla | string | SLA级别 | 1:L1;2:L2;3:L3 |
| bk_service_term | int | 质保年限 | 1-10 |
| bk_os_type | string | 操作系统类型 | 1:Linux;2:Windows;3:AIX |
| bk_os_version | string | 操作系统版本 | |
| bk_os_bit | int | 操作系统位数 | |
| bk_mem | string | 内存容量| |
| bk_mac | string | 内网MAC地址 | |
| bk_host_outerip | string | 外网IP | |
| bk_host_name | string | 主机名称 |  |
| bk_host_innerip | string | 内网IP | |
| bk_host_id | int | 主机ID | |
| bk_disk | int | 磁盘容量 | |
| bk_cpu_module | string | CPU型号 | |
| bk_cpu_mhz | int | CPU频率 | |
| bk_cpu | int | CPU逻辑核心数 | 1-1000000 |
| bk_comment | string | 备注 | |
| bk_cloud_id | int | 云区域 | |
| bk_bak_operator | string | 备份维护人 | |
| bk_asset_id | string | 固资编号 | |
