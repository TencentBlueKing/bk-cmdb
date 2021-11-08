### Functional description

search cloud area list 

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description       | 
|----------------------|------------|--------|-------------|
|condition|object|No|search condition|
|bk_cloud_id|int|No|cloud id|
|bk_cloud_name|string|No|cloud name|
| page| object| Yes |page info |



#### page params

| Field                 |  Type      | Required	   |  Description       | 
|--------|------------|--------|------------|
|start|int|No|get the data offset location|
|limit|int|Yes|The number of data points in the past is limited, suggest 200|


### Request Parameters Example

``` python
{
    "condition": {
        "bk_cloud_id": 12,
        "bk_cloud_name" "aws",
    },
    "page":{
        "start":0,
        "limit":10
    }
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
    "count": 10,
    "info": [
         {
            "bk_cloud_id": 0,
            "bk_cloud_name": "aws",
            "bk_supplier_account": "0",
            "create_time": "2019-05-20T14:59:48.354+08:00",
            "last_time": "2019-05-20T14:59:48.354+08:00"
        },
        .....
    ]   
  }
}
```

### Return Result Parameters Description

#### data


| Field       | Type     | Description         |
|------------|----------|--------------|
| count| int| the num of record|
| info| object array | the queried cloud list|



#### data.info 

| Field       | Type     | Description         |
|---|---|---|---|
| bk_cloud_id | int | cloud area id |
| bk_cloud_name | string  | cloud area name  | 
| create_time | string | cloud area create time |
| last_time | string | cloud arae last modify time | 


