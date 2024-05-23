### Description

Query model classification

### Parameters

### Request Example

```python
{
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
     "data": [
         {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "主机管理",
            "bk_classification_type": "inner",
            "bk_supplier_account": "0",
            "id": 1
         }
     ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name                   | Type   | Description                                                                                                                           |
|------------------------|--------|---------------------------------------------------------------------------------------------------------------------------------------|
| bk_classification_id   | string | Classification ID, used for internal use in the system in English description                                                         |
| bk_classification_name | string | Classification name                                                                                                                   |
| bk_classification_type | string | Used to classify the classification (such as: inner code for built-in classification, empty string for custom classification)         |
| bk_classification_icon | string | Icon of the model classification, the value can refer to [(classIcon.json)](https://chat.openai.com/c/resource_define/classIcon.json) |
| id                     | int    | Data record ID                                                                                                                        |
| bk_supplier_account    | string | Developer account                                                                                                                     |
