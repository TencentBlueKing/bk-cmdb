### Function Description

Get the business topology of the mainline model.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type | Required | Description |
| ----- | ---- | -------- | ----------- |
|       |      |          |             |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx"
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_obj_id": "biz",
            "bk_obj_name": "Business",
            "bk_next_obj": "set",
            "bk_next_name": "Set",
            "bk_pre_obj_id": "",
            "bk_pre_obj_name": ""
        },
        {
            "bk_obj_id": "set",
            "bk_obj_name": "Set",
            "bk_next_obj": "module",
            "bk_next_name": "Module",
            "bk_pre_obj_id": "biz",
            "bk_pre_obj_name": "Business"
        },
        {
            "bk_obj_id": "module",
            "bk_obj_name": "Module",
            "bk_next_obj": "host",
            "bk_next_name": "Host",
            "bk_pre_obj_id": "set",
            "bk_pre_obj_name": "Set"
        },
        {
            "bk_obj_id": "host",
            "bk_obj_name": "Host",
            "bk_next_obj": "",
            "bk_next_name": "",
            "bk_pre_obj_id": "module",
            "bk_pre_obj_name": "Module"
        }
    ]
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |

#### data

| Field               | Type   | Description                                           |
| ------------------- | ------ | ----------------------------------------------------- |
| bk_obj_id           | string | Unique ID of the model                                |
| bk_obj_name         | string | Model name                                            |
| bk_next_obj         | string | Unique ID of the next model for the current model     |
| bk_next_name        | string | Name of the next model for the current model          |
| bk_pre_obj_id       | string | Unique ID of the previous model for the current model |
| bk_pre_obj_name     | string | Name of the previous model for the current model      |