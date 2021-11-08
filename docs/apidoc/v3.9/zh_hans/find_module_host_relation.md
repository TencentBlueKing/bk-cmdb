### 功能描述

根据模块ID查询主机和模块的关系(v3.8.7)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段          | 类型         | 必选 | 描述                                         |
| ------------- | ------------ | ---- | -------------------------------------------- |
| bk_biz_id     | int          | 是   | 业务ID                                       |
| bk_module_ids | int array    | 是   | 模块ID数组，最多200条                        |
| module_fields | string array | 是   | 模块属性列表，控制返回结果的模块里有哪些字段 |
| host_fields   | string array | 是   | 主机属性列表，控制返回结果的主机里有哪些字段 |
| page          | object       | 是   | 分页参数                                     |

#### page

| 字段  | 类型 | 必选 | 描述                  |
| ----- | ---- | ---- | --------------------- |
| start | int  | 否   | 记录开始位置,默认值0  |
| limit | int  | 是   | 每页限制条数,最大1000 |

**注: 一个模块下的主机关系可能会拆分多次返回，分页方式是按主机ID排序进行分页。**

### 请求参数示例

```json
{
    "bk_module_ids": [
        1,
        2,
        3
    ],
    "module_fields": [
        "bk_module_id",
        "bk_module_name"
    ],
    "host_fields": [
        "bk_host_innerip",
        "bk_host_id"
    ],
    "page": {
        "start": 0,
        "limit": 500
    }
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "count": 2,
    "relation": [
      {
        "host": {
          "bk_host_id": 1,
          "bk_host_innerip": "127.0.0.1",
        },
        "modules": [
          {
            "bk_module_id": 1,
            "bk_module_name": "m1",
          },
          {
            "bk_module_id": 2,
            "bk_module_name": "m2",
          }
        ]
      },
      {
        "host": {
          "bk_host_id": 2,
          "bk_host_innerip": "127.0.0.2",
        },
        "modules": [
          {
            "bk_module_id": 3,
            "bk_module_name": "m3",
          }
        ]
      }
    ]
  }
}
```

### 返回结果参数说明

| 名称    | 类型   | 说明                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| data    | object | 请求返回的数据                             |

data 字段说明：

| 名称     | 类型         | 说明               |
| -------- | ------------ | ------------------ |
| count    | int          | 记录条数           |
| relation | object array | 主机和模块实际数据 |


info 字段说明:

| 名称    | 类型         | 说明               |
| ------- | ------------ | ------------------ |
| host    | object       | 主机数据           |
| modules | object array | 主机所属的模块信息 |
