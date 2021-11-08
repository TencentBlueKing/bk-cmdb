### 功能描述

对指定业务id下通过集群id批量删除集群

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | 是     | 业务ID     |
| inst_ids  | int array  | 是     | 集群ID集合 |

### 请求参数示例

```python
{
    "bk_biz_id":0,
    "delete": {
    "inst_ids": [123]
    }
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}
```
