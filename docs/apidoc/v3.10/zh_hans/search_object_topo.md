### 功能描述

通过对象模型的分类ID查询普通模型拓扑

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                  |  类型      | 必选   |  描述                                    |
|----------------------|------------|--------|------------------------------------------|
| bk_classification_id |string      |是      | 对象模型的分类ID，只能用英文字母序列命名 |


### 请求参数示例

```python
{
    "bk_classification_id": "test"
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": [
        {
           "arrows": "to",
           "from": {
               "bk_classification_id": "bk_host_manage",
               "bk_obj_id": "host",
               "bk_obj_name": "主机",
               "position": "{\"bk_host_manage\":{\"x\":-357,\"y\":-344},\"lhmtest\":{\"x\":163,\"y\":75}}",
               "bk_supplier_account": "0"
           },
           "label": "switch_to_host",
           "label_name": "",
           "label_type": "",
           "to": {
               "bk_classification_id": "bk_network",
               "bk_obj_id": "bk_switch",
               "bk_obj_name": "交换机",
               "position": "{\"bk_network\":{\"x\":-172,\"y\":-160}}",
               "bk_supplier_account": "0"
           }
        }
   ]
}
```

### 返回结果参数说明

#### data

| 字段       | 类型      | 描述                               |
|------------|-----------|------------------------------------|
| arrows     | string    | 取值 to（单向） 或 to,from（双向） |
| label_name | string    | 关联关系的名字                     |
| label      | string    | 表明From通过哪个字段关联到To的     |
| from       | string    | 对象模型的英文id，拓扑关系的发起方 |
| to         | string    | 对象模型的英文ID，拓扑关系的终止方 |
