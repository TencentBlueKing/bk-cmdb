### Description

Query the association relationships between models. (Permission: Model view permission)

### Parameters

| Name      | Type       | Required | Description      |
|-----------|------------|----------|------------------|
| condition | string map | Yes      | Query conditions |

Condition Parameters

| Name           | Type   | Required | Description                                                |
|----------------|--------|----------|------------------------------------------------------------|
| bk_asst_id     | string | No       | Unique ID of the model's association type                  |
| bk_obj_id      | string | No       | Source model ID, either this or bk_asst_obj_id is required |
| bk_asst_obj_id | string | No       | Target model ID, either this or bk_obj_id is required      |

**Note: Without the condition limit of bk_asst_id, if only the condition of bk_obj_id is filled, it will query all
association relationships where the model acts as the source model; if only the condition of bk_asst_obj_id is filled,
it will query all association relationships where the model acts as the target model**

### Request Example

```json
{
    "condition": {
        "bk_asst_id": "belong",
        "bk_obj_id": "bk_switch",
        "bk_asst_obj_id": "bk_host"
    }
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": [
        {
           "id": 27,
           "bk_supplier_account": "0",
           "bk_obj_asst_id": "test1_belong_biz",
           "bk_obj_asst_name": "1",
           "bk_obj_id": "test1",
           "bk_asst_obj_id": "biz",
           "bk_asst_id": "belong",
           "mapping": "n:n",
           "on_delete": "none",
           "ispre": null
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name                | Type   | Description                                                                                                                                                                                                              |
|---------------------|--------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                  | int64  | Identity ID of the model association relationship                                                                                                                                                                        |
| bk_obj_asst_id      | string | Unique ID of the model association relationship.                                                                                                                                                                         |
| bk_obj_asst_name    | string | Alias of the association relationship.                                                                                                                                                                                   |
| bk_asst_id          | string | ID of the association type                                                                                                                                                                                               |
| bk_obj_id           | string | Source model ID                                                                                                                                                                                                          |
| bk_asst_obj_id      | string | Target model ID                                                                                                                                                                                                          |
| mapping             | string | Mapping relationship between the source model and the target model, one of [1:1, 1:n, n:n]                                                                                                                               |
| on_delete           | string | Action when deleting the association relationship, one of [none, delete_src, delete_dest]. "none" does nothing, "delete_src" deletes instances of the source model, "delete_dest" deletes instances of the target model. |
| bk_supplier_account | string | Developer account                                                                                                                                                                                                        |
| ispre               | bool   | true: pre-installed field, false: non-built-in field                                                                                                                                                                     |
