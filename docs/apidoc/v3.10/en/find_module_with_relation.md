### Functional description

find module with relation (v3.9.7)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_biz_id           | int    | Yes      | Business ID           |
| bk_set_ids  | int array  | Yes     | bk_set_id array，the max length is 200 |
| bk_service_template_ids  | int array  | Yes     |  bk_service_template_id array |
| fields  |  string array   | Yes     | module property list, the specified module property feilds will be returned |
| page                | object | Yes       | page info             |

#### page

| Field | Type   | Required | Description                                       |
| ----- | ------ | -------- | ------------------------------------------------- |
| start | int    | Yes       | start record                                      |
| limit | int    | Yes       | page limit, maximum value is 500                 |

### Request Parameters Example

```json
{
    "bk_biz_id": 2,
    "bk_set_ids":[1,2],
    "bk_service_template_ids": [3,4],
    "fields":["bk_module_id", "bk_module_name"],
    "page": {
        "start": 0,
        "limit": 10
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
    "data": {
        "count": 2,
        "info": [
            {
                "bk_module_id": 8,
                "bk_module_name": "license"
            },
            {
                "bk_module_id": 12,
                "bk_module_name": "gse_proc"
            }
        ]
    }
}
```