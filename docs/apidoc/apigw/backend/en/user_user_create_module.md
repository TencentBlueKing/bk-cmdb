### Description

Create a module (Permission: Business Topology Creation Permission)

### Parameters

| Name            | Type   | Required | Description                                                                                                                                                                            |
|-----------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id       | int    | Yes      | Business ID                                                                                                                                                                            |
| bk_set_id       | int    | Yes      | Cluster id                                                                                                                                                                             |
| bk_parent_id    | int    | Yes      | The ID of the parent instance node, the last level instance node of the current instance node, in the topology structure, for modules, it generally refers to the bk_set_id of the set |
| bk_module_name  | string | Yes      | Module name                                                                                                                                                                            |
| bk_module_type  | string | No       | Module type                                                                                                                                                                            |
| operator        | string | No       | Main maintainer                                                                                                                                                                        |
| bk_bak_operator | string | No       | Backup maintainer                                                                                                                                                                      |
| bk_created_at   | string | No       | Creation time                                                                                                                                                                          |
| bk_updated_at   | string | No       | Update time                                                                                                                                                                            |
| bk_created_by   | string | No       | Creator                                                                                                                                                                                |
| bk_updated_by   | string | No       | Last updater                                                                                                                                                                           |

**Note: The input parameters here only explain the required and system-built parameters. The rest of the parameters to
be filled in depend on the user's own defined attribute fields. The parameter values are set according to the
configuration of the module's attribute fields.**

### Request Example

```json
{
  "bk_parent_id": 4,
  "bk_module_name": "redis-1",
  "bk_module_type": "2",
  "operator": "admin",
  "bk_bak_operator": "admin",
  "bk_created_at": "",
  "bk_updated_at": "",
  "bk_created_by": "admin",
  "bk_updated_by": "admin"
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
    "bk_bak_operator": "admin",
    "bk_biz_id": 3,
    "bk_created_at": "2023-11-14T17:11:21.225+08:00",
    "bk_created_by": "admin",
    "bk_module_id": 20,
    "bk_module_name": "redis-1",
    "bk_module_type": "2",
    "bk_parent_id": 4,
    "bk_set_id": 4,
    "bk_supplier_account": "0",
    "bk_updated_at": "2023-11-14T17:11:21.225+08:00",
    "create_time": "2023-11-14T17:11:21.225+08:00",
    "default": 0,
    "host_apply_enabled": false,
    "last_time": "2023-11-14T17:11:21.225+08:00",
    "operator": "admin",
    "service_category_id": 2,
    "service_template_id": 0,
    "set_template_id": 0
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |

#### data

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_bak_operator     | string  | Backup maintainer                                          |
| bk_module_id        | int     | Module ID                                                  |
| bk_biz_id           | int     | Business ID                                                |
| bk_module_id        | int     | Module ID                                                  |
| bk_module_name      | string  | Module name                                                |
| bk_module_type      | string  | Module type                                                |
| bk_parent_id        | int     | Parent node ID                                             |
| bk_set_id           | int     | Cluster id                                                 |
| bk_supplier_account | string  | Developer account                                          |
| create_time         | string  | Creation time                                              |
| last_time           | string  | Update time                                                |
| default             | int     | Module type                                                |
| host_apply_enabled  | bool    | Whether to enable automatic application of host attributes |
| operator            | string  | Main maintainer                                            |
| service_category_id | integer | Service category ID                                        |
| service_template_id | int     | Service template ID                                        |
| set_template_id     | int     | Cluster template ID                                        |
| bk_created_at       | string  | Creation time                                              |
| bk_updated_at       | string  | Update time                                                |
| bk_created_by       | string  | Creator                                                    |
| bk_updated_by       | string  | Last updater                                               |
