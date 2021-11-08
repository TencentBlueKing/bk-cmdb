### Functional description

find association between object's instance.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description          |
|----------------------|------------|--------|-----------------------------|
| metadata           | object     | Yes    | request meta data             |
| condition | string map     | Yes   | query condition |
| bk_obj_id           | string     | YES     | the association's source object's id(v3.10+)|

condition params
| Field                 |  Type      | Required	   |  Description         |
|---------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           | string     | Yes     | object's association unique id|
| bk_asst_id           | string     | NO     | object's association kind unique name|
| bk_asst_obj_id           | string     | NO     | the association's destination object's id|


### Request Parameters Example

``` json
{
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_asst_obj_id": ""
    },
    "bk_object_id": "xxx",
    "metadata":{
        "label":{
            "bk_biz_id":"3"
        }
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": [{
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
}

```


### Return Result Parameters Description

#### data ：

| Field       | Type     | Description         |
|------------|----------|--------------|
|id|int64|the association's unique id|
| bk_obj_asst_id| string|  auto generated id, which represent this association.|
| bk_obj_id| string| the association source object's id |
| bk_asst_obj_id| string| the association destination object's id|
| bk_inst_id| int64| the association source object's instance id|
| bk_asst_inst_id| int64| the association destination object's instance id|

