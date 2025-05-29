### Description

Query all relationships of a certain instance (including its situation as the original model and the target model of the
relationship, Permission: Model Instance Query Permission)

### Parameters

| Name       | Type   | Required | Description                      |
|------------|--------|----------|----------------------------------|
| bk_inst_id | int    | Yes      | Instance ID                      |
| bk_obj_id  | string | Yes      | Model ID                         |
| fields     | array  | Yes      | Fields to be returned            |
| start      | int    | No       | Record start position            |
| limit      | int    | No       | Page size, maximum value is 500. |

#### page

| Name  | Type | Required | Description                                |
|-------|------|----------|--------------------------------------------|
| start | int  | Yes      | Record start position                      |
| limit | int  | Yes      | Number of records per page, maximum is 200 |

### Request Example

```json
{
    "condition": {
        "bk_inst_id": 16,
        "bk_obj_id": "bk_router"
    },
    "fields": [
        "id",
        "bk_inst_id",
        "bk_obj_id",
        "bk_asst_inst_id",
        "bk_asst_obj_id",
        "bk_obj_asst_id",
        "bk_asst_id"
        ],
    "page": {
        "start":0,
        "limit":2
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
    "data": [
        {
            "id": 4,
            "bk_inst_id": 1,
            "bk_obj_id": "bk_switch",
            "bk_asst_inst_id": 16,
            "bk_asst_obj_id": "bk_router",
            "bk_obj_asst_id": "bk_switch_default_bk_router",
            "bk_asst_id": "default"
        },
        {
            "id": 6,
            "bk_inst_id": 2,
            "bk_obj_id": "bk_switch",
            "bk_asst_inst_id": 16,
            "bk_asst_obj_id": "bk_router",
            "bk_obj_asst_id": "bk_switch_default_bk_router",
            "bk_asst_id": "default"
        }
    ]
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

| Name            | Type   | Description                                   |
|-----------------|--------|-----------------------------------------------|
| id              | int64  | Relationship ID                               |
| bk_inst_id      | int64  | Source model instance ID                      |
| bk_obj_id       | string | Source model ID                               |
| bk_asst_inst_id | int64  | Target model instance ID                      |
| bk_asst_obj_id  | string | Target model ID                               |
| bk_obj_asst_id  | string | Automatically generated model relationship ID |
| bk_asst_id      | string | Relationship name                             |
