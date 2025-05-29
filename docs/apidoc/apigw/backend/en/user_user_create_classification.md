### Description

Add model classification (Permission: Model Group New Permission)

### Parameters

| Name                   | Type   | Required | Description                                                           |
|------------------------|--------|----------|-----------------------------------------------------------------------|
| bk_classification_id   | string | Yes      | Classification ID, English description for internal use in the system |
| bk_classification_name | string | Yes      | Classification name                                                   |
| bk_classification_icon | string | No       | Model classification icon                                             |

### Request Example

```python
{
    "bk_classification_id": "cs_test",
    "bk_classification_name": "test_name",
    "bk_classification_icon": "icon-cc-business"
}
```

### Response Example

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
    }
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Request return data                                                         |

#### data

| Name                   | Type   | Description                                                                                                                  |
|------------------------|--------|------------------------------------------------------------------------------------------------------------------------------|
| id                     | int    | ID of the newly added data record                                                                                            |
| bk_classification_id   | string | Classification ID, English description for internal use in the system                                                        |
| bk_classification_name | string | Classification name                                                                                                          |
| bk_classification_icon | string | Model classification icon                                                                                                    |
| bk_classification_type | string | Used to classify classifications (e.g., inner code for built-in classifications, an empty string for custom classifications) |
| bk_supplier_account    | string | Developer account                                                                                                            |
