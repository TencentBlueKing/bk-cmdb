
### 添加模型主关联
- API POST /api/{version}/topo/model/mainline
- API 名称：create_mainline_object
- 功能说明：
	- 中文：添加主线模型
	- English：create the main line model

- input body

``` json
{
	"bk_classification_id": "XXX",
	"bk_obj_id": "cc_test",
	"bk_obj_name": "cc_test",
	"bk_supplier_account": "0",
	"bk_asst_obj_id": "id-XXX",
	"bk_obj_icon": "icon-XXX"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

- 输入参数

| 字段                 | 类型   | 必填 | 默认值 | 说明                                                     | Description                       |
| -------------------- | ------ | ---- | ------ | -------------------------------------------------------- | --------------------------------- |
| bk_classification_id | string | 是   | 无     | 对象模型的分类ID，只能用英文字母序列命名                 | the classification identifier     |
| bk_obj_id            | string | 是   | 无     | 对象模型的ID，只能用英文字母序列命名                     | the object identifier             |
| bk_obj_name          | string | 是   | 无     | 对象模型的名字，用于展示，可以使用人类可以阅读的任何语言 | the object name                   |
| bk_supplier_account  | string | 是   | 无     | 开发商账号                                               | supplier account code             |
| bk_asst_obj_id       | string | 是   | 无     | 主线模型关联的父对象模型的ID（bk_obj_id）                | the association object identifier |
| bk_obj_icon          | string | 是   | 无     | 模型的图标                                               | the icon of the object            |

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

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | string | 请求返回的数据                             | the data response                                          |

### 删除模型主关联

- API: DELETE  /api/{version}/topo/model/mainline/owners/{bk_supplier_account}/objectids/{bk_obj_id}
- API 名称：delete_mainline_object
- 功能说明：
	- 中文：删除主线模型
	- English：delete the mainline object

- input body

    无


- input 字段说明

| 字段                | 类型   | 必填 | 默认值 | 说明         | Description           |
| ------------------- | ------ | ---- | ------ | ------------ | --------------------- |
| bk_supplier_account | string | 是   | 无     | 开发商账号   | supplier account code |
| bk_obj_id           | string | 是   | 无     | 对象模型的ID | the object identifier |


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

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | string | 请求返回的数据                             | the data response                                          |

### 查询模型拓扑

- API: GET/api/{version}/topo/model/{bk_supplier_account}  
- API 名称：search_mainline_object
- 功能说明：
	- 中文：搜索主线模型
	- English：search the main line model

- input body

    无

-  input字段说明

| 字段                | 类型   | 必填 | 默认值 | 说明       | Description           |
| ------------------- | ------ | ---- | ------ | ---------- | --------------------- |
| bk_supplier_account | string | 是   | 无     | 开发商账号 | supplier account code |


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": [{
		"bk_next_name": "",
		"bk_next_obj": "",
		"bk_obj_id": "biz",
		"bk_obj_name": "业务",
		"bk_pre_obj_id": "",
		"bk_pre_obj_name": "",
		"bk_supplier_account": "0"
	}]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | array  | 请求返回的数据                             | the data response                                          |

data 字段说明：

| 名称                | 类型   | 说明             | Description                   |
| ------------------- | ------ | ---------------- | ----------------------------- |
| bk_next_name        | string | 下一个模型的名字 | the next object name          |
| bk_next_obj         | string | 下一个模型的ID   | the next object identifier    |
| bk_obj_id           | string | 当前的模型ID     | the current object identifier |
| bk_obj_name         | string | 当前模型的名字   | the current object name       |
| bk_pre_obj_id       | string | 上一个模型的ID   | the pre object identifier     |
| bk_pre_obj_name     | string | 上一个模型的名字 | the pre object name           |
| bk_supplier_account | string | 开发商账号       | supplier account code         |



### 获取实例拓扑

- API: GET /api/{version}/topo/inst/{bk_supplier_account}/{bk_biz_id}
- API 名称：get_inst_topo
- 功能说明：
	- 中文：获取实例拓扑
	- English：get the  topo of the inst
	
- input body

    无


- input 输入参数

| 字段                | 类型   | 必填 | 默认值 | 说明       | Description           |
| ------------------- | ------ | ---- | ------ | ---------- | --------------------- |
| bk_biz_id           | int    | 是   | 无     | 业务id     | the business id       |
| bk_supplier_account | string | 是   | 无     | 开发商账号 | supplier account code |


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": [{
		"default": 0,
		"bk_inst_id": 96,
		"bk_inst_name": "cc_biz_test",
		"bk_obj_id": "biz",
		"bk_obj_name": "业务",
		"child": [{
			"default": 0,
			"bk_inst_id": 58,
			"bk_inst_name": "obj_id_name",
			"bk_obj_id": "obj_id",
			"bk_obj_name": "obj_id_name",
			"child": [{
				"default": 0,
				"bk_inst_id": 59,
				"bk_inst_name": "obj_inst_name",
				"bk_obj_id": "obj_inst",
				"bk_obj_name": "obj_inst",
				"child": []
			}]
		}]
	}]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称          | 类型   | 说明                                    | Description                                                |
| ------------- | ------ | --------------------------------------- | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:成功；false:失败     | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误 | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                  | error message from failed request                          |
| data          | array  | 请求返回的数据                          | the data response                                          |

data 字段说明：

| 名称         | 类型   | 说明       | Description           |
| ------------ | ------ | ---------- | --------------------- |
| bk_inst_id   | int    | 实例ID     | the inst identifier   |
| bk_inst_name | string | 实例名字   | the inst name         |
| bk_obj_id    | string | 模型的标识 | the object identifier |
| bk_obj_name  | string | 模型名     | the object name       |
| child        | array  | 实例集合   | the inst array        |

**注:child节点下包含的字段于data节点包含的字段一致。**

###  获取子节点实例

- API: GET /api/{version}/topo/inst/child/{bk_supplier_account}/{bk_obj_id}/{bk_biz_id}/{bk_inst_id}
- API名称：search_inst_topo
- 功能说明：
	- 中文：获取子节点实例拓扑
	- English：search inst topo

- input body

    无

- input 输入参数

| 字段                | 类型   | 必填 | 默认值 | Description  |
| ------------------- | ------ | ---- | ------ | ------------ |
| bk_biz_id           | int    | 是   | 无     | 业务id       | the business id       |
| bk_supplier_account | string | 是   | 无     | 开发商账号   | supplier account code |
| bk_obj_id           | string | 是   | 无     | 对象模型的ID | the object identifier |
| bk_inst_id          | string | 是   | 无     | 实例ID       | the inst id           |

- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": [{
		"default": 0,
		"bk_inst_id": 96,
		"bk_inst_name": "cc_biz_test",
		"bk_obj_id": "biz",
		"bk_obj_name": "业务",
		"child": [{
			"default": 0,
			"bk_inst_id": 58,
			"bk_inst_name": "obj_id_name",
			"bk_obj_id": "obj_id",
			"bk_obj_name": "obj_id_name",
			"child": [{
				"default": 0,
				"bk_inst_id": 59,
				"bk_inst_name": "obj_inst_name",
				"bk_obj_id": "obj_inst",
				"bk_obj_name": "obj_inst",
				"child": []
			}]
		}]
	}]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | array  | 请求返回的数据                             | the data response                                          |

data 字段说明：

| 名称         | 类型   | 说明                                       | Description                                                   |
| ------------ | ------ | ------------------------------------------ | ------------------------------------------------------------- |
| default      | int    | 1-资源模块（空闲机），2-故障模块（故障机） | 1-Resource Module(Idle Machine),2-Fault Module(Fault Machine) |
| bk_inst_id   | int    | 实例ID                                     | the inst identifier                                           |
| bk_inst_name | string | 实例名字                                   | the inst name                                                 |
| bk_obj_id    | string | 模型的标识                                 | the object identifier                                         |
| bk_obj_name  | string | 模型名                                     | the object name                                               |
| child        | array  | 实例集合                                   | the inst array                                                |

**注:child节点下包含的字段于data节点包含的字段一致。**

###  查询内置模块集
- API: GET /api/{version}/topo/internal/{bk_supplier_account}/{bk_biz_id}
- API名称： get_internal_topo
- 功能说明：
	- 中文：获取业务的空闲机和故障机模块
	- English：get the internal idle-cluster and the fault-cluster


- input body

    无


- input 字段说明

| 字段                | 类型   | 必填 | 默认值 | 说明Description |
| ------------------- | ------ | ---- | ------ | --------------- |
| bk_supplier_account | string | 是   | 无     | 开发商账号      | supplier account code |
| bk_biz_id           | int    | 是   | 无     | 业务ID          | the business id       |


- output
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
        “module":[
            {
                “bk_module_id":503,
                “bk_module_name":"空闲机"
            },
            {
                “bk_module_id":504,
                “bk_module_name":"故障机"
            }
        ],
        “bk_set_id":214,
        “bk_set_name":"内置模块集"
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | the data response                                          |

data 字段说明:

| 名称        | 类型   | 说明     | Description  |
| ----------- | ------ | -------- | ------------ |
| bk_set_id   | int    | 集群ID   | the set id   |
| bk_set_name | string | 集群名字 | the set name |

module 字段说明:

| 名称           | 类型   | 说明       | Description               |
| -------------- | ------ | ---------- | ------------------------- |
| bk_module_id   | int    | 模块记录ID | the module data record id |
| bk_module_name | string | 模块名     | the module name           |
