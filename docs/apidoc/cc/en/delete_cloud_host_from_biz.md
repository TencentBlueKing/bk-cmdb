### Function description

delete cloud host from biz idle set (cloud host management dedicated interface, version: v3.10.19+, permission: edit business host)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| field       | type      | mandatory | description                                                                                                         |
|-------------|-----------|-----------|---------------------------------------------------------------------------------------------------------------------|
| bk_biz_id   | int       | yes       | business id                                                                                                         |
| bk_host_ids | array | yes       | to be deleted cloud host ids, array length is limited to 200, these hosts can only succeed or fail at the same time |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 123,
    "bk_host_ids": [
        1,
        2
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807"
}
```

### Return Result Parameters Description

#### response

| name       | type   | description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. true:request successful; false request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                  |
| message    | string | The error message returned by the failed request.                                         |
| permission | object | Permission information                                                                    |
| request_id | string | Request chain id                                                                          |