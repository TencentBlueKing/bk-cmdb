### Description

Update Object Model Property (Permission: Model Editing Permission)

### Parameters

| Name              | Type   | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
|-------------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                | int    | Yes      | Record ID of the target data                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| description       | string | No       | Description information of the data                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| isonly            | bool   | No       | Indicates uniqueness                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| isreadonly        | bool   | No       | Indicates whether it is read-only                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| isrequired        | bool   | No       | Indicates whether it is required                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| bk_property_group | string | No       | Name of the field column                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| option            | string | No       | User-defined content, the content and format stored are determined by the caller. For example, using numeric content ({"min":"1","max":"2"})                                                                                                                                                                                                                                                                                                                                                                                                             |
| bk_property_name  | string | No       | Model property name, used for display                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| unit              | string | No       | Unit                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| bk_property_type  | string | Yes      | Defined property field for storing data type (singlechar (short character), longchar (long character), int (integer), enum (enumeration type), date (date), time (time), objuser (user), enummulti (enumeration multiple), enumquote (enumeration reference), timezone (time zone), bool (boolean), organization (organization))                                                                                                                                                                                                                         |
| placeholder       | string | No       | Placeholder                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| ismultiple        | bool   | No       | Whether it can be selected multiple times. For field types such as short character, long character, number, floating point, enumeration, date, time, time zone, boolean, multiple selection is not supported temporarily. When updating the property, if the field type is one of the above types, ismultiple cannot be updated to true. If updated to true, it will prompt that this type does not support multiple selection temporarily. Enumeration multiple selection, enumeration reference, user, organization fields support multiple selection. |
| default           | object | No       | Add a default value to the attribute. When updating, the value of default is passed according to the actual type of the field. If you want to clear the default value of the field, you need to pass default: null                                                                                                                                                                                                                                                                                                                                       |

### Request Example

Update Default Value Scenario

```python
{
    "id":1,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":"{\"min\":\"1\",\"max\":\"4\"}",
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0",
    "default":3
}
```

Do Not Update Default Value Scenario

```python
{
    "id":1,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":"{\"min\":\"1\",\"max\":\"4\"}",
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0"
}
```

Clear Default Value Scenario

```python
{
    "id":1,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":"{\"min\":\"1\",\"max\":\"4\"}",
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0",
    "default":null
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
