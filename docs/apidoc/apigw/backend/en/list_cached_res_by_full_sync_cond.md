### Description

List cached resource details by full sync cache cond (version: v3.14.1+, permission: general cache query permission)

### Parameters

| Name    | Type         | Required | Description                                                                                                     |
|---------|--------------|----------|-----------------------------------------------------------------------------------------------------------------|
| cond_id | int          | Yes      | The full sync cache cond ID                                                                                     |
| cursor  | int          | Yes      | The starting ID, returns the resource details with resource IDs greater than the starting ID in ascending order |
| limit   | int          | Yes      | The paging limit, maximum 500                                                                                   |
| fields  | string array | No       | Return field list, controls which fields are returned                                                           |

### Request Example

```json
{
  "cond_id": 111,
  "cursor": 222,
  "limit": 10,
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

**Note: The return value here only uses the scenario of listing some fields of switch as an example to illustrate its attribute fields. The specific return value depends on the resource type and user-defined attribute fields**
