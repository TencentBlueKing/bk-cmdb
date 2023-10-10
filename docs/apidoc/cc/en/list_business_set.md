### Functional description

Query business set (v3.10.12+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_filter | object  |no   | Business set condition range|
| time_condition    |  object  |no   | Business set time range|
| fields            |  array   | no   | Query criteria. The parameter is any attribute of the business. If it is not written, it means to search all data.|
| page              |  object  |Yes.   | Paging condition|

#### bk_biz_set_filter

This parameter is a combination of filtering rules for business set attribute fields, and is used to search business sets according to business set attribute fields. The combination supports AND and OR, allowing nesting, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  | yes      | Rule operator|
| rules |  array  |yes     | Scope rule for filtering business|


#### rules
The filtering rule is triplet`field`,`operator`,`value`

| Name     | Type   | Required| Default value|  Description                                                  |
| -------- | ------ | ---- | ------ | ------------------------------------------------------------ |
| field    |  string |yes   | None     | Field name|                                                              |
| operator | string |yes   | None     | Operator| Optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between|
| value    | -      |no   | None     | Operand| Different values correspond to different value formats                            |

Assembly rules can be found at: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### time_condition

| Field   | Type   | Required| Description              |
|-------|--------|-----|--------------------|
| oper  | string |yes| Operator, currently only and is supported|
| rules | array  | yes      | Time query criteria         |

#### rules

| Field   | Type   | Required| Description                             |
|-------|--------|-----|----------------------------------|
| field | string |yes| The value is the field name of the model                  |
| start | string |yes| Start time in the format yyyy MM dd hh: mm:ss|
| end   |  string |yes| End time in the format yyyy MM dd hh: mm:ss|

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|
| enable_count |  bool |yes| Whether this request is a token to obtain quantity or details|
| sort     |   string |no     | Sort the field. By adding sort in front of the field, for example, sort&#34;: sort field&#34; can indicate descending order by field field|

**Note:**
- `enable_count`If this flag is true, this request is a get quantity. The remaining fields must be initialized, start is 0, and limit is: 0, sort is "."
- `sort`If the caller does not specify it, the background specifies it as the business set ID by default.

### Request Parameters Example

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"bk_biz_set_id",
                "operator":"equal",
                "value":10
            },
            {
                "field":"bk_biz_set_maintainer",
                "operator":"equal",
                "value":"admin"
            }
        ]
    },
    "time_condition":{
        "oper":"and",
        "rules":[
            {
                "field":"create_time",
                "start":"2021-05-13 01:00:00",
                "end":"2021-05-14 01:00:00"
            }
        ]
    },
    "fields": [
        "bk_biz_id"
    ],
    "page":{
        "start":0,
        "limit":500,
        "enable_count":false,
        "sort":"bk_biz_set_id"
    }
}
```

### Return Result Example

### Details interface response
```python

{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":0,
        "info":[
            {
                "bk_biz_set_id":10,
                "bk_biz_set_name":"biz_set",
                "bk_biz_set_desc":"dba",
                "biz_set_maintainer":"tom",
                "create_time":"2021-09-06T08:10:50.168Z",
                "last_time":"2021-10-15T02:30:01.867Z",
                "bk_scope":{
                    "match_all":true
                }
            },
            {
                "bk_biz_set_id":11,
                "bk_biz_set_name":"biz_set1",
                "bk_biz_set_desc":"dba",
                "biz_set_maintainer":"tom",
                "create_time":"2021-09-06T08:10:50.168Z",
                "last_time":"2021-10-15T02:30:01.867Z",
                "bk_scope":{
                    "match_all":false,
                    "filter":{
                        "condition":"AND",
                        "rules":[
                            {
                                "field":"bk_sla",
                                "operator":"equal",
                                "value":"3"
                            },
                            {
                                "field":"bk_biz_maintainer",
                                "operator":"equal",
                                "value":"admin"
                            }
                        ]
                    }
                }
            }
        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### Business set quantity interface response
```python
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":2,
        "info":[
        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### Return Result Parameters Description
#### response

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
| info      |  array     | Actual business data|

#### info

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_id   |   int  |yes   | Business set ID|
| create_time   |   string  |no   | Business set creation time|
| last_time   |   string  |no   | Business set modification time|
| bk_biz_set_name   |   string  |yes   | Business set name|
| bk_biz_maintainer |  string  |no   | Operation and maintenance personnel|
| bk_biz_set_desc   |   string  |no   | Business set description|
| bk_scope   |   object  |no   | Business set selected business scope|

#### bk_scope

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| match_all |  bool  |yes    | Selected business scope tag|
| filter |  object  |no     | Scope criteria for the selected business|

#### filter

This parameter is a combination of filtering rules for service attribute fields, and is used to search for services according to the service attribute fields. The combination only supports AND operation and can be nested, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  |yes     | Scope condition rule for selected business|


#### rules

| Name     | Type   | Required| Default value|  Description                                                  |
| -------- | ------ | ---- | ------ | ------------------------------------------------------------ |
| field    |  string |yes   | None     | Field name|                                                              |
| operator | string |yes   | None     | Operator| Optional value equal,in|
| value    | -      |no   | None     | Operand| Different values correspond to different value formats                            |

**Note:**
- If this request is to query details, count is 0. If the query is quantity, info is empty.
- The input here`info` only describes the required and built-in parameters for parameters, and the rest of the parameters to be filled in depend on the attribute fields defined by the user
