## datacollection
> 用于数据采集

## 支持的采集类型
- `middleware` 可以录入模型/模型属性/实例
- `hostsnap`  用于采集/更新主机信息
- `netcollect` 用于录入/更新数据到 `cc_NetcollectReport` collection

## 注意事项
- 实例录入必须有 `bk_inst_key` 字段, 否则实例无法录入

## 测试
> `datacollection` 模块提供了一个 `mockmanager` 服务用于测试。
> 改模块的实现原理是启动一个http服务(绑定地址localhost:12140)，收到http请求后，
> 将消息放入对应消息类型的队列中 (scene_server/datacollection/datacollection/manager.go:65）

- 编译参数配置 `go build -i -ldflags "-X configcenter/src/common/version.CCRunMode=dev"`
- 向 datacollection 模块发送数据

```python
# -*- coding: utf8 -*-
import json

import requests

msg = {
  "host": {
    "bk_supplier_account": "0"
  },
  "data": {
    "meta": {
      "model": {
        "bk_classification_id": "middelware",
        "bk_obj_id": "test1",
        "bk_obj_name": "test1n",
        "bk_supplier_account": "0"
      },
      "fields": {
        "bk_inst_name": {
          "bk_property_name": "实例名",
          "bk_property_type": "longchar"
        },
        "field1": {
          "bk_property_name": "field1",
          "bk_property_type": "longchar"
        }
      }
    },
    "data": {
      "bk_inst_key": "test1",
      "field1": "field 1",
      "bk_inst_name": "inst1"
    }
  }
}
data = {
    "name": "middleware",
    "mesg": json.dumps(msg)
}
url = "http://127.0.0.1:12140"
response = requests.request("POST", url, data=json.dumps(data))
print(response.status_code, response.text)
```
