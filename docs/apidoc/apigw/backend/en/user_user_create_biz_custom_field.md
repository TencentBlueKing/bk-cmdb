### Description

Create Business Custom Model Property (Permission: Business Custom Field Edit Permission)

### Parameters

| Name              | Type   | Required | Description                                                                                                                                                                                                                                                                                                                                                    |
|-------------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id         | string | Yes      | Model ID                                                                                                                                                                                                                                                                                                                                                       |
| bk_property_id    | string | Yes      | Model property ID                                                                                                                                                                                                                                                                                                                                              |
| bk_property_name  | string | Yes      | Model property name used for display                                                                                                                                                                                                                                                                                                                           |
| bk_property_type  | string | Yes      | Defined attribute field used to store data types, with a value range (singlechar(short character), longchar(long character), int(integer), enum(enum type), date(date), time(time), objuser(user), enummulti(enum multiple choice), enumquote(enum reference), timezone(time zone), bool(boolean), organization(organization))                                 |
| bk_biz_id         | int    | Yes      | Business ID                                                                                                                                                                                                                                                                                                                                                    |
| creator           | string | No       | Data creator                                                                                                                                                                                                                                                                                                                                                   |
| description       | string | No       | Data description                                                                                                                                                                                                                                                                                                                                               |
| editable          | bool   | No       | Indicates whether the data is editable                                                                                                                                                                                                                                                                                                                         |
| isonly            | bool   | No       | Indicates uniqueness                                                                                                                                                                                                                                                                                                                                           |
| ispre             | bool   | No       | true: Preset field, false: Non-built-in field                                                                                                                                                                                                                                                                                                                  |
| isreadonly        | bool   | No       | true: Read-only, false: Non-read-only                                                                                                                                                                                                                                                                                                                          |
| isrequired        | bool   | No       | true: Required, false: Optional                                                                                                                                                                                                                                                                                                                                |
| option            | string | No       | User-defined content, the content and format stored are determined by the calling party, as an example of a numeric type ({"min":1,"max":2})                                                                                                                                                                                                                   |
| unit              | string | No       | Unit                                                                                                                                                                                                                                                                                                                                                           |
| placeholder       | string | No       | Placeholder                                                                                                                                                                                                                                                                                                                                                    |
| bk_property_group | string | No       | Field column name                                                                                                                                                                                                                                                                                                                                              |
| bk_asst_obj_id    | string | No       | If there is a relation to other models, then this field must be set, otherwise it does not need to be set                                                                                                                                                                                                                                                      |
| default           | object | No       | Add default value to property field, the value of default is passed according to the actual type of the field, for example, when creating an int type field, if you want to set a default value for this field, you can pass default:5, if it is a short character type, then default:"aaa", if you do not want to set a default value, do not pass this field |

**Note:**

- The `create_biz_custom_field` interface is used to create business custom fields, which are only valid within the
  business. The difference between business custom fields and other model fields is that the `bk_biz_id` of business
  custom fields is the actual business ID, while the `bk_biz_id` of other model fields is 0.
- When calling this interface, the `bk_biz_id` parameter in the parameters should be the actual business ID, and
  the `bk_obj_id` can only be set to "set", "module", and "host".

### Request Example

```json
{
    "bk_biz_id": 2,
    "creator": "user",
    "description": "test",
    "editable": true,
    "isonly": false,
    "ispre": false,
    "isreadonly": false,
    "isrequired": false,
    "option": {"min":1,"max":2},
    "unit": "1",
    "placeholder": "test",
    "bk_property_group": "default",
    "bk_obj_id": "set",
    "bk_property_id": "cc_test",
    "bk_property_name": "cc_test",
    "bk_property_type": "singlechar",
    "bk_asst_obj_id": "test"
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
		"bk_biz_id": 2,
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
		"option": {"min":1,"max":2},
		"description": "test",
		"creator": "user",
		"create_time": "2020-03-25 17:12:08",
		"last_time": "2020-03-25 17:12:08"
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
| bk_biz_id              | int    | Business ID of the business custom field                                                |
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
| bk_property_group_name | string | Field column name                                                                       |
| bk_obj_id              | string | Model ID                                                                                |
| bk_supplier_account    | string | Vendor account                                                                          |
| bk_property_id         | string | Model property ID                                                                       |
| bk_property_name       | string | Model property name used for display                                                    |
| bk_property_type       | string | Defined attribute field used to store data types                                        |
| bk_asst_obj_id         | string | If there is a relation to other models, then this field must be set                     |
| create_time            | string | Creation time                                                                           |
| last_time              | string | Update time                                                                             |
| id                     | int    | Primary key ID                                                                          |
