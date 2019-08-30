### 添加业务

- API: POST /api/{version}/biz/{bk_supplier_account}
- API 名称：create_business
- 功能说明：
	- 中文：新建业务
	- English：create a new business

- input body

``` json
	{
		"bk_biz_name": "cc_app_test",
		"bk_biz_maintainer": "admin",
		"bk_biz_productor": "admin",
		"bk_biz_developer": "admin",
		"bk_biz_tester": "admin",
    }
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_biz_name|string|是|无|业务名|the business name|
|bk_biz_maintainer|string|否|无|运维人员|operation staff|
|bk_biz_developer|string|否|无|开发人员|the developer|
|bk_biz_tester|string|否|无|测试人员|the tester|
**注：此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段。**

### 删除业务

- API: DELETE /api/{version}/biz/{bk_supplier_account}/{bk_biz_id}
- API 名称：delete_business
- 功能说明：
	- 中文：删除业务
	- English：delete a business

- input body

	无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id|the business id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|


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
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|

### 修改业务启用状态

- API: PUT /api/{version}/biz/status/{flag}/{bk_supplier_account}/{bk_biz_id}
- API 名称：update_business_enable_status
- 功能说明：
	- 中文：修改业务启用状态
	- English：update the business enable status

- input body

``` json
无
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id|the business id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|flag|string|是|无|启用状态，为disabled 或者enable|the business name, it should be disabled or enable|




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
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|



### 修改业务

- API: PUT /api/{version}/biz/{bk_supplier_account}/{bk_biz_id}
- API 名称：update_business
- 功能说明：
	- 中文：更新业务信息
	- English：update the business

- input body

``` json
{
    "bk_biz_developer": "",
    "bk_biz_maintainer": "admin,jobdevelop,cmdbdevelop",
    "bk_biz_name": "example_biz",
    "bk_biz_productor": "admin",
    "bk_biz_tester": "",
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id|the business id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_biz_name|string|否|无|业务名称|the business name|
|bk_biz_developer|string|否|无|开发人员|the developer|
|bk_biz_maintainer|string|否|无|运维人员|the maintainers|
|bk_biz_productor|string|否|无|产品人员|the productor|
|bk_biz_tester|string|否|无|测试人员|the tester|

**注：此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段。**

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
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|

### 查询业务

- API: POST /api/{version}/biz/search/{bk_supplier_account}
- API 名称：search_business
- 功能说明：
	- 中文：查询业务
	- English：search the business

- input body

```  json
{
	"page": {
		"start": 0,
		"limit": 10,
		"sort": ""
	},
	"fields": [],
	"condition": {}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明:

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|fields|array|否|无|指定查询的字段|need to show|
|condition|object|否|无|查询条件|search condition|
|page|object|否|无|分页条件|page condition|

page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

**注：sort 字段通过在字段前面增加 -，如 sort:"-field" 可以表示按照字段field降序。**

fields参数说明：

参数为业务的任意属性，如果不填写字段信息，系统会返回业务的所有字段。

**注：需要填写的参数取决于用户自己定义的属性字段。**

condition 参数说明：

condition 参数为业务的任意属性，如果不写代表搜索全部数据。

**注：需要填写的参数取决于用户自己定义的属性字段。**

