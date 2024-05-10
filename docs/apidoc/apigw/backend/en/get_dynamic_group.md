### Description

Get details of a dynamic group (Version: v3.9.6, Permission: Business access permission)

### Parameters

| Name      | Type   | Required | Description                         |
|-----------|--------|----------|-------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                         |
| id        | string | Yes      | Target dynamic group primary key ID |

### Request Example

```json
{
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
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
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### Explanation of data Parameters

| Name        | Type   | Description                                                                    |
|-------------|--------|--------------------------------------------------------------------------------|
| bk_biz_id   | int    | Business ID                                                                    |
| id          | string | Dynamic group primary key ID                                                   |
| bk_obj_id   | string | Target resource object type of dynamic group, which can be host or set for now |
| name        | string | Dynamic group naming                                                           |
| info        | object | Dynamic group rule information                                                 |
| last_time   | string | Update time                                                                    |
| modify_user | string | Modifier                                                                       |
| create_time | string | Creation time                                                                  |
| create_user | string | Creator                                                                        |

#### Explanation of info Parameters

| Name      | Type  | Description      |
|-----------|-------|------------------|
| condition | array | Query conditions |

#### Explanation of condition Parameters

| Name      | Type   | Description                                                                                                                                                              |
|-----------|--------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id | string | Condition object resource type, the dynamic group of the host type supports info.conditon:set,module,host; the dynamic group of the set type supports info.condition:set |
| condition | array  | Query conditions                                                                                                                                                         |

#### Explanation of condition.condition Parameters

| Name     | Type   | Description                                                                                  |
|----------|--------|----------------------------------------------------------------------------------------------|
| field    | string | Object field                                                                                 |
| operator | string | Operator, op value can be eq (equal)/ne (not equal)/in (belongs to)/nin (does not belong to) |
| value    | object | Value corresponding to the field                                                             |
