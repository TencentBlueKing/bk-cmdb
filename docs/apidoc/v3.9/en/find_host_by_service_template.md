### Functional description

find host by service template (v3.8.6)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_biz_id           | int    | Yes      | Business ID           |
| inst_ids  | int array  | Yes     | bk_set_id array，the max length is 500 |
| bk_service_template_ids  | int array  | Yes     |  bk_service_template_id array，the max length is 500 |
| bk_module_ids  | int array  | No     |  bk_module_id array，the max length is 500 |
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
    "bk_service_template_ids": [
        48,
        49
    ],
    "bk_module_ids": [
        65,
        68
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
        "count": 6,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 2
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
