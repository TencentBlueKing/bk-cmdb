### Description

Update Object Instance (Permission: Model Instance Editing Permission)

- This interface is only applicable to custom hierarchical models and general model instances, not applicable to model
  instances such as business, cluster, module, host, etc.

### Parameters

| Name         | Type   | Required | Description                                                                      |
|--------------|--------|----------|----------------------------------------------------------------------------------|
| bk_obj_id    | string | Yes      | Model ID                                                                         |
| bk_inst_id   | int    | Yes      | Instance ID                                                                      |
| bk_inst_name | string | No       | Instance name, can also be other custom fields                                   |
| bk_biz_id    | int    | No       | Business ID, required when deleting custom mainline hierarchical model instances |

Note: When operating on custom mainline hierarchical model instances, and if using Permission Center, for CMDB versions
less than 3.9, the metadata parameter containing the business id of the instance must be passed, otherwise it will
result in Permission Center authentication failure. The format is "metadata": { "label": { "bk_biz_id": "64" } }

### Request Example

```json
{
    "bk_obj_id": "1",
    "bk_inst_id": 0,
    "bk_inst_name": "test"
 }
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | No data returned                                                    |
