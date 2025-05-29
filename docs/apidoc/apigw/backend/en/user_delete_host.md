### Description

Delete Host (Permission: Host Pool Host Deletion Permission)

### Parameters

| Name                | Type   | Required | Description                              |
|---------------------|--------|----------|------------------------------------------|
| bk_supplier_account | string | No       | Developer account                        |
| bk_host_id          | string | Yes      | Host ID, separated by commas if multiple |

### Request Example

```json
{
    "bk_host_id": "1,2,3"
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

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |
