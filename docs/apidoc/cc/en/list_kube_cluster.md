### Functional description

list container clusters (v3.10.23+, permissions: no permissions required)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id    |  int  | yes     | business ID|
| filter | object  | no   | container cluster query scope|
| fields | array   | no   | the container cluster attribute to be queried, if not written, it means to search all data |
| page | object  | yes   | paging condition |

#### filter

- This parameter is a combination of filtering rules for container cluster attribute fields, and is used to search container 
clusters according to container cluster attribute fields. The combination supports AND and OR, allowing nesting, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  | yes      | rule operator|
| rules |  array  |yes     | scope rule for filtering cluster|


#### rules
The filtering rule is triplet`field`,`operator`,`value`

| Name     | Type   | Required| Default value|  Description                                                  |
| -------- | ------ | ---- | ------ | ------------------------------------------------------------ |
| field    |  string |yes   | None     | Field name|                                                              |
| operator | string |yes   | None     | Operator| Optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between|
| value    | -      |no   | None     | Operand| Different values correspond to different value formats                            |

Assembly rules can be found at: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|
| enable_count |  bool |yes| Whether this request is a token to obtain quantity or details|
| sort     |   string |no     | Sort the field. By adding sort in front of the field, for example, sort&#34;: sort field&#34; can indicate descending order by field field|

**Note:**

- `enable_count`If this flag is true, this request is a get quantity. The remaining fields must be initialized, start is 0, and limit is: 0, sort is "."
- `sort`If the caller does not specify it, the background specifies it as the container cluster ID by default.
- Paging parameters must be set, and the maximum query data at one time does not exceed 500.

### Request Parameters Example

### Request Details Request Parameters

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

### get quantity request parameters

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

### Return Result Example

### Details interface response
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

### kube cluster quantity interface response
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

### Return Result Parameters Description

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| data    |  object |Data returned by request                           |
| request_id    |  string |Request chain id    |

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int       | Number of records|
| info      |  array     | Actual cluster data|

#### info

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| name    |  string  | no     | cluster|
| scheduling_engine |  string  | no  | scheduling engine |
| uid   |  string  | no   | cluster own ID|
| xid |  string  | no   | associated cluster ID |
| version   |  string  | no   |  cluster version |
| network_type   |  string  | no   | network type |
| region |  string  | no    | the region where the cluster is located|
| vpc |  string  | no    | vpc network|
| network |  array  | no    | cluster network|
| type |  string  | no     | cluster type |


**Note:**
- If this request is to query details, count is 0. If the query is quantity, info is empty.
