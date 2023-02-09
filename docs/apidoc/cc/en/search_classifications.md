### Functional description

Query model classification

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|

### Request Parameters Example

``` python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
     "data": [
         {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "hosts manage",
            "bk_classification_type": "inner",
            "bk_supplier_account": "0",
            "id": 1
         }
     ]
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

#### data

| Field                   | Type     | Description                                                                                          |
|------------------------|----------|-----------------------------------------------------------------------------------------------|
| bk_classification_id   |  string   | Classification ID, English description for internal use of the system                                                              |
| bk_classification_name | string   | Class name                                                                                        |
| bk_classification_type | string   | Used to classify a classification (for example: Internal code is built-in classification, empty string is user-defined classification)                           |
| bk_classification_icon | string   | Icon of model classification, value can be referred to, value can be referred to [(classIcon.json)](resource_define/classIcon.json)|
| id                     |  int      | Data record ID                                                                                    |
| bk_supplier_account|  string| Developer account|