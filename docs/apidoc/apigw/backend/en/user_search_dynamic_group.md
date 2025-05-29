### Description

Query dynamic group list (Version: v3.9.6, Permission: Business access permission)

### Parameters

| Name            | Type   | Required | Description                                                                                                                  |
|-----------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id       | int    | Yes      | Business ID                                                                                                                  |
| condition       | object | No       | Query conditions, the condition field is the attribute field for custom queries, which can be create_user, modify_user, name |
| disable_counter | bool   | No       | Whether to not return the total number of records, default is to return                                                      |
| page            | object | Yes      | Paging settings                                                                                                              |

#### page

| Name  | Type   | Required | Description                                            |
|-------|--------|----------|--------------------------------------------------------|
| start | int    | Yes      | Record start position                                  |
| limit | int    | Yes      | Number of restrictions per page, maximum is 200        |
| sort  | string | No       | Retrieval sorting, default is to sort by creation time |

### Request Example

```json
{
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

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name  | Type  | Description                                                                                                                                                                                                          |
|-------|-------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| count | int   | Total number of records that the current rule can match (used for the caller to pre-pagination, the actual number of returns and whether the data is completely pulled as JSON Array parsing quantity is subject to) |
| info  | array | Custom query data                                                                                                                                                                                                    |

#### data.info

| Name        | Type   | Description                                                    |
|-------------|--------|----------------------------------------------------------------|
| bk_biz_id   | int    | Business ID                                                    |
| id          | string | Dynamic group primary key ID                                   |
| name        | string | Dynamic group naming                                           |
| bk_obj_id   | string | Target resource object type of dynamic group, can be host, set |
| info        | object | Dynamic group information                                      |
| last_time   | string | Update time                                                    |
| modify_user | string | Modifier                                                       |
| create_time | string | Creation time                                                  |
| create_user | string | Creator                                                        |

#### data.info.info.condition

| Name      | Type   | Description                           |
|-----------|--------|---------------------------------------|
| bk_obj_id | string | Object name, can be set, module, host |
| condition | array  | Query condition                       |

#### data.info.info.condition.condition

| Name     | Type   | Description                                                                                                         |
|----------|--------|---------------------------------------------------------------------------------------------------------------------|
| field    | string | Object field                                                                                                        |
| operator | string | Operator, op value is eq (equal) / ne (not equal) / in (belongs to) / nin (does not belong to) / like (fuzzy match) |
| value    | object | Value corresponding to the field                                                                                    |
