### Functional description

Query business in business set (v3.10.12+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_id | int    | yes  | Business set ID|
| filter      |   object  |no     | Business attribute combination query criteria|
| fields      |   array   | no     | Specify the fields to query. The parameter is any attribute of the business. If you do not fill in the field information, the system will return all the fields of the business|
| page        |   object  |yes     | Paging condition|

#### filter

Query criteria. The combination supports AND and OR. Can be nested, up to 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  |yes     | Scope rule for filtering business|


#### rules
The filtering rule is triplet`field`,`operator`,`value`

| Name     | Type   | Required| Default value|  Description                                                  |
| -------- | ------ | ---- | ------ | ------------------------------------------------------------ |
| field    |  string |yes   | None     | Field name|                                                              |
| operator | string |yes   | None     | Operator| Optional value equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between|
| value    | -      |no   | None     | Operand| Different values correspond to different value formats                            |

Assembly rules can be found at: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|
| enable_count |  bool  |yes| Whether to get the flag of the number of query objects|
| sort     |   string |no     | Sort the field. By adding sort in front of the field, for example, sort&#34;: sort field&#34; can indicate descending order by field field|

**Note:**
- `enable_count`If this flag is true, then the request is to get the quantity. The remaining fields must be initialized,start is 0, and limit is: 0, sort is "."

### Request Parameters Example

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_id":2,
    "filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"xxx",
                "operator":"equal",
                "value":"xxx"
            },
            {
                "field":"xxx",
                "operator":"in",
                "value":[
                    "xxx"
                ]
            }
        ]
    },
    "fields":[
        "bk_biz_id",
        "bk_biz_name"
    ],
    "page":{
        "start":0,
        "limit":10,
        "enable_count":false,
        "sort":"bk_biz_id"
    }
}
```

### Details return result example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission":null,
    "data": {
        "count": 0,
        "info": [
            {
                "bk_biz_id": 1,
                "bk_biz_name": "xxx"
            }
        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### Example of return result of query business quantity

```python
{
    "result":true,
    "code":0,
    "message":"",
    "permission":null,
    "data":{
        "count":10,
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
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
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

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| bk_biz_id     |  int       | Business ID |
| bk_biz_name      |  string     | Business name|


**Note:**
- If this request is to query details, count is 0. If the query is quantity, info is empty.
