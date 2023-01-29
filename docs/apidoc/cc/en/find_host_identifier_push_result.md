### Function description

Get the result from the identity of the push host to the machine (you can only get the tasks pushed within 30 minutes).
(version: v3.10.23+ permission: when the host included in the task is under business, the access permission of the corresponding business is required; when the host is under the host pool, the update permission of the host is required.)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| field | type | required | description |
| ---- | ---- | ---- | ---------- |
| task_id | string | yes | task_id |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "task_id": "GSETASK:F:202201251046313618521052:198"
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "msg": "success",
    "permission": null,
    "request_id": "c11aasdadadadsadasdadasd1111ds",
    "data": {
            "success_list": [
                1,
                2
            ],
            "pending_list": [
                3,
                4
            ],
            "failed_list": [
                5,
                6
            ]
        }
}
```

### Return Result Parameter Description

#### response

| name | type | description |
| ------- | ------ | ------------------------------------------ |
| result | bool | Whether the request was successful or not. true:request successful; false request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | The data returned by the request.

#### data Field Description
| name | type | description |
| ------- | ------ | ------------------------------------------ |
| success_list | array | List of host ids that executed successfully |
| failed_list | array | list of failed host ids |
| pending_list | array |List of host ids for which gse was invoked to send down the host identity and the result is not yet available |
