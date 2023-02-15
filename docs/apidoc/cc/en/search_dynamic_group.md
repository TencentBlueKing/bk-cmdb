### Functional description

Query dynamic group list (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | yes  | Business ID |
| condition |  object    | no     | Query condition: the condition field is the attribute field of the user-defined query, which can be create_user, modify_user, name|
| disable_counter |  bool |no     | Return total number of records; default|
| page     |   object   | yes  | Paging settings|

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start     |   int     | yes  | Record start position|
| limit     |   int     | yes  | Limit bars per page, Max. 200|
| sort      |   string  |no     | Retrieve sort, by default by creation time|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "disable_counter": true,
    "condition": {
        "name": "my-dynamic-group"
    },
    "page":{
        "start": 0,
        "limit": 200
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 0,
        "info": [
            {
                "bk_biz_id": 1,
                "id": "XXXXXXXX",
                "name": "my-dynamic-group",
                "bk_obj_id": "host",
                "info": {
                    "condition":[
                			{
                				"bk_obj_id":"set",
                				"condition":[
                					{
                						"field":"default",
                						"operator":"$ne",
                						"value":1
                					}
                				]
                			},
                			{
                				"bk_obj_id":"module",
                				"condition":[
                					{
                						"field":"default",
                						"operator":"$ne",
                						"value":1
                					}
                				]
                			},
                			{
                				"bk_obj_id":"host",
                				"condition":[
                					{
                						"field":"bk_host_innerip",
                						"operator":"$eq",
                						"value":"127.0.0.1"
                					}
                				]
                			}
                    ]
                },
                "name": "test",
                "bk_obj_id": "host",
                "id": "1111",
                "create_user": "admin",
                "create_time": "2018-03-27T16:22:43.271+08:00",
                "modify_user": "admin",
                "last_time": "2018-03-27T16:29:26.428+08:00"
            }
        ]
    }
}
```

### Return result parameter
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int |The total number of records that can be matched by the current rule (used for pre-paging by the caller, the actual number of returns from a single request and whether all data are pulled are subject to the number of JSON Array parsing)|
| info      |  array        | Custom query data|

#### data.info

| Field      | Type       | Description      |
|-----------|------------|-----------|
| bk_biz_id    |  int     | Business ID |
| id           |  string  |Dynamic grouping pk ID|
| name         |  string  |Dynamic group naming|
| bk_obj_id    |  string  |The target resource object type of dynamic grouping can be host,set at present|
| info         |  object  |Dynamic grouping information|
| last_time    |  string  |Update time|
| modify_user  | string  |Modifier|
| create_time  | string  |Settling time|
| create_user  | string  |Creator|

#### data.info.info.condition

| Field      | Type     | Description      |
|-----------|-----------|------------|
| bk_obj_id |  string   | Object name, which can be set,module,host|
| condition |  array    | Query criteria|

#### data.info.info.condition.condition

| Field      | Type     | Description      |
|-----------|------------|---------------|
| field     |   string    | Fields of the object|
| operator  |  string    | Operator, op values are eq(equal)/ne(unequal)/in(of)/nin(not of)/like(fuzzy match)|
| value     |   object    | The value for the field|
