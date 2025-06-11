### 描述

查询动态分组列表 (版本：v3.9.6，权限：业务访问权限)

### 输入参数

| 参数名称            | 参数类型   | 必选 | 描述                                                              |
|-----------------|--------|----|-----------------------------------------------------------------|
| bk_biz_id       | int    | 是  | 业务ID                                                            |
| condition       | object | 否  | 查询条件，condition 字段为自定义查询的属性字段, 可以是create_user, modify_user, name |
| disable_counter | bool   | 否  | 是否不返回总记录条数，默认返回                                                 |
| page            | object | 是  | 分页设置                                                            |

#### page

| 参数名称  | 参数类型   | 必选 | 描述               |
|-------|--------|----|------------------|
| start | int    | 是  | 记录开始位置           |
| limit | int    | 是  | 每页限制条数,最大200     |
| sort  | string | 否  | 检索排序， 默认按照创建时间排序 |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "disable_counter": true,
    "condition": {
        "name": "my-dynamic-group"
    },
    "page":{
        "start": 0,
        "limit": 200
    }
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {
        "count": 0,
        "info": [
            {
                "bk_biz_id": 1,
                "id": "XXXXXXXX",
                "name": "my-dynamic-group",
                "bk_obj_id": "host",
                "info": {
                    "condition":[
                			{
                				"bk_obj_id":"set",
                				"condition":[
                					{
                						"field":"default",
                						"operator":"$ne",
                						"value":1
                					}
                				]
                			},
                			{
                				"bk_obj_id":"module",
                				"condition":[
                					{
                						"field":"default",
                						"operator":"$ne",
                						"value":1
                					}
                				]
                			},
                			{
                				"bk_obj_id":"host",
                				"condition":[
                					{
                						"field":"bk_host_innerip",
                						"operator":"$eq",
                						"value":"127.0.0.1"
                					}
                				]
                			}
                    ]
                },
                "name": "test",
                "bk_obj_id": "host",
                "id": "1111",
                "create_user": "admin",
                "create_time": "2018-03-27T16:22:43.271+08:00",
                "modify_user": "admin",
                "last_time": "2018-03-27T16:29:26.428+08:00"
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

| 参数名称  | 参数类型  | 描述                                                                 |
|-------|-------|--------------------------------------------------------------------|
| count | int   | 当前规则能匹配到的总记录条数（用于调用者进行预分页，实际单次请求返回数量以及数据是否全部拉取完毕以JSON Array解析数量为准） |
| info  | array | 自定义查询数据                                                            |

#### data.info

| 参数名称        | 参数类型   | 描述                          |
|-------------|--------|-----------------------------|
| bk_biz_id   | int    | 业务ID                        |
| id          | string | 动态分组主键ID                    |
| name        | string | 动态分组命名                      |
| bk_obj_id   | string | 动态分组的目标资源对象类型,目前可以为host,set |
| info        | object | 动态分组信息                      |
| last_time   | string | 更新时间                        |
| modify_user | string | 修改者                         |
| create_time | string | 创建时间                        |
| create_user | string | 创建者                         |

#### data.info.info.condition

| 参数名称      | 参数类型   | 描述                     |
|-----------|--------|------------------------|
| bk_obj_id | string | 对象名,可以为set,module,host |
| condition | array  | 查询条件                   |

#### data.info.info.condition.condition

| 参数名称     | 参数类型   | 描述                                                |
|----------|--------|---------------------------------------------------|
| field    | string | 对象的字段                                             |
| operator | string | 操作符, op值为eq(相等)/ne(不等)/in(属于)/nin(不属于)/like(模糊匹配) |
| value    | object | 字段对应的值                                            |
