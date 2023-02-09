### 功能描述

查询容器集群(v3.10.23+, 权限:不需要权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id    |  int  | 是     | 业务ID|
| filter | object  | 否   | 容器集群查询范围 |
| fields | array   | 否   | 所要查询的容器集群属性，如果不写代表搜索全部数据 |
| page | object  | 是   | 分页条件 |

#### filter

该参数为容器集群属性字段过滤规则的组合，用于根据容器集群属性字段搜索容器集群。组合支持AND 和 OR 两种方式，允许嵌套，最多嵌套2层。

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition |  string  | 是    | 规则操作符|
| rules |  array  | 是     | 过滤集群的范围规则 |


#### rules
过滤规则为三元组 `field`, `operator`, `value`

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| -------- | ------ | ---- | ------ | ------ | ------------------------------------------------------------ |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大500 |
| enable_count |  bool | 是 | 本次请求是否为获取数量还是详情的标记 |
| sort     |  string | 否     | 排序字段，通过在字段前面增加 -，如 sort:&#34;-field&#34; 可以表示按照字段 field降序 |

**注意：**
- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- `sort`如果调用方没有指定，后台默认指定为容器集群ID。
- 必须设置分页参数，一次最大查询数据不超过500个。


### 请求参数示例

### 获取详细信息请求参数
```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"scheduling_engine",
                "operator":"equal",
                "value":"k8s"
            },
            {
                "field":"version",
                "operator":"equal",
                "value":"1.1.0"
            }
        ]
    },
    "page":{
        "start":0,
        "limit":500,
        "enable_count":false
    }
}
```

### 获取数量请求示例

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"scheduling_engine",
                "operator":"equal",
                "value":"k8s"
            },
            {
                "field":"version",
                "operator":"equal",
                "value":"1.1.0"
            }
        ]
    },
    "page":{
        "start":0,
        "limit":0,
        "enable_count":true
    }
}
```

### 返回结果示例

### 详细信息接口响应

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":0,
        "info":[
            {
                "name":"cluster",
                "scheduling_engine":"k8s",
                "uid":"xxx",
                "xid":"xxx",
                "version":"1.1.0",
                "network_type":"underlay",
                "region":"xxx",
                "vpc":"xxx",
                "network":"127.0.0.0/21",
                "type":"public-cluster"
            }
        ]
    },
    "request_id":"87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 获取容器集群数量接口响应

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":1,
        "info":[
        ]
    },
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 返回结果参数说明

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| data    | object | 请求返回的数据                           |
| request_id    | string | 请求链id    |

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int       | 记录条数 |
| info      | array     | 集群实际数据 |

#### info

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| name    |  string  | 否     | 集群名称|
| scheduling_engine |  string  | 否  | 调度引擎 |
| uid   |  string  | 否   | 集群自有ID|
| xid |  string  | 否   | 关联集群ID |
| version   |  string  | 否   | 集群版本 |
| network_type   |  string  | 否   | 网络类型 |
| region |  string  | 否    | 地域|
| vpc |  string  | 否    | vpc网络|
| network |  array  | 否    | 集群网络|
| type |  string  | 否     | 集群类型 |


**注意：**
- 如果本次请求是查询详细信息那么count为0，如果查询的是数量，那么info为空。
