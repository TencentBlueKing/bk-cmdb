### Description

Create a dynamic group (Version: v3.9.6+, Permission: Dynamic Group Creation Permission)

### Parameters

| Name      | Type   | Required | Description                                                                  |
|-----------|--------|----------|------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                  |
| bk_obj_id | string | Yes      | Target resource object type of the dynamic group, currently can be host, set |
| info      | object | Yes      | Common query conditions                                                      |
| name      | string | Yes      | Dynamic group name                                                           |

#### info.condition

| Name      | Type   | Required | Description                                                                                                                                             |
|-----------|--------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id | string | Yes      | Condition object resource type, host type dynamic group supports info.condition: set, module, host; set type dynamic group supports info.condition: set |
| condition | array  | Yes      | Query conditions                                                                                                                                        |

#### info.condition.condition

| Name     | Type   | Required | Description                                                                                        |
|----------|--------|----------|----------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Object field                                                                                       |
| operator | string | Yes      | Operator, op value can be eq (equal) / ne (not equal) / in (belongs to) / nin (does not belong to) |
| value    | object | Yes      | Value of the field                                                                                 |

### Request Example

```json
{
    "bk_biz_id": 1,
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
    "message": "success",
    "permission": null,
    "data": {
        "id": "XXXXXXXX"
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Request return data                                                         |

#### data

| Name | Type   | Description                                                                   |
|------|--------|-------------------------------------------------------------------------------|
| id   | string | Newly created dynamic group primary key ID returned after successful creation |
