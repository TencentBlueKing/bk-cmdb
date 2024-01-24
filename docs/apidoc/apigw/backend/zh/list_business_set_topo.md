### 描述

查询业务集拓扑(版本：v3.10.12+，权限：业务集访问)

### 输入参数

| 参数名称             | 参数类型   | 必选 | 描述                |
|------------------|--------|----|-------------------|
| bk_biz_set_id    | int    | 是  | 业务集ID             |
| bk_parent_obj_id | string | 是  | 需要查询模型的parent对象ID |
| bk_parent_id     | int    | 是  | 需要查询模型的parent ID  |

### 调用示例

```json
{
  "bk_biz_set_id":3,
  "bk_parent_obj_id":"bk_biz_set_obj",
  "bk_parent_id":344
}
```

### 响应示例

```json
{
  "result":true,
  "code":0,
  "message":"success",
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
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | array  | 请求返回的数据                    |

#### data

| 参数名称         | 参数类型   | 描述     |
|--------------|--------|--------|
| bk_obj_id    | string | 模型对象ID |
| bk_inst_id   | int    | 模型实例ID |
| bk_inst_name | string | 模型实例名称 |
| default      | int    | 模型实例分类 |
