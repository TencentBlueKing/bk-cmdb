### Description

You can use optional parameters to query object model properties based on the model id or business id (Permission: Model
View Permission)

### Parameters

| Name      | Type   | Required | Description                                                           |
|-----------|--------|----------|-----------------------------------------------------------------------|
| bk_obj_id | string | Yes      | Model ID                                                              |
| bk_biz_id | int    | No       | Business id, if set, the query result contains business custom fields |

### Request Example

```python
{
    "bk_obj_id": "test",
    "bk_biz_id": 2
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
           "bk_biz_id": 0,
           "bk_asst_obj_id": "",
           "bk_asst_type": 0,
           "create_time": "2018-03-08T11:30:27.898+08:00",
           "creator": "cc_system",
           "description": "",
           "editable": false,
           "id": 51,
           "isapi": false,
           "isonly": true,
           "ispre": true,
           "isreadonly": false,
           "isrequired": true,
           "last_time": "2018-03-08T11:30:27.898+08:00",
           "bk_obj_id": "process",
           "option": "",
           "placeholder": "",
           "bk_property_group": "default",
           "bk_property_group_name": "基础信息",
           "bk_property_id": "bk_process_name",
           "bk_property_index": 0,
           "bk_property_name": "进程名称",
           "bk_property_type": "singlechar",
           "bk_supplier_account": "0",
           "unit": ""
       },
       {
            "bk_biz_id": 2,
            "id": 7,
            "bk_supplier_account": "0",
            "bk_obj_id": "process",
            "bk_property_id": "biz_custom_field",
            "bk_property_name": "业务自定义字段",
            "bk_property_group": "biz_custom_group",
            "bk_property_index": 4,
            "unit": "",
            "placeholder": "",
            "editable": true,
            "ispre": true,
            "isrequired": false,
            "isreadonly": false,
            "isonly": false,
            "bk_issystem": false,
            "bk_isapi": false,
            "bk_property_type": "singlechar",
            "option": "",
            "description": "",
            "creator": "admin",
            "create_time": "2020-03-25 17:12:08",
            "last_time": "2020-03-25 17:12:08",
            "bk_property_group_name": "业务自定义分组"
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

| Name                | Type   | Description                                                                                                                                                                                                                                                                                     |
|---------------------|--------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| creator             | string | Creator of the data                                                                                                                                                                                                                                                                             |
| description         | string | Description information of the data                                                                                                                                                                                                                                                             |
| editable            | bool   | Indicates whether the data is editable                                                                                                                                                                                                                                                          |
| isonly              | bool   | Indicates uniqueness of the data                                                                                                                                                                                                                                                                |
| ispre               | bool   | true: pre-installed field, false: non-built-in field                                                                                                                                                                                                                                            |
| isreadonly          | bool   | true: read-only, false: non-read-only                                                                                                                                                                                                                                                           |
| isrequired          | bool   | true: required, false: optional                                                                                                                                                                                                                                                                 |
| option              | string | User-defined content, the content and format stored is determined by the caller                                                                                                                                                                                                                 |
| unit                | string | Unit                                                                                                                                                                                                                                                                                            |
| placeholder         | string | Placeholder                                                                                                                                                                                                                                                                                     |
| bk_property_group   | string | Name of the field column                                                                                                                                                                                                                                                                        |
| bk_obj_id           | string | Model ID                                                                                                                                                                                                                                                                                        |
| bk_supplier_account | string | Vendor account                                                                                                                                                                                                                                                                                  |
| bk_property_id      | string | Model property ID                                                                                                                                                                                                                                                                               |
| bk_property_name    | string | Model property name used for display                                                                                                                                                                                                                                                            |
| bk_property_type    | string | Defined property field for storing data types (singlechar(short character), longchar(long character), int(integer), enum(enum type), date(date), time(time), objuser(user), enummulti(enum multiple), enumquote(enum reference), timezone(timezone), bool(boolean), organization(organization)) |
| bk_asst_obj_id      | string | If there is a relationship with other models, this field must be set, otherwise it does not need to be set                                                                                                                                                                                      |
| bk_biz_id           | int    | Business id of business custom field                                                                                                                                                                                                                                                            |
| create_time         | string | Creation time                                                                                                                                                                                                                                                                                   |
| last_time           | string | Update time                                                                                                                                                                                                                                                                                     |
| id                  | int    | Query object id value                                                                                                                                                                                                                                                                           |
