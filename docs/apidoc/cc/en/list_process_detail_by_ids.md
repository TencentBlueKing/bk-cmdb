### Functional description

Query the process details corresponding to the process ID under a business (v3.9.8)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|Field| Type| Required| Description|
|---|---|---|---|
|bk_biz_id| int| yes | The business ID of the process|
|bk_process_ids| array| yes | Process ID list, up to 500|
|fields| array| No| Process attribute list, which controls which fields are in the process instance information that returns the result, can speed up interface requests and reduce network traffic transmission<br> If blank, all fields of the process are returned, and bk_process_id is a required field|


### Request Parameters Example

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":1,
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
``` json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

| Name| Type| Description|
|---|---|--- |
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | array |Data returned by request|

#### data
| Name| Type| Description|
|---|---|--- |
|bk_func_id| string| Function ID|
|bk_func_name| string| Process name|
|bk_process_id| int| Process id|
|bk_process_name| string| Process alias|