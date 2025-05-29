### Description

Query model instances based on association relationship instances (Permission: Model Instance Query Permission)

- This interface is only applicable to custom hierarchical models and general model instances, not applicable to
  business, cluster, module, host, etc.

### Parameters

| Name           | Type   | Required | Description                                                                                                                                                 |
|----------------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id      | string | Yes      | Model ID                                                                                                                                                    |
| page           | object | Yes      | Pagination parameters                                                                                                                                       |
| condition      | object | No       | Query conditions for model instances with association relationships                                                                                         |
| time_condition | object | No       | Query conditions for model instances based on time                                                                                                          |
| fields         | object | No       | Specify the fields to be returned for the queried model instances, where the key is the model ID, and the value is the model property fields to be returned |

#### page

| Name  | Type   | Required | Description                           |
|-------|--------|----------|---------------------------------------|
| start | int    | Yes      | Record start position                 |
| limit | int    | Yes      | Each page limit, maximum value is 200 |
| sort  | string | No       | Sorting field                         |

#### condition

In the example, user is the model

| Name     | Type   | Required | Description                                              |
|----------|--------|----------|----------------------------------------------------------|
| field    | string | Yes      | The field value is the field name of the model           |
| operator | string | Yes      | The value is: $regex $eq $ne                             |
| value    | string | Yes      | The value corresponding to the field configured by field |

#### time_condition

| Name  | Type   | Required | Description                           |
|-------|--------|----------|---------------------------------------|
| oper  | string | Yes      | Operator, currently only supports and |
| rules | array  | Yes      | Time query conditions                 |

#### rules

| Name  | Type   | Required | Description                                   |
|-------|--------|----------|-----------------------------------------------|
| field | string | Yes      | The value is the field name of the model      |
| start | string | Yes      | Start time, in the format yyyy-MM-dd hh:mm:ss |
| end   | string | Yes      | End time, in the format yyyy-MM-dd hh:mm:ss   |

### Request Example

```json
{
    "bk_obj_id": "bk_switch",
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

### Response Example

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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name  | Type  | Description                    |
|-------|-------|--------------------------------|
| count | int   | Record count                   |
| info  | array | Actual data of model instances |

#### data.info[n]

| Name         | Type   | Description   |
|--------------|--------|---------------|
| bk_asset_id  | string | Asset ID      |
| bk_inst_id   | int    | Instance ID   |
| bk_inst_name | string | Instance name |
| bk_obj_id    | string | Model ID      |
