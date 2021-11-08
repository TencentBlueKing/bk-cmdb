### Functional description

list difference between the service template and service instances (v3.9.19)

- only used for GSEKit，is hidden in ESB doc

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field      | Type      | Required | Description                                                  |
| ---------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id  | int64       | Yes      | Business ID                                                  |
|bk_module_ids|int64 array|No|Module ID arrary, the max length is 20|
|service_template_ids|int64 array|No|Service template ID arrary, the max length is 20|
|is_partial|bool|Yes|If true, take service_template_ids as request variables and return service template status. Otherwise, take bk_module_ids as request variables and return module status |

### Request Parameters Example

- Example 1
``` json
{
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

#### data

| Field | Type  | Description       |
| ----- | ----- | ----------------- |
|service_templates|object array|Sevice template info array|
|modules|object array|Module info array|

#### service_templates

| Field | Type  | Description       |
| ----- | ----- | ----------------- |
|service_template_id|int|Service template ID|
|need_sync|bool|Whether the service instances associated with the service template has any difference with the service template|

#### modules

| Field | Type  | Description       |
| ----- | ----- | ----------------- |
|bk_module_id|int|Module ID|
|need_sync|bool|Whether the service instances under the module has any difference with service template|

