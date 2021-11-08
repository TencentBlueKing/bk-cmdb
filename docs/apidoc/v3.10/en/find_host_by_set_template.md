### Functional description

find host by set template (v3.8.6)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_biz_id           | int    | Yes      | Business ID           |
| bk_set_ids  | int array  | Yes     | bk_set_id array，the max length is 500 |
| bk_set_template_ids  | int array  | Yes     |  bk_set_template_id array，the max length is 500 |
| bk_set_ids  | int array  | No     |  bk_set_id array，the max length is 500 |
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
    "bk_biz_id": 5,
    "bk_set_template_ids": [
        1,
        3
    ],
    "bk_set_ids": [
        13,
        14
    ],
    "fields": [
        "bk_host_id",
        "bk_cloud_id"
    ],
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
    "code": 0,
    "message": "success",
    "data": {
        "count": 7,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 3
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 4
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 5
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 6
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 7
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 8
            }
        ]
    }
}
```
