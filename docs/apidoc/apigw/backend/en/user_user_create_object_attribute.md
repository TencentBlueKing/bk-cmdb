### Description

Create Model Property (Permission: Model Edit Permission)

### Parameters

| Name              | Type   | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
|-------------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| creator           | string | No       | Data creator                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| description       | string | No       | Data description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| editable          | bool   | No       | Indicates whether the data is editable                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| isonly            | bool   | No       | Indicates uniqueness                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| ispre             | bool   | No       | true: Preset field, false: Non-built-in field                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| isreadonly        | bool   | No       | true: Read-only, false: Non-read-only                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| isrequired        | bool   | No       | true: Required, false: Optional                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| option            | string | No       | User-defined content, the content and format stored are determined by the calling party, as an example of a numeric type ({"min":1,"max":2})                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| unit              | string | No       | Unit                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| placeholder       | string | No       | Placeholder                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| bk_property_group | string | No       | Field column name                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| bk_obj_id         | string | Yes      | Model ID                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| bk_property_id    | string | Yes      | Model property ID                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| bk_property_name  | string | Yes      | Model property name used for display                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| bk_property_type  | string | Yes      | Defined attribute field used to store data types, with a value range (singlechar(short character), longchar(long character), int(integer), enum(enum type), date(date), time(time), objuser(user), enummulti(enum multiple choice), enumquote(enum reference), timezone(time zone), bool(boolean), organization(organization))                                                                                                                                                                                                                                                                         |
| ismultiple        | bool   | No       | Whether it can be selected multiple times, where the field types are short character, long character, number, float, enum, date, time, time zone, boolean, and the list does not support multiple selections. When creating a property, the field types above do not need to pass the `ismultiple` parameter, and the default is false. If true is passed, it will prompt that this type does not support multiple selections for now. Enum multiple selection, enum reference, user, and organization fields support multiple selections, with user fields and organization fields defaulting to true |
| default           | object | No       | Add default value to the property field, the value of `default` is passed according to the actual type of the field. For example, when creating an int type field, if you want to set a default value for this field, you can pass `default:5`, if it is a short character type, then `default:"aaa"`, if you do not want to set a default value, do not pass this field                                                                                                                                                                                                                               |

### Request Example

```json
{
    "creator": "user",
    "description": "test",
    "editable": true,
    "isonly": false,
    "ispre": false,
    "isreadonly": false,
    "isrequired": false,
    "option": "^[0-9a-zA-Z_]{1,}$",
    "unit": "1",
    "placeholder": "test",
    "bk_property_group": "default",
    "bk_obj_id": "cc_test_inst",
    "bk_property_id": "cc_test",
    "bk_property_name": "cc_test",
    "bk_property_type": "singlechar",
    "bk_asst_obj_id": "test",
    "ismultiple": false,
    "default":"aaaa"
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
	"data": {
		"id": 7,
		"bk_supplier_account": "0",
		"bk_obj_id": "cc_test_inst",
		"bk_property_id": "cc_test",
		"bk_property_name": "cc_test",
		"bk_property_group": "default",
		"bk_property_index": 4,
		"unit": "1",
		"placeholder": "test",
		"editable": true,
		"ispre": false,
		"isrequired": false,
		"isreadonly": false,
		"isonly": false,
		"bk_issystem": false,
		"bk_isapi": false,
		"bk_property_type": "singlechar",
		"option": "",
		"description": "test",
		"creator": "user",
		"create_time": "2020-03-25 17:12:08",
		"last_time": "2020-03-25 17:12:08",
		"bk_property_group_name": "default",
        	"ismultiple": false,
        	"default":"aaaa"
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
| data       | object | Data returned in the request                                                |

#### data

| Name                   | Type   | Description                                                                             |
|------------------------|--------|-----------------------------------------------------------------------------------------|
| creator                | string | Data creator                                                                            |
| description            | string | Data description                                                                        |
| editable               | bool   | Indicates whether the data is editable                                                  |
| isonly                 | bool   | Indicates uniqueness                                                                    |
| ispre                  | bool   | true: Preset field, false: Non-built-in field                                           |
| isreadonly             | bool   | true: Read-only, false: Non-read-only                                                   |
| isrequired             | bool   | true: Required, false: Optional                                                         |
| option                 | string | User-defined content, the content and format stored are determined by the calling party |
| unit                   | string | Unit                                                                                    |
| placeholder            | string | Placeholder                                                                             |
| bk_property_group      | string | Field column name                                                                       |
| bk_obj_id              | string | Model ID                                                                                |
| bk_supplier_account    | string | Vendor account                                                                          |
| bk_property_id         | string | Model property ID                                                                       |
| bk_property_name       | string | Model property name used for display                                                    |
| bk_property_type       | string | Defined attribute field used to store data types                                        |
| bk_property_group_name | string | Field column name                                                                       |
| ismultiple             | bool   | Whether the field supports multiple selections                                          |
| default                | object | Property default value                                                                  |
