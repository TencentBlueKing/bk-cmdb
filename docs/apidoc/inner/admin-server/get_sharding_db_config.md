### 请求方式

GET /migrate/v3/find/system/sharding_db_config

### 描述

查询DB分库配置

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": {
    "master_db": "masteruuid",
    "for_new_data": "slave1uuid",
    "slave_db": {
      "slave1uuid": {
        "name": "slave1",
        "disabled": false,
        "config": {
          "address": "127.0.0.1",
          "port": "27017",
          "user": "user",
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

| 参数名称           | 参数类型              | 描述                         |
|----------------|-------------------|----------------------------|
| master_db      | string            | 主库唯一标识                     |
| for_new_data | string            | 指定新增租户数据写入哪个库，存储这个数据库的唯一标识 |
| slave_db       | map[string]object | 从库唯一标识->从库配置的映射            |

#### data.slave_db[key]

| 参数名称     | 参数类型   | 描述               |
|----------|--------|------------------|
| name     | string | 从库名称，不可以重复       |
| disabled | bool   | 是否禁用，禁用的数据库不允许操作 |
| config   | object | 数据库配置            |

#### data.slave_db[key].config

| 参数名称           | 参数类型   | 描述                                         |
|----------------|--------|--------------------------------------------|
| address        | string | 地址，可以填写多个                                  |
| port           | string | 端口，如果地址里没有包含端口信息则会和地址拼接成实际地址               |
| user           | string | 用户名                                        |
| database       | string | 数据库                                        |
| mechanism      | string | 身份验证机制                                     |
| rs_name        | string | 副本集名称                                      |
| max_open_conns | int    | 最大开启连接数                                    |
| max_idle_conns | int    | 最大空闲连接数                                    |
| socket_timeout | int    | mongo的socket连接的超时时间，以秒为单位，默认10s，最小5s，最大30s |
