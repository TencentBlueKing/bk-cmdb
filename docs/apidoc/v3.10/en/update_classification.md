### Functional description

update classification

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                   |  Type    | Required	   |  Description                                      |
|------------------------|----------|--------|--------------------------------------------|
| id                     | int      | No     | Record ID of the target data, as a condition for update        |
| bk_classification_name | string   | No     | Classification name  |
| bk_classification_icon | string   | No     | Classfication icon, that can refer to [(classIcon.json)](resource_define/classIcon.json) |




### Request Parameters Example

```python
{
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
    "data": "success"
}
```
