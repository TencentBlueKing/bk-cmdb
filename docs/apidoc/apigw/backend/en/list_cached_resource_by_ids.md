### Description

List cached resource details by ID list (version: v3.14.1+, permission: general cache query permission)

### Parameters

| Name         | Type         | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
|--------------|--------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| resource     | string       | Yes      | The resource type to be queried. Enumeration: host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project. Among them, host represents host, biz represents biz, set represents set, module represents module, process represents process, object_instance represents common model instance, mainline_instance represents mainline model instance, biz_set represents business set, plat represents cloud area, project represents project. |
| sub_resource | string       | No       | The subordinate resource type to be queried. It needs to be specified when resource is object_instance or mainline_instance, which represents bk_obj_id of the model that needs to be synchronized                                                                                                                                                                                                                                                                          |
| ids          | int array    | yes      | ID list to be queried, up to 500                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| fields       | string array | No       | Return field list, controls which fields are returned                                                                                                                                                                                                                                                                                                                                                                                                                       |

### Request Example

```json
{
   "resource": "object_instance",
   "sub_resource": "bk_switch",
   "ids": [
     123,
     456
   ],
   "fields": [
     "bk_asset_id",
     "bk_inst_id",
     "bk_inst_name",
     "bk_obj_id"
   ]
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "info": [
      {
        "bk_asset_id": "sw00001",
        "bk_inst_id": 1,
        "bk_inst_name": "sw1",
        "bk_obj_id": "bk_switch"
      },
      {
        "bk_asset_id": "sw00002",
        "bk_inst_id": 2,
        "bk_inst_name": "sw2",
        "bk_obj_id": "bk_switch"
      }
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error    |
| message    | string | Error message returned in case of request failure                |
| permission | object | Permission information                                           |
| data       | object | Data returned in the request                                     |

#### data

| Name | Type  | Description                  |
|------|-------|------------------------------|
| info | array | Resource caching detail list |

#### data.info

| Name         | Type   | Description   |
|--------------|--------|---------------|
| bk_asset_id  | string | Asset ID      |
| bk_inst_id   | int    | Instance ID   |
| bk_inst_name | string | Instance name |
| bk_obj_id    | string | Model ID      |

**Note: The return value here only uses the scenario of listing some fields of switch as an example to illustrate its
attribute fields. The specific return value depends on the resource type and user-defined attribute fields**
