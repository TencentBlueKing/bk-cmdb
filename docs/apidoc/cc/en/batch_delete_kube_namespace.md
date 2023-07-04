### Function Description

Batch delete namespace (version: v3.10.23+, auth: delete container namespace)

### Request parameters

{{ common_args_desc }}

#### Interface parameters

| field | type | required | description |
|----------------------------|------------|--------|--------------------------------------------|
| bk_biz_id | int| Yes |business id|
| ids | array| Yes |the namespace ids array to be deleted, limited to 200 at a time|

### Example request parameters

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "ids": [
        1
     ]
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
| result | bool | Whether the request was successful or not. true:request successful; false request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | The data returned by the request.
