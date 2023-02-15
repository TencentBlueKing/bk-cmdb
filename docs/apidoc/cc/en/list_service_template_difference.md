### Functional description

List the differences between service templates and service instances (v3.9.19)

- This interface is intended for use by GSEKit and is hidden in the ESB documentation

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|Field| Type| Required| Description|
|---|---|---|---|
| bk_biz_id  | int64       |  yes   | Business ID |
|bk_module_ids| int64 array| no | Module ID list, no more than 20|
|service_template_ids| int64 array| no | List of service template IDs, up to 20|
|is_partial| bool| yes | If true, use service_template_ids parameter to return the state of service_template; When false, returns the status of the module using the bk_module_ids parameter|


### Request Parameters Example

- Example 1
``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "service_template_ids": [
        1,
        2
    ],
    "is_partial": true
}
```
- Example 2
```
{
    "bk_biz_id": 3,
    "bk_module_ids": [
        11,
        12
    ],
    "is_partial": false
}
```

### Return Result Example
- Example 1
``` json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "service_templates": [
            {
                "service_template_id": 1,
                "need_sync": true
            },
            {
                "service_template_id": 2,
                "need_sync": false
            }
        ]
    }
}
```
- Example 2
```
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "modules": [
            {
                "bk_module_id": 11,
                "need_sync": false
            },
            {
                "bk_module_id": 12,
                "need_sync": true
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
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

- Data field Description

| Name| Type| Description|
|---|---|--- |
|service_templates| object array| Service template Info list|
|modules| object array| Module info list|

- Service_templates Field Description

| Name| Type| Description|
|---|---|--- |
|service_template_id| int| Service template ID|
|need_sync| bool| Is there any difference between the service instance and the service template under the module to which the service template applies|

- Modules Field Description

| Name| Type| Description|
|---|---|--- |
|bk_module_id| int| Module ID|
|need_sync| bool| Is there any difference between service instance and service template under module|
