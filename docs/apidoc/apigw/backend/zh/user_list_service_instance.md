### 描述

根据业务id查询服务实例列表,也可以加上模块id等信息查询

### 输入参数

| 参数名称         | 参数类型      | 必选 | 描述                                                        |
|--------------|-----------|----|-----------------------------------------------------------|
| bk_biz_id    | int       | 是  | 业务id                                                      |
| bk_module_id | int       | 否  | 模块ID                                                      |
| bk_host_ids  | int array | 否  | 主机id列表，最多支持1000个主机id                                      |
| selectors    | int       | 否  | label过滤功能，operator可选值: `=`,`!=`,`exists`,`!`,`in`,`notin` |
| page         | object    | 否  | 分页参数                                                      |
| search_key   | string    | 否  | 名字过滤参数，可填写进程名称包含的字符用于模糊搜索                                 |

#### page

| 参数名称  | 参数类型 | 必选 | 描述           |
|-------|------|----|--------------|
| start | int  | 是  | 记录开始位置       |
| limit | int  | 是  | 每页限制条数,最大500 |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 1
  },
  "bk_module_id": 56,
  "bk_host_ids":[26,10],
  "search_key": "",
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  },{
    "key": "key1",
    "operator": "in",
    "values": ["value1", "value2"]
  }]
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
        "bk_biz_id": 1,
        "id": 72,
        "name": "t1",
        "bk_host_id": 26,
        "bk_module_id": 62,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-06-20T22:46:00.69+08:00",
        "last_time": "2019-06-20T22:46:00.69+08:00",
        "bk_supplier_account": "0"
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

#### data 字段说明

| 参数名称  | 参数类型    | 描述   |
|-------|---------|------|
| count | integer | 总数   |
| info  | array   | 返回结果 |

#### info 字段说明

| 参数名称                | 参数类型   | 描述          |
|---------------------|--------|-------------|
| id                  | int    | 服务实例ID      |
| name                | string | 服务实例名称      |
| bk_biz_id           | int    | 业务ID        |
| bk_module_id        | int    | 模块ID        |
| bk_host_id          | int    | 主机ID        |
| creator             | string | 本条数据创建者     |
| modifier            | string | 本条数据的最后修改人员 |
| create_time         | string | 创建时间        |
| last_time           | string | 更新时间        |
| bk_supplier_account | string | 开发商账号       |
