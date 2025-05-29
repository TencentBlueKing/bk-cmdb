### Description

Host Transfer to Module within Business (Permission: Service Instance Edit Permission)

### Parameters

| Name                | Type   | Required | Description                                                                                                                                 |
|---------------------|--------|----------|---------------------------------------------------------------------------------------------------------------------------------------------|
| bk_supplier_account | string | No       | Developer account                                                                                                                           |
| bk_biz_id           | int    | Yes      | Business ID                                                                                                                                 |
| bk_host_id          | array  | Yes      | Host ID                                                                                                                                     |
| bk_module_id        | array  | Yes      | Module ID                                                                                                                                   |
| is_increment        | bool   | No       | Whether to cover or append, will delete the original relationship. True is to append, false is to cover, not filling in is default to false |

### Request Example

```json
{
    "bk_biz_id": 1,
    "bk_host_id": [
        9,
        10
    ],
    "bk_module_id": [
        10
    ],
    "is_increment": true
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
