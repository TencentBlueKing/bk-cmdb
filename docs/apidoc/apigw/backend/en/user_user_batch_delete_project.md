### Description

Delete project (Version: v3.10.23+, Permission: Project deletion permission)

### Parameters

| Name | Type  | Required | Description                                                                  |
|------|-------|----------|------------------------------------------------------------------------------|
| ids  | array | Yes      | Array of unique identifiers of projects in CC, with a limit of 200 at a time |

### Request Example

```json
{
    "ids":[
        1, 2, 3
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
