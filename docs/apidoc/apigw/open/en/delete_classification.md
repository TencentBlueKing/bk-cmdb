### Description

Delete a model category by model category ID. (Permission: Model category deletion permission)

### Parameters

| Name | Type | Required | Description             |
|------|------|----------|-------------------------|
| id   | int  | Yes      | Category data record ID |

**Note** Can only delete an empty model category. Deletion will fail if the category has models.

### Request Example

```python
{
    "id": 13
}
```

### Response Example

#### Deletion Successful

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": "success"
}
```

#### Deletion Failed (Category has models)

```python
{
    "result": false,
    "code": 1101029,
    "data": null,
    "message": "There is a model under the category, not allowed to delete",
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
