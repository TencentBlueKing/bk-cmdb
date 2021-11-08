### 功能描述

获取主线模型的业务拓扑

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|

### 请求参数示例

```python
{
    "bk_supplier_account":"0"
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": [
    {
      "bk_obj_id": "biz",
      "bk_obj_name": "业务",
      "bk_supplier_account": "0",
      "bk_next_obj": "set",
      "bk_next_name": "集群",
      "bk_pre_obj_id": "",
      "bk_pre_obj_name": ""
    },
    {
      "bk_obj_id": "set",
      "bk_obj_name": "集群",
      "bk_supplier_account": "0",
      "bk_next_obj": "module",
      "bk_next_name": "模块",
      "bk_pre_obj_id": "biz",
      "bk_pre_obj_name": "业务"
    },
    {
      "bk_obj_id": "module",
      "bk_obj_name": "模块",
      "bk_supplier_account": "0",
      "bk_next_obj": "host",
      "bk_next_name": "主机",
      "bk_pre_obj_id": "set",
      "bk_pre_obj_name": "集群"
    },
    {
      "bk_obj_id": "host",
      "bk_obj_name": "主机",
      "bk_supplier_account": "0",
      "bk_next_obj": "",
      "bk_next_name": "",
      "bk_pre_obj_id": "module",
      "bk_pre_obj_name": "模块"
    }
  ]
}
```

### 返回结果参数说明

#### data
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_obj_id | string | 模型的唯一ID |
|bk_obj_name | string |模型名称|
|bk_supplier_account | string |开发商帐户名称|
|bk_next_obj | string |当前模型的下一个模型唯一ID|
|bk_next_name | string |当前模型的下一个模型名称|
|bk_pre_obj_id | string |当前模型的前一个模型的唯一ID|
|bk_pre_obj_name | string |当前模型的前一个模型的名称|
