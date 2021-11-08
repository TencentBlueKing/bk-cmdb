### Functional description

create dynamic group (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type    | Required  | Description                            |
|---------------------|---------|-----------|----------------------------------------|
| bk_biz_id           | int     | Yes       | Business ID                            |
| bk_obj_id           | string  | Yes       | object name, it can be set,host        |
| info                | object  | Yes       | common search query parameters         |
| name                | string  | Yes       | the name of dynamic group              |

#### info.condition

| Field     | Type    | Required | Description                                                                                                                |
|-----------|---------|----------|----------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id | string  | Yes      | object name, it can be set,module,host object type for host type dynamic group, only set object type for set dynamic group |
| condition | array   | Yes      | search condition                                                                                                           |

#### info.condition.condition

| Field     |  Type    | Required  | Description                                                                            |
|-----------|----------|-----------|----------------------------------------------------------------------------------------|
| field     |  string  | Yes       | field of object                                                                        |
| operator  |  string  | Yes       | condition operator, $eq is equal, $ne is not equal, $in is belongs, $nin is not belong |
| value     |  object  | Yes       | the value of field                                                                     |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
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

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": "XXXXXXXX"
    }
}
```

### Return Result Parameters Description

#### data

| Field  | Type    | Description                                                   |
|--------|---------|---------------------------------------------------------------|
| id     | string  | Primary key ID returned when new dynamic group create success |
