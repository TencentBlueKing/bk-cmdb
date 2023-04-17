#### 添加模型唯一约束

- API: POST /api/v3/create/objectunique/object/{bk_obj_id}
- API 名称: create_object_unique_constraints
- 功能说明: 
    - 中文: 添加模型唯一约束
    - English: create an object's unique constraint
- input body

```json
{
    "must_check": true,
    "keys": [
        {
            "key_kind": "property",
            "key_id": 3
        },
        {
            "key_kind": "association",
            "key_id": 4
        }
    ]
}
```

- input 字段说明
  
|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_obj_id|string|是|无|模型英文ID||
|must_check|bool|是|否|是否必须检查，若为否则只有该实例所有的keys字段/关联均非空时才检查||
|keys|object array|是|无|受约束的关键字段/关联||

keys 字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|key_kind|string|是|无|key类型，可选 property, association(未启用), 分别指这个key是模型属性还是模型关联||
|key_id|int64|是|无|模型属性或模型关联的自增ID||

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
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
|id|int|模型唯一约束的自增ID|auto-increment id|

#### 编辑模型唯一约束

- API: PUT /api/v3/update/objectunique/object/{bk_obj_id}/unique/{id}
- API 名称: update_object_unique_constraints
- 功能说明: 
    - 中文: 编辑模型唯一约束
    - English: update an object's unique constraint
  
    若变更keys，我们会重查是否已有的数据是否符合约束，若不符合，会报错

- input body

```json
{
    "must_check": true,
    "keys": [
        {
            "key_kind": "property",
            "key_id": 3,
        },
        {
            "key_kind": "association",
            "key_id": 4,
        }
    ]
    
}
```

- input 字段说明
  
|字段名|类型|必填|默认值|说明|Description|
|-|-|-|-|-|-|
|id|int|是|无|模型唯一约束的自增ID|auto-increment id|
|must_check|bool|否|无|是否必须检查，若为否则只有该实例所有的keys字段/关联均非空时才检查|
|keys|object array|否|无|受约束的关键字段/关联，必须传入全量的keys，更新时会覆盖原来的keys||

keys 字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|key_kind|string|是|无|key类型，可选 property, association(未启用), 分别指这个key是模型属性还是模型关联||
|key_id|int64|是|无|模型属性或模型关联的自增ID||

- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
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

#### 删除模型唯一约束

- API: POST /api/v3/delete/objectunique/object/{bk_obj_id}/unique/{id}
- API 名称: delete_object_unique_constraints
- 功能说明: 
    - 中文: 删除模型唯一约束
    - English: delete an object's unique constraint
- input body

```json
{
    
}
```

- input 字段说明


- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
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

#### 查询模型唯一约束

- API: POST /api/v3/find/objectunique/object/{bk_obj_id}
- API 名称: search_object_unique_constraints
- 功能说明: 
    - 中文: 查询模型唯一约束
    - English: search an object's unique constraint
- input body

```json
{
    
}
```

- input 字段说明



- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": [{
            "id": 1,
            "must_check": true,
            "keys": [
                {
                    "key_kind": "property",
                    "key_id": 3,
                },
                {
                    "key_kind": "association",
                    "key_id": 4,
                }
            ],
            "last_time":"2018-03-08 15:10:42",
            "ispre": false,
            "bk_supplier_account": "0"
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
|data|object array|结果数据|the result|

data 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int64|自增ID||
|bk_obj_id|string|模型英文ID||
|must_check|bool|是否必须检查，若为否则只有该实例所有的keys字段/关联均非空时才检查|
|keys|object array|受约束的关键字段/关联||
|last_time|time|最后更新时间|last update time|
|ispre|bool|是否内置，若为内置则不允许修改删除|
|bk_supplier_account|string|开发商账号|supplier account code|

keys 字段说明

|字段|类型|说明|Description|
|---|---|---|---|
|key_kind|string|key类型，可选 property, association(未启用), 分别指这个key是模型属性还是模型关联||
|key_id|int64|模型属性或模型关联的自增ID||

