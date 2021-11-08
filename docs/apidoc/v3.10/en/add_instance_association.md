### Functional description

create association between object's instance.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description          |
|----------------------|------------|--------|-----------------------------|
| metadata           | object     | Yes    | request meta data             |
| condition | string map     | Yes   | query condition |


metadata params

| Field                 |  Type      | Required	   |  Description         |
|---------------------|------------|--------|-----------------------------|
| label           | string map     | Yes     |the label data request should with, such as biz info |


label params

| Field                 |  Type      | Required	   |  Description         |
|---------------------|------------|--------|-----------------------------|
| bk_biz_id           | string      | Yes     | business's ID |


condition params

| Field                 |  Type      | Required	   |  Description         |
|---------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           | string     | Yes     | object's association unique id |
| bk_inst_id           | int64     | Yes     | association's source object's instance id |
| bk_asst_inst_id           | int64     | Yes     | association's destination object's instance id |


### Request Parameters Example

``` json
{
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 11,
    "bk_asst_inst_id": 21,
    "metadata":{
        "label":{
            "bk_biz_id":"1"
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
    "data": {
        "id": 1038
    }
}

```

### Return Result Parameters Description

#### data ：

| Field       | Type     | Description         |
|------------|----------|--------------|
|id|int64|the instance association's unique id|

