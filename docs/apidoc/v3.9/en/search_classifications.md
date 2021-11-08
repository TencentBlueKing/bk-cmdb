### Functional description

search classifications

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description                                                    |
|----------------------|------------|--------|----------------------------------------------------------|
| bk_supplier_account  | string     | No     | Supplier account                                               |

### Request Parameters Example

``` python
{
    "bk_supplier_account": "0"
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
     "data": [
         {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "host managment",
            "bk_classification_type": "inner",
            "id": 1
         }
     ]
}
```

### Return Result Parameters Description

#### data

| Field                   | Type     | Description                                                                                          |
|------------------------|----------|-----------------------------------------------------------------------------------------------|
| bk_classification_id   | string   | Classification ID，English description is used in system                                                              |
| bk_classification_name | string   | Classification name                                                                                        |
| bk_classification_type | string   | For classification （example：inner code is inner classification, null string is custom classification）                           |
| bk_classification_icon | string   | Classification icon, that can refer to[(classIcon.json)](resource_define/classIcon.json) |
| id                     | int      | Data record ID                                                                                   |
