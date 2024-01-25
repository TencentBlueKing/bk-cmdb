### Function Description

Transfer hosts across businesses. You can only transfer hosts from the source business's idle host pool cluster to the target business's idle host pool cluster (Version: v3.10.27+, Permission: Transfer hosts to another business)

### Request Parameters

#### Common Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field         | Type  | Required | Description                                                  |
| ------------- | ----- | -------- | ------------------------------------------------------------ |
| src_bk_biz_id | int   | Yes      | The business ID to which the hosts to be transferred belong  |
| bk_host_id    | array | Yes      | List of host IDs to be transferred, with a maximum length of 500 |
| dst_bk_biz_id | int   | Yes      | The business ID to which the hosts will be transferred       |
| bk_module_id  | int   | Yes      | The module ID to which the hosts will be transferred. This module ID must be under the idle host pool set. |

### Request Parameter Example

```json
{
    "bk_app_code": "xxx",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "src_bk_biz_id": 2,
    "dst_bk_biz_id": 3,
    "bk_host_id": [
        9,
        10
    ],
    "bk_module_id": 10
}
```

### Return Result Example

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

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |