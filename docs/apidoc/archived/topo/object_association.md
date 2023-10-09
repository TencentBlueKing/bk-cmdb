### 关联类型
#### 查询关联类型

- API: POST /api/v3/find/associationtype
- API 名称: search_association_type
- 功能说明: 
    - 中文: 查询关联类型
    - English: search a association type
- input body

```json
{
    "page":{
        "start":0,
        "limit":100,
        "sort":"bk_asst_id"
    },
    "condition":{
        "bk_asst_name": {
            "$eq": "belong"
        }
    }
}
```

- input 字段说明

page 参数说明: 

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|
  
|字段名|类型|必填|默认值|说明|Description|
|---|---| --- |---  | --- | ---|
|bk_asst_id|string|是|无|唯一标识|Unique Identification|
|bk_asst_name|string|是|无|显示的名称|name of the association|
|src_des|string|是|无|"源->目标"描述|description from source to destination|
|dest_des|string|是|无|"目标->源"描述|description from destination to source|
|direction|enum|否|无|连线的方向，默认none。枚举值: [none, src_to_dest, dest_to_src, bidirectional(双向)]|direction for theline|

- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "count": 1,
        "info":[{
            "id": 1,
            "bk_asst_id": "belong",
            "bk_asst_name": "属于",
            "src_des": "属于",
            "dest_des": "被属于",
            "direction": "none",
        }]
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
|bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
|bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|结果数据|the result|

data 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|count|int|数据条数|the data item count|
|info|object array|数据集合|the data array|

info 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|自增ID|auto-increment id|
|bk_asst_id|string|唯一标识|Unique Identification|
|bk_asst_name|string|显示的名称|name of the association|
|src_des|string|"源->目标"描述|description from source to destination|
|dest_des|string|"目标->源"描述|description from destination to source|
|direction|enum|连线的方向，默认none。枚举值: [none, src_to_dest, dest_to_src, bidirectional(双向)]|direction for theline|
|ispre|bool|是否是预置，预置的关联类型不允许删除|is preset,if this association type is preset, then it is forbidden to delete|

#### 添加关联类型

- API: POST /api/v3/create/associationtype
- API 名称: create_association_type
- 功能说明: 
    - 中文: 创建关联类型
    - English: create a association type
- input body

```json
{
    "bk_asst_id": "belong",
    "bk_asst_name": "属于",
    "src_des": "属于",
    "dest_des": "被属于",
    "direction": "none",
}
```

- input 字段说明
  
|字段名|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|bk_asst_id|string|是|无|唯一标识|Unique Identification|
|bk_asst_name|string|是|无|显示的名称|name of the association|
|src_des|string|是|无|"源->目标"描述|description from source to destination|
|dest_des|string|是|无|"目标->源"描述|description from destination to source|
|direction|enum|否|无|连线的方向，默认none。枚举值: [none, src_to_dest, dest_to_src, bidirectional(双向)]|direction for theline|

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "id": 1038
    }
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|自增ID|auto-increment id|

#### 编辑关联类型

- API: PUT /api/v3/update/associationtype/{id}
- API 名称: update_association_type
- 功能说明: 
    - 中文: 编辑关联类型
    - English: update a association type
- input body

```json
{
    "bk_asst_name": "属于",
    "src_des": "属于",
    "dest_des": "被属于",
    "direction": "none",
}
```

- input 字段说明
  
|字段名|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|id|int|是|无|自增ID|auto-increment id|
|bk_asst_name|string|是|无|显示的名称|name of the association|
|src_des|string|是|无|"源->目标"描述|description from source to destination|
|dest_des|string|是|无|"目标->源"描述|description from destination to source|
|direction|enum|否|无|连线的方向，默认none。枚举值: [none, src_to_dest, dest_to_src, bidirectional(双向)]|direction for theline|

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|结果数据|the result|

#### 删除关联类型

- API: DELETE /api/v3/delete/associationtype/{id}
- API 名称: delete_association_type
- 功能说明: 
    - 中文: 删除关联类型
    - English: delete a association type
- input body

无

- input 字段说明
  
无

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|结果数据|the result|

### 模型关联

#### 查询模型关联
- API: POST /api/v3/find/objectassociation
- API 名称: search_object_association
- 功能说明: 
    - 中文: 查询模型关联
    - English: search a association between object
- input body

```json
{
    "condition": {
        "bk_asst_id": "belong",
        "bk_obj_id": "bk_switch",
        "bk_asst_obj_id": "bk_host"
    }
}
```

- input 字段说明
  
|字段名|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|bk_asst_id|string|否|无|关联类型的唯一标识|
|bk_obj_id|string|否|无|源模型ID|
|bk_asst_obj_id|string|否|无|目标模型ID|


- output

