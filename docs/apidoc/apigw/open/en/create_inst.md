### Description

Create an Instance (Permission: Model Instance Creation Permission)

- This interface is only applicable to custom hierarchical models and general model instances, not applicable to
  business, set, module, host, and other model instances.

### Parameters

| Name         | Type   | Required | Description                                                                                                                |
|--------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id    | string | Yes      | Model ID                                                                                                                   |
| bk_inst_name | string | Yes      | Instance name                                                                                                              |
| bk_biz_id    | int    | No       | Business ID, required when creating custom mainline hierarchical model instances                                           |
| bk_parent_id | int    | No       | Required when creating custom mainline hierarchical model instances, represents the ID of the parent hierarchical instance |

Note: When operating on custom mainline hierarchical model instances and using permission center, for CMDB versions less
than 3.9, you also need to pass the metadata parameter containing the business ID of the instance in the metadata
parameter; otherwise, it will cause permission center authentication failure. The format is:

"metadata": { "label": { "bk_biz_id": "64" } }

Other fields belonging to instance properties can also be passed as parameters. For table-type attributes, the value is
a list of IDs of associated instances of the table-type model (needs to be created first using the
batch_create_quoted_inst interface), with a maximum of 50.

### Request Example

```json
{
    "bk_obj_id":"test3",
    "bk_inst_name":"example18",
    "bk_biz_id":0
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "data": {
        "bk_biz_id": 0,
        "bk_inst_id": 1177099,
        "bk_inst_name": "example18",
        "bk_obj_id": "test3",
        "bk_supplier_account": "0",
        "create_time": "2022-01-05T17:28:27.069+08:00",
        "last_time": "2022-01-05T17:28:27.069+08:00",
        "test4": ""
    },
    "message": "success",
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Request return data                                                         |

#### data

| Name                | Type   | Description                                    |
|---------------------|--------|------------------------------------------------|
| bk_inst_id          | int    | Instance ID returned after successful creation |
| bk_biz_id           | int    | Business ID                                    |
| bk_inst_name        | string | Instance name                                  |
| bk_obj_id           | string | Model ID                                       |
| bk_supplier_account | string | Supplier account                               |
| create_time         | string | Creation time                                  |
| last_time           | string | Update time                                    |
