### 描述

根据业务拓扑实例节点，查询该实例节点下的主机关系信息(权限：业务访问权限)

### 输入参数

| 参数名称        | 参数类型   | 必选 | 描述                                                                                               |
|-------------|--------|----|--------------------------------------------------------------------------------------------------|
| page        | dict   | 是  | 查询条件                                                                                             |
| fields      | array  | 是  | 主机属性列表，控制返回结果的主机里有哪些字段，请按需求填写，可以为bk_biz_id,bk_host_id,bk_module_id,bk_set_id,bk_supplier_account |
| bk_obj_id   | string | 是  | 拓扑节点的模型ID，可以是自定义层级模型ID，set，module等，但不能是业务                                                        |
| bk_inst_ids | array  | 是  | 拓扑节点的实例ID，最多支持50个实例节点                                                                            |
| bk_biz_id   | int    | 是  | 业务id                                                                                             |

#### page

| 参数名称  | 参数类型   | 必选 | 描述             |
|-------|--------|----|----------------|
| start | int    | 是  | 记录开始位置         |
| limit | int    | 是  | 每页限制条数,最大值为500 |
| sort  | string | 否  | 排序字段           |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "page": {
        "start": 0,
        "limit": 10
    },
    "fields": [
        "bk_module_id",
        "bk_host_id"
    ],
    "bk_obj_id": "province",
    "bk_inst_ids": [10,11]
}
```

### 响应示例

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission": null,
  "data":  {
      "count": 1,
      "info": [
          {
              "bk_host_id": 2,
              "bk_module_id": 51
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
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称  | 参数类型  | 描述     |
|-------|-------|--------|
| count | int   | 记录条数   |
| info  | array | 主机关系信息 |

#### info

| 参数名称         | 参数类型 | 描述   |
|--------------|------|------|
| bk_host_id   | int  | 主机id |
| bk_module_id | int  | 模块id |
