### Functional description

Query the service instance list according to the service id, or add information such as module id to query

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            |  int  |yes   | Business ID |
| bk_module_id         |  int  |no   | Module ID|
| selectors            |  int  |no   | Label filtering function, operator optional value: `=`,`!=`,` exists`,`!`,` in`,`notin`|
| page         |  object  | No   | Paging parameter|
| search_key         |  string  | No   | Name filter parameter|

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|

### Request Parameters Example

```python

{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 1
  },
  "bk_module_id": 56,
  "search_key": "",
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  },{
    "key": "key1",
    "operator": "in",
    "values": ["value1", "value2"]
  }]
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
    "count": 1,
    "info": [
      {
        "bk_biz_id": 1,
        "id": 72,
        "name": "t1",
        "bk_host_id": 26,
        "bk_module_id": 62,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-06-20T22:46:00.69+08:00",
        "last_time": "2019-06-20T22:46:00.69+08:00",
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
|count| integer| Total||
|info| array| Return result||

#### Info Field Description

| Field| Type| Description|
|---|---|---|
|id| int| Service instance ID|
|name| string| Service instance name|
|bk_biz_id| int| Business ID |
|bk_module_id| int| Module ID|
|bk_host_id| int| Host ID|
| creator              |  string             | Creator of this data                                                                                 |
| modifier             |  string             | The last person to modify this piece of data            |
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string       | Developer account number|
