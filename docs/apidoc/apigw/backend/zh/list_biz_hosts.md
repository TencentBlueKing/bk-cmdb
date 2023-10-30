### 描述

根据业务ID查询业务下的主机，可附带其他的过滤信息，如集群id,模块id等

### 输入参数

| 参数名称                 | 参数类型   | 必选 | 描述                                                  |
|----------------------|--------|----|-----------------------------------------------------|
| page                 | object | 是  | 查询条件                                                |
| bk_biz_id            | int    | 是  | 业务id                                                |
| bk_set_ids           | array  | 否  | 集群ID列表，最多200条 **bk_set_ids和set_cond只能使用其中一个**       |
| set_cond             | array  | 否  | 集群查询条件 **bk_set_ids和set_cond只能使用其中一个**              |
| bk_module_ids        | array  | 否  | 模块ID列表，最多500条 **bk_module_ids和module_cond只能使用其中一个** |
| module_cond          | array  | 否  | 模块查询条件 **bk_module_ids和module_cond只能使用其中一个**        |
| host_property_filter | object | 否  | 主机属性组合查询条件                                          |
| fields               | array  | 是  | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输            |

#### host_property_filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 参数名称      | 参数类型   | 必选 | 描述     |
|-----------|--------|----|--------|
| condition | string | 否  | 组合查询条件 |
| rules     | array  | 否  | 规则     |

#### rules

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------| 
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符，可选值：equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数，不同的operator对应不同的value格式                                                                       |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### set_cond

| 参数名称     | 参数类型   | 必选 | 描述                |
|----------|--------|----|-------------------|
| field    | string | 是  | 取值为集群的字段          |
| operator | string | 是  | 取值为：$eq $ne       |
| value    | string | 是  | field配置的集群字段所对应的值 |

#### module_cond

| 参数名称     | 参数类型   | 必选 | 描述                |
|----------|--------|----|-------------------|
| field    | string | 是  | 取值为模块的字段          |
| operator | string | 是  | 取值为：$eq $ne       |
| value    | string | 是  | field配置的模块字段所对应的值 |

#### page

| 参数名称  | 参数类型   | 必选 | 描述           |
|-------|--------|----|--------------|
| start | int    | 是  | 记录开始位置       |
| limit | int    | 是  | 每页限制条数,最大500 |
| sort  | string | 否  | 排序字段         |

### 调用示例

```json
{
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
  "bk_module_ids": [
    54,
    56
  ],
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

### 响应示例

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

### 响应参数说明

#### response

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称  | 参数类型  | 描述     |
|-------|-------|--------|
| count | int   | 记录条数   |
| info  | array | 主机实际数据 |

#### data.info

| 参数名称            | 参数类型   | 描述      |
|-----------------|--------|---------| 
| bk_os_type      | string | 操作系统类型  | 
| bk_mac          | string | 内网MAC地址 | 
| bk_host_innerip | string | 内网IP    | 
| bk_host_id      | int    | 主机ID    | 
| bk_cloud_id     | int    | 云区域     |
