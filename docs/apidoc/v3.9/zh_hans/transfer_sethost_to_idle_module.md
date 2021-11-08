### 功能描述

根据业务id,集群id,模块id,将指定业务集群模块下的主机上交到业务的空闲机模块

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段          |  类型      | 必选     |  描述    |
|---------------|------------|----------|----------|
| bk_biz_id     | int        | 是       | 业务id   |
| bk_set_id     | int        | 是       | 集群id   |
| bk_module_id  | int        | 是       | 模块id   |


### 请求参数示例

```python
{
    "bk_biz_id":10,
    "bk_module_id":58,
    "bk_set_id":1
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "sucess"
}
```
