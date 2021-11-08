### Functional description

search instance by the associated instance

- the api is just suitable for instances of self-defined mainline model and common model, not suitable for instances of business, set, module, host model, etc.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
| bk_obj_id           | string     | Yes     | Object ID                      |
| bk_supplier_account | string     | Yes     | Supplier account,please fill '0' by independent deployment  |
| page                | object     | Yes     | Page parameters                    |
| condition           | object     | No     | the associated model instance condition                    |
| time_condition           | object     | No     | the model instance time condition                    |
| fields              | map | No     | the model attribution fields to return,key is Object ID，value is the fields to return |

#### page

| Field      |  Type      | Required	   |  Description                |
|-----------|------------|--------|----------------------|
| start     |  int       | Yes     | The record of start position         |
| limit     |  int       | Yes     | Limit number of each page,maximum 200 |
| sort      |  string    | No     | Sort fields             |

#### condition

| Field      |  Type      | Required	   |  Description      |
|-----------|------------|--------|------------|
| field     |string      |Yes      | Value of model field                                                |
| operator  |string      |Yes      | value : $regex $eq $ne                                           |
| value     |string      |Yes      | Value of model field                                   |

#### time_condition

| Field | Type   | Required | Description                          |
| ----- | ------ | -------- | ------------------------------------ |
| oper  | string | Yes      | operator, only support "and" for now |
| rules | array  | Yes      | search time condition                |

#### rules

| Field | Type   | Required | Description                                   |
| ----- | ------ | -------- | --------------------------------------------- |
| field | string | Yes      | Value of model field                          |
| start | string | Yes      | start time in the form of yyyy-MM-dd hh:mm:ss |
| end   | string | Yes      | end time in the form of yyyy-MM-dd hh:mm:ss   |


### Request Parameters Example

```json
{
    "bk_obj_id": "bk_switch",
    "bk_supplier_account": "0",
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_inst_id"
    },
    "fields": {
        "bk_switch": [
            "bk_asset_id",
            "bk_inst_id",
            "bk_inst_name",
            "bk_obj_id"
        ]
    },
    "condition": {
        "user": [
            {
                "field": "operator",
                "operator": "$regex",
                "value": "admin"
            }
        ]
    },
    "time_condition": {
        "oper": "and",
        "rules": [
            {
                "field": "create_time",
                "start": "2021-05-13 01:00:00",
                "end": "2021-05-14 01:00:00"
            }
        ]
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "count": 2,
        "info": [
            {
                "bk_asset_id": "sw00001",
                "bk_inst_id": 1,
                "bk_inst_name": "sw1",
                "bk_obj_id": "bk_switch"
            },
            {
                "bk_asset_id": "sw00002",
                "bk_inst_id": 2,
                "bk_inst_name": "sw2",
                "bk_obj_id": "bk_switch"
            }
        ]
    }
}
```

### Return Result Parameters Description

#### data

| Field      | Type      | Description         |
|-----------|-----------|--------------|
| count     | int       | Count number     |
| info      | array     | The real model instance data |
