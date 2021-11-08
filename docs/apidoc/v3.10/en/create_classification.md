### Functional description

create classification

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                       |  Type      | Required	   |  Description                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_classification_id       | string     | Yes     | Classification ID, English description is used in system            |
| bk_classification_name     | string     | Yes     | Classification name      |
| bk_classification_icon     | string     | No     | Classification icon, that can refer to '[(classIcon.json)](resource_define/classIcon.json)'|



### Request Parameters Example

```python
{
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
    "message": "",
    "data": {
        "id": 18
    }
}
```

### Return Result Parameters Description

#### data

| Field       | Type      | Description                |
|----------- |-----------|--------------------|
| id         | int       |  ID of the new data record   |
