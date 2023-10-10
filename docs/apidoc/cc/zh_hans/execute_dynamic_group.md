### 功能描述

根据指定动态分组规则查询获取数据 (V3.9.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | 是     | 业务ID |
| id        |  string     | 是     | 动态分组主键ID |
| fields    |  array   | 是     | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输,目标资源不具备指定的字段时该字段将被忽略 |
| disable_counter |  bool | 否     | 是否返回总记录条数，默认返回 |
| page     |  object     | 是     | 分页设置 |

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start     |  int     | 是     | 记录开始位置 |
| limit     |  int     | 是     | 每页限制条数,最大200 |
| sort     |  string   | 否     | 检索排序， 默认按照创建时间排序 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "disable_counter": true,
    "id": "XXXXXXXX",
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_host_name"
    ],
    "page":{
        "start": 0,
        "limit": 10
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 1,
        "info": [
            {
                "bk_obj_id": "host",
                "bk_host_id": 1,
                "bk_host_name": "nginx-1",
                "bk_host_innerip": "10.0.0.1",
                "bk_cloud_id": 0
            }
        ]
    }
}
```

### 返回结果参数

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                          |

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int | 当前规则能匹配到的总记录条数（用于调用者进行预分页，实际单次请求返回数量以及数据是否全部拉取完毕以JSON Array解析数量为准） |
| info      | array        | dict数组，主机实际数据, 当动态分组为host查询时返回host自身属性信息,当动态分组为set查询时返回set信息 |

#### data.info
| 名称  | 类型  | 说明 |
| ---------------- | ------ | ---------------|
| bk_obj_id       | string | 模型id  |               
| bk_host_name           | string | 主机名   |                                                               |                              |
| bk_host_innerip  | string | 内网IP        |                                 
| bk_host_id       | int    | 主机ID        |                                 
| bk_cloud_id      | int    | 云区域        |  