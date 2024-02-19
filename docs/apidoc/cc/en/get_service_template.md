### Function Description

Get service template by service template ID.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type | Required | Description         |
| ------------------- | ---- | -------- | ------------------- |
| service_template_id | int  | Yes      | Service template ID |

### Request Parameter Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "service_template_id": 51
}
```

### Response Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "bk_biz_id": 3,
        "id": 51,
        "name": "mm2",
        "service_category_id": 12,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-05-26T09:46:15.259Z",
        "last_time": "2020-05-26T09:46:15.259Z",
        "bk_supplier_account": "0",
        "host_apply_enabled": false
    }
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

#### data Field Description

| Field               | Type    | Description                                                |
| ------------------- | ------- | ---------------------------------------------------------- |
| bk_biz_id           | int     | Business ID                                                |
| id                  | int     | Service template ID                                        |
| name                | array   | Service template name                                      |
| service_category_id | integer | Service category ID                                        |
| creator             | string  | Creator of the service template                            |
| modifier            | string  | Last modifier of the service template                      |
| create_time         | string  | Creation time                                              |
| last_time           | string  | Last update time                                           |
| bk_supplier_account | string  | Supplier account                                           |
| host_apply_enabled  | bool    | Whether to enable automatic application of host properties |