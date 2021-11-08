### Functional description

execute dynamic group (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               |  Type    | Required | Description                            |
|---------------------|----------|----------|----------------------------------------|
| bk_biz_id           |  int     | Yes      | Business ID                            |
| id                  |  string  | Yes      | Primary key ID of target dynamic group |
| fields              |  array   | Yes      | fields of object                       |
| disable_counter     |  bool    | No       | disable counter flag                   |
| page                |  object  | Yes      | query page settings                    |

#### page

| Field  | Type    | Required  | Description            |
|--------|---------|-----------|------------------------|
| start  | int     | Yes       | start record           |
| limit  | int     | Yes       | page limit, max is 200 |
| sort   | string  | No        | query order by         |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "disable_counter": true,
    "id": "XXXXXXXX",
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_host_name"
    ],
    "page":{
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
    "message": "",
    "data": {
        "count": 1,
        "info": [
            {
                "bk_obj_id": "host",
                "bk_host_id": 1,
                "bk_host_name": "nginx-1",
                "bk_host_innerip": "10.0.0.1",
                "bk_cloud_id": 0
            }
        ]
    }
}
```

### Return Result Parameters Description

#### data

| Field  | Type  | Description       |
|--------|-------|-------------------|
| count  | int   | the num of record |
| info   | array | detail of record  |
