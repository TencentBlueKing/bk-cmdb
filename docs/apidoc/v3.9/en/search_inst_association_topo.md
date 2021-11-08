### Functional description

search instance association topology

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
|bk_supplier_account  |string|Yes |Supplier account|
|bk_obj_id            |string|Yes |Object ID|
|bk_inst_id           |int|Yes |Instance ID|


### Request Parameters Example

``` python
{
    "bk_supplier_account":"0",
    "bk_obj_id":"test",
    "bk_inst_id":"test"
}
```


### Return Result Example

```python
{
    "result":true,
    "code":0,
    "message":"",
    "data":[
        {
            "bk_inst_id":0,
            "bk_inst_name":"",
            "bk_obj_icon":"icon-cc-business",
            "bk_obj_id":"biz",
            "bk_obj_name":"业务",
            "count":1,
            "children":[
                {
                    "bk_inst_id":2,
                    "bk_inst_name":"蓝鲸",
                    "bk_obj_icon":"",
                    "bk_obj_id":"biz",
                    "bk_obj_name":"业务"
                }
            ]
        }
    ]
}
```

### Return Result Parameters Description

#### data

| Field         | Type         | Description                          |
|--------------|--------------|-------------------------------|
| bk_inst_id   | int          | Instance ID                        |
| bk_inst_name | string       | Instance name for display            |
| bk_obj_icon  | string       | Object icon name                |
| bk_obj_id    | string       | Object ID                        |
| bk_obj_name  | string       | Object name for display            |
| children     | object array | The set of associated instances under this model|
| count        | int          | Children include node's number   |

#### children

| Field         | Type      | Description               |
|--------------|-----------|--------------------|
| bk_inst_id   |int        | Instance ID             |
| bk_inst_name |string     | Instance name for display |
| bk_obj_icon  |string     | Object icon name     |
| bk_obj_id    |string     | Object ID             |
| bk_obj_name  |string     | Object name for display |
