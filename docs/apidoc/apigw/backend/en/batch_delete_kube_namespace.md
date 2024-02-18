### Description

Batch delete namespace (Version: v3.12.1+, Permission: Container namespace deletion permission)

### Parameters

| Name      | Type  | Required | Description                                                                         |
|-----------|-------|----------|-------------------------------------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                                                         |
| ids       | array | Yes      | Unique identifiers of namespaces to be deleted in CC, with a limit of 200 at a time |

### Request Example

```json
{
  "bk_biz_id": 3,
  "ids": [
    1
  ]
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "data": null,
  "message": "success",
  "permission": null,
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
