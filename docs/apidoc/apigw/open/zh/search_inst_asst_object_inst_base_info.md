### 描述

查询实例关联模型实例基本信息

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                                        |
|-----------|--------|----|-------------------------------------------|
| fields    | array  | 否  | 指定查询的字段，参数为业务的任意属性，如果不填写字段信息，系统会返回业务的所有字段 |
| condition | object | 是  | 查询条件                                      |
| page      | object | 否  | 分页条件                                      |

#### condition

| 参数名称               | 参数类型   | 必选 | 描述                                                                                |
|--------------------|--------|----|-----------------------------------------------------------------------------------|
| bk_obj_id          | string | 是  | 实例模型ID                                                                            |
| bk_inst_id         | int    | 是  | 实例ID                                                                              |
| association_obj_id | string | 是  | 关联对象的模型ID， 返回association_obj_id模型与bk_inst_id实例有关联的实例基本数据（bk_inst_id,bk_inst_name） |
| is_target_object   | bool   | 否  | bk_obj_id 是否为目标模型， 默认false， 关联关系中的源模型，否则是目标模型                                     |

#### page

| 参数名称  | 参数类型 | 必选 | 描述                 |
|-------|------|----|--------------------|
| start | int  | 否  | 记录开始位置,默认值0        |
| limit | int  | 否  | 每页限制条数,默认值20,最大200 |

### 调用示例

```json
{
  "condition": {
    "bk_obj_id": "bk_switch",
    "bk_inst_id": 12,
    "association_obj_id": "host",
    "is_target_object": true
  },
  "page": {
    "start": 0,
    "limit": 10
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
    "count": 4,
    "info": [
      {
        "bk_inst_id": 1,
        "bk_inst_name": "127.0.0.3"
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

| 参数名称  | 参数类型         | 描述                                                |
|-------|--------------|---------------------------------------------------|
| count | int          | 记录条数                                              |
| info  | object array | 关联对象的模型ID， 实例关联模型的实例基本数据（bk_inst_id,bk_inst_name） |
| page  | object       | 分页信息                                              |

#### data.info 字段说明：

| 参数名称         | 参数类型   | 描述   |
|--------------|--------|------|
| bk_inst_id   | int    | 实例ID |
| bk_inst_name | string | 实例名  |

##### data.info.bk_inst_id,data.info.bk_inst_name 字段说明

不同模型bk_inst_id, bk_inst_name 对应的值

| 模型   | bk_inst_id    | bk_inst_name     |
|------|---------------|------------------|
| 业务   | bk_biz_id     | bk_biz_name      |
| 集群   | bk_set_id     | bk_set_name      |
| 模块   | bk_module_id  | bk_module_name   |
| 进程   | bk_process_id | bk_process_name  |
| 主机   | bk_host_id    | bk_host_inner_ip |
| 通用模型 | bk_inst_id    | bk_inst_name     |
