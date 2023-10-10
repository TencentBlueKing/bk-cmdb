### 功能描述

查询业务集拓扑(v3.10.12+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_id    | int    | 是 | 业务集ID |
| bk_parent_obj_id | string | 是 | 需要查询模型的parent对象ID |
| bk_parent_id     | int    | 是 | 需要查询模型的parent ID |

### 请求参数示例

```json
{
  "bk_app_code":"esb_test",
  "bk_app_secret":"xxx",
  "bk_username":"xxx",
  "bk_token":"xxx",
  "bk_biz_set_id":3,
  "bk_parent_obj_id":"bk_biz_set_obj",
  "bk_parent_id":344
}
```

### 返回结果示例

```json
{
  "result":true,
  "code":0,
  "message":"",
  "permission":null,
  "data":[
    {
      "bk_obj_id":"bk_biz_set_obj",
      "bk_inst_id":5,
      "bk_inst_name":"xxx",
      "default":0
    },
    {
      "bk_obj_id":"bk_biz_set_obj",
      "bk_inst_id":6,
      "bk_inst_name":"xxx",
      "default":0
    }
  ],
  "request_id": "dsda1122adasadadada2222"
}
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| data    | array | 请求返回的数据                           |
| request_id    | string | 请求链id    |

#### data

| 名称    | 类型   | 描述              |
| ------- | ------ | --------------- |
| bk_obj_id  | string   | 模型对象ID  |
| bk_inst_id    | int    | 模型实例ID   |
| bk_inst_name | string | 模型实例名称   |
| default    | int | 模型实例分类    |


