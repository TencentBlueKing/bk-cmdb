### Description

Batch Update Instances of Referenced Models (Version: v3.10.30+, Permission: Edit Permission of Source Model Instances)

### Parameters

| Name           | Type        | Required | Description                                                 |
|----------------|-------------|----------|-------------------------------------------------------------|
| bk_obj_id      | string      | Yes      | Source model ID                                             |
| bk_property_id | string      | Yes      | Property ID of the source model that references this model  |
| ids            | int64 array | Yes      | Array of instance IDs to be updated, up to a maximum of 500 |
| data           | object      | Yes      | Information of instances to be updated                      |

#### data

| Name        | Type   | Required                                     | Description                                                                            |
|-------------|--------|----------------------------------------------|----------------------------------------------------------------------------------------|
| name        | string | At least one field in data must be filled in | Name, this is just an example, the actual fields depend on the model properties        |
| operator    | string | At least one field in data must be filled in | Operator, this is just an example, the actual fields depend on the model properties    |
| description | string | At least one field in data must be filled in | Description, this is just an example, the actual fields depend on the model properties |

### Request Example

```json
{
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "ids": [
    1,
    2
  ],
  "data": {
    "name": "test",
    "operator": "user",
    "description": "test instance"
  }
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
