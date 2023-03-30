### Function description

batch delete project (version: v3.10.23+, permission: delete permission for project)

### Request parameters

{{ common_args_desc }}


#### Interface Parameters

| field | type | required | description                                                         |
| ----------------------------|------------|----------|---------------------------------------------------------------------|
| ids | array| yes      | an array of ids uniquely identified in cc, limited to 200 at a time |

### Request parameter examples

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "ids":[
        1, 2, 3
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
| result | bool | Whether the request was successful or not. true:request successful; false request failed.|
| code | int | The error code. 0 means success, >0 means failure error.|
| message | string | The error message returned by the failed request.|
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | The data returned by the request.|
