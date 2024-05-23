### Function Description

Update dynamic group (Version: v3.9.6, Permission: Dynamic group editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_id | int    | Yes      | Business ID                                                  |
| id        | string | Yes      | Primary key ID                                               |
| bk_obj_id | string | No       | Target resource object type of dynamic group, can be host, set. When updating rules, both this field and the info field must be provided |
| info      | object | No       | General query conditions                                     |
| name      | string | No       | Dynamic group name                                           |

#### info
| Field     | Type   | Required | Description                                                                        |
|-----------|--------|----------|------------------------------------------------------------------------------------|
| condition | object | No       | dynamic group locking condition, which is at least one from the variable condition |
| variable_condition | object | No | dynamic group variable condition, which is at least one from the locking condition  |

#### info.condition

| Field     | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| bk_obj_id | string | Yes      | Type of condition object resource, host type dynamic group supports info.conditon:set,module,host; set type dynamic group supports info.condition:set |
| condition | array  | Yes      | Query condition                                              |

#### info.condition.condition

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | Yes      | Object field                                                 |
| operator | string | Yes      | Operator, op value can be $eq (equal) / $ne (not equal) / $in (belongs to) / $nin (does not belong to)/ $regex (fuzzy match) |
| value    | object | Yes      | Field corresponding value                                    |

#### info.variable_condition

| Field     | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| bk_obj_id | string | Yes      | Condition object resource type, host type dynamic group supports info.condition: set, module, host; set type dynamic group supports info.condition: set |
| condition | array  | Yes      | Query conditions                                             |

#### info.variable_condition.condition

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | Yes      | Object field                                                 |
| operator | string | Yes      | Operator, op value can be $eq (equal) / $ne (not equal) / $in (belongs to) / $nin (does not belong to)/ $regex (fuzzy match) |
| value    | object | Yes      | Value of the field                                           |

### Request Parameter Example

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

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Data returned by the request                                 |