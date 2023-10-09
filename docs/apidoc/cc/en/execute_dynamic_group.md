### Functional description

Query to obtain data according to specified dynamic grouping rules (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | yes  | Business ID |
| id        |   string     | yes  | Dynamic grouping pk ID|
| fields    |   array   | yes  | Host attribute list, which controls the fields in the host that returns the result, can speed up interface requests and reduce network traffic transmission. If the target resource does not have the specified field, this field will be ignored|
| disable_counter |  bool |no     | Return total number of records; default|
| page     |   object     | yes  | Paging settings|

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start     |   int     | yes  | Record start position|
| limit     |   int     | yes  | Limit number of bars per page, maximum 200|
| sort     |   string   | no     | Retrieve sort, by default by creation time|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
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
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### Return result parameter

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
| count     |  int |The total number of records that can be matched by the current rule (used for pre-paging by the caller, the actual number of returns from a single request and whether all data are pulled are subject to the number of JSON Array parsing)|
| info      |  array        | Dict array, host actual data, returns host own attribute information when dynamic grouping is host query, and returns set information when dynamic grouping is set query|

#### data.info
| Name| Type| Description|
| ---------------- | ------ | ---------------|
| bk_obj_id       |  string |Model id|
| bk_host_name           |  string |Host name   |
| bk_host_innerip  | string |Intranet IP        |
| bk_host_id       |  int    | Host ID        |
| bk_cloud_id      |  int    | Cloud area    |