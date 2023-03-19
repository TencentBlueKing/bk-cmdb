### Functional description

Update dynamic grouping (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | yes  | Business ID |
| id        |   string  |yes     | Primary key ID|
| bk_obj_id |  string  |no     | The target resource object type of dynamic grouping can be host, set at present. When updating rules, both this field and info field shall be provided.|
| info      |   object  |no     | General query criteria|
| name      |   string  |no     | Dynamic group name|

#### info.condition

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_obj_id |  string   | yes  | Conditional object resource type, info.conditon supported for dynamic grouping of host type: set,module,host; Info.conditions supported for dynamic grouping of type set: set|
| condition |  array    | yes  | Query criteria|

#### info.condition.condition

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| field     |   string    | yes  | The fields of the object|
| operator  |  string    | yes  | Operator with op values eq(equal)/ne(unequal)/in(of)/nin(not of)|
| value     |   object    | yes  | The value corresponding to the field|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "id": "XXXXXXXX",
    "bk_obj_id": "host",
    "name": "my-dynamic-group",
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
    "data": {}
}
```

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |
