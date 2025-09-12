### Description

Update Module (Permission: Business Topology Editing Permission)

### URL Parameters

| Name                | Type   | Required | Description       |
|---------------------|--------|----------|-------------------|
| bk_supplier_account | string | No       | Developer account |
| bk_biz_id           | int    | Yes      | Business ID       |
| bk_set_id           | int    | Yes      | Cluster ID        |
| bk_module_id        | int    | Yes      | Module ID         |

#### Parameters

| Name            | Type   | Required | Description       |
|-----------------|--------|----------|-------------------|
| bk_module_name  | string | No       | Module name       |
| bk_module_type  | string | No       | Module type       |
| operator        | string | No       | Main maintainer   |
| bk_bak_operator | string | No       | Backup maintainer |

**Note: The parameter here only explains the system-built editable parameters, and the rest of the parameters to be
filled depend on the user's own defined attribute fields. Modules created through service templates can only be modified
through service templates.**

### Request Example

```json
{
    "bk_module_name": "test",
    "bk_module_type": "1",
    "operator": "admin",
    "bk_bak_operator": "admin"
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
