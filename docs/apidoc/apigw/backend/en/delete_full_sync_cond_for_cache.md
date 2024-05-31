### Description

Delete full synchronization cache condition (version: v3.14.1+, permission: delete permission for full sync cache cond)

### Parameters

| Name | Type | Required | Description                                             |
|------|------|----------|---------------------------------------------------------|
| id   | int  | yes      | ID of the full sync cache cond that needs to be deleted |

### Request Example

```json
{
  "id": 123
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

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error    |
| message    | string | Error message returned in case of request failure                |
| permission | object | Permission information                                           |
| data       | object | Data returned in the request                                     |
