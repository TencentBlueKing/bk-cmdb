### Function Description

Get details of a dynamic group (Version: v3.9.6, Permission: Business access permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                         |
| --------- | ------ | -------- | ----------------------------------- |
| bk_biz_id | int    | Yes      | Business ID                         |
| id        | string | Yes      | Target dynamic group primary key ID |

### Request Parameter Example

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

### Response Example

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
    		],
            "variable_condition":[
              {
                "bk_obj_id":"set",
                "condition":[
                  {
                    "field":"bk_parent_id",
                    "operator":"$ne",
                    "value":1
                  }
                ]
              },
              {
                "bk_obj_id":"module",
                "condition":[
                  {
                    "field":"bk_parent_id",
                    "operator":"$ne",
                    "value":1
                  }
                ]
              },
              {
                "bk_obj_id":"host",
                "condition":[
                  {
                    "field":"bk_host_outerip",
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

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |

#### data

| Field       | Type   | Description                                                  |
| ----------- | ------ | ------------------------------------------------------------ |
| bk_biz_id   | int    | Business ID                                                  |
| id          | string | Dynamic group primary key ID                                 |
| bk_obj_id   | string | Target resource object type of dynamic group, which can be host or set for now |
| name        | string | Dynamic group naming                                         |
| info        | object | Dynamic group rule information                               |
| last_time   | string | Update time                                                  |
| modify_user | string | Modifier                                                     |
| create_time | string | Creation time                                                |
| create_user | string | Creator                                                      |

#### data.info
| Field     | Type   | Description                    |
|-----------|-------|-------------------------|
| condition | object   | dynamic group locking condition |
| variable_condition | object | dynamic group variable condition  |

#### data.info.condition

| Field     | Type   | Description                                                  |
| --------- | ------ | ------------------------------------------------------------ |
| bk_obj_id | string | Condition object resource type, the dynamic group of the host type supports info.conditon:set,module,host; the dynamic group of the set type supports info.condition:set |
| condition | array  | Query conditions                                             |

#### data.info.condition.condition

| Field    | Type   | Description                                                                                      |
| -------- | ------ |--------------------------------------------------------------------------------------------------|
| field    | string | Object field                                                                                     |
| operator | string | Operator, op value can be $eq (equal)/$ne (not equal)/$in (belongs to)/$nin (does not belong to)/ $regex (fuzzy match) |
| value    | object | Value corresponding to the field                                                                 |

#### data.info.variable_condition

| Field     | Type   | Description                           |
| --------- | ------ | ------------------------------------- |
| bk_obj_id | string | Object name, can be set, module, host |
| condition | array  | Query condition                       |

#### data.info.variable_condition.condition

| Field    | Type   | Description                                                                                                              |
| -------- | ------ |--------------------------------------------------------------------------------------------------------------------------|
| field    | string | Object field                                                                                                             |
| operator | string | Operator, op value is $eq (equal) / $ne (not equal) / $in (belongs to) / $nin (does not belong to) / $regex (fuzzy match) |
| value    | object | Value corresponding to the field                                                                                         |
