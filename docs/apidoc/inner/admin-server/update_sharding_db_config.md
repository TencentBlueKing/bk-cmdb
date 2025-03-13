### 请求方式

PUT /migrate/v3/update/system/sharding_db_config

### 描述

更新DB分库配置

### 输入参数

| 参数名称            | 参数类型              | 必选 | 描述                                          |
|-----------------|-------------------|----|---------------------------------------------|
| for_new_data  | string            | 否  | 指定新增租户数据写入哪个库。对于存量数据库指定它的唯一标识。对于新增的从库指定它的名称 |
| create_slave_db | object array      | 否  | 新增的从库配置数组                                   |
| update_slave_db | map[string]object | 否  | 更新的从库唯一标识->从库配置的映射                          |

#### create_slave_db[n] 和 update_slave_db[key]

| 参数名称     | 参数类型   | 必选    | 描述               |
|----------|--------|-------|------------------|
| name     | string | 新增时必填 | 从库名称，不可以重复       |
| disabled | bool   | 否     | 是否禁用，禁用的数据库不允许操作 |
| config   | object | 新增时必填 | 数据库配置            |

#### data.slave_db[key].config

| 参数名称           | 参数类型   | 必选 | 描述                                         |
|----------------|--------|----|--------------------------------------------|
| address        | string | 是  | 地址，可以填写多个                                  |
| port           | string | 是  | 端口，如果地址里没有包含端口信息则会和地址拼接成实际地址               |
| user           | string | 是  | 用户名                                        |
| password       | string | 是  | 密码                                         |
| database       | string | 是  | 数据库                                        |
| mechanism      | string | 是  | 身份验证机制                                     |
| rs_name        | string | 是  | 副本集名称                                      |
| max_open_conns | int    | 是  | 最大开启连接数                                    |
| max_idle_conns | int    | 是  | 最大空闲连接数                                    |
| socket_timeout | int    | 是  | mongo的socket连接的超时时间，以秒为单位，默认10s，最小5s，最大30s |

### 调用示例

```json
{
  "for_new_data": "slave1uuid",
  "create_slave_db": [
    {
      "name": "slave2",
      "disabled": false,
      "config": {
        "address": "127.0.0.1",
        "port": "27017",
        "user": "user",
        "password": "password",
        "database": "cmdb",
        "mechanism": "SCRAM-SHA-1",
        "max_open_conns": 100,
        "max_idle_conns": 100,
        "rs_name": "rs0",
        "socket_timeout": 10
      }
    }
  ],
  "update_slave_db": {
    "slave1uuid": {
      "name": "slave1",
      "disabled": false,
      "config": {
        "address": "127.0.0.1",
        "port": "27018",
        "user": "user",
        "password": "password",
        "database": "cmdb",
        "mechanism": "SCRAM-SHA-1",
        "max_open_conns": 100,
        "max_idle_conns": 100,
        "rs_name": "rs0",
        "socket_timeout": 10
      }
    }
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
  "data": null
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
