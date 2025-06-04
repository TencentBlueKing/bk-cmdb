### Description

Batch Update Object Instances (Permission: Model Instance Editing Permission)

### Parameters

| Name    | Type   | Required | Description                                            |
|---------|--------|----------|--------------------------------------------------------|
| datas   | object | Yes      | Fields and values to be updated for instances          |
| inst_id | int    | Yes      | Specific instance for which datas is used for updating |

#### datas

| Name         | Type   | Required | Description                                    |
|--------------|--------|----------|------------------------------------------------|
| bk_inst_name | string | No       | Instance name, can also be other custom fields |

**datas is a map-type object, where the key is the field defined in the model for the instance, and the value is the
value of the field**

### Request Example

```python
{
    "bk_obj_id": "test",
    "update": [
        {
            "datas": {
                "bk_inst_name": "batch_update"
            },
            "inst_id": 46
        }
    ]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": "success"
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
