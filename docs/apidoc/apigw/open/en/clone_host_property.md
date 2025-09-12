### Description

Clone Host Properties (Permission: Business Host Editing Permission)

### Parameters

| Name        | Type   | Required | Description               |
|-------------|--------|----------|---------------------------|
| bk_org_ip   | string | Yes      | Source host's internal IP |
| bk_dst_ip   | string | Yes      | Target host's internal IP |
| bk_org_id   | int    | Yes      | Source host ID            |
| bk_dst_id   | int    | Yes      | Target host ID            |
| bk_biz_id   | int    | Yes      | Business ID               |
| bk_cloud_id | int    | Yes      | Control area ID           |

Note: Cloning using the internal IP of the host and cloning using the identity ID of the host, these two methods can
only be used separately and cannot be mixed.

### Request Example

```json
{
    "bk_biz_id": 2,
    "bk_org_ip": "127.0.0.1",
    "bk_dst_ip": "127.0.0.2",
    "bk_cloud_id": 0
}
```

or

```json
{
    "bk_biz_id": 2,
    "bk_org_id": 10,
    "bk_dst_id": 11,
    "bk_cloud_id": 0
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
