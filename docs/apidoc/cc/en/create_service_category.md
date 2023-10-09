#### Functional description

Create service classification

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| name            |  string  |yes   | Service class name|
| parent_id         |  int  |no   | Parent node ID|
| bk_biz_id         |  int  |yes   | Business ID |

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "parent_id": 0,
  "bk_biz_id": 1,
  "name": "test101"
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
    "id": 6,
    "name": "test101",
    "root_id": 5,
    "parent_id": 5,
    "bk_supplier_account": "0",
    "is_built_in": false
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
| data | object |New service classification information|

#### Data field Description

| Field| Type| Description|
|---|---|---|
|id| integer| Service class ID|
|root_id| integer| Service classification root node ID|
|parent_id| integer| Service classification parent node ID|
|is_built_in| bool| Is it a built-in node (built-in node can not be edited)|
| bk_biz_id    |  int     | Service ID|
| name    |  string     | Service class name|
| bk_supplier_account|  string| Developer account number|
