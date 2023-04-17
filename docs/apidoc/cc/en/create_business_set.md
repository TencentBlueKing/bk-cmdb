### Functional description

New business set (v3.10.12+)

### Request Parameters

{{ common_args_desc }}


#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_attr    |   object  |yes     | Business set model attribute |
| bk_scope |  object  |yes     | Selected business scope|

#### bk_biz_set_attr

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_name   |   string  |yes   | Business set name|
| bk_biz_maintainer |  string  |no   | Operation and maintenance personnel|
| bk_biz_set_desc   |   string  |no   | Business set description|

#### bk_scope

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| match_all |  bool  |yes    | Selected business scope tag|
| filter |  object  |no     | Scope criteria for the selected business|

#### filter

This parameter is a combination of filtering rules for the service attribute field, and is used to search for services according to the service attribute field. The combination only supports AND operation and can be nested, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  |yes     | Scope condition rule for selected business|


#### rules

| Name     | Type   | Required| Default value| Description   |  Description                                                  |
| -------- | ------ | ---- | ------ | ------ | ------------------------------------------------------------ |
| field    |  string |yes   | None     | Field name|                                                              |
| operator | string |yes   | None     | Operator| Optional value equal,in|
| value    | -      |no   | None     | Operand| Different values correspond to different value formats                            |

**Note:**
- The input here`bk_biz_set_attr` only describes the required and built-in parameters for parameters, and the rest of the parameters to be filled in depend on the attribute fields defined by the user
- `bk_scope`If the field in`match_all` is true, it means that the selected business range of the business set is all. In this case, the parameter`filter` is blank. If the`match_all` field is false`filter`, it needs to be non-empty, and the user needs to explicitly point to
Scope of business selection
- The circled type of the selected business attribute in the business set is organization and enumeration
### Request Parameters Example

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_attr":{
        "bk_biz_set_name":"biz_set",
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
                    "value":"3"
                },
                {
                    "field":"life_cycle",
                    "operator":"equal",
                    "value":1
                }
            ]
        }
    }
}
```

### Return Result Example

```python

{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":5,
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
| data    |  int |Business set id created                           |
| request_id    |  string |Request chain id    |
