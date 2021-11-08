### 功能描述

根据业务ID获取业务空闲机, 故障机和待回收模块

### 请求参数

{{ common_args_desc }}


#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | 是     | 业务ID     |

### 请求参数示例

```python

{
    "bk_biz_id":0,
    "bk_supplier_account":"0"
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "bk_set_id": 2,
    "bk_set_name": "空闲机池",
    "module": [
      {
        "bk_module_id": 3,
        "bk_module_name": "空闲机"
      },
      {
        "bk_module_id": 4,
        "bk_module_name": "故障机"
      },
      {
        "bk_module_id": 5,
        "bk_module_name": "待回收"
      }
    ]
  }
}
```

### 返回结果参数说明

#### data说明
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_set_id | int64 | 空闲机, 故障机和待回收模块所属的set的实例ID |
|bk_set_name | string |空闲机, 故障机和待回收模块所属的set的实例名称|

#### module说明
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_module_id | int64 | 空闲机, 故障机或待回收模块的实例ID |
|bk_module_name | string |空闲机, 故障机或待回收模块的实例名称|

