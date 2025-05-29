### Description

Delete object instances under a specified model based on model ID and instance ID (Permission: Model instance deletion
permission)

- This interface is only applicable to custom hierarchy models and general model instances, not applicable to business,
  cluster, module, host, and other model instances.

### Parameters

| Name       | Type   | Required | Description                                                                       |
|------------|--------|----------|-----------------------------------------------------------------------------------|
| bk_obj_id  | string | Yes      | Model ID                                                                          |
| bk_inst_id | int    | Yes      | Instance ID                                                                       |
| bk_biz_id  | int    | No       | Business ID, required when deleting instances of custom mainline hierarchy models |

Note: When operating on instances of custom mainline hierarchy models and using permission center, for CMDB versions
less than 3.9, you also need to pass the metadata parameter containing the business ID of the instance, otherwise
permission center authentication will fail. The format is "metadata": { "label": { "bk_biz_id": "64" } }

### Request Example

```json
{ 
    "bk_obj_id": "test",
    "bk_inst_id": 44
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

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
