### Description

Create full synchronization cache condition (version: v3.14.1+, permission: creation permission for full sync cache
cond)

### Parameters

| Name         | Type   | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
|--------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| resource     | string | Yes      | The resource type that needs to be cached for full sync data. Enumeration: host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project. Among them, host represents host, biz represents biz, set represents set, module represents module, process represents process, object_instance represents common model instance, mainline_instance represents mainline model instance, biz_set represents business set, plat represents cloud area, project represents project. |
| sub_resource | string | No       | The subordinate resource type. It needs to be specified when resource is object_instance or mainline_instance, which represents bk_obj_id of the model that needs to be synchronized                                                                                                                                                                                                                                                                                                                     |
| is_all       | bool   | No       | Whether to synchronize all data, each resource can have only one condition whose is_all is true                                                                                                                                                                                                                                                                                                                                                                                                          |
| condition    | object | No       | Used to specify sync condition when is_all is false. Its format can be referred to: https://github.com/TencentBlueKing/bk-cmdb/blob/master/pkg/filter/README.md                                                                                                                                                                                                                                                                                                                                          |
| interval     | int    | Yes      | Sync period, in hours, used to specify the cache expiration time, the minimum is 6 hours, the maximum is 7 days                                                                                                                                                                                                                                                                                                                                                                                          |

**Notice:**

- The maximum number of customized conditions is 100
- Only resources that have created corresponding full sync cache cond will be fully cached.

### Request Example

```json
{
  "resource": "object_instance",
  "sub_resource": "bk_switch",
  "is_all": true,
  "interval": 24
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
    "id": 123
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

| Name | Type | Description                            |
|------|------|----------------------------------------|
| id   | int  | ID of the created full sync cache cond |
