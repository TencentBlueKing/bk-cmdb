### Functional description

list process related info according to condition (v3.9.13)

- only used for GSEKit，is hidden in ESB doc

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field      | Type      | Required | Description                                                  |
| ---------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id  | int64       | Yes      | Business ID                                                  |
|bk_set_ids|int64 array|No|set id array, empty represent anyone|
|bk_module_ids|int64 array|No|set id array, empty represent anyone|
|ids|int64 array|No|set id array, empty represent anyone|
|bk_process_names|string array|No|process name array,empty represent anyone, `bk_process_name，bk_func_id can only use one`|
|bk_func_ids|string array|No|func id array, empty represent anyone, `bk_process_name，bk_func_id can only use one`|
|bk_process_ids|int64 array|No|process id array, empty represent anyone,|
| fields     | array     | No      | process property list, the specified process property feilds will be returned <br>it can speed up the request and reduce the network payload |
| page       | object    | Yes      | page info                                                    |


#### page

| Field | Type | Required | Description                      |
| ----- | ---- | -------- | -------------------------------- |
| start | int  | Yes      | start record                     |
| limit | int  | Yes      | page limit, maximum value is 500 |
| sort  | string | No       | the field for sort, '-' represent decending order, default sorted by bk_process_id |

### Request Parameters Example

```json
{
    "set": {
        "bk_set_ids": [
            11,
            12
        ]
    },
    "module": {
        "bk_module_ids": [
            60,
            61
        ]
    },
    "service_instance": {
        "ids": [
            4,
            5
        ]
    },
    "process": {
        "bk_process_names": [
            "pr1",
            "alias_pr2"
        ],
        "bk_func_ids": [],
        "bk_process_ids": [
            45,
            46,
            47
        ]
    },
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_id",
        "bk_func_name"
    ],
    "page": {
        "start": 0,
        "limit": 100,
        "sort": "bk_process_id"
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "count": 2,
        "info": [
            {
                "set": {
                    "bk_set_id": 11,
                    "bk_set_name": "set1",
                    "bk_set_env": "3"
                },
                "module": {
                    "bk_module_id": 60,
                    "bk_module_name": "mm1"
                },
                "host": {
                    "bk_host_id": 4,
                    "bk_cloud_id": 0,
                    "bk_host_innerip": "192.168.15.22"
                },
                "service_instance": {
                    "id": 4,
                    "name": "192.168.15.22_pr1_3333"
                },
                "process_template": {
                    "id": 48
                },
                "process": {
                    "bk_func_id": "",
                    "bk_func_name": "pr1",
                    "bk_process_id": 45,
                    "bk_process_name": "pr1"
                }
            },
            {
                "set": {
                    "bk_set_id": 11,
                    "bk_set_name": "set1",
                    "bk_set_env": "3"
                },
                "module": {
                    "bk_module_id": 60,
                    "bk_module_name": "mm1"
                },
                "host": {
                    "bk_host_id": 4,
                    "bk_cloud_id": 0,
                    "bk_host_innerip": "192.168.15.22"
                },
                "service_instance": {
                    "id": 4,
                    "name": "192.168.15.22_pr1_3333"
                },
                "process_template": {
                    "id": 49
                },
                "process": {
                    "bk_func_id": "",
                    "bk_func_name": "pr2",
                    "bk_process_id": 46,
                    "bk_process_name": "alias_pr2"
                }
            }
        ]
    }
}
```

### Return Result Parameters Description

#### data

| Field | Type  | Description       |
| ----- | ----- | ----------------- |
| count | int   | the num of record |
| info  | array | process related info         |
|set|object|set info|
|module|object|module info|
|host|object|host info|
|service_instance|object|service_instance info|
|process|object|process info|
