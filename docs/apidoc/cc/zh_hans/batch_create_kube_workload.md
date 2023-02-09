### 功能描述

批量创建workload (版本：v3.10.23+，权限：容器工作负载新建)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_biz_id | int| 是 |业务id|
|kind | string |是 |workload类型，目前支持的workload类型有deployment、daemonSet、statefulSet、gameStatefulSet、gameDeployment、cronJob、job、pods(放不通过workload而直接创建Pod)|
| data | array| 是 |数组, 一次限制创建200|

#### data

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
|bk_namespace_id | int |是 |namespace在cc中的唯一标识|
|name | string |是 |workload名称|
| labels| map | 否 |标签|
| selector| object | 否 |工作负载选择器|
| replicas| 否 | 否 |工作负载实例个数|
| strategy_type| string | 否 |工作负载更新机制|
| min_ready_seconds| int | 否 |指定新创建的 Pod 在没有任意容器崩溃情况下的最小就绪时间， 只有超出这个时间 Pod 才被视为可用|
| rolling_update_strategy| object | 否 |滚动更新策略|

#### selector
| 字段  | 类型  | 必选   |描述         |
| ----- | ----- | ------------|------------ |
|match_labels | map |否|根据label匹配|
|match_expressions | array |否|匹配表达式|

#### match_expressions[0]
| 字段  | 类型  | 必选   |描述         |
| ----- | ----- | ------------|------------ |
|key | string |是|标签的key|
|operator | string |是|操作符，可选值："In"、"NotIn"、"Exists"、"DoesNotExist"|
|values | array |否|字符串数组，如果操作符为"In"或"NotIn",不能为空，如果为"Exists"或"DoesNotExist"，必须为空|

#### rolling_update_strategy
当strategy_type为RollingUpdate，不为空，其他情况为空

| 字段  | 类型  |必选   | 描述         |
| ----- | ----- | ------------|------------ |
|max_unavailable | object |否|最大不可用|
|max_surge | object |否|最大溢出|

#### max_unavailable
| 字段  | 类型  |必选   | 描述         |
| ----- | ----- | ------------|------------ |
|type | int |是|可选值为0(表示int类型)或1(表示string类型)|
|int_val | int |否|当type为0(表示int类型)，不能为空，对应的的int值|
|str_val | string |否|当type为1(表示string类型),不能为空，对应的string值|

#### max_surge
| 字段  | 类型  |必选   | 描述         |
| ----- | ----- | ------------|------------ |
|type | int |是|可选值为0(表示int类型)或1(表示string类型)|
|int_val | int |否|当type为0(表示int类型)，不能为空，对应的的int值|
|str_val | string |否|当type为1(表示string类型),不能为空，对应的string值|

### 请求参数示例
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "kind": "deployment",
    "data": [
      {
        "bk_namespace_id": 1,
        "name": "test",
        "labels": {
            "test": "test",
            "test2": "test2"  
        },
        "selector": {
            "match_labels": {
                "test": "test",
                "test2": "test2" 
            },
            "match_expressions": [
                {
                    "key": "tier",
                    "operator": "In", 
                    "values": ["cache"]
                }
            ]
        },
        "replicas": 1,
        "strategy_type": "RollingUpdate",
        "min_ready_seconds": 1,
        "rolling_update_strategy": {
            "max_unavailable": {
                "type": 0,
                "int_val": 1
            },
            "max_surge": {
                "type": 0,
                "int_val": 1
            }
        }
      }  
    ]   
}
```

### 返回结果示例

```json

{
    "result": true,
    "code": 0,
    "data": {
      "ids": [1]
    },
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```
**注意：**
- 返回的data中的节点ID数组顺序与参数中的数组数据顺序保持一致。

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

#### data

| 字段       | 类型      | 描述     |
|----------- |-----------|----------|
| ids |    array    |  在cc中的唯一标识数组  |
