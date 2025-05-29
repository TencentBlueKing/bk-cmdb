### Description

Update custom model attributes for a business (Permission: Business custom field editing permission)

### Parameters

| Name              | Type                                                | Required | Description                                                                                                                                                                                                                                                                          |
|-------------------|-----------------------------------------------------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                | int                                                 | Yes      | Record ID of the target data                                                                                                                                                                                                                                                         |
| bk_biz_id         | int                                                 | Yes      | Business ID                                                                                                                                                                                                                                                                          |
| description       | string                                              | No       | Description of the data                                                                                                                                                                                                                                                              |
| isonly            | bool                                                | No       | Indicates uniqueness                                                                                                                                                                                                                                                                 |
| isreadonly        | bool                                                | No       | Indicates if it is read-only                                                                                                                                                                                                                                                         |
| isrequired        | bool                                                | No       | Indicates if it is required                                                                                                                                                                                                                                                          |
| bk_property_group | string                                              | No       | Name of the field column                                                                                                                                                                                                                                                             |
| option            | object                                              | No       | User-defined content, the format and content are determined by the caller, using numeric content as an example (`{"min":1,"max":2}`)                                                                                                                                                 |
| bk_property_name  | string                                              | No       | Model attribute name for display                                                                                                                                                                                                                                                     |
| unit              | string                                              | No       | Unit                                                                                                                                                                                                                                                                                 |
| bk_property_type  | string                                              | Yes      | Defined property field for storing data types (`singlechar(short string),longchar(long string),int(integer),enum(enum type),date(date),time(time),objuser(user),enummulti(multi-select enum),enumquote(enum reference),timezone(timezone),bool(boolean),organization(organization)`) |
| placeholder       | string                                              | No       | Placeholder                                                                                                                                                                                                                                                                          |
| bk_asst_obj_id    | string                                              | No       | If there is a relationship with other models, then this field must be set; otherwise, it does not need to be set                                                                                                                                                                     |
| default           | Depends on the type specified by `bk_property_type` | No       | Default value                                                                                                                                                                                                                                                                        |

### Request Example

```json
{

    "id": 1,
    "bk_biz_id": 2,
    "description": "test",
    "placeholder": "test",
    "unit": "1",
    "isonly": false,
    "isreadonly": false,
    "isrequired": false,
    "bk_property_group": "default",
    "option": {"min":1,"max":4},
    "bk_property_name": "aaa",
    "bk_property_type": "int",
    "bk_asst_obj_id": "0"
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
