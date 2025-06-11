### 描述

根据业务id查询集群模板

### 输入参数

| 参数名称             | 参数类型   | 必选 | 描述       |
|------------------|--------|----|----------|
| bk_biz_id        | int    | 是  | 业务ID     |
| set_template_ids | array  | 否  | 集群模板ID数组 |
| page             | object | 否  | 分页信息     |

#### page 字段说明

| 参数名称  | 参数类型   | 必选 | 描述            |
|-------|--------|----|---------------|
| start | int    | 否  | 记录开始位置        |
| limit | int    | 否  | 每页限制条数,最大1000 |
| sort  | string | 否  | 排序字段，'-'表示倒序  |

### 调用示例

```json
{
  "bk_biz_id": 10,
  "set_template_ids":[1, 11],
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "-name"
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
    "count": 2,
    "info": [
      {
        "id": 1,
        "name": "zk1",
        "bk_biz_id": 10,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:09:23.859+08:00",
        "last_time": "2020-03-25T18:59:00.167+08:00",
        "bk_supplier_account": "0"
      },
      {
        "id": 11,
        "name": "q",
        "bk_biz_id": 10,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:10:05.176+08:00",
        "last_time": "2020-03-16T15:10:05.176+08:00",
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

| 参数名称  | 参数类型  | 描述   |
|-------|-------|------|
| count | int   | 总数   |
| info  | array | 返回结果 |

#### info 字段说明

| 参数名称                | 参数类型   | 描述     |
|---------------------|--------|--------|
| id                  | int    | 集群模板ID |
| name                | array  | 集群模板名称 |
| bk_biz_id           | int    | 业务ID   |
| creator             | string | 创建者    |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间   |
| last_time           | string | 更新时间   |
| bk_supplier_account | string | 开发商账号  |
