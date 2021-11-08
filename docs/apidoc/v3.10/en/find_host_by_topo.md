### Functional description

find host by topo node (v3.8.13)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field      | Type      | Required | Description                                                  |
| ---------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id  | int       | Yes      | Business ID                                                  |
| bk_obj_id  | int array | Yes      | topo node object ID                                          |
| bk_inst_id | int array | Yes      | topo node instance ID                                        |
| fields     | array     | Yes      | host property list, the specified host property feilds will be returned <br>it can speed up the request and reduce the network payload |
| page       | object    | Yes      | page info                                                    |

#### page

| Field | Type | Required | Description                      |
| ----- | ---- | -------- | -------------------------------- |
| start | int  | Yes      | start record                     |
| limit | int  | Yes      | page limit, maximum value is 500 |

### Request Parameters Example

```json
{
    "bk_biz_id": 5,
    "bk_obj_id": "xxx",
    "bk_inst_id": 10,
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
        "count": 2,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 2
            }
        ]
    }
}
```

### Return Result Parameters Description

#### data

| Field | Type  | Description       |
| ----- | ----- | ----------------- |
| count | int   | the num of record |
| info  | array | host data         |
