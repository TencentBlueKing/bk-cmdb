### 功能描述

没有业务信息的主机查询

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | 否     | 业务id |
| page       |  object    | 是     | 查询条件 |
| host_property_filter| object| 否| 主机属性组合查询条件 |
| fields  |  array   | 是     | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输   |

#### host_property_filter
该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition       |  string    | 否     |  组合查询条件|
| rules      |  array    | 否     | 规则 |


#### rules
| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| field|string|是|无|字段名 |字段名|
| operator|string|是|无|操作符 |可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value| - | 否| 无|操作数|不同的operator对应不同的value格式|

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>



#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大500 |



### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "page": {
        "start": 0,
        "limit": 3
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
    "count": 30,
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
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 3,
        "bk_host_innerip": "192.168.15.12",
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

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int       | 记录条数 |
| info      | array     | 主机实际数据 |

#### data.info
| 名称  | 类型  | 说明 |Description|
| ---------------- | ------ | ------------- | -------------------------------  |
| bk_os_type       | string | 操作系统类型  | 1:Linux;2:Windows;3:AIX         |                            |
| bk_mac           | string | 内网MAC地址   |                                 |                               |                              |
| bk_host_innerip  | string | 内网IP        |                                 |
| bk_host_id       | int    | 主机ID        |                                 |
| bk_cloud_id      | int    | 云区域        |                                 |
