### 功能描述

批量更新namespace (版本：v3.10.23+，权限：容器命名空间编辑)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_biz_id | int| 是 |业务id|
| data |  object | 是 | 包含需要更新的字段|
| ids | array| 是 |在cc中的id唯一标识数组|

#### data

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| labels| map | 否 |标签|
| resource_quotas| array | 否 |命名空间CPU与内存的requests与limits|

### 请求参数示例

```json
{ 
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "ids":[
        1, 2, 3
     ],       
    "data":{ 
      "labels": {
          "test": "test",
          "test2": "test2"
      },
      "resource_quotas": [  
        {
          "hard":{ 
            "memory": "20000Gi", 
            "pods": "100",
            "cpu": "10k" 
          },
          "scopeSelector": {
            "matchExpressions": [
              {
                "values": [ "high" ],
                "operator": "In",
                "scopeName": "PriorityClass"
              }
            ]
          }
        }
      ]
    }
}
```

### 返回结果示例

```json

{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
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
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |
