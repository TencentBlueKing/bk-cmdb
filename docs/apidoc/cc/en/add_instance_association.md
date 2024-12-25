### Function Description

Create a new association relationship between model instances. (Permission: Model instances editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field           | Type   | Required | Description                                  |
| --------------- | ------ | -------- | -------------------------------------------- |
| bk_obj_asst_id  | string | Yes      | Unique ID of the relationship between models |
| bk_inst_id      | int64  | Yes      | Source model instance ID                     |
| bk_asst_inst_id | int64  | Yes      | Target model instance ID                     |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 11,
    "bk_asst_inst_id": 21
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "id": 1038
    },
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807"
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error |
| message    | string | Error message returned for a failed request                  |
| data       | object | Data returned by the request                                 |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |

#### data Field Description

| Field | Type  | Description                                             |
| ----- | ----- | ------------------------------------------------------- |
| id    | int64 | ID of the newly added instance association relationship |