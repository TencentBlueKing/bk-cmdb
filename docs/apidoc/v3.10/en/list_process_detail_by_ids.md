### Functional description

list process detail by their ids in one biz (v3.9.8)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field      | Type      | Required | Description                                                  |
| ---------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id  | int64       | Yes      | Business ID                                                  |
|bk_process_ids|int64 array|Yes|process id array,its max length is 500|
| fields     | array     | No      | process property list, the specified process property feilds will be returned <br>it can speed up the request and reduce the network payload |


### Request Parameters Example

```json
{
    "bk_process_ids": [
        43,
        44
    ],
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_id",
        "bk_func_name"
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "bk_func_id": "",
            "bk_func_name": "pr1",
            "bk_process_id": 43,
            "bk_process_name": "pr1"
        },
        {
            "bk_func_id": "",
            "bk_func_name": "pr2",
            "bk_process_id": 44,
            "bk_process_name": "pr2"
        }
    ]
}
```

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |