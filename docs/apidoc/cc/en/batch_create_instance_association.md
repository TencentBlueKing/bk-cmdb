### Function Description

Batch create relationships between common model instances (Version: v3.10.2+, Permission: Editing permissions for source model instances and target model instances)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Parameter      | Type   | Required | Description                                                  |
| -------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_obj_id      | string | Yes      | Source model ID                                              |
| bk_asst_obj_id | string | Yes      | Target model ID                                              |
| bk_obj_asst_id | string | Yes      | Unique ID for the relationship between models                |
| details        | array  | Yes      | Content of batch-created relationship, up to 200 relationships allowed |

#### details

| Parameter       | Type | Required | Description              |
| --------------- | ---- | -------- | ------------------------ |
| bk_inst_id      | int  | Yes      | Source model instance ID |
| bk_asst_inst_id | int  | Yes      | Target model instance ID |

#### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "bk_switch",
    "bk_asst_obj_id": "host",
    "bk_obj_asst_id": "bk_switch_belong_host",
    "details": [
        {
            "bk_inst_id": 11,
            "bk_asst_inst_id": 21
        },
        {
            "bk_inst_id": 12,
            "bk_asst_inst_id": 22
        }
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "success_created": {
            "0": 73
        },
        "error_msg": {
            "1": "Associated instance does not exist"
        }
    }
}
```

### Response Parameters Description

#### response

| Name       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error |
| message    | string | Error message returned for a failed request                  |
| data       | object | Data returned by the request                                 |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |

#### data

| Field           | Type | Description                                                  |
| --------------- | ---- | ------------------------------------------------------------ |
| success_created | map  | Key is the index of the instance relationship in the parameter details array, value is the ID of the successfully created instance relationship |
| error_msg       | map  | Key is the index of the instance relationship in the parameter details array, value is the error message for failure |