### Functional description

find modules in one biz (v3.8.6)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_biz_id           | int    | Yes      | Business ID           |
| bk_ids  | int array  | Yes     | bk_set_id array，the max length is 500 |
| fields  |  string array   | Yes     | module property list, the specified module property feilds will be returned |

### Request Parameters Example

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        56,
        57,
        58,
        59,
        60
    ],
    "fields": [
        "bk_module_id",
        "bk_module_name",
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
            "bk_module_id": 60,
            "bk_module_name": "sm1",
            "create_time": "2020-05-15T22:15:51.725+08:00",
            "default": 0
        },
        {
            "bk_module_id": 59,
            "bk_module_name": "m1",
            "create_time": "2020-05-12T21:04:47.286+08:00",
            "default": 0
        },
        {
            "bk_module_id": 58,
            "bk_module_name": "recycle host",
            "create_time": "2020-05-12T21:03:37.238+08:00",
            "default": 3
        },
        {
            "bk_module_id": 57,
            "bk_module_name": "fault host",
            "create_time": "2020-05-12T21:03:37.183+08:00",
            "default": 2
        },
        {
            "bk_module_id": 56,
            "bk_module_name": "idle host",
            "create_time": "2020-05-12T21:03:37.122+08:00",
            "default": 1
        }
    ]
}
```
