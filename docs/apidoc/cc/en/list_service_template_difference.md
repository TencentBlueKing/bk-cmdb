### Function Description

List the differences between service templates and service instances (v3.9.19).

- This interface is specifically designed for GSEKit and is in a hidden state in the ESB documentation.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type        | Required | Description                                                  |
| -------------------- | ----------- | -------- | ------------------------------------------------------------ |
| bk_biz_id            | int64       | Yes      | Business ID                                                  |
| bk_module_ids        | int64 array | No       | List of module IDs, up to 20                                 |
| service_template_ids | int64 array | No       | List of service template IDs, up to 20                       |
| is_partial           | bool        | Yes      | When true, use the `service_template_ids` parameter to return the status of service templates; when false, use the `bk_module_ids` parameter to return the status of modules |

### Request Parameter Example

- Example 1

```json
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

```json
{
    "bk_biz_id": 3,
    "bk_module_ids": [
        11,
        12
    ],
    "is_partial": false
}
```

### Response Example

- Example 1

```json
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

```json
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

### Response Result Explanation

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Data returned by the request                                 |

- data Field Explanation

| Field              | Type         | Description                          |
| ----------------- | ------------ | ------------------------------------ |
| service_templates | object array | List of service template information |
| modules           | object array | List of module information           |

- service_templates Field Explanation

| Field                | Type | Description                                                  |
| ------------------- | ---- | ------------------------------------------------------------ |
| service_template_id | int  | Service template ID                                          |
| need_sync           | bool | Whether there are differences between service instances and service templates under the module where the service template is applied |

- modules Field Explanation

| Field         | Type | Description                                                  |
| ------------ | ---- | ------------------------------------------------------------ |
| bk_module_id | int  | Module ID                                                    |
| need_sync    | bool | Whether there are differences between service instances and service templates under the module |