### 描述

查询容器节点(v3.12.1+，权限：业务访问)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                       |
|-----------|--------|----|--------------------------|
| bk_biz_id | int    | 是  | 业务ID                     |
| filter    | object | 否  | 容器节点查询范围                 |
| fields    | array  | 否  | 所要查询的容器节点属性，如果不写代表搜索全部数据 |
| page      | object | 是  | 分页条件                     |

#### filter

该参数为容器节点属性字段过滤规则的组合，用于根据容器节点属性字段搜索容器集群。组合支持AND 和 OR 两种方式，允许嵌套，最多嵌套2层。

| 参数名称      | 参数类型   | 必选 | 描述        |
|-----------|--------|----|-----------|
| condition | string | 是  | 规则操作符     |
| rules     | array  | 是  | 过滤节点的范围规则 |

#### rules

过滤规则为三元组 `field`, `operator`, `value`

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### page

| 参数名称         | 参数类型   | 必选 | 描述                                                        |
|--------------|--------|----|-----------------------------------------------------------|
| start        | int    | 是  | 记录开始位置                                                    |
| limit        | int    | 是  | 每页限制条数,最大500                                              |
| enable_count | bool   | 是  | 本次请求是否为获取数量还是详情的标记                                        |
| sort         | string | 否  | 排序字段，通过在字段前面增加 -，如 sort:&#34;-field&#34; 可以表示按照字段 field降序 |

**注意：**

- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- `sort`如果调用方没有指定，后台默认指定为节点ID。
- 必须设置分页参数，一次最大查询数据不超过500个。

### 调用示例

#### 详细信息请求参数

```json
{
  "bk_biz_id": 2,
  "filter": {
    "condition": "OR",
    "rules": [
      {
        "field": "id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "bk_cluster_id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "hostname",
        "operator": "equal",
        "value": "name"
      }
    ]
  },
  "page": {
    "enable_count": false,
    "start": 0,
    "limit": 500
  }
}
```

### 响应示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "filter": {
    "condition": "OR",
    "rules": [
      {
        "field": "id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "bk_cluster_id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "hostname",
        "operator": "equal",
        "value": "name"
      }
    ]
  },
  "page": {
    "enable_count": true,
    "start": 0,
    "limit": 0
  }
}
```

### 响应参数说明
