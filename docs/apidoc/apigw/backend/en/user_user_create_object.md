### Description

Create Model (Permission: Create Model)

### Parameters

| Name                 | Type   | Required | Description                                                                           |
|----------------------|--------|----------|---------------------------------------------------------------------------------------|
| creator              | string | No       | Creator of this data                                                                  |
| bk_classification_id | string | Yes      | ID of the classification for the object model, can only be named with English letters |
| bk_obj_id            | string | Yes      | ID of the object model, can only be named with English letters                        |
| bk_obj_name          | string | Yes      | Name of the object model, used for display, can be in any language readable by humans |
| bk_obj_icon          | string | No       | ICON information of the object model, used for frontend display                       |
| obj_sort_number      | int    | No       | Sorting order of the object model under the corresponding model group                 |

### Request Example

```python
{
    "creator": "admin",
    "bk_classification_id": "test",
    "bk_obj_name": "test",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "test",
    "obj_sort_number": 1
}
```

### Response Example

```python
{
    "code": 0,
    "permission": null,
    "result": true,
    "message": "success",
    "data": {
        "description": "",
        "bk_ishidden": false,
        "bk_classification_id": "test",
        "creator": "admin",
        "bk_obj_name": "test",
        "bk_ispaused": false,
        "last_time": null,
        "bk_obj_id": "test",
        "create_time": null,
        "bk_supplier_account": "0",
        "position": "",
        "bk_obj_icon": "icon-cc-business",
        "modifier": "",
        "id": 2000002118,
        "ispre": false,
        "obj_sort_number": 1
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
