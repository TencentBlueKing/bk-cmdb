### 功能描述

查询容器节点(v3.10.23+，权限：不需要权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id    |  int  | 是     | 业务ID|
| filter | object  | 否   | 容器节点查询范围 |
| fields | array   | 否   | 所要查询的容器节点属性，如果不写代表搜索全部数据 |
| page | object  | 是   | 分页条件 |

#### filter

该参数为容器节点属性字段过滤规则的组合，用于根据容器节点属性字段搜索容器集群。组合支持AND 和 OR 两种方式，允许嵌套，最多嵌套2层。

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition |  string  | 是    | 规则操作符|
| rules |  array  | 是     | 过滤节点的范围规则 |


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
- `sort`如果调用方没有指定，后台默认指定为节点ID。
- 必须设置分页参数，一次最大查询数据不超过500个。

### 请求参数示例

### 详细信息请求参数

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_id":2,
    "filter":{
        "condition":"OR",
        "rules":[
            {
                "field":"id",
                "operator":"equal",
                "value":10
            },
            {
                 "field":"bk_cluster_id",
                 "operator":"equal",
                 "value":10
            },
            {
                "field":"hostname",
                "operator":"equal",
                "value":"name"
            }
        ]
    },
    "page":{
        "enable_count":false,
        "start":0,
        "limit":500
    }
}
```

### 获取节点数量请求参数

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_id":2,
    "filter":{
        "condition":"OR",
        "rules":[
            {
                "field":"id",
                "operator":"equal",
                "value":10
            },
            {
               "field":"bk_cluster_id",
                "operator":"equal",
                "value":10
            },
            {
                "field":"hostname",
                "operator":"equal",
                "value":"name"
            }
        ]
    },
    "page":{
        "enable_count":true,
        "start":0,
        "limit":0
    }
}
```

### 返回结果示例

### 详细信息接口响应

```json
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"success",
    "permission":null,
    "data":{
        "count":0,
        "info":[
            {
                "name":"k8s",
                "roles":"master",
                "labels":{
                    "env":"test"
                },
                "taints":{
                    "type":"gpu"
                },
                "unschedulable":false,
                "internal_ip":[
                    "127.0.0.1"
                ],
                "external_ip":[
                    "127.0.0.1"
                ],
                "hostname":"name",
                "runtime_component":"runtime_component",
                "kube_proxy_mode":"ipvs",
                "pod_cidr":"127.0.0.128/26"
            }
        ]
    }
}
```

### 获取容器节点数量接口响应

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
| info      | array     | 节点的实际数据 |

#### info

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| name   |  string  | 否   | 节点名称 |
| roles   |  string  | 否   | 节点类型 |
| labels |  object  | 否    | 标签|
| taints |  object  | 否    | 污点|
| unschedulable |  bool| 否 | 是否关闭可调度，true为不可调度，false代表可调度|
| internal_ip |  array  | 否 | 内网IP |
| external_ip |  array  | 否  | 外网IP |
| hostname |  string  | 否     | 主机名 |
| runtime_component |  string  | 否 | 运行时组件 |
| kube_proxy_mode |  string  | 否 | kube-proxy 代理模式 |
| pod_cidr |  string  | 否 | 此节点Pod地址的分配范围 |


**注意：**
- 如果本次请求是查询详细信息那么count为0，如果查询的是数量，那么info为空。
