### Functional description

Update object model properties

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type   | Required   | Description                                   |
|---------------------|---------|--------|-----------------------------------------|
| id                  |  int     | yes  | Record ID of the target data                        |
| description         |  string  |no     | Description information of data                          |
| isonly              |  bool    | no     | Show uniqueness                              |
| isreadonly          |  bool    | no     | Indicates whether it is read-only                            |
| isrequired          |  bool    | no     | Indicates whether it is required                            |
| bk_property_group   |  string  |no     | Name of the field column                          |
| option              |  string  |no     | User-defined content, stored content and format determined by the caller, take digital content as an example ({"min":"1","max":"2"}）|
| bk_property_name    |  string  |no     | Model property name, used to show                    |
| bk_property_type    |  string  |no     | The data type of the defined attribute field used to store the data (singlechar,longchar,int,enum,date,time,objUser,singleasst,multiasst,timezone,bool)|
| unit                |  string  |no     | Unit                                    |
| placeholder         |  string  |no     | Placeholder                                  |
| bk_asst_obj_id      |  string  |no     | This field must be set if there are other models associated with it, otherwise it is not required|

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
| singleasst |Simple correlation   |
| multiasst  |Multiple correlation   |
| timezone   | Time zone     |
| bool       | Bull     |


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```
### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request succeeded or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |No data return|
