### 功能描述

查询资源池中的主机(权限：主机池主机查看权限)

### 请求参数
{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| page       |  dict    | 否     | 查询条件 |
| host_property_filter| object| 否| 主机属性组合查询条件 |
| fields  |  array   | 否     | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输，不填默认返回所有字段   |

#### host_property_filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition       |  string    | 否     |  |
| rules      |  array    | 否     | 规则 |

#### rules
过滤规则为四元组 `field`, `operator`, `value`

| 字段     | 类型   | 必填 | 描述                                                                                                |
| -------- | ------ |  ------ |---------------------------------------------------------------------------------------------------|
| field    | string | 是   | 字段名                                                                                               |
| operator | string | 是   | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>



#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大500 |
| sort     |  string | 否     | 排序字段 |



### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_host_id"
    },
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
        {
            "field": "bk_host_outerip",
            "operator": "equal",
            "value": "127.0.0.1"
        }, {
            "condition": "OR",
            "rules": [{
                "field": "bk_os_type",
                "operator": "not_in",
                "value": ["3"]
            }, {
                "field": "bk_sla",
                "operator": "equal",
                "value": "1"
            }]
        }]
    }
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "count": 1,
    "info": [
      {
        "bk_cloud_id": "0",
        "bk_host_id": 17,
        "bk_host_innerip": "192.168.1.1",
        "bk_mac": "",
        "bk_os_type": "1"
      }
    ]
  }
}
```

### 返回结果参数说明
#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | array | 请求返回的数据 |

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int       | 记录条数 |
| info      | array     | 主机实际数据 |

#### data.info
| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| bk_host_name         | string | 主机名               |    
| bk_host_innerip      | string | 内网IP              |                                 
| bk_host_id           | int    | 主机ID              |                                 
| bk_cloud_id          | int    | 管控区域               |  
| import_from          | string | 主机导入来源,以api方式导入为3 |
| bk_asset_id          | string | 固资编号              |
| bk_cloud_inst_id     | string | 云主机实例ID           |
| bk_cloud_vendor      | string | 云厂商               |
| bk_cloud_host_status | string | 云主机状态             |
| bk_comment           | string | 备注                |
| bk_cpu               | int    | CPU逻辑核心数          |
| bk_cpu_architecture  | string | CPU架构             |
| bk_cpu_module        | string | CPU型号             |
| bk_disk              | int    | 磁盘容量（GB）          |
| bk_host_outerip      | string | 主机外网IP            |
| bk_host_innerip_v6   | string | 主机内网IPv6          |
| bk_host_outerip_v6   | string | 主机外网IPv6          |
| bk_isp_name          | string | 所属运营商             |
| bk_mac               | string | 主机内网MAC地址         |
| bk_mem               | int    | 主机名内存容量（MB）       |
| bk_os_bit            | string | 操作系统位数            |
| bk_os_name           | string | 操作系统名称            |
| bk_os_type           | string | 操作系统类型            |
| bk_os_version        | string | 操作系统版本            |
| bk_outer_mac         | string | 主机外网MAC地址         |
| bk_province_name     | string | 所在省份              |
| bk_service_term      | int    | 质保年限              |
| bk_sla               | string | SLA级别             |
| bk_sn                | string | 设备SN              |
| bk_state             | string | 当前状态              |
| bk_state_name        | string | 所在国家              |
| operator             | string | 主要维护人             |
| bk_bak_operator      | string | 备份维护人             |
**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
