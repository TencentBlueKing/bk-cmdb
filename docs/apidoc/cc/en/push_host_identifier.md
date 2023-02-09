### Function description

push the host identity to the machine and return the gse task id of this task，according to this task, id can go to gse to query the push result of the task.(v3.10.18+, for machines in business, business access is required, and for machines in host pool, host pool host editing permission is required)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| field | type | required | description |
| ---- | ---- | ---- | ---------- |
| bk_host_ids | array | Yes | Array of host ids, the number cannot exceed 200 |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_ids": [1,2]
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
        "task_id": "GSETASK:F:202206222053523618521052:393",
        "host_infos": [
            {
                "bk_host_id": 2,
                "identification": "0:127.0.0.1"
            }
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
| task_id | string | task_id，this id is the task_id from the gse  |

#### host_infos Field Description
| name    | type   | description                                       |
| ------- | ------ | ------------------------------------------ |
|  bk_host_id |  int  |host id |
|  identification |  string  |the identity of the pushed host in the task |
