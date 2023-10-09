### Functional description

Obtain service template according to service template ID

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| service_template_id | int  |yes   | Service template ID|


### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "service_template_id": 51
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
        "bk_biz_id": 3,
        "id": 51,
        "name": "mm2",
        "service_category_id": 12,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-05-26T09:46:15.259Z",
        "last_time": "2020-05-26T09:46:15.259Z",
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
| data | object |Data returned by request|

#### Data field Description

| Field| Type| Description|
|---|---|---|
|bk_biz_id| int| Business ID |
|id| int| Service template ID|
|name| array| Service template name|
|service_category_id| integer| Service class ID|
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
