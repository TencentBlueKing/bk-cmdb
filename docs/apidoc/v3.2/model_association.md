# 模型关联

### 查询模型关联

* API：POST / api / {version} / find / objectassociation
* API名称：search_object_association
* 功能说明：
  * 中文：查询模型关联
  * 英语：搜索对象之间的关联
* 输入体
```
{
    "condition": {
        "bk_asst_id": "belong",
        "bk_obj_id": "bk_switch",
        "bk_asst_obj_id": "bk_host"
    },
    "metadata":{
        "label":{
            "bk_biz_id":"1"
        }
    }
}
```
* 输入字段说明

| 字段名         | 类型 | 必填 | 默认值 | 说明               | 描述 |
| -------------- | ---- | ---- | ------ | ------------------ | ---- |
| bk_asst_id     | 串   | 否   | 无     | 关联类型的唯一标识 |
| bk_obj_id      | 串   | 否   | 无     | 源模型ID           |
| bk_asst_obj_id | 串   | 否   | 无     | 目标模型ID         |

* 产量
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": [
        {
            "id": 1,
            "bk_obj_asst_id": "bk_switch_belong_bk_host",
            "bk_obj_asst_name": "",
            "bk_asst_id": "belong",
            "bk_asst_name": "属于",
            "bk_obj_id": "bk_switch",
            "bk_obj_name": "交换机",
            "bk_asst_obj_id": "bk_host",
            "bk_asst_obj_name": "主机",
            "mapping": "1:n",
            "on_delete": "none"
        }
    ]
}
```
注：以上JSON数据中各字段的取值仅为示例数据。

* 输出字段说明

| 字段          | 类型 | 说明                                | 描述                                 |
| ------------- | ---- | ----------------------------------- | ------------------------------------ |
| 结果          | 布尔 | ture：成功，假：失败                | true：成功，错误：失败               |
| bk_error_code | INT  | 错误编码.0表示成功，> 0表示失败错误 | 错误代码。0表示成功，> 0表示失败代码 |
| bk_error_msg  | 串   | 请求失败返回的错误信息              | 失败请求的错误消息                   |
| 数据          | 宾语 | 结果数据                            | 结果                                 |

* 数据说明

| 字段             | 类型 | 说明                                                                                                      | 描述       |
| ---------------- | ---- | --------------------------------------------------------------------------------------------------------- | ---------- |
| ID               | INT  | 自增ID                                                                                                    | 自动递增ID |
| bk_obj_asst_id   | 串   | 唯一标识，自动生成。规则：源模型英文ID +关联类型英文标识+目标模型英文ID。由前端生成传入，后端只做唯一校验 |            |
| bk_obj_asst_name | 串   | 别名                                                                                                      |            |
| bk_asst_id       | 串   | 关联类型                                                                                                  |            |
| bk_asst_name     | 串   | 显示的名称                                                                                                |            |
| bk_obj_id        | 串   | 源模型ID                                                                                                  |            |
| bk_obj_name      | 串   | 源模型ID                                                                                                  |            |
| bk_asst_obj_id   | 串   | 目标模型名称                                                                                              |            |
| bk_asst_obj_name | 串   | 源模型名称                                                                                                |            |
| 制图             | 枚举 | 关联映射，任选：[1：1,1：n，n：n]                                                                         |            |
| on_delete        | 枚举 | 删除时对实例的动作，可选none，delete_src，delete_dest，分别表示不处理，删除源实例，删除目标实             |            |

### 添加模型关联


* API：POST / api / {version} / create / objectassociation
* API名称：create_object_association
* 功能说明：
  * 中文：添加模型关联
  * 英语：在对象之间创建关联
* 输入体

```
{
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_obj_asst_name": "",
    "bk_asst_id": "belong",
    "bk_obj_id": "bk_switch",
    "bk_asst_obj_id": "bk_host",
    "mapping": "1:n",
    "on_delete": "none",
    "metadata":{
        "label":{
            "bk_biz_id":"1"
        }
    }
}
```
* 输入字段说明

| 字段名           | 类型 | 必填 | 默认值 | 说明                                                                                                      | 描述 |
| ---------------- | ---- | ---- | ------ | --------------------------------------------------------------------------------------------------------- | ---- |
| bk_obj_asst_id   | 串   | 是   | 无     | 唯一标识，自动生成。规则：源模型英文ID +关联类型英文标识+目标模型英文ID。由前端生成传入，后端只做唯一校验 |      |
| bk_obj_asst_name | 串   | 否   | 无     | 别名                                                                                                      |      |
| bk_asst_id       | 串   | 是   | 无     | 关联类型                                                                                                  |      |
| bk_obj_id        | 串   | 是   | 无     | 源模型ID                                                                                                  |      |
| bk_asst_obj_id   | 串   | 是   | 无     | 目标模型ID                                                                                                |      |
| 制图             | 枚举 | 是   | 无     | 关联映射，任选：[1：1,1：n，n：n]                                                                         |      |
| on_delete        | 枚举 | 否   | 没有   | 删除时的动作，可选无，delete_src，delete_dest                                                             |      |

* 产量
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "id": 1038
    }
}
```
* 输出字段说明

