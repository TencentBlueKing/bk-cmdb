### Description

Update Project ID. This interface is dedicated to BCS for project data migration. Other platforms should not use it (
Version: v3.10.23+, Permission: Project Update Permission)

### Parameters

| Name          | Type   | Required | Description                                                |
|---------------|--------|----------|------------------------------------------------------------|
| id            | int    | Yes      | Unique identifier of the project in cc                     |
| bk_project_id | string | Yes      | The final value that needs to be updated for bk_project_id |

### Request Example

```json
{
    "id": 1,
    "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2"
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

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Data returned by the request                                        |
