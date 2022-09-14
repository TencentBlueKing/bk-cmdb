### 功能描述

查询资源池中的主机

### 请求参数
{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| page       |  dict    | 否     | 查询条件 |
| host_property_filter| object| 否| 主机属性组合查询条件 |
| fields  |  array   | 是     | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输   |

#### host_property_filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition       |  string    | 否     |  |
| rules      |  array    | 否     | 规则 |

#### rules
过滤规则为四元组 `field`, `operator`, `value`

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| field|string|是|无|字段名 |字段名|
| operator|string|是|无|操作符 |可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value| string | 否| 无|操作数|不同的operator对应不同的value格式|

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
| 名称             | 类型   | 说明          | Description                     |
| ---------------- | ------ | ------------- | -------------------------------  |
| bk_os_type       | string | 操作系统类型  | 1:Linux;2:Windows;3:AIX         |                            |
| bk_mac           | string | 内网MAC地址   |                                 |                               |                              |
| bk_host_innerip  | string | 内网IP        |                                 |
| bk_host_id       | int    | 主机ID        |                                 |
| bk_cloud_id      | int    | 云区域        |  
