### Functional description

Query the service template list according to the service id, and add the service classification id for further query

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id           |  int    | yes | Business ID |
| service_category_id         |  int  |no   | Service class ID|
| search         |  string  |no   | Query by service template name. It is blank by default|
| is_exact         |  bool  |no   | Whether to exactly match the service template name. The default value is no. It is used in combination with the search parameter. It is valid when the search parameter is not empty (v3.9.19) |

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|
| sort     |   string |no     | Sort field|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "service_category_id": 1,
    "search": "test2",
    "is_exact": true,
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-name"
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
        "count": 1,
        "info": [
            {
                "bk_biz_id": 1,
                "id": 50,
                "name": "test2",
                "service_category_id": 1,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2019-09-18T20:31:29.607+08:00",
                "last_time": "2019-09-18T20:31:29.607+08:00",
                "bk_supplier_account": "0"
            }
        ]
    }
}
```

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request succeeded or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

#### Data field Description

| Field| Type| Description|
|---|---|---|
|count| int| Total||
|info| array| Return result||

#### Info Field Description

| Field| Type| Description|
|---|---|---|
|bk_biz_id| int| Business ID ||
|id| int| Service template ID||
|name| array| Service template name||
|service_category_id| integer| Service class ID||
|creator| string| Founder||
|modifier| string| Modified by||
|create_time| string| Settling time||
|last_time| string| Repair time||
|bk_supplier_account| string| Vendor ID||
