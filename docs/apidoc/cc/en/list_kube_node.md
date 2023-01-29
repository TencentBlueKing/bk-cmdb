### Functional description

list container nodes (v3.10.23+, permission: no permission required)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id    |  int  | yes     | business ID|
| filter | object  | no   | Container node query scope |
| fields | array   | no   | The attribute of the container node to be queried, if not written, it means to search all data |
| page | object  | yes   | Paging condition |

#### filter

- This parameter is a combination of filtering rules for container node attribute fields, and is used to search container
  node according to container cluster attribute fields. The combination supports AND and OR, allowing nesting, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  | yes      | Rule operator|
| rules |  array  |yes     | Scope rule for filtering node|


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
- `sort`If the caller does not specify it, the background specifies it as the container node ID by default.
- Paging parameters must be set, and the maximum query data at one time does not exceed 500.
- bk_cluster_id and cluster_uid cannot be empty or filled at the same time.

### Request Parameters Example

### Request Details Request Parameters

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

### get quantity request parameters

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

### Return Result Example

### Details interface response
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

### kube node quantity interface response
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
| info      |  array     | Actual node data|

#### info

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| name   |  string  | yes   | node name |
| roles   |  string  | no   | node roles |
| labels |  object  | no    | label|
| taints |  object  | no    | taints|
| unschedulable |  bool| no | Whether to turn off schedulable, true means not schedulable, false means schedulable|
| internal_ip |  array  | no | internal ip |
| external_ip |  array  | no  | external ip |
| hostname |  string  | no     | hostname |
| runtime_component |  string  | no | runtime components |
| kube_proxy_mode |  string  | no | kube-proxy proxy mode |
| pod_cidr |  string  | no | The allocation range of the Pod address of this node  |


**Note:**
- If this request is to query details, count is 0. If the query is quantity, info is empty.

