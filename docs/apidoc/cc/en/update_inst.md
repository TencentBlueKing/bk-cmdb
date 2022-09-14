### Functional description

Update object instance

- This interface is only applicable to user-defined hierarchical model and general model instances, not applicable to model instances such as business, set, module, host, etc

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                            |
|---------------------|------------|--------|----------------------------------|
| bk_obj_id           |  string     | yes  | Model ID       |
| bk_inst_id          |  int        | yes  | Instance ID|
| bk_inst_name        |  string     | no     | Instance name, or any other custom field   |
| bk_biz_id                  |  int        | no     | Business ID, which must be transferred when deleting a user-defined mainline level model instance|

 Note: If the operation is a user-defined mainline hierarchy model instance and permission Center is used, for the version with CMDB less than 3.9, the metadata parameter containing the service id of the instance needs to be transferred. Otherwise, the permission Center authentication will fail. The format is
"metadata": {
    "label": {
        "bk_biz_id": "64"
    }
}

### Request Parameters Example (generic example)

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "bk_obj_id": "1",
    "bk_inst_id": 0,
    "bk_inst_name": "test"
 }
```

### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |No data return|
