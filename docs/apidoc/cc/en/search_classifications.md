### Function Description

Query model classification

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type | Required | Description |
| ----- | ---- | -------- | ----------- |
|       |      |          |             |

### Request Parameter Example

```python
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
            "bk_classification_name": "主机管理",
            "bk_classification_type": "inner",
            "id": 1
         }
     ]
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |

#### data

| Field                  | Type   | Description                                                                                                                           |
|------------------------|--------|---------------------------------------------------------------------------------------------------------------------------------------|
| bk_classification_id   | string | Classification ID, used for internal use in the system in English description                                                         |
| bk_classification_name | string | Classification name                                                                                                                   |
| bk_classification_type | string | Used to classify the classification (such as: inner code for built-in classification, empty string for custom classification)         |
| bk_classification_icon | string | Icon of the model classification, the value can refer to [(classIcon.json)](https://chat.openai.com/c/resource_define/classIcon.json) |
| id                     | int    | Data record ID                                                                                                                        |