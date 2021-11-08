### Functional description

 Query all the association relationships of an instance by relation name (including the situation that it is used as the original model of association relationship and the target model of association relationship)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field      | Type   | Required | Description                 |
| ---------- | ------ | -------- | --------------------------- |
| bk_inst_id | int    | Yes      | Instance ID                 |
| bk_obj_id  | string | Yes      | Object ID                   |
| fields     | array  | Yes      | Fields to be returned       |
| start      | int    | No       | Record start position       |
| limit      | int    | No       | Page size,  maximum is 500. |

### Request Parameters Example

```json
{
    "condition": {
        "bk_inst_id": 16,
        "bk_obj_id": "bk_router"
    },
    "fields": [
        "id",
        "bk_inst_id",
        "bk_obj_id",
        "bk_asst_inst_id",
        "bk_asst_obj_id",
        "bk_obj_asst_id",
        "bk_asst_id"
        ],
    "page": {
        "start":0,
        "limit":2
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "permission": null,
    "data": [
        {
            "id": 4,
            "bk_inst_id": 1,
            "bk_obj_id": "bk_switch",
            "bk_asst_inst_id": 16,
            "bk_asst_obj_id": "bk_router",
            "bk_obj_asst_id": "bk_switch_default_bk_router",
            "bk_asst_id": "default"
        },
        {
            "id": 6,
            "bk_inst_id": 2,
            "bk_obj_id": "bk_switch",
            "bk_asst_inst_id": 16,
            "bk_asst_obj_id": "bk_router",
            "bk_obj_asst_id": "bk_switch_default_bk_router",
            "bk_asst_id": "default"
        }
    ]
}
```

### Return Result Parameters Description

#### data

| Field           | Type   | Description                                   |
| --------------- | ------ | --------------------------------------------- |
| id              | int64  | Association ID                                |
| bk_inst_id      | int64  | Source instance ID                            |
| bk_obj_id       | string | Source object ID                              |
| bk_asst_inst_id | int64  | Target instance ID                            |
| bk_asst_obj_id  | string | Target object ID                              |
| bk_obj_asst_id  | string | Automatically generated object association ID |
| bk_asst_id      | string | Relationship name                             |

