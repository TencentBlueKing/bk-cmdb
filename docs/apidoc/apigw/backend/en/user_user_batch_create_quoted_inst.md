### Description

Batch create quoted model instance (Version: v3.10.30+, permission: Model instance editing permission)

### Parameters

| Name           | Type         | Required | Description                                                |
|----------------|--------------|----------|------------------------------------------------------------|
| bk_obj_id      | string       | Yes      | Source model ID                                            |
| bk_property_id | string       | Yes      | Property ID of the source model referencing this model     |
| data           | object array | Yes      | Information of instances to be created, up to 50 instances |

#### data[n]

| Name        | Type   | Required                                                | Description                                                                                                                                                    |
|-------------|--------|---------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_inst_id  | int64  | No                                                      | Source model instance ID, if not filled in, it needs to be associated with the source model instance through the interface for creating source model instances |
| name        | string | Depends on the "Required" configuration in the property | Name, this is just an example, the actual field depends on the model property                                                                                  |
| operator    | string | Depends on the "Required" configuration in the property | Maintainer, this is just an example, the actual field depends on the model property                                                                            |
| description | string | Depends on the "Required" configuration in the property | Description, this is just an example, the actual field depends on the model property                                                                           |

### Request Example

```json
{
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "data": [
    {
      "bk_inst_id": 123,
      "name": "test",
      "operator": "user",
      "description": "test instance"
    }
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
    "ids": [
      1,
      2
    ]
  },
}
```

**Note:**

- The order of the IDs array in the returned data is consistent with the order of the array data in the parameters.

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned for the request                                     |

#### data

| Name | Type        | Description                   |
|------|-------------|-------------------------------|
| ids  | int64 array | Unique identifier array in cc |
