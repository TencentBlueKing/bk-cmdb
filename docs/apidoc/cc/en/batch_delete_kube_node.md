### Function Description

Delete container node (v3.12.1+, Permission: Container node deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type  | Required | Description                            |
| --------- | ----- | -------- | -------------------------------------- |
| bk_biz_id | int   | Yes      | Business ID of the container node      |
| ids       | array | Yes      | List of IDs of the nodes to be deleted |

**Note:**

- Users need to ensure that there are no associated resources (such as pods) under the nodes to be deleted, otherwise, deletion will fail.
- The number of nodes to be deleted in one go should not exceed 100.

### Request Parameter Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "ids": [
    1,
    2
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
  "data": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Response Parameter Description

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| data       | object | No data returned                                             |
| request_id | string | Request chain ID                                             |