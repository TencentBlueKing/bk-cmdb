### Functional description

update business custom object attribute

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type   | Required	   |  Description                                   |
|---------------------|---------|--------|-----------------------------------------|
| id                  | int     | Yes     |   ID of target data record                        |
| bk_biz_id           | int     | Yes    | business's ID                                              |
| description         | string  | No     |  Description information of datas                          |
| isonly              | bool    | No     | Uniqueness data                              |
| isreadonly          | bool    | No     | Read-only, true or not                            |
| isrequired          | bool    | No     | Required, true or not                            |
| bk_property_group   | string  | No     | Property group name                          |
| option              | string  | No     | User's custom content，the content and format of memory is determined by caller, example for digital content({"min":"1","max":"2"})|
| bk_property_name    | string  | No     | Property name, for display                    |
| bk_property_type    | string  | No     | The storage data type of defined property field,rang list（singlechar,longchar,int,enum,date,time,objuser,singleasst,multiasst,timezone,bool)|
| unit                | string  | No     | Unit                                    |
| placeholder         | string  | No     | Placeholder                                  |
| bk_asst_obj_id      | string  | No     | If there are other models associated with the object, then must be set this field, otherwise, it doesn't to be set |

#### bk_property_type

| identifier       | name     |
|------------|----------|
| singlechar | Single character   |
| longchar   | Long character   |
| int        | Integer     |
| enum       | Enumeration |
| date       | Date     |
| time       | Time     |
| objuser    | Object user     |
| singleasst | Single association   |
| multiasst  | Multiple association   |
| timezone   | Timezone     |
| bool       | Bool     |


### Request Parameters Example

```json
{
    "id":1,
    "bk_biz_id": 2,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":{"min":1,"max":4},
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0"
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null
}
```
