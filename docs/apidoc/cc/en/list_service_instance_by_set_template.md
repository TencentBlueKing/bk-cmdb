### Functional description

Obtaining the service instance list according to the set template id

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            |  int  |yes   | Business ID |
| set_template_id            |  int  |yes   | Set template ID |
| page       |   object    | no     | Query criteria|

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500, suggest 200.|

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",  
  "bk_biz_id": 1,
  "set_template_id":1,
  "page": {
    "start": 0,
    "limit": 10
  }
}
```

### Return Result Example

```python
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
                "bk_biz_id": 3,
                "id": 1,
                "name": "10.0.0.1_lgh-process-1",
                "labels": null,
                "service_template_id": 50,
                "bk_host_id": 1,
                "bk_module_id": 59,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2020-10-09T02:46:25.002Z",
                "last_time": "2020-10-09T02:46:25.002Z",
                "bk_supplier_account": "0"
            },
            {
                "bk_biz_id": 3,
                "id": 3,
                "name": "127.0.122.2_lgh-process-1",
                "labels": null,
                "service_template_id": 50,
                "bk_host_id": 3,
                "bk_module_id": 59,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2020-10-09T03:04:19.859Z",
                "last_time": "2020-10-09T03:04:19.859Z",
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
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

#### Data field Description

| Field| Type| Description|
|---|---|---|
|count| int| Total|
|info| array| Return result|

#### Info Field Description

| Field| Type| Description|
|---|---|---|
|id| int| Service instance ID|
|name| string| Service instance name|
|bk_biz_id| int| Business ID |
|bk_module_id| int| Model id|
|service_template_id| int| Service template ID|
| labels           |  map  |Tag information|
|bk_host_id| int| Host id|
| creator              |  string             | Creator of this data                                                                                 |
| modifier             |  string             | The last person to modify this piece of data            |
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string       | Developer account number|
