### 描述

查询给定模型的实例信息(权限：模型实例查询权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述      |
|---------------------|--------|----|---------|
| bk_supplier_account | string | 否  | 开发商账号   |
| bk_obj_id           | string | 是  | 自定义模型ID |
| fields              | array  | 否  | 指定查询的字段 |
| condition           | dict   | 否  | 查询条件    |
| page                | dict   | 否  | 分页条件    |

#### page

| 参数名称  | 参数类型   | 必选 | 描述           |
|-------|--------|----|--------------|
| start | int    | 是  | 记录开始位置       |
| limit | int    | 是  | 每页限制条数,最大200 |
| sort  | string | 否  | 排序字段         |

#### fields参数说明

参数为查询的目标实例对应的模型定义的所有字段

#### condition 参数说明

condition 参数为查询的目标实例对应的模型定义的所有字段

### 调用示例

```json
{
  "bk_obj_id": "bk_switch",
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "bk_inst_id"
  },
  "fields": [
    "bk_asset_id",
    "bk_inst_id",
    "bk_inst_name",
    "bk_obj_id"
  ],
  "condition": {
    "bk_inst_name": "aaa"
  }
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 1,
    "info": [
      {
        "bk_asset_id": "aaa",
        "bk_inst_id": 3,
        "bk_inst_name": "aaa",
        "bk_obj_id": "bk_switch"
      }
    ]
  }
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data

| 参数名称  | 参数类型  | 描述            |
|-------|-------|---------------|
| count | int   | info 集合中元素的数量 |
| info  | array | 查询的模型的实例集合    |

#### info

| 参数名称                | 参数类型   | 描述      |
|---------------------|--------|---------|
| bk_inst_id          | int    | 实例ID    |
| bk_inst_name        | string | 实例名     |
| bk_supplier_account | string | 开发商账号   |
| bk_obj_id           | string | 自定义模型ID |
| bk_created_at       | string | 创建时间    |
| bk_updated_at       | string | 更新时间    |
| bk_created_by       | string | 创建人     |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
