### Functional description

find association between object.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description          |
|----------------------|------------|--------|-----------------------------|
| condition | string map     | Yes   | query condition |


condition params
| Field                 |  Type      | Required	   |  Description         |
|---------------------|------------|--------|-----------------------------|
| bk_asst_id           | string     | Yes     | object's association kind unique name|
| bk_obj_id           | string     | Yes     | the association's source object's id|
| bk_asst_id           | string     | Yes     | the association's destination object's id|


### Request Parameters Example

``` json
{
    "condition": {
        "bk_asst_id": "belong",
        "bk_obj_id": "bk_switch",
        "bk_asst_obj_id": "bk_host"
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": [
        {
            "id": 1,
            "bk_obj_asst_id": "bk_switch_belong_bk_host",
            "bk_obj_asst_name": "",
            "bk_asst_id": "belong",
            "bk_asst_name": "belong",
            "bk_obj_id": "bk_switch",
            "bk_obj_name": "switch",
            "bk_asst_obj_id": "bk_host",
            "bk_asst_obj_name": "host",
            "mapping": "1:n",
            "on_delete": "none"
        }
    ]
}

```


### Return Result Parameters Description

#### data ：

| Field       | Type     | Description         |
|------------|----------|--------------|
|id|int64|the association's unique id|
| bk_obj_asst_id| string|  auto generated id, which represent this association.|
| bk_obj_asst_name| string| the alias name for this association. |
| bk_asst_id| string| association kind id |
| bk_asst_name| string| association kind name |
| bk_obj_id| string| the association source object's id |
| bk_obj_name| string| the association source object's name |
| bk_asst_obj_id| string| the association destination object's id|
| bk_asst_obj_name| string| the association destination object's name|
| mapping| string| association between object's instance type, could be one of [1:1, 1:n, n:n] |
| on_delete| string| the action when this association is delete, could be one of [none, delete_src, delete_dest], "none" means do nothing, "delete_src" means delete source object's instance, "delete_dest" means delete destination object's instance.|
