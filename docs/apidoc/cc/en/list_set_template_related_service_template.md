### Function Description

Retrieve a list of service templates for a specified business and cluster template ID.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field           | Type | Required | Description         |
| --------------- | ---- | -------- | ------------------- |
| set_template_id | int  | Yes      | Cluster template ID |
| bk_biz_id       | int  | Yes      | Business ID         |

### Request Parameter Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "set_template_id": 1,
  "bk_biz_id": 3
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
    "data": [
        {
            "bk_biz_id": 3,
            "id": 48,
            "name": "sm1",
            "service_category_id": 2,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:14:57.691Z",
            "last_time": "2020-05-15T14:14:57.691Z",
            "host_apply_enabled": false
        },
        {
            "bk_biz_id": 3,
            "id": 49,
            "name": "sm2",
            "": 16,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:19:09.813Z",
            "last_time": "2020-05-15T14:19:09.813Z",
            "host_apply_enabled": false
        }
    ]
}
```

### Response Result Explanation

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | array  | Data returned by the request                                 |

#### data

| Field               | Type   | Description                                                |
| ------------------- | ------ | ---------------------------------------------------------- |
| bk_biz_id           | int    | Business ID                                                |
| id                  | int    | Service template ID                                        |
| name                | string | Service template name                                      |
| service_category_id | int    | Service category ID                                        |
| creator             | string | Creator                                                    |
| modifier            | string | Last modifier                                              |
| create_time         | string | Creation time                                              |
| last_time           | string | Update time                                                |
| host_apply_enabled  | bool   | Whether to enable automatic application of host properties |