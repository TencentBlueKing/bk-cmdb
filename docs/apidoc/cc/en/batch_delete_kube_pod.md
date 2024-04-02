### Function Description

Batch delete pods (Version: v3.12.1+, Permission: Container pod deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type         | Required | Description                                                  |
| ----- | ------------ | -------- | ------------------------------------------------------------ |
| data  | object array | Yes      | Array of pod information to be deleted, with a maximum of 200 pods in total in the 'data' array |

#### data

| Field     | Type      | Required | Description                                                  |
| --------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id | int       | Yes      | Business ID                                                  |
| ids       | int array | Yes      | Array of cc IDs of pods to be deleted, with a maximum of 200 IDs in total in the 'data' array |

### Request Parameter Example

```json
{
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
  "data": [
    {
      "bk_biz_id": 123,
      "ids": [
        5,
        6
      ]
    }
  ]
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "",
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