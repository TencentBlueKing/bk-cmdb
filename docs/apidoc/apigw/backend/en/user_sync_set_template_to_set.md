### Description

Synchronize cluster templates to clusters based on business ID, cluster template ID, and a list of clusters to be
synchronized (Permission: Cluster editing permission)

### Parameters

| Name            | Type  | Required | Description                         |
|-----------------|-------|----------|-------------------------------------|
| bk_biz_id       | int   | Yes      | Business ID                         |
| set_template_id | int   | Yes      | Cluster template ID                 |
| bk_set_ids      | array | Yes      | List of clusters to be synchronized |

### Request Example

```json
{

    "bk_biz_id": 20,
    "set_template_id": 6,
    "bk_set_ids": [46]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
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
