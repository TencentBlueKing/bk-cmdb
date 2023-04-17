### Functional description

Update model classification

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                   | Type    | Required   | Description                                      |
|------------------------|----------|--------|--------------------------------------------|
| id                     |  int      | no     | Record ID of the target data as a condition for the update operation       |
| bk_classification_name | string   | no     | Class name|
| bk_classification_icon | string   | no     | Icon of model classification, value can be referred to, value can be referred to [(classIcon.json)](resource_define/classIcon.json)|




### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id": 1,
    "bk_classification_name": "cc_test_new",
    "bk_classification_icon": "icon-cc-business"
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |
