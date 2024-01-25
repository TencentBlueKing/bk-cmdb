### Function Description

Delete project (Version: v3.10.23+, Permission: Project deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type  | Required | Description                                                  |
| ----- | ----- | -------- | ------------------------------------------------------------ |
| ids   | array | Yes      | Array of unique identifiers of projects in CC, with a limit of 200 at a time |

### Request Parameter Example

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
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
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Response Parameter Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned in the request                                 |
