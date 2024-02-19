### Function Description

Update Business Information (Permission: Business Edit Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type   | Required | Description       |
| ------------------- | ------ | -------- | ----------------- |
| bk_supplier_account | string | No       | Developer account |
| bk_biz_id           | int    | Yes      | Business ID       |
| data                | dict   | Yes      | Business data     |

#### data

| Field             | Type   | Required | Description                                             |
| ----------------- | ------ | -------- | ------------------------------------------------------- |
| bk_biz_name       | string | No       | Business name                                           |
| bk_biz_developer  | string | No       | Developer                                               |
| bk_biz_maintainer | string | No       | Maintainer                                              |
| bk_biz_productor  | string | No       | Productor                                               |
| bk_biz_tester     | string | No       | Tester                                                  |
| operator          | string | No       | Operator                                                |
| life_cycle        | string | No       | Life cycle: Testing(1), Online(2, default), Shutdown(3) |
| language          | string | No       | Language, "1" for Chinese, "2" for English              |

**Note: The data parameter here only explains the system-built editable parameters, and the rest of the parameters to be filled depend on the user's own defined attribute fields.**

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 1,
    "data": {
        "bk_biz_name": "cc_app_test",
        "bk_biz_maintainer": "admin",
        "bk_biz_productor": "admin",
        "bk_biz_developer": "admin",
        "bk_biz_tester": "admin",
        "language": "1",
        "operator": "admin",
        "life_cycle": "2"
    }
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
    "data": null
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure        |
| message    | string | Error message returned in case of request failure            |
| data       | object | Request returned data                                        |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |