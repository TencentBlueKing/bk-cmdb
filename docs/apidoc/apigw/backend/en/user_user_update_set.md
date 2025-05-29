### Description

Update Cluster (Permission: Business Topology Editing Permission)

### Parameters

| Name              | Type   | Required | Description                                                                      |
|-------------------|--------|----------|----------------------------------------------------------------------------------|
| bk_biz_id         | int    | Yes      | Business ID                                                                      |
| bk_set_id         | int    | Yes      | Cluster ID                                                                       |
| bk_set_name       | string | No       | Cluster name                                                                     |
| default           | int    | No       | 0 - Normal cluster, 1 - Built-in module set, default is 0                        |
| set_template_id   | int    | No       | Cluster template ID, required when creating a cluster through a cluster template |
| bk_capacity       | int    | No       | Design capacity                                                                  |
| description       | string | No       | Remarks, description of data                                                     |
| bk_set_desc       | string | No       | Cluster description                                                              |
| bk_set_env        | string | No       | Environment type: Test (1), Experience (2), Formal (3, default)                  |
| bk_service_status | string | No       | Service status: Open (1, default), Close (2)                                     |

**Note: The input parameter here only explains the system-built editable parameters, and the rest of the parameters to be
filled depend on the user's own defined attribute fields. Clusters created through cluster templates can only be
modified through cluster templates.**

### Request Example

```json
{
  "bk_set_name": "test",
  "default": 0,
  "bk_capacity": 500,
  "bk_set_desc": "Cluster description",
  "description": "Cluster remarks",
  "bk_set_env": "3",
  "bk_service_status": "1"
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
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |
