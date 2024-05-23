### Description

Create a control area based on the control area name (Permission: Control Area Creation Permission)

### Parameters

| Name          | Type   | Required | Description       |
|---------------|--------|----------|-------------------|
| bk_cloud_name | string | Yes      | Control area name |

### Request Example

```python
{
    
    "bk_cloud_name": "test1"
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "created": {
            "origin_index": 0,
            "id": 6
        }
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Request return data                                                         |

#### data

| Name    | Type   | Description                              |
|---------|--------|------------------------------------------|
| created | object | Created successfully, return information |

#### data.created

| Name         | Type | Description                                       |
|--------------|------|---------------------------------------------------|
| origin_index | int  | Corresponding to the order of the request results |
| id           | int  | Control area id, bk_cloud_id                      |
