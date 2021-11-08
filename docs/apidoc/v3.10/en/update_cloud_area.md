### Functional description

update cloud area 

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description       | 
|----------------------|------------|--------|-------------|
| bk_cloud_id  | int      | Yes     | cloud area id       |


### Request Parameters Example

``` json
{
    "bk_cloud_id": 5,
	"bk_cloud_name": "cloudname1"
}

```

### Return Result Example


```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}

```
