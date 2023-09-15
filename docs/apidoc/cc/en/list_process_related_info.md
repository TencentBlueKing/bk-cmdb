### Functional description

Click Five-digit query process instance related information (v3.9.13)

- This interface is intended for use by GSEKit and is hidden in the ESB documentation

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|Field| Type| Required| Description|
|---|---|---|---|
| bk_biz_id  | int64       |  yes   | Business ID |
|bk_set_ids| int64 array| no | The set ID list, if empty, represents any set |
|bk_module_ids| int64 array| no | Module ID list, if empty, it represents any module|
|ids| int64 array| no | Service instance ID list. If empty, it represents any instance|
|bk_process_names| string array| no | List of process names. If empty, it represents any process. This field is mutually exclusive with bk_func_id. Only one of them can be selected `it can not have value at the same time`|
|bk_func_ids| string array| no | Function ID list of process. If empty, it represents any process. bk_process_name `only one of the two can be selected, and can not have value at the same time|
|bk_process_ids| int64 array| no | Process ID list, if empty, represents any process|
|fields| string array| no | Process attribute list, which controls which fields are in the process instance information that returns the result, can speed up interface requests and reduce network traffic transmission<br> If it is null, all fields of the process are returned, and bk_process_id,bk_process_name and bk_func_id are required fields to be returned|
|page| dict| yes | Paging condition|

The conditional relationship for these fields is relationship and (&amp;&amp;), and only process instances that meet the criteria you fill in are queried<br>
For example, if both bk_set_ids and bk_module_ids are filled in, and neither bk_module_ids belongs to bk_set_ids, the query result is empty

#### page

| Field| Type| Required| Description|
| ---  | ---  | ---  | --- |
| start| int| no | Record start position, default is 0|
| limit| int| yes | Limit bars per page, Max. 500|
| sort  | string |no   | Sort field,'backward' means reverse order, can only be the field of the process, and sort by bk_process_id by default|


### Request Parameters Example

``` json
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
``` json
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
                    "bk_host_innerip_v6":"1::1",
                    "bk_addressing":"dynamic",
                    "bk_agent_id":"xxxxxx"
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

| Name| Type| Description|
|---|---|--- |
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|

- Data field Description

| Name| Type| Description|
|---|---|--- |
|count| int| Total number of eligible process instances|
|set| object| Set information to which the process belongs |
|module| object| Module information to which the process belongs|
|host| object| Host information to which the process belongs|
|service_instance| object| Service instance information to which the process belongs|
|process| object| Details of the process itself|
