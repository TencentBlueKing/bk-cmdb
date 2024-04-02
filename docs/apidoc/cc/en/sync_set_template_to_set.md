### Function Description

Synchronize cluster templates to clusters based on business ID, cluster template ID, and a list of clusters to be synchronized (Permission: Cluster editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field           | Type  | Required | Description                         |
| --------------- | ----- | -------- | ----------------------------------- |
| bk_biz_id       | int   | Yes      | Business ID                         |
| set_template_id | int   | Yes      | Cluster template ID                 |
| bk_set_ids      | array | Yes      | List of clusters to be synchronized |

### Request Parameter Example

```json
{

    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 20,
    "set_template_id": 6,
    "bk_set_ids": [46]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": null
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |