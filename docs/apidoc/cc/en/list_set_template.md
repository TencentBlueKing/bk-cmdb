### Functional description

Query set template by business id

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type   | Required| Description           |
| ------------------- | ------ | ---- | -------------- |
| bk_biz_id           |  int    | yes | Business ID |
| set_template_ids    |  array  |no   | Set template ID array |
| page                |  object |no   | Paging information       |

#### Page field Description

| Field| Type   | Required| Description                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | no   | Record start position          |
| limit | int    | no   | Limit bars per page, Max. 1000|
| sort  | string |no   | Sort field,'inverted' for reverse order|


### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_supplier_account": "0",
  "bk_biz_id": 10,
  "set_template_ids":[1, 11],
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
    "count": 2,
    "info": [
      {
        "id": 1,
        "name": "zk1",
        "bk_biz_id": 10,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:09:23.859+08:00",
        "last_time": "2020-03-25T18:59:00.167+08:00",
        "bk_supplier_account": "0"
      },
      {
        "id": 11,
        "name": "q",
        "bk_biz_id": 10,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:10:05.176+08:00",
        "last_time": "2020-03-16T15:10:05.176+08:00",
        "bk_supplier_account": "0"
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

#### Data field Description

| Field| Type| Description     |
| ----- | ----- | -------- |
| count | int   | Total     |
| info  | array |Return result|

#### Info Field Description

| Field                | Type   | Description         |
| ------------------- | ------ | ------------ |
| id                  |  int    | Set template ID |
| name                |  array  |Set template name|
| bk_biz_id           |  int    | Business ID |
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
