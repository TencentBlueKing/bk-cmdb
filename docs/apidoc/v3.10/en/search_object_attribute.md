### Functional description

search object attribute

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
|bk_obj_id            | string     | No     | Object ID                      |
|bk_supplier_account  | string     | No     | Supplier account                  |
| bk_biz_id           | int        | No    | business's ID, set this and query result wil include biz custom field |


### Request Parameters Example

``` python
{
    "bk_obj_id": "test",
    "bk_supplier_account": "0",
    "bk_biz_id": 2
}
```


### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
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
           "bk_property_group_name": "base information",
           "bk_property_id": "bk_process_name",
           "bk_property_index": 0,
           "bk_property_name": "process name",
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
            "bk_property_name": "biz custom field",
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
            "bk_property_group_name": "biz custom group"
       }
   ]
}
```

### Return Result Parameters Description

#### data

| Field                | Type         | Description                                                       |
|---------------------|--------------|------------------------------------------------------------|
| creator             | string       | The creator of data                                               |
| description         | string       | Description information of data                                              |
| editable            | bool         | Editable data                                         |
| isonly              | bool         | Uniqueness data                                                 |
| ispre               | bool         | true:preset field, false:non preset field                             |
| isreadonly          | bool         | true:read-only, false:non read-only                                    |
| isrequired          | bool         | true:required, false:optional                                      |
| option              | string       | User's custom content，the content and format of memory is determined by caller               |
| unit                | string       | Unit                                                       |
| placeholder         | string       | Placeholder                                                     |
| bk_property_group   | string       | Object property group name                                             |
| bk_obj_id           | string       | Object ID                                                     |
| bk_supplier_account | string       | Supplier account                                                 |
| bk_property_id      | string       | Object Property ID                                               |
| bk_property_name    | string       | Object property name                                       |
| bk_property_type    | string       | The storage data type of defined property field,range list(singlechar,longchar,int,enum,date,time,objuser,singleasst,multiasst,timezone,bool)|
| bk_asst_obj_id      | string       | If there are other models associated with the object, then must be set this field, otherwise, it doesn't to be set|
| bk_biz_id           | int          | business's ID of biz custom field                          |

#### bk_property_type

| identifier       | name     |
|------------|----------|
| singlechar | Single character   |
| longchar   | Long character   |
| int        | Integer     |
| enum       | Enumeration |
| date       | Date     |
| time       | Time      |
| objuser    | Object user      |
| singleasst | Single association   |
| multiasst  | Multiple association   |
| timezone   | Timezone     |
| bool       | Bool     |