``` json
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

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
|bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
|bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|结果数据|the result|

data 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|自增ID|auto-increment id|
|bk_obj_asst_id|string|唯一标识，自动生成。规则: 源模型英文ID+关联类型英文标识+目标模型英文ID。由前端生成传入，后端只做唯一校验|
|bk_obj_asst_name|string|别名|
|bk_asst_id|string|关联类型|
|bk_asst_name|string|显示的名称|
|bk_obj_id|string|源模型ID|
|bk_obj_name|string|源模型ID|
|bk_asst_obj_id|string|目标模型名称|
|bk_asst_obj_name|string|源模型名称|
|mapping|enum|关联映射，可选: [1:1, 1:n, n:n]|
|on_delete|enum|删除时对实例的动作，可选none, delete_src, delete_dest，分别表示不处理，删除源实例，删除目标实例|




#### 添加模型关联
- API: POST /api/v3/create/objectassociation
- API 名称: create_object_association
- 功能说明: 
    - 中文: 添加模型关联
    - English: create a association between object
    - English: create a association between object
- input body

```json
{
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_obj_asst_name": "",
    "bk_asst_id": "belong",
    "bk_obj_id": "bk_switch",
    "bk_asst_obj_id": "bk_host",
    "mapping": "1:n",
    "on_delete": "none"
}
```

- input 字段说明
  
|字段名|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|bk_obj_asst_id|string|是|无|唯一标识，自动生成。规则: 源模型英文ID+关联类型英文标识+目标模型英文ID。由前端生成传入，后端只做唯一校验|
|bk_obj_asst_name|string|否|无|别名|
|bk_asst_id|string|是|无|关联类型|
|bk_obj_id|string|是|无|源模型ID|
|bk_asst_obj_id|string|是|无|目标模型ID|
|mapping|enum|是|无|关联映射，可选: [1:1, 1:n, n:n]|
|on_delete|enum|否|none|删除时的动作，可选none, delete_src, delete_dest|

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "id": 1038
    }
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|自增ID|auto-increment id|

#### 编辑模型关联

- API: PUT /api/v3/update/objectassociation/{id}
- API 名称: update_object_association
- 功能说明: 
    - 中文: 编辑关联类型，只有input body里的任一一个可以更新。
    - English: update a association between object
- input body

```json
{
    "bk_asst_name": "属于",
    "bk_asst_id":"belong",
    "on_delete":""// 具体枚举值见上。
}
```

- input 字段说明
  
|字段名|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|id|int|是|无|自增ID|auto-increment id|
|bk_asst_name|string|否|无|显示的名称|name of the association|
|bk_asst_id|string|否|无|关联类型||
|on_delete|string|否|无|删除时的动作||

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|结果数据|the result|

#### 删除模型关联

- API: DELETE /api/v3/delete/objectassociation/{id}
- API 名称: delete_object_association
- 功能说明: 
    - 中文: 删除关联类型
    - English: delete association between object
- input body

无

- input 字段说明
  
无

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|结果数据|the result|

# 根据关联类型查询使用这些关联类型的关联关系列表
* API: POST /api/v3/find/topoassociationtype
* API 名称：serch_association_list_with_association_kind_list
* 功能说明：
  - 根据关联类型查询使用这些关联类型的关联关系列表

* input body

```json
{
	"asst_ids": ["run","group"]
}

```
* 字段说明
 
|名称|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|asst_ids|string数组|是|无|要查询的关联类型bk_asst_id列表||

* output
```json
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

* output说明
	- data.association中包含了所有查到的每个关联类型所包含的模型关联关系的信息。
	- bk_asst_id: 查询时所用的关联关系id名称。
	- assts: 使用关联类型的所有关联关系列表。



### 实例关联

#### 查询实例关联
- API: POST /api/v3/find/instassociation
- API 名称: search_inst_association
- 功能说明: 
    - 中文: 查询实例之间的关联信息
    - English: search a association between inst
- input body

```json
{
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_object_id": "",
        "bk_asst_obj_id": ""
    }
}
```

- input 字段说明

|字段名|类型|必填|说明|
|-|-|-|-|
|bk_obj_asst_id|string|否|关联唯一标识|
|bk_asst_id|string|否|关联类型|
|bk_object_id|string|否|源模型ID|
|bk_asst_obj_id|string|否|目标实例ID|

- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": [{
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
|bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
|bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|结果数据|the result|

data 说明（结构待定）

| 字段|类型|说明|Description|
|---|---|---|---|
|bk_obj_asst_id|string|模型关联唯一标识|Object association unique identification|
|bk_obj_id|string|源模型ID，冗余字段|source object id|
|bk_asst_obj_id|string|目标模型ID| target object id |
|bk_inst_id|int|源实例ID| source inst id |
|bk_asst_inst_id|int|目标实例ID|target inst id|
|bk_supplier_account|string|开发商账号|supplier account code|


#### 添加实例关联
- API: POST /api/v3/create/instassociation
- API 名称: create_inst_association
- 功能说明: 
    - 中文: 添加实例关联
    - English: create a association between inst
- input body

```json
{
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 1,
    "bk_asst_inst_id": 2
}
```

- input 字段说明
  
|字段名|类型|必填|说明|
|-|-|-|-|
|bk_obj_asst_id|string|是|唯一标识|
|bk_inst_id|int|是|源实例ID|
|bk_asst_inst_id|int|是|目标实例ID|

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "id": 1038
    }
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|自增ID|auto-increment id|

#### 删除实例关联

- API: DELETE /api/v3/delete/instassociation/{bk_obj_id}/{id}
- API 名称: delete_inst_association
- 功能说明: 
    - 中文: 删除实例关联
    - English: delete association between inst
- association_id: 实联关联关系的自增id值。

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture: 成功，false: 失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|
