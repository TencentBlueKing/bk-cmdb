### Functional description

Query the service classification list. Query according to the business ID, and the common service classification will also be returned

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id           |  int    | yes | Business ID |

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
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
    "count": 20,
    "info": [
      {
	"bk_biz_id": 0,
        "id": 16,
        "name": "Apache",
        "bk_root_id": 14,
        "bk_parent_id": 14,
        "bk_supplier_account": "0",
        "is_built_in": true
      },
      {
	"bk_biz_id": 0,
        "id": 19,
        "name": "Ceph",
        "bk_root_id": 18,
        "bk_parent_id": 18,
        "bk_supplier_account": "0",
        "is_built_in": true
      },
      {
	"bk_biz_id": 1,
        "id": 1,
        "name": "Default",
        "bk_root_id": 1,
        "bk_supplier_account": "0",
        "is_built_in": true
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
|count| int| Total|
|info| array| Return result|

#### Info Field Description

| Field| Type| Description|
|---|---|---|
|id| int| Service class ID|
|name| string| Service class name|
|bk_root_id| int| Root service class ID|
|bk_parent_id| int| Parent service class ID|
|is_built_in| bool| Built-in|
|bk_supplier_account | string |Developer account name|
