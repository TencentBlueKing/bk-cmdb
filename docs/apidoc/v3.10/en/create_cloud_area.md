### Functional description

create cloud area

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description       | 
|----------------------|------------|--------|-------------|
| bk_cloud_name  | string     | Yes     |    cloud area name |

### Request Parameters Example

``` python
{
	"bk_cloud_name": "test1"
}

```


### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "created": {
            "id": 2
        }
    }
}
```

### Return Result Parameters Description

#### data

| Field                 |  Type    	   |  Description       | 
|---------------|----------|----------|----------|
| created      | object   |  create sucess, return information  |


#### data.created

| Field                 |  Type    	   |  Description       | 
|---------|--------|------------|
| id| int | cloud area id, bk_cloud_id |


