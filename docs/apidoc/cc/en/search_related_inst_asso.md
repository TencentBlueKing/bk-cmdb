### Functional description

 Query all Association relationships of an instance (including the case that it is the original model of Association relationship and the target model of Association relationship)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field       | Type   | Required| Description                  |
| ---------- | ------ | ---- | --------------------- |
| bk_inst_id | int    | yes | Instance id                |
| bk_obj_id  | string |yes   | Model id                |
| fields     |  array  |yes   | Fields to be returned        |
| start      |  int    | no   | Record start position          |
| limit      |  int    | no   | Page size, maximum 500. |

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start     |   int     | yes  | Record start position|
| limit     |   int     | yes  | Limit bars per page, Max. 200|


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data

| Name            | Type   | Description                     |
| :-------------- | :----- | :----------------------- |
| id              |  int64  |Association id                   |
| bk_inst_id      |  int64  |Source model instance id             |
| bk_obj_id       |  string |Association relationship source model id         |
| bk_asst_inst_id | int64  |Association relation target model id       |
| bk_asst_obj_id  | string |Target model instance id           |
| bk_obj_asst_id  | string |Auto-generated model association id|
| bk_asst_id      |  string |Relationship name                 |