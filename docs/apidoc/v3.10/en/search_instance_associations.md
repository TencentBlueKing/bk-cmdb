### Functional description

search object instance associations (v3.10.1+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|  Field     |  Type   | Required | Description                                                                                                            |
|------------|---------|----------|------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id  | integer |  No      | application business id, only used for mainline object search                                                          |
| bk_obj_id  | string  |  Yes     | object id                                                                                                              |
| conditions | object  |  No      | conditions, support AND/OR types，max conditions deep 3, max OR conditions rules count is 20, empty means matching all（as is conditions is null value） |
| fields     | array   |  No      | fields of object, empty means all fields                                                                               |
| page       | object  |  Yes     | query page settings                                                                                                    |

#### conditions

|  Field   |  Type  | Required | Description                                                                                                                |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------|
| field    | string |  Yes     | condition field, support id, bk_inst_id, bk_obj_id, bk_asst_inst_id, bk_asst_obj_id, bk_obj_asst_id, bk_asst_id            |
| operator | string |  Yes     | condition operator, support like equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    |   -    |  No      | condition value, max slice(array) elements count is 500                                                                    |

condition rules detail: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Field  | Type    | Required  | Description            |
|--------|---------|-----------|------------------------|
| start  | int     | Yes       | start record           |
| limit  | int     | Yes       | page limit, max is 500 |
| sort   | string  | No        | query order by         |

### Request Parameters Example

```json
{
    "bk_app_code":"code",
    "bk_app_secret":"secret",
    "bk_token":"xxxx",
    "bk_obj_id":"bk_switch",
    "conditions":{
        "condition": "AND",
        "rules": [
            {
                "field": "bk_obj_asst_id",
                "operator": "equal",
                "value": "bk_switch_connect_host"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                         "field": "bk_inst_id",
                         "operator": "in",
                         "value": [2,4,6]
                    },
                    {
                        "field": "bk_asst_id",
                        "operator": "equal",
                        "value": 3
                    }
                ]
            }
        ]
    },
    "fields":[
        "bk_inst_id",
        "bk_asst_inst_id",
        "bk_asst_obj_id",
        "bk_asst_id",
        "bk_obj_asst_id"
    ],
    "page":{
        "start":0,
        "limit":500
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
        "info": [
            {
                "bk_inst_id": 2,
                "bk_asst_inst_id": 8,
                "bk_asst_obj_id": "host",
                "bk_asst_id": "connect",
                "bk_obj_asst_id": "bk_switch_connect_host"
            }
        ]
    }
}
```

### Return Result Parameters Description

#### data

| Field  |  Type   | Description       |
|--------|---------|-------------------|
| info   | array   | data of record    |
