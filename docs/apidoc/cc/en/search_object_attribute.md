### Functional description

You can query object model properties based on model id or business id with optional parameters

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
|bk_obj_id            |  string     | no     | Model ID                      |
| bk_biz_id           |  int        | no     | Business id: after setting, the query result contains the business user-defined field|


### Request Parameters Example

``` python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "test",
    "bk_biz_id": 2
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
           "bk_property_group_name": "Basic Info",
           "bk_property_id": "bk_process_name",
           "bk_property_index": 0,
           "bk_property_name": "Process name",
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
            "bk_property_name": "Business Custom Fields",
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
            "bk_property_group_name": "Business Custom Grouping"
       }
   ]
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data

| Field                | Type         | Description                                                       |
|---------------------|--------------|------------------------------------------------------------|
| creator             |  string       | The creator of the data                                               |
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
| bk_property_name    |  string       | Model property name, used to show                                       |
| bk_property_type    |  string       | The data type of the defined attribute field used to store the data (singlechar,longchar,int,enum,date,time,objUser,singleasst,multiasst,timezone,bool)|
| bk_asst_obj_id      |  string       | This field must be set if there are other models associated with it, otherwise it is not required|
| bk_biz_id           |  int          | Business id of business custom field                                       |
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| id                  |  int    | The id value of the query object   |
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
