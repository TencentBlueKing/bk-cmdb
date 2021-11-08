### 功能描述

根据云区域名字创建云区域

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选   |  描述       |
|----------------------|------------|--------|-------------|
| bk_cloud_name  | string     | 是     |    云区域名字|

### 请求参数示例

``` python
{
	"bk_cloud_name": "test1"
}

```

### 返回结果示例

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

### 返回结果参数说明

#### data

| 字段          | 类型     | 描述     |
|---------------|----------|----------|
| created      | object   |  创建成功，返回信息  |


#### data.created

| 名称    | 类型   | 描述       |
|---------|--------|------------|
| id| int | 云区域id, bk_cloud_id |


