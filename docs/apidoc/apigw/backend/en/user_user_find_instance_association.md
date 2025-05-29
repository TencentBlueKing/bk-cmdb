### Description

Query the instance association relationship of the model. (Permission: Model instance query permission)

### Parameters

| Name      | Type   | Required | Description              |
|-----------|--------|----------|--------------------------|
| condition | object | Yes      | Query conditions         |
| bk_obj_id | string | Yes      | Source model id (v3.10+) |

#### condition

| Name           | Type   | Required | Description                                     |
|----------------|--------|----------|-------------------------------------------------|
| bk_obj_asst_id | string | Yes      | Unique id of the model association relationship |
| bk_asst_id     | string | No       | Unique id of the association type               |
| bk_asst_obj_id | string | No       | Target model id                                 |

### Request Example

```json
{
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_asst_obj_id": ""
    },
    "bk_obj_id": "xxx"
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [{
        "id": 481,
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
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

#### data Field Explanation:

| Name            | Type   | Description                                                |
|-----------------|--------|------------------------------------------------------------|
| id              | int    | The association's unique id                                |
| bk_obj_asst_id  | string | Automatically generated model association relationship id. |
| bk_obj_id       | string | Source model id of the association relationship            |
| bk_asst_obj_id  | string | Target model id of the association relationship            |
| bk_inst_id      | int    | Source model instance id                                   |
| bk_asst_inst_id | int    | Target model instance id                                   |
