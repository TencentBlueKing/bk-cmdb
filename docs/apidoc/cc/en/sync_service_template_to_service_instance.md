### Functional description

Synchronize the service template information to the corresponding service instance(v3.12.3+, permission: service instance create, edit, and delete permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type  | Required | Description                              |
|---------------------|-------|----------|------------------------------------------|
| bk_biz_id           | int   | yes      | Business ID                              |
| service_template_id | int   | yes      | Service template ID                      |
| bk_module_ids       | array | yes      | ID list of the module to be synchronized |

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "service_template_id": 1,
  "bk_module_ids": [
    28
  ]
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": null
}
```

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | object | Data returned by request                                                                |
