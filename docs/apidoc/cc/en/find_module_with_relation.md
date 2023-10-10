### Functional description

Query modules under Business by criteria (v3.9.7)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  |yes     | Business ID |
| bk_set_ids  |  array  |no     | List of set IDs, up to 200 |
| bk_service_template_ids  |  array  |no     | Service template ID list|
| fields  |   array   | yes  | Module attribute list, which controls the fields in the module information that returns the result|
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
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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
### Return Result Parameters Description
#### response
| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

Data field Description:

| Name     | Type         | Description               |
| -------- | ------------ | ------------------ |
| count    |  int          | Number of records           |
| info | object array |Module actual data|