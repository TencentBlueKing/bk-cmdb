### Functional description

Query hosts under topology nodes (v3.8.13)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field        | Type| Required   | Description      |
|------------|--------|--------|------------|
| bk_biz_id  | int    | yes  | Business ID |
| bk_obj_id  | string |yes     | Topology node model ID, which can not be biz|
| bk_inst_id | int    | yes  | Topology node instance ID|
| fields     |  array  |yes     | Host attribute list, which controls which fields are in the host that returns the result, can speed up interface requests and reduce network traffic transmission   |
| page       |  object |yes     | Paging information|

#### Page field Description

| Field| Type   | Required| Description                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | yes | Record start position          |
| limit | int    | yes | Limit bars per page, Max. 500|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int       | Number of records|
| info      |  array     | Host actual data|
