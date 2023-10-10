### Functional description

Update service template information

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| name            |  string  |And service_category_id are required, or both   | Service template name|
| service_category_id            |  int  |And name are required, or both   | Service class id|
| id         |  int  |yes   | Service template ID|
| bk_biz_id     |   int     | yes  | Business ID |

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "name": "test1",
  "id": 50,
  "service_category_id": 3
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
    "bk_biz_id": 1,
    "id": 50,
    "name": "test1",
    "service_category_id": 3,
    "creator": "admin",
    "modifier": "admin",
    "create_time": "2019-06-05T11:22:22.951+08:00",
    "last_time": "2019-06-05T11:22:22.951+08:00",
    "bk_supplier_account": "0"
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
| data | object |Updated service template information|

#### Data field Description

| Field                | Type   | Description         |
| ------------------- | ------ | ------------ |
| id                  |  int    | Service template ID   |
| name                |  string  |Service template name|
| bk_biz_id           |  int    | Business ID |
| service_category_id | int    | Service class id|
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
