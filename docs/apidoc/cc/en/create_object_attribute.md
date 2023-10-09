### Functional description

Create model properties

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                  | Type      | Required   | Description                                                    |
|-----------------------|------------|--------|----------------------------------------------------------|
| creator               |  string     | no     | Who created the data                                             |
| description           |  string     | no     | Description information of data                                           |
| editable              |  bool       | no     | Indicates whether the data is editable                                       |
| isonly                |  bool       | no     | Show uniqueness                                               |
| ispre                 |  bool       | no     | True: preset field,false: Non-built-in field                           |
| isreadonly            |  bool       | no     | True: read-only, false: Not read-only                                  |
| isrequired            |  bool       | no     | True: required, false: Optional                                    |
| option                |  string     | no     | User-defined content, stored content and format determined by the caller, taking numeric type as an example ({"min":"1","max":"2"}）|
| unit                  |  string     | no     | Unit                                                     |
| placeholder           |  string     | no     | Placeholder                                                   |
| bk_property_group     |  string     | no     | Name of the field column                                           |
| bk_obj_id             |  string     | yes     | Model ID                                                   |
| bk_property_id        |  string     | yes     | The property ID of the model                                             |
| bk_property_name      |  string     | yes     | Model attribute name, used to show                                     |
| bk_property_type      |  string     | yes     | The defined attribute field is used to store the data type of the data, and the value range can be (singlechar,longchar,int,enum,date,time,objUser,singleasst,multiasst,timezone,bool)|
| ismultiple        |  bool     | no     | Whether multiple choices are allowed, where the field type is singlechar, longchar, int, float, enum, date, time, timezone, bool, and the list, temporarily does not support multiple choices. When creating an attribute, the field type is the above type, and the ismultiple parameter can not be passed. The default is false. If you pass true, you will be prompted that the type does not support multiple choices. enummulti, enumquote , user and organization fields support multiple choices, among which the user field and organization field are true by default |
| default | object | no | Add a default value to the attribute field. The default value is passed according to the actual type of the field. For example, create an int type field. If you want to set the default value for this field, you can pass default: 5. If it is a short character type, then  default: "aaa". If you do not want to set the default value, you can not pass this field |

#### bk_property_type

| Identification       | Name     |
|------------|----------|
| singlechar |Short character   |
| longchar   | Long character   |
| int        | Reshaping     |
| enum       | Enumeration type|
| date       | Date     |
| time       | Time     |
| objuser    | User     |
| timezone   | Time zone     |
| bool       | Bull     |
| enummulti | Enumerate multiple |
| enumquote | Enumeration References |
| organization | Organization |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "default":"aaaa"
}
```


### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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
        	"default":"aaaa"
	}
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field                | Type         | Description                                                       |
|---------------------|--------------|------------------------------------------------------------|
| creator             |  string       | Who created the data                                               |
| description         |  string       | Description information of data                                             |
| editable            |  bool         | Indicates whether the data is editable                                         |
| isonly              |  bool         | Show uniqueness                                                 |
| ispre               |  bool         | True: preset field,false: Non-built-in field                             |
| isreadonly          |  bool         | True: read-only, false: Not read-only                                    |
| isrequired          |  bool         | True: required, false: Optional                                      |
| option              |  string       | User-defined content, stored content and format determined by the caller               |
| unit                |  string       | Unit                                                       |
| placeholder         |  string       | Placeholder                                                     |
| bk_property_group   |  string       | Name of the field column                                             |
| bk_obj_id           |  string       | Model ID                                                     |
| bk_supplier_account | string       | Developer account number                                                 |
| bk_property_id      |  string       | The property ID of the model                                               |
| bk_property_name    |  string       | Model attribute name, used to show                                       |
| bk_property_type    |  string       | The data type of the defined attribute field used to store the data (singlechar,longchar,int,enum,date,time,objUser,singleasst,multiasst,timezone,bool)|
| bk_biz_id           |  int          | Business id of business custom field                                       |
| bk_property_group_name           |  string          | Name of the field column                                       |
| ismultiple | bool | Can multiple fields be selected |
| default | object | attribute default vaule |

#### bk_property_type

| Identification       | Name     |
|------------|----------|
| singlechar |Short character   |
| longchar   | Long character   |
| int        | Reshaping     |
| enum       | Enumeration type|
| date       | Date     |
| time       | Time     |
| objuser    | User     |
| timezone   | Time zone     |
| bool       | Bull     |
| enummulti | Enumerate multiple |
| enumquote | Enumeration References |
| organization | Organization |
