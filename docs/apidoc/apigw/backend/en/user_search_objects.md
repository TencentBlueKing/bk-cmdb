### Description

Query models based on optional conditions (Permission: Model View Permission)

### Parameters

| Name                 | Type   | Required | Description                                                                         |
|----------------------|--------|----------|-------------------------------------------------------------------------------------|
| creator              | string | No       | Creator of this data                                                                |
| modifier             | string | No       | Last modifier of this data                                                          |
| bk_classification_id | string | No       | Classification ID of the object model, can only be named with alphabetical sequence |
| bk_obj_id            | string | No       | ID of the object model, can only be named with alphabetical sequence                |
| bk_obj_name          | string | No       | Name of the object model, used for display, can be any language readable by humans  |
| obj_sort_number      | int    | No       | Sorting order of the object model under the corresponding model group               |

### Request Example

```python
{
    "creator": "user",
    "modifier": "user",
    "bk_classification_id": "test",
    "bk_obj_id": "biz",
    "bk_obj_name": "aaa",
    "obj_sort_number": 1
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": [
        {
            "bk_classification_id": "bk_organization",
            "create_time": "2018-03-08T11:30:28.005+08:00",
            "creator": "cc_system",
            "description": "",
            "id": 4,
            "bk_ispaused": false,
            "ispre": true,
            "last_time": null,
            "modifier": "",
            "bk_obj_icon": "icon-XXX",
            "bk_obj_id": "XX",
            "bk_obj_name": "XXX",
            "position": "{\"test_obj\":{\"x\":-253,\"y\":137}}",
            "bk_supplier_account": "0",
            "obj_sort_number": 1
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name                 | Type               | Description                                                                                                                                                                 |
|----------------------|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                   | int                | Data record ID                                                                                                                                                              |
| creator              | string             | Creator of this data                                                                                                                                                        |
| modifier             | string             | Last modifier of this data                                                                                                                                                  |
| bk_classification_id | string             | Classification ID of the object model, can only be named with alphabetical sequence                                                                                         |
| bk_obj_id            | string             | ID of the object model, can only be named with alphabetical sequence                                                                                                        |
| bk_obj_name          | string             | Name of the object model, used for display                                                                                                                                  |
| bk_supplier_account  | string             | Vendor account                                                                                                                                                              |
| bk_ispaused          | bool               | Whether it is paused, true or false                                                                                                                                         |
| ispre                | bool               | Whether it is predefined, true or false                                                                                                                                     |
| bk_obj_icon          | string             | ICON information of the object model, used for front-end display, values can refer to [(modleIcon.json)](https://chat.openai.com/static/esb/api_docs/res/cc/modleIcon.json) |
| position             | json object string | Coordinates used for front-end display                                                                                                                                      |
| description          | string             | Description of the data                                                                                                                                                     |
| obj_sort_number      | int                | Sorting order of the object model under the corresponding model group                                                                                                       |
