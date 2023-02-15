# 更新服务实例

### 接口地址
PUT /api/v3/updatemany/proc/service_instance/biz/{bk_biz_id}

### 请求参数

  | 字段         | 类型   | 必选 | 描述         |
  | ------------ | ------ | ---- | ------------ |
  | bk_biz_id   | int | 是   | 服务实例所属的业务ID   |
  | data | obj array| 是   | 所有要更新的服务实例数据，类型为数组  |

#### obj字段
  | 字段         | 类型   | 必选 | 描述         |
  | ------------ | ------ | ---- | ------------ |
  | service_instance_id   | int | 是   | 服务实例ID   |
  | update | dict| 是   | 单个服务实例要更新的数据，当前支持更新的服务实例字段有name|
  
### 请求参数示例

```json
{
    "data": [
        {
            "service_instance_id": 1,
            "update": {
                "name": "serviceInstance001"
            }
        },
                {
            "service_instance_id": 2,
            "update": {
                "name": "serviceInstance002"
            }
        }
    ]
}
``` 

### 返回结果示例
```json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "permission": null,
    "data": null
}
```