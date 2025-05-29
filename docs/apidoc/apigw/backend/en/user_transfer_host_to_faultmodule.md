### Description

Transfer Hosts to the Fault Module within Business (Permission: Service Instance Edit Permission)

### Parameters

| Name                | Type   | Required | Description       |
|---------------------|--------|----------|-------------------|
| bk_supplier_account | string | No       | Developer account |
| bk_biz_id           | int    | Yes      | Business ID       |
| bk_host_id          | array  | Yes      | Host ID           |

### Request Example

```json
{
    "bk_biz_id": 1,
    "bk_host_id": [
        9,
        10
    ]
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
