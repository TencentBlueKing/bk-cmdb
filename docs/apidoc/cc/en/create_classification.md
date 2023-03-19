### Functional description

Add model classification

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                       | Type      | Required   | Description                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_classification_id       |  string     | yes  | Classification ID, English description for internal use of the system           |
| bk_classification_name     |  string     | yes     | Class name     |
| bk_classification_icon     |  string     | no     | Icon for model classification|



### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_classification_id": "cs_test",
    "bk_classification_name": "test_name",
    "bk_classification_icon": "icon-cc-business"
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "data": {
        "id": 11,
        "bk_classification_id": "cs_test",
        "bk_classification_name": "test_name",
        "bk_classification_type": "",
        "bk_classification_icon": "icon-cc-business",
        "bk_supplier_account": ""
    },
    "message": "success",
    "permission": null,
    "request_id": "76e9134a953b4055bb55853bb248dcb7"
    }
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field       | Type      | Description                |
|----------- |-----------|--------------------|
| id         |  int       | ID of new data record   |
| bk_classification_id       |  string          | Classification ID, English description for internal use of the system           |
| bk_classification_name     |  string        | Class name     |
| bk_classification_icon     |  string         | Icon for model classification|
| bk_classification_type | string   | Used to classify a classification (for example: Internal code is built-in classification, empty string is user-defined classification)                           |
| bk_supplier_account|  string| Developer account number|