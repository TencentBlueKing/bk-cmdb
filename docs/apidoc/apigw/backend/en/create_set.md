### Description

Create a cluster (Permission: Business Topology Creation Permission)

### Parameters

| Name              | Type   | Required | Description                                                                                                                                                                |
|-------------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id         | int    | Yes      | Business ID                                                                                                                                                                |
| bk_parent_id      | int    | Yes      | The ID of the parent instance node, the last level instance node of the current instance node, in the topology structure, for sets, it generally refers to the business ID |
| bk_set_name       | string | Yes      | Cluster name                                                                                                                                                               |
| default           | int    | No       | 0-ordinary cluster, 1-built-in module collection, default is 0                                                                                                             |
| set_template_id   | int    | No       | Cluster template ID, required when creating a cluster through a cluster template                                                                                           |
| bk_capacity       | int    | No       | Design capacity                                                                                                                                                            |
| description       | string | No       | Remark, description of the data                                                                                                                                            |
| bk_set_desc       | string | No       | Cluster description                                                                                                                                                        |
| bk_set_env        | string | No       | Environment type: test(1), experience(2), formal(3, default)                                                                                                               |
| bk_service_status | string | No       | Service status: open(1, default), close(2)                                                                                                                                 |
| bk_created_at     | string | No       | Creation time                                                                                                                                                              |
| bk_updated_at     | string | No       | Update time                                                                                                                                                                |
| bk_created_by     | string | No       | Creator                                                                                                                                                                    |
| bk_updated_by     | string | No       | Last updater                                                                                                                                                               |

**Note: The input parameters here only explain the required and system-built parameters. The rest of the parameters to
be filled in depend on the user's own defined attribute fields. The parameter values are set according to the
configuration of the cluster's attribute fields.**

### Request Example

```json
{
  "bk_parent_id": 3,
  "bk_set_name": "set_a1",
  "set_template_id": 0,
  "default": 0,
  "bk_capacity": 1000,
  "bk_set_desc": "test-set",
  "bk_set_env": "1",
  "bk_service_status": "1",
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
    "bk_biz_id": 3,
    "bk_capacity": 1000,
    "bk_created_at": "2023-11-14T17:30:43.048+08:00",
    "bk_created_by": "admin",
    "bk_parent_id": 3,
    "bk_service_status": "1",
    "bk_set_desc": "test-set",
    "bk_set_env": "1",
    "bk_set_id": 10,
    "bk_set_name": "set_a1",
    "bk_supplier_account": "0",
    "bk_updated_at": "2023-11-14T17:30:43.048+08:00",
    "create_time": "2023-11-14T17:30:43.048+08:00",
    "default": 0,
    "description": "",
    "last_time": "2023-11-14T17:30:43.048+08:00",
    "set_template_id": 0,
    "set_template_version": null
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

| Name                 | Type   | Description                                                    |
|----------------------|--------|----------------------------------------------------------------|
| bk_biz_id            | int    | Business ID                                                    |
| bk_capacity          | int    | Design capacity                                                |
| bk_parent_id         | int    | Parent node ID                                                 |
| bk_set_id            | int    | Cluster ID                                                     |
| bk_service_status    | string | Service status: 1/2 (1: open, 2: close)                        |
| bk_set_desc          | string | Cluster description                                            |
| bk_set_env           | string | Environment type: 1/2/3 (1: test, 2: experience, 3: formal)    |
| bk_set_name          | string | Cluster name                                                   |
| create_time          | string | Creation time                                                  |
| last_time            | string | Update time                                                    |
| bk_supplier_account  | string | Developer account                                              |
| default              | int    | 0-ordinary cluster, 1-built-in module collection, default is 0 |
| description          | string | Data description information                                   |
| set_template_version | array  | Current version of the cluster template                        |
| set_template_id      | int    | Cluster template ID                                            |
| bk_created_at        | string | Creation time                                                  |
| bk_updated_at        | string | Update time                                                    |
| bk_created_by        | string | Creator                                                        |
| bk_updated_by        | string | Last updater                                                   |
