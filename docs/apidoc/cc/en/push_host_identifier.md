### Function Description

Push host identity information to the machines and return the task ID of this GSE task. You can use this GSE task ID to query the results of the push task in GSE. (v3.10.18+, for hosts in a business, business access permission is required, for hosts in a host pool, host pool host editing permission is required)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field       | Type  | Required | Description                  |
| ----------- | ----- | -------- | ---------------------------- |
| bk_host_ids | array | Yes      | Array of host IDs, up to 200 |

### Request Parameter Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_ids": [1, 2]
}
```

### Response Example

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

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Data returned by the request                                 |

#### data

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| task_id    | string | Task ID, this ID is the task_id on the GSE side              |
| host_infos | array  | Host information pushed in the task, only contains information of successfully pushed hosts |

#### host_infos

| Field           | Type   | Description                            |
| -------------- | ------ | -------------------------------------- |
| bk_host_id     | int    | Host ID                                |
| identification | string | Identification of the host in the task |