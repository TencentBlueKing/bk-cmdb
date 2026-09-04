### Description

Create Model (Permission: Create Model)

### Parameters

| Name                 | Type         | Required | Description                                                                           |
|----------------------|--------------|----------|---------------------------------------------------------------------------------------|
| creator              | string       | No       | Creator of this data                                                                  |
| bk_classification_id | string       | Yes      | ID of the classification for the object model, can only be named with English letters |
| bk_obj_id            | string       | Yes      | ID of the object model, can only be named with English letters                        |
| bk_obj_name          | string       | Yes      | Name of the object model, used for display, can be in any language readable by humans |
| bk_obj_icon          | string       | No       | ICON information of the object model, used for frontend display                       |
| obj_sort_number      | int          | No       | Sorting order of the object model under the corresponding model group                 |
| bk_labels            | string array | No       | Labels for model categorization and grouping, e.g. ["network"]                        |

### Request Example

```python
{
    "creator": "admin",
    "bk_classification_id": "test",
    "bk_obj_name": "test",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "test",
    "obj_sort_number": 1,
    "bk_labels": ["network"]
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
        "bk_labels": ["network"],
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

#### data

| Name                 | Type               | Description                                         |
|----------------------|--------------------|-----------------------------------------------------|
| id                   | int                | Data record ID                                      |
| creator              | string             | Creator of this data                                |
| modifier             | string             | Last modifier of this data                          |
| create_time          | string             | Creation time                                       |
| last_time            | string             | Update time                                         |
| bk_supplier_account  | string             | Vendor account                                      |
| bk_obj_id            | string             | Object model ID                                     |
| bk_obj_name          | string             | Object model name                                   |
| bk_obj_icon          | string             | ICON information of the object model                |
| position             | json object string | Coordinates used for front-end display              |
| ispre                | bool               | Whether it is predefined, true or false             |
| obj_sort_number      | int                | Sorting order of the object model under the model group |
| bk_labels            | string array       | Labels for model categorization and grouping        |