| 字段          | 类型 | 说明                                | 描述                                 |
| ------------- | ---- | ----------------------------------- | ------------------------------------ |
| 结果          | 布尔 | ture：成功，假：失败                | true：成功，错误：失败               |
| bk_error_code | INT  | 错误编码.0表示成功，> 0表示失败错误 | 错误代码。0表示成功，> 0表示失败代码 |
| bk_error_msg  | 串   | 请求失败返回的错误信息              | 失败请求的错误消息                   |
| 数据          | 宾语 | 操作结果                            | 结果                                 |

* 数据字段说明

| 字段 | 类型 | 说明   | 描述       |
| ---- | ---- | ------ | ---------- |
| ID   | INT  | 自增ID | 自动递增ID |

### 编辑模型关联

* API：PUT / api / {version} / update / objectassociation / {id}
* API名称：update_object_association
* 功能说明：
  * 中文：编辑关联类型，只有输入体里的任一一个可以更新。
  * 英语：更新对象之间的关联

* 输入体
```
{
    "bk_asst_name": "属于",
    "bk_asst_id":"belong",
    "on_delete":""// 具体枚举值见上。
    "metadata":{
        "label":{
            "bk_biz_id":"1"
        }
    }
}
```
* 输入字段说明

| 字段名       | 类型 | 必填 | 默认值 | 说明         | 描述       |
| ------------ | ---- | ---- | ------ | ------------ | ---------- |
| ID           | INT  | 是   | 无     | 自增ID       | 自动递增ID |
| bk_asst_name | 串   | 否   | 无     | 显示的名称   | 协会的名称 |
| bk_asst_id   | 串   | 否   | 无     | 关联类型     |
| on_delete    | 串   | 否   | 无     | 删除时的动作 |

* 产量
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```
注：以上JSON数据中各字段的取值仅为示例数据。

* 输出字段说明

| 字段          | 类型 | 说明                                | 描述                                 |
| ------------- | ---- | ----------------------------------- | ------------------------------------ |
| 结果          | 布尔 | ture：成功，假：失败                | true：成功，错误：失败               |
| bk_error_code | INT  | 错误编码.0表示成功，> 0表示失败错误 | 错误代码。0表示成功，> 0表示失败代码 |
| bk_error_msg  | 串   | 请求失败返回的错误信息              | 失败请求的错误消息                   |
| 数据          | 串   | 结果数据                            | 结果                                 |

### 删除模型关联

* API：DELETE / api / {version} / delete / objectassociation / {id}
* API名称：delete_object_association
* 功能说明：
  * 中文：删除关联类型
  * 英语：删除对象之间的关联
* 输入体

无

* 输入字段说明

无

* 产量
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```
注：以上JSON数据中各字段的取值仅为示例数据。

* 输出字段说明

| 字段          | 类型 | 说明                                | 描述                                 |
| ------------- | ---- | ----------------------------------- | ------------------------------------ |
| 结果          | 布尔 | ture：成功，假：失败                | true：成功，错误：失败               |
| bk_error_code | INT  | 错误编码.0表示成功，> 0表示失败错误 | 错误代码。0表示成功，> 0表示失败代码 |
| bk_error_msg  | 串   | 请求失败返回的错误信息              | 失败请求的错误消息                   |
| 数据          | 串   | 结果数据                            | 结果                                 |

### 根据关联类型查询使用这些关联类型的关联关系列表

* API：POST / api / v3 / find / topoassociationtype
* API名称：serch_association_list_with_association_kind_list
* 功能说明：
  * 根据关联类型查询使用这些关联类型的关联关系列表
* 输入体
```
{
    "asst_ids": ["run","group"]，
    "metadata":{
        "label":{
            "bk_biz_id":"1"
        }
    }
}
```
* 字段说明

| 名称     | 类型       | 必填 | 默认值 | 说明                           | 描述 |
| -------- | ---------- | ---- | ------ | ------------------------------ | ---- |
| asst_ids | 字符串数组 | 是   | 无     | 要查询的关联类型bk_asst_id列表 |      |

* 产量
```
{
"result": true,
"bk_error_code": 0,
"bk_error_msg": "",
"data": {
  "associations": [
    {
      "bk_asst_id": "run",
      "assts": [
        {
          "id": 8,
          "bk_supplier_account": "0",
          "bk_obj_asst_id": "set_default_nation",
          "bk_obj_asst_name": "test",
          "bk_obj_id": "set",
          "bk_asst_obj_id": "nation",
          "bk_asst_id": "group",
          "mapping": "1:1",
          "on_delete": "none",
          "ispre": false
        }
      ]
    },
    {
      "bk_asst_id": "group",
      "assts": [
        {
          "id": 20,
          "bk_supplier_account": "0",
          "bk_obj_asst_id": "moduel_default_nation",
          "bk_obj_asst_name": "test",
          "bk_obj_id": "moduel",
          "bk_asst_obj_id": "nation",
          "bk_asst_id": "default",
          "mapping": "1:1",
          "on_delete": "none",
          "ispre": false
        }
      ]
    }
  ]
}
}
```

* 输出说明
  * data.association中包含了所有查到的每个关联类型所包含的模型关联关系的信息。
  * bk_asst_id：查询时所用的关联关系id名称。
  * assts：使用关联类型的所有关联关系列表。
