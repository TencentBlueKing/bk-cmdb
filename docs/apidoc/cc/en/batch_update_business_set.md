### Functional description

Update business set information (v3.10.12+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_ids | array  |yes| Business set ID list|
| data           |  object |Yes.| Business set data|

#### data

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_attr |  object  |no     | Business set model attribute |
| bk_scope  |  object  |no     | Selected business scope|

#### bk_biz_set_attr

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_name   |   string  |yes     | Business set name|
| bk_biz_maintainer |  string  |no     | Operation and maintenance personnel|
| bk_biz_set_desc   |   string  |no     | Business set description|

#### bk_scope

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| match_all |  bool  |yes     | Selected business scope tag|
| filter    |   object| no     | Scope criteria for the selected business|

#### filter

This parameter is a combination of filtering rules for business attribute fields, and is used to search for hosts according to host attribute fields. The combination only supports AND operation and can be nested, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  |yes     | Scope condition rule for selected business|


#### rules

| Name     | Type   | Required| Default value  |  Description                                                  |
| -------- | ------ | ---- | ------  | ------------------------------------------------------------ |
| field    |  string |yes   | None     | Field name|                                                              |
| operator | string |yes   | None     | Operator| Optional value equal,in|
| value    | -      |no   | None     | Operand| Different values correspond to different value formats                            |


**Note:**
- The input parameters here only describe the required and built-in parameters, and the rest of the parameters to be filled in depend on the attribute fields defined by the user
- The and fields are not allowed to change for batch scenarios (number of IDs in bk_biz_set_ids is greater than 1`bk_biz_set_name``bk_scope`
- The maximum number of batch updates is 200.

### Request Parameters Example

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_ids":[
        2
    ],
    "data":{
        "bk_biz_set_attr":{
            "bk_biz_set_name": "test",
            "bk_biz_set_desc":"xxx",
            "biz_set_maintainer":"xxx"
        },
        "bk_scope":{
            "match_all":false,
            "filter":{
                "condition":"AND",
                "rules":[
                    {
                        "field":"bk_sla",
                        "operator":"equal",
                        "value":"2"
                    }
                ]
            }
        }
    }
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission":null,
    "data": {},
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
