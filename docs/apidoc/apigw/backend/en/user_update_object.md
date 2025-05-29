### Description

Update Model Definition (Permission: Model Editing Permission)

### Parameters

| Name                 | Type               | Required | Description                                                                                                                                                                                                                                                                                                                                                               |
|----------------------|--------------------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                   | int                | No       | ID of the object model, used as a condition for the update operation                                                                                                                                                                                                                                                                                                      |
| modifier             | string             | No       | Last modifier of this data                                                                                                                                                                                                                                                                                                                                                |
| bk_classification_id | string             | Yes      | Classification ID of the object model, can only be named with an alphabetical sequence                                                                                                                                                                                                                                                                                    |
| bk_obj_name          | string             | No       | Name of the object model                                                                                                                                                                                                                                                                                                                                                  |
| bk_obj_icon          | string             | No       | ICON information of the object model, used for frontend display, values can be referred to [(modleIcon.json)](https://chat.openai.com/static/esb/api_docs/res/cc/modleIcon.json)                                                                                                                                                                                          |
| position             | json object string | No       | Coordinates for frontend display                                                                                                                                                                                                                                                                                                                                          |
| obj_sort_number      | int                | No       | Sorting number of the object model under its model group; when updating this value, if the set value exceeds the maximum value of this value in the group model, the updated value will be the maximum value plus one. For example, if the set value is 999, and the current maximum value of this value in the group model is 6, then the updated value will be set to 7 |

### Request Example

```python
{
    "id": 1,
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc2_test_inst",
    "bk_obj_icon": "icon-cc-business",
    "position":"{\"ff\":{\"x\":-863,\"y\":1}}",
    "obj_sort_number": 1
}
```

### Response Example

```python
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
| permission | object | Permission information                                              |
| data       | object | No data returned                                                    |
