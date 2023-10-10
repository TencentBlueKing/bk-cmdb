### Functional description

Get host under service template (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  |yes     | Business ID |
| bk_service_template_ids  |  array  |yes     | List of service template IDs, up to 500|
| bk_module_ids  |  array  |no     | Module ID list, up to 500|
| fields  |   array   | yes  | Host attribute list, which controls what fields are in the module information that returns the result|
| page       |   object    | yes  | Paging information|

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
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

#### data.info

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| bk_cloud_id     |  int       | Cloud area id |
| bk_host_id      |  int     | Host id|

