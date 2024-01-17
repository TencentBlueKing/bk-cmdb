### 功能描述

根据业务ID查询业务下的主机，可附带其他的过滤信息，如集群id,模块id等(权限：业务访问权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 | 类型   | 必选 | 描述                                                         |
| -------------------- | ------ | ---- | ------------------------------------------------------------ |
| page                 | object   | 是   | 查询条件                                                     |
| bk_biz_id            | int    | 是   | 业务id                                                       |
| bk_set_ids           | array  | 否   | 集群ID列表，最多200条 **bk_set_ids和set_cond只能使用其中一个** |
| set_cond             | array  | 否   | 集群查询条件 **bk_set_ids和set_cond只能使用其中一个**        |
| bk_module_ids        | array  | 否   | 模块ID列表，最多500条 **bk_module_ids和module_cond只能使用其中一个** |
| module_cond          | array  | 否   | 模块查询条件 **bk_module_ids和module_cond只能使用其中一个**  |
| host_property_filter | object | 否   | 主机属性组合查询条件                                         |
| fields               | array  | 是   | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输 |

#### host_property_filter
该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition       |  string    | 否     |  组合查询条件|
| rules      |  array    | 否     | 规则 |


#### rules
| 名称     | 类型   | 必填 | 默认值 | 说明   | 描述                                                  |
| -------- | ------ | ---- | ------ | ------ | ------------------------------------------------------------ |
| field    | string | 是   | 无     | 字段名 |         字段名                                                     |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### set_cond
| 字段     | 类型   | 必选 | 描述                          |
| -------- | ------ | ---- | ----------------------------- |
| field    | string | 是   | 取值为集群的字段              |
| operator | string | 是   | 取值为：$eq $ne               |
| value    | string | 是   | field配置的集群字段所对应的值 |

#### module_cond
| 字段     | 类型   | 必选 | 描述                          |
| -------- | ------ | ---- | ----------------------------- |
| field    | string | 是   | 取值为模块的字段              |
| operator | string | 是   | 取值为：$eq $ne               |
| value    | string | 是   | field配置的模块字段所对应的值 |

#### page

| 字段  | 类型   | 必选 | 描述                 |
| ----- | ------ | ---- | -------------------- |
| start | int    | 是   | 记录开始位置         |
| limit | int    | 是   | 每页限制条数,最大500 |
| sort  | string | 否   | 排序字段             |



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
    "set_cond": [
        {
            "field": "bk_set_name",
            "operator": "$eq",
            "value": "set1"
        }
    ],
    "bk_biz_id": 3,
    "bk_module_ids": [54,56],
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
                "field": "bk_host_innerip",
                "operator": "equal",
                "value": "127.0.0.1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_os_type",
                        "operator": "not_in",
                        "value": [
                            "3"
                        ]
                    },
                    {
                        "field": "bk_cloud_id",
                        "operator": "equal",
                        "value": 0
                    }
                ]
            }
        ]
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
    "count": 2,
    "info": [
      {
        "bk_cloud_id": 0,
        "bk_host_id": 1,
        "bk_host_innerip": "192.168.15.18",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 2,
        "bk_host_innerip": "192.168.15.4",
        "bk_mac": "",
        "bk_os_type": null
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
| data | object | 请求返回的数据 |

#### data

| 字段  | 类型  | 描述         |
| ----- | ----- | ------------ |
| count | int   | 记录条数     |
| info  | array | 主机实际数据 |

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