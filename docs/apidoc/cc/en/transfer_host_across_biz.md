### Functional description

transfer hosts from one business to another business. can only transfer hosts between resource sets(v3.10.27+, permissions: host transferred to other business)

### Request Parameters

#### General Parameters
{{ common_args_desc }}

#### Interface Parameters

| Field         | Type  | Required | Description                                       |
| ------------- | ----- | -------- | ------------------------------------------------- |
| src_bk_biz_id | int   | Yes      | the source business id these hosts belongs to     |
| bk_host_id    | array | Yes      | to be transfered hosts id list, max length is 500 |
| dst_bk_biz_id | int   | Yes      | the target business id                            |
| bk_module_id  | int   | Yes      | the target module id，must be one of idle set's module |

### Request Parameters Example

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

### Return result parameter description
#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | The success or failure of the request. true: the request was successful; false: the request failed.|
| code | int | The error code. 0 means success, >0 means failure error.|
| message | string | The error message returned by the failed request.|
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | data returned by the request|
