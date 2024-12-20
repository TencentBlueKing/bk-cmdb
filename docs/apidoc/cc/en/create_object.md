### Function Description

Create Model (Permission: Create Model)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type   | Required | Description                                                                           |
|----------------------|--------|----------|---------------------------------------------------------------------------------------|
| creator              | string | No       | Creator of this data                                                                  |
| bk_classification_id | string | Yes      | ID of the classification for the object model, can only be named with English letters |
| bk_obj_id            | string | Yes      | ID of the object model, can only be named with English letters                        |
| bk_obj_name          | string | Yes      | Name of the object model, used for display, can be in any language readable by humans |
| bk_obj_icon          | string | No       | ICON information of the object model, used for frontend display                       |
| obj_sort_number      | int    | No       | Sorting order of the object model under the corresponding model group                 |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "request_id": "b529879b85c74e3c91b3d8119df8dbc7",
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

### Response Parameters Description

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| data       | object | Request return data                                          |
| request_id | string | Request chain ID                                             |

### data

| Field                | Type               | Description                                                           |
|----------------------|--------------------|-----------------------------------------------------------------------|
| id                   | int                | New ID of data record                                                 |
| bk_classification_id | int                | ID of the classification for the object model                         |
| creator              | string             | Creator                                                               |
| modifier             | string             | Last modifier                                                         |
| create_time          | string             | Creation time                                                         |
| last_time            | string             | Update time                                                           |
| bk_obj_id            | string             | Model type                                                            |
| bk_obj_name          | string             | Model name                                                            |
| bk_obj_icon          | string             | ICON information of the object model, used for frontend display       |
| position             | json object string | Coordinates for frontend display                                      |
| ispre                | bool               | Whether it is predefined, true or false                               |
| obj_sort_number      | int                | Sorting order of the object model under the corresponding model group |