### Functional description

Get dynamic grouping details (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | yes  | Business ID |
| id        |   string  |yes     | Target dynamic grouping pk ID|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
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
    	"bk_biz_id": 1,
    	"name": "my-dynamic-group",
    	"id": "XXXXXXXX",
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
       "create_user": "admin",
       "create_time": "2018-03-27T16:22:43.271+08:00",
       "modify_user": "admin",
       "last_time": "2018-03-27T16:29:26.428+08:00"
    },
    "permission": null,
    "request_id": "dsda1122adasadadada2222",
}
```

### Return result parameter
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| bk_biz_id    |  int     | Business ID |
| id           |  string  |Dynamic grouping pk ID|
| bk_obj_id    |  string  |Target resource object type of dynamic grouping, which can be host,set at present|
| name         |  string  |Dynamic group naming|
| info         |  object  |Dynamic grouping rule information|
| last_time    |  string  |Update time|
| modify_user  | string  |Modifier|
| create_time  | string  |Settling time|
| create_user  | string  |Creator|

#### data.info.condition

| Field      | Type     | Description      |
|-----------|-----------|------------|
| bk_obj_id |  string   | Conditional object resource type, info.conditon supported for dynamic grouping of host type: set,module,host; Info.conditions supported for dynamic grouping of type set: set|
| condition |  array    | Query criteria|

#### data.info.condition.condition

| Field      | Type     | Description       |
|-----------|------------|------------|
| field     |   string    | Fields of the object|
| operator  |  string    | Operator with op values eq(equal)/ne(unequal)/in(of)/nin(not of)|
| value     |   object    | The value corresponding to the field|
