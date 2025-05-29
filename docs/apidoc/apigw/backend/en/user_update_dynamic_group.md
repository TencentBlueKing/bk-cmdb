### Description

Update dynamic group (Version: v3.9.6, Permission: Dynamic group editing permission)

### Parameters

| Name      | Type   | Required | Description                                                                                                                              |
|-----------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                                              |
| id        | string | Yes      | Primary key ID                                                                                                                           |
| bk_obj_id | string | No       | Target resource object type of dynamic group, can be host, set. When updating rules, both this field and the info field must be provided |
| info      | object | No       | General query conditions                                                                                                                 |
| name      | string | No       | Dynamic group name                                                                                                                       |

#### info.condition

| Name      | Type   | Required | Description                                                                                                                                           |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id | string | Yes      | Type of condition object resource, host type dynamic group supports info.conditon:set,module,host; set type dynamic group supports info.condition:set |
| condition | array  | Yes      | Query condition                                                                                                                                       |

#### info.condition.condition

| Name     | Type   | Required | Description                                                                                        |
|----------|--------|----------|----------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Object field                                                                                       |
| operator | string | Yes      | Operator, op value can be eq (equal) / ne (not equal) / in (belongs to) / nin (does not belong to) |
| value    | object | Yes      | Field corresponding value                                                                          |

### Request Example

```json
{
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

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {}
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
