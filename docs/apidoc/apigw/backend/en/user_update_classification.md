### Description

Update model classification (Permission: Model group editing permission)

### Parameters

| Name                   | Type   | Required | Description                                                                                                                    |
|------------------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------|
| id                     | int    | No       | The record ID of the target data, used as a condition for the update operation                                                 |
| bk_classification_name | string | No       | Classification name                                                                                                            |
| bk_classification_icon | string | No       | Model classification icon, the value can refer to [(classIcon.json)](https://chat.openai.com/c/resource_define/classIcon.json) |

### Request Example

```python
{
    "id": 1,
    "bk_classification_name": "cc_test_new",
    "bk_classification_icon": "icon-cc-business"
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": "success"
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
