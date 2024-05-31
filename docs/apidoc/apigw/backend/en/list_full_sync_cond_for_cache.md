### Description

List full synchronization cache conditions (version: v3.14.1+, permission: Query permission for full sync cache cond)

### Parameters

| Name         | Type      | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|--------------|-----------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| resource     | string    | Yes      | The resource type to be queried. Enumeration: host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project. Among them, host represents host, biz represents biz, set represents set, module represents module, process represents process, object_instance represents common model instance, mainline_instance represents mainline model instance, biz_set represents business set, plat represents cloud area, project represents project. One of resource and ids must be set. |
| sub_resource | string    | No       | The subordinate resource type to be queried. It needs to be specified when resource is object_instance or mainline_instance, which represents bk_obj_id of the model that needs to be synchronized                                                                                                                                                                                                                                                                                                               |
| ids          | int array | yes      | ID list to be queried, up to 500. One of resource and ids must be set.                                                                                                                                                                                                                                                                                                                                                                                                                                           |

### Request Example

```json
{
   "ids": [
     123,
     456
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
         "id": 123,
         "resource": "object_instance",
         "sub_resource": "bk_switch",
         "is_all": true,
         "interval": 24
       },
       {
         "id": 456,
         "resource": "host",
         "is_all": false,
         "interval": 6,
         "condition": {
           "condition": "AND",
           "rules": [
             {
               "field": "bk_host_innerip",
               "operator": "not_equal",
               "value": "127.0.0.1"
             },
             {
               "condition": "OR",
               "rules": [
                 {
                   "field": "bk_os_type",
                   "operator": "in",
                   "value": [
                     "3"
                   ]
                 },
                 {
                   "field": "bk_cloud_id",
                   "operator": "equal",
                   "value": 0
                 }
               ]
             }
           ]
         }
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

| Name | Type  | Description                        |
|------|-------|------------------------------------|
| info | array | The full sync cache cond data list |

#### info[x]

| Name         | Type   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
|--------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id           | int    | The auto-incremented ID of full sync cache cond                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| resource     | string | The resource type that needs to be cached for full sync data. Enumeration: host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project. Among them, host represents host, biz represents biz, set represents set, module represents module, process represents process, object_instance represents common model instance, mainline_instance represents mainline model instance, biz_set represents business set, plat represents cloud area, project represents project. |
| sub_resource | string | The subordinate resource type. It needs to be specified when resource is object_instance or mainline_instance, which represents bk_obj_id of the model that needs to be synchronized                                                                                                                                                                                                                                                                                                                     |
| is_all       | bool   | Whether to synchronize all data, each resource can have only one condition whose is_all is true                                                                                                                                                                                                                                                                                                                                                                                                          |
| condition    | object | Used to specify sync condition when is_all is false. Its format can be referred to: https://github.com/TencentBlueKing/bk-cmdb/blob/master/pkg/filter/README.md                                                                                                                                                                                                                                                                                                                                          |
| interval     | int    | Sync period, in hours, used to specify the cache expiration time, the minimum is 6 hours, the maximum is 7 days                                                                                                                                                                                                                                                                                                                                                                                          |
