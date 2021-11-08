### Functional description

find sets in one biz (v3.8.6)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_biz_id           | int    | Yes      | Business ID           |
| bk_ids  | int array  | Yes     | bk_set_id array，the max length is 500 |
| fields  |  string array   | Yes     | set property list, the specified set property feilds will be returned |

### Request Parameters Example

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        11,
        12
    ],
    "fields": [
        "bk_set_id",
        "bk_set_name",
        "create_time"
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "bk_set_id": 12,
            "bk_set_name": "ss1",
            "create_time": "2020-05-15T22:15:51.67+08:00",
            "default": 0
        },
        {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "create_time": "2020-05-12T21:04:36.644+08:00",
            "default": 0
        }
    ]
}
```
