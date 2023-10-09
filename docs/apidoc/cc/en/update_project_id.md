### Function description

update the project id. This API is a dedicated interface for BCS to migrate project data. Do not use it on other platforms (version: v3.10.23+, permissions: update permissions of the project)

### Request parameters

{{ common_args_desc }}


#### Interface parameters

| field         | type   | required | description                                         |
|---------------|--------|----------|-----------------------------------------------------|
| id            | int    | yes      | The unique identification of the project's id in cc |
| bk_project_id | string | yes      | The final value of bk_project_id to be updated      |


### Request parameter examples
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
