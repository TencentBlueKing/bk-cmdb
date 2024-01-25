### Function Description

Update Project ID. This interface is dedicated to BCS for project data migration. Other platforms should not use it (Version: v3.10.23+, Permission: Project Update Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field         | Type   | Required | Description                                                |
| ------------- | ------ | -------- | ---------------------------------------------------------- |
| id            | int    | Yes      | Unique identifier of the project in cc                     |
| bk_project_id | string | Yes      | The final value that needs to be updated for bk_project_id |

### Request Parameters Example

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
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
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure        |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |