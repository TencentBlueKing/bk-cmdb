### Functional description

query dynamic group (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type    | Required | Description                            |
|---------------------|---------|----------|----------------------------------------|
| bk_biz_id           | int     | Yes      | Business ID                            |
| id                  | string  | Yes      | Primary key ID of target dynamic group |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
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
    }
}
```

### Return Result Parameters Description

#### data

| Field        | Type    | Description                     |
|--------------|---------|---------------------------------|
| bk_biz_id    | int     | Business ID                     |
| id           | string  | Primary key ID of dynamic group |
| bk_obj_id    | string  | object name, it can be set,host |
| name         | string  | the name of dynamic group       |
| info         | object  | common search query parameters  |
| last_time    | string  | last update timestamp           |
| modify_user  | string  | last modify user                |
| create_time  | string  | create timestamp                |
| create_user  | string  | creator                         |

#### data.info.condition

| Field     | Type    | Description                                                                                                                |
|-----------|---------|----------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id |  string | object name, it can be set,module,host object type for host type dynamic group, only set object type for set dynamic group |
| condition |  array  | search condition                                                                                                           |

#### data.info.condition.condition

| Field     | Type    | Description                                                                            |
|-----------|---------|----------------------------------------------------------------------------------------|
| field     | string  | field of object                                                                        |
| operator  | string  | condition operator, $eq is equal, $ne is not equal, $in is belongs, $nin is not belong |
| value     | object  | the value of field                                                                     |
