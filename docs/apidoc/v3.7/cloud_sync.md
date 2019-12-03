#### 云资源发现

##### 新增发现任务

- API: POST /create/cloud/sync/task

- API名称：create_cloud_sync_task

- 功能说明：

  - 中文：创建云同步任务
  - English：create cloud sync task

- input:

  ```
  {
  	"task_name": "王者荣耀发现任务",
  	"account_id": 4,
  	"resource_type": "host",
  	"sync_to": {
  		"bk_biz_name": true,
  		"bk_biz_id": 1
  	},
  	"need_confirm": false,
  	"sync_vpc": {
  		"all": true,
  		"vpc_ids": [],
  	}
  }
  ```

- input字段说明：

  | 名称          | 类型   | 必填 | 默认值 | 说明                                    | Description                                        |
  | ------------- | ------ | ---- | ------ | --------------------------------------- | -------------------------------------------------- |
  | task_name     | string | 是   | 无     | 云同步任务名称                          | cloud sync task name                               |
  | account_id    | int64  | 是   | 无     | 云账户id（cc_CloudAccount中的唯一标识） | cloud account id                                   |
  | resource_type | string | 是   | 主机   | 同步的资源类型                          | cloud sync resource                                |
  | sync_to       | object | 是   | 无     | 主机录入的资源池目录                    | Directory of the resource pool entered by the host |
  | need_confirm  | bool   | 是   | false  | 资源录入是否需要确认                    | Does resource entry require confirmation           |
  | sync_vpcs     | object | 是   | 无     | 云账户中同步的目标vpc                   | Target vpc should sync in cloud account            |

  sync_vpcs字段说明：

  | 名称    | 类型  | 必填 | 默认值 | 说明                      | Description                                                  |
  | ------- | ----- | ---- | ------ | ------------------------- | ------------------------------------------------------------ |
  | all     | bool  | 是   | false  | 是否同步云账户下全部vpc   | Whether to synchronize all vpc under the cloud account       |
  | vpc_ids | array | 是   | 无     | 云账户下需要同步的vpc的id | The id of the vpc that needs to be synchronized under the cloud account |

- output:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": {
  		"task_id": 1
  	}
  	"permission": null,
  	"result": true
  }
  ```


* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

  data字段说明：

| 名称    | 类型  | 说明                 | Description                      |
| ------- | ----- | -------------------- | -------------------------------- |
| task_id | int64 | 云同步任务的唯一标识 | Unique ID of the cloud sync task |



##### 查询账户VPC

- API: POST /search/cloud/account/vpc

- API名称：search_cloud_account_vpc

- 功能说明：

  - 中文：查询云账户下的vpc信息
  - English：Query vpc information under a cloud account

- input:

  ```
  {
  	"account_id": 5
  }
  ```

- input字段说明：

  | 名称       | 类型  | 必填 | 默认值 | 说明         | Description                    |
  | ---------- | ----- | ---- | ------ | ------------ | ------------------------------ |
  | account_id | int64 | 是   | 无     | 云账户唯一id | unique id of the cloud account |

- output:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": {
  		"count": 2,
  		"info": [
  			{
  				"vpc_name": "vpc-25c909ft（Default-VPC)",
  				"region": "广东一区",
  				"host_count": 25，
  			},
  			{
  				"vpc_name": "vpc-25c909ft（Default-VPC)",
  				"region": "广东二区",
  				"host_count": 15，
  			}
  		]
  	}
  	"permission": null,
  	"result": true
  }
  ```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

  data字段说明：

| 名称  | 类型   | 说明            | Description                       |
| ----- | ------ | --------------- | --------------------------------- |
| count | int64  | 云账户下vpc数量 | Number of vpc under cloud account |
| info  | object | vpc具体的信息   | vpc specific information          |

  info字段说明：

| 名称       | 类型   | 说明          | Description               |
| ---------- | ------ | ------------- | ------------------------- |
| vpc_name   | string | vpc的名称     | name of vpc               |
| region     | string | vpc所属区域   | region of vpc             |
| host_count | int64  | vpc下主机数量 | Number of hosts under vpc |



##### 查询发现任务

- API： POST /search/cloud/sync/task

- API名称：search_cloud_sync_task

- 功能说明：

  - 中文：查询云同步任务
  - English：search cloud sync task

- input:

  ```
  {
  	"condition": {
  		"field": "task_name",
  		"operator": "$regex",
  		"value": "x"
  	}
  	"page": {
  		"start": 0,
  		"limit": 10,
  		"sort": "-create_time"
  	}
  }
  ```

- input字段说明：

  | 名称      | 类型   | 必填 | 默认值 | 说明     | Description      |
  | --------- | ------ | ---- | ------ | -------- | ---------------- |
  | condition | object | 是   | 无     | 查询条件 | Query conditions |
  | page      | object | 否   | 无     | 查询条件 | Query conditions |
  condition字段说明：

  | 名称     | 类型   | 必填 | 默认值 | 说明            | Description       |
  | -------- | ------ | ---- | ------ | --------------- | ----------------- |
  | field    | string | 是   | 无     | 查询的字段      | Query fields      |
  | operator | string | 是   | $regex | 操作符          | Operator          |
  | value    | string | 是   | ""     | 查询字段的value | Query field value |

  page 参数说明：

  | 名称  | 类型 |必填| 默认值 | 说明 | Description|
  | ---  | ---  | --- |---  | --- | ---|
  | start|int|是|无|记录开始位置 |start record|
  | limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
  | sort| string| 否| 无|排序字段|the field for sort|

- ouput:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": {
  		"count": 2,
  		"info": [{
  				"task_id": 1,
  				"task_name": "测试任务一",
  				"resource_type": "主机",
  				"account_id": 1,
  				"status": "成功",
  				"status_description": "同步耗时1s",
  				"last_sync_time": "2019-11-04T17:10:47.819+08:00",
  				"sync_to": {
  					"bk_biz_name": "蓝鲸",
  					"bk_biz_id": 2
  				}
  				"need_confirm": false,
  				"sync_vpc": {
  					"all": true,
  					"vpc_ids": [],
  				}
  				"creator": "admin",
  				"create_time": "2019-11-04T17:10:47.819+08:00",
  				"last_editor": "admin"
      			"last_time": "2019-11-04T17:10:47.819+08:00"
  		}, {
  				"task_id": 2,
  				"task_name": "测试任务二",
  				"resource_type": "主机",
  				"account_id": 1,
  				"status": "成功",
  				"status_description": "同步耗时1s",
  				"last_sync_time": "2019-11-04T17:10:47.819+08:00",
  				"sync_to": {
  					"bk_biz_name": "蓝鲸",
  					"bk_biz_id": 2
  				}
  				"need_confirm": false,
  				"sync_vpc": {
  					"all": false,
  					"vpc_ids": [1, 2, 3],
  				}
  				"creator": "admin",
  				"create_time": "2019-11-04T17:10:47.819+08:00",
  				"last_editor": "admin"
      			"last_time": "2019-11-04T17:10:47.819+08:00"
  		}]
  	}
  	"permission": null,
  	"result": true
  }
  ```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

data字段说明：

| 名称  | 类型   | 说明                 | Description                                          |
| ----- | ------ | -------------------- | ---------------------------------------------------- |
| count | int64  | 云同步任务的数量     | Number of cloud synchronization tasks                |
| info  | object | 云同步任务的具体信息 | Specific information for cloud synchronization tasks |

info字段说明：

| 名称               | 类型   | 说明                         | Description                                                  |
| ------------------ | ------ | ---------------------------- | ------------------------------------------------------------ |
| task_id            | int64  | 云同步任务唯一标识           | Cloud sync task unique identifier                            |
| task_name          | string | 云同步任务名称               | cloud sync task name                                         |
| resource_type      | string | 同步资源类型                 | cloud sync resource type                                     |
| account_id         | int64  | 云同步对应的云账户的唯一标识 | Unique ID of the cloud account corresponding to cloud synchronization |
| status             | string | 同步的状态                   | Status of sync                                               |
| status_description | string | 同步状态说明                 | sync status description                                      |
| last_sync_time     | time   | 上次同步的时间               | Last sync time                                               |
| sync_to            | object | 主机录入的资源池目录         | Directory of the resource pool entered by the host           |
| need_confirm       | bool   | 录入需要确认                 | Entry needs confirmation                                     |
| sync_vpc           | object | 云账户中同步的目标vpc        | Target vpc synced in cloud account                           |
| creator            | string | 创建者                       | creator                                                      |
| create_time        | time   | 创建时间                     | create time                                                  |
| last_editor        | string | 最近编辑人                   | last editor                                                  |
| last_time          | time   | 最近编辑时间                 | last time                                                    |

  sync_vpcs字段说明：

| 名称    | 类型  | 必填 | 默认值 | 说明                      | Description                                                  |
| ------- | ----- | ---- | ------ | ------------------------- | ------------------------------------------------------------ |
| all     | bool  | 是   | false  | 是否同步云账户下全部vpc   | Whether to synchronize all vpc under the cloud account       |
| vpc_ids | array | 是   | 无     | 云账户下需要同步的vpc的id | The id of the vpc that needs to be synchronized under the cloud account |

#####  更新发现任务

- API： POST /update/cloud/sync/task

- API名称：update_cloud_sync_task

- 功能说明：

  - 中文：更新云同步任务
  - English：updata cloud sync task

- input:

  ```
  {
      "task_id": 2,
      "task_name": "测试任务二",
      "resource_type": "主机",
      "account_id": 1,
      "status": "成功",
      "status_description": "同步耗时1s",
      "last_sync_time": "2019-11-04T17:10:47.819+08:00",
      "sync_to": {
          "bk_biz_name": "蓝鲸",
          "bk_biz_id": 2
      }
      "need_confirm": false,
      "sync_vpc": {
  		"all": true,
  		"vpc_ids": [],
      }
      "creator": "admin",
      "create_time": "2019-11-04T17:10:47.819+08:00",
      "last_editor": "admin"
      "last_time": "2019-11-04T17:10:47.819+08:00"
  }
  ```

- input字段说明：


| 名称               | 类型   | 说明                         | Description                                                  |
| ------------------ | ------ | ---------------------------- | ------------------------------------------------------------ |
| task_id            | int64  | 云同步任务唯一标识           | Cloud sync task unique identifier                            |
| task_name          | string | 云同步任务名称               | cloud sync task name                                         |
| resource_type      | string | 同步资源类型                 | cloud sync resource type                                     |
| account_id         | int64  | 云同步对应的云账户的唯一标识 | Unique ID of the cloud account corresponding to cloud synchronization |
| status             | string | 同步的状态                   | Status of sync                                               |
| status_description | string | 同步状态说明                 | sync status description                                      |
| last_sync_time     | time   | 上次同步的时间               | Last sync time                                               |
| sync_to            | object | 主机录入的资源池目录         | Directory of the resource pool entered by the host           |
| need_confirm       | bool   | 录入需要确认                 | Entry needs confirmation                                     |
| sync_vpc           | object | 云账户中同步的目标vpc        | Target vpc synced in cloud account                           |
| creator            | string | 创建者                       | creator                                                      |
| create_time        | time   | 创建时间                     | create time                                                  |
| last_editor        | string | 最近编辑人                   | last editor                                                  |
| last_time          | time   | 最近编辑时间                 | last time                                                    |

  sync_vpcs字段说明：

| 名称    | 类型  | 必填 | 默认值 | 说明                      | Description                                                  |
| ------- | ----- | ---- | ------ | ------------------------- | ------------------------------------------------------------ |
| all     | bool  | 是   | false  | 是否同步云账户下全部vpc   | Whether to synchronize all vpc under the cloud account       |
| vpc_ids | array | 是   | 无     | 云账户下需要同步的vpc的id | The id of the vpc that needs to be synchronized under the cloud account |

- ouput:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": null,
  	"permission": null,
  	"result": true
  }
  ```



##### 删除发现任务

- API： DELETE /delete/cloud/sync/task/{taskID}

- API名称：delete_cloud_sync_task

- 功能说明：

  - 中文：删除云同步任务
  - English：delete cloud sync task

- input:

  ```

  ```

- input字段说明：

  | 名称   | 类型  | 必填 | 默认值 | 说明                 | Description                                 |
  | ------ | ----- | ---- | ------ | -------------------- | ------------------------------------------- |
  | taskID | int64 | 是   | 无     | 云同步任务的唯一标识 | Unique ID of the cloud synchronization task |



- ouput:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": null,
  	"permission": null,
  	"result": true
  }
  ```



##### 查询录入历史

- API： POST /search/cloud/sync/history

- API名称：search_cloud_sync_history

- 功能说明：

  - 中文：查询录入历史
  - English：search cloud sync history

- input:

  ```
  {
  	"task_id"：5,
  	"condition": {
  		"inst_name": "xx"
  		"create_time": {
  			"0": "2019-11-03 00:00:00"
  			"1": "2019-11-06 23:59:59"
  		}
  	},
  	"page": {
  		"start": 0,
  		"limit": 10,
  		"sort": "-create_time"
  	}
  }
  ```

- input字段说明：

  | 名称      | 类型   | 必填 | 默认值 | 说明                 | Description                                 |
  | --------- | ------ | ---- | ------ | -------------------- | ------------------------------------------- |
  | task_id   | int64  | 是   | 无     | 云同步任务的唯一标识 | Unique ID of the cloud synchronization task |
  | condition | object | 是   | 无     | 查询组合条件         | Query combination conditions                |
  | page      | object | 否   | 无     | 查询条件             | Query conditions                            |

  conditon字段说明：

  | 名称        | 类型   | 必填 | 默认值 | 说明           | Description      |
  | ----------- | ------ | ---- | ------ | -------------- | ---------------- |
  | inst_name   | string | 否   | 无     | 实例名         | instance name    |
  | create_time | time   | 否   | 无     | 查询的时间范围 | Query time range |

  page 参数说明：

  | 名称  | 类型 |必填| 默认值 | 说明 | Description|
  | ---  | ---  | --- |---  | --- | ---|
  | start|int|是|无|记录开始位置 |start record|
  | limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
  | sort| string| 否| 无|排序字段|the field for sort|

- output:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": {
  		"count": 1,
  		"info": [{
  			"history_id": 1,
  			"task_id": 2,
  			"operation": "录入"
  			"inst_name": "xxx"
  			"create_time": "2019-11-04T17:10:47.819+08:00"
  			"description": "
                  IP：192.168.1.1
                  云区域：腾讯云广州1
                  MAC地址：DXLKSJDFLSJDFLKSDJFSD
                  云资产状态：正常
                  操作系统类型：linux
                  内存容量：6 GB"
  		}]
  	},
  	"permission": null,
  	"result": true
  }
  ```

- output字段说明：

  | 名称        | 类型   | 说明                     | Description                                 |
  | ----------- | ------ | ------------------------ | ------------------------------------------- |
  | history_id  | int64  | 录入历史的唯一标识       | Unique ID for entry history                 |
  | task_id     | int64  | 云同步任务的唯一标识     | Unique ID of the cloud synchronization task |
  | operation   | string | 实例的操作（录入or更新） | Instance operation (entry or update)        |
  | inst_name   | string | 实例名                   | instance name                               |
  | create_time | time   | 创建时间                 | create time                                 |
  | description | string | 录入详情                 | Entry details                               |


#### 云账户

##### 连接测试

- API： POST /try/connect/cloud/account

- API名称：try_to_connect_cloud_account

- 功能说明：

  - 中文：云账户连接测试
  - English：try to connect cloud account

- input:

  ```
  {
  	"secret_id": "qweqwe",
  	"secret_key": "asasdas",
  	"type": "aws"
  }
  ```

- input字段说明：

  | 名称       | 类型   | 必填 | 默认值 | 说明                              | Description        |
  | ---------- | ------ | ---- | ------ | --------------------------------- | ------------------ |
  | secret_id  | string | 是   | 无     | 云账户id                          | cloud account id   |
  | secret_key | string | 是   | 无     | 云账户key                         | cloud account key  |
  | type       | string | 是   | 无     | 账户类型(可选 aws、tencent_cloud) | cloud account type |



- output:

  ```
  {
  	"bk_error_code": 0
  	"bk_error_msg": "success"
  	"data": {
  		"connected": false,
  		"error_msg": "连接超时"
		}
  	"permission": null
  	"result": true
  }
  ```

- output字段说明：
    | 名称  | 类型  | 说明 |Description|
    |---|---|---|---|
    | result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
    | bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
    | bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
    | data | object| 请求返回的数据 |return data|

    data字段说明：

    | 名称      | 类型   | 说明               | Description                                         |
    | --------- | ------ | ------------------ | --------------------------------------------------- |
    | connected | bool   | 是否成功连接云账户 | Whether the cloud account is successfully connected |
    | error_msg | string | 连接失败的错误信息 | Connection failed error message                     |



##### 新建云账户

- API POST /create/cloud/account

- API名称：create_cloud_account

- 功能说明：

  - 中文：新建云账户
  - English：create cloud account

- input:

  ```
  {
  	"name": "LPL腾讯云",
  	"account_type": "腾讯云",
  	"secret_id": "qweqwe",
  	"secret_key": "asasdas",
  	"description": "接口测试",
  	"creator": "admin"
  }
  ```

- input字段说明：

  | 名称         | 类型   | 必填 | 默认值 | 说明         | Description              |
  | ------------ | ------ | ---- | ------ | ------------ | ------------------------ |
  | name         | string | 是   | 无     | 云账户名称   | cloud account name       |
  | account_type | string | 是   | 无     | 云账户类型   | cloud account type       |
  | secret_id    | string | 是   | 无     | 云账户id     | cloud account id         |
  | secret_key   | string | 是   | 无     | 云账户key    | cloud account key        |
  | description  | string | 是   | 无     | 云账户的备注 | Notes for cloud accounts |
  | creator      | string | 是   | 无     | 创建者       | creator                  |



- output:

  ```
  {
  	"bk_error_code": 0
  	"bk_error_msg": "success"
  	"data": {
  		"account_id": 1
  	}
  	"permission": null
		"result": true
  }
  ```

 - output字段说明：
    | 名称  | 类型  | 说明 |Description|
    |---|---|---|---|
    | result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
    | bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
    | bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
    | data | object| 请求返回的数据 |return data|

    data字段说明：

    | 名称       | 类型  | 说明           | Description             |
    | ---------- | ----- | -------------- | ----------------------- |
    | account_id | int64 | 云账户唯一标识 | Cloud account unique id |



##### 查询云账户

- API POST /search/cloud/account

- API名称：search_cloud_account

- 功能说明：

  - 中文：查询云账户
  - English：search cloud account

- input:

  ```
  {
  	"condition": {
  		"field": "name",
  		"operator": "$regex",
  		"value": "x"
  	}
  	"page": {
  		"start": 0,
  		"limit": 10,
  		"sort": "-create_time"
  	}
  }
  ```

- input字段说明：

  | 名称      | 类型   | 必填 | 默认值 | 说明     | Description      |
  | --------- | ------ | ---- | ------ | -------- | ---------------- |
  | condition | object | 是   | 无     | 查询条件 | Query conditions |
  | page      | object | 否   | 无     | 查询条件 | Query conditions |
  condition字段说明：

  | 名称     | 类型   | 必填 | 默认值 | 说明            | Description       |
  | -------- | ------ | ---- | ------ | --------------- | ----------------- |
  | field    | string | 是   | 无     | 查询的字段      | Query fields      |
  | operator | string | 是   | $regex | 操作符          | Operator          |
  | value    | string | 是   | ""     | 查询字段的value | Query field value |

  page 参数说明：

  | 名称  | 类型 |必填| 默认值 | 说明 | Description|
  | ---  | ---  | --- |---  | --- | ---|
  | start|int|是|无|记录开始位置 |start record|
  | limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
  | sort| string| 否| 无|排序字段|the field for sort|

- output:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": {
  		"count" : 2,
  		"info": [{
  			"name": "LPL腾讯云",
  			"account_type": "腾讯云",
  			"account_id": 1,
  			"secret_id": "12312312"
  			"secret_key": "asasdas",
  			"description": "接口测试"
  			"creator": "admin"，
  			"last_editor": "admin",
  			"create_time": "2019-11-04T17:10:47.819+08:00",
  			"last_time": "2019-11-04T17:10:47.819+08:00"
  		}, {
  			"name": "LPL腾讯云",
  			"account_type": "腾讯云",
  			"secret_id": 2,
  			"secret_key": "asasdas",
  			"description": "接口测试"
  			"creator": "admin"，
  			"last_editor": "admin",
  			"create_time": "2019-11-04T17:10:47.819+08:00",
  			"last_time": "2019-11-04T17:10:47.819+08:00"
  		}]
  	}
  	"permission": null,
  	"result": true,
  }
  ```

- output字段说明：
    | 名称  | 类型  | 说明 |Description|
    |---|---|---|---|
    | result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
    | bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
    | bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
    | data | object| 请求返回的数据 |return data|

    data字段说明：

    | 名称         | 类型   | 说明           | Description              |
    | ------------ | ------ | -------------- | ------------------------ |
    | name         | string | 云账户名称     | cloud account name       |
    | account_type | string | 云账户类型     | cloud account type       |
    | account_id   | string | 云账户唯一标识 | Cloud account unique id  |
    | secret_id    | string | 云账号id       | cloud account id         |
    | secret_key   | string | 云账号key      | cloud account key        |
    | description  | string | 云账户备注     | Notes for cloud accounts |
    | creator      | string | 创建者         | creator                  |
    | create_time  | time   | 创建时间       | create time              |
    | last_editor  | string | 最近编辑人     | last editor              |
    | last_time    | time   | 最近编辑时间   | last edit time           |



##### 更新云账户

- API： POST /update/cloud/account

- API名称：update_cloud_account

- 功能说明：

  - 中文：更新云账户
  - English：update cloud account

- input:

  ```
  {
  	"name": "LPL腾讯云",
  	"account_type": "腾讯云",
  	"account_id": 1,
  	"secret_id": "qweqwe",
  	"secret_key": "asasdas",
  	"description": "",
  	"creator": "admin"
  	"create_time": "2019-10-23 20:12:22"
  }
  ```

- input字段说明：


| 名称         | 类型   | 必填 | 默认值 | 说明         | Description              |
| ------------ | ------ | ---- | ------ | ------------ | ------------------------ |
| name         | string | 是   | 无     | 云账户名称   | cloud account name       |
| account_type | string | 是   | 无     | 云账户类型   | cloud account type       |
| secret_id    | string | 是   | 无     | 云账户id     | cloud account id         |
| secret_key   | string | 是   | 无     | 云账户key    | cloud account key        |
| description  | string | 是   | 无     | 云账户的备注 | Notes for cloud accounts |
| creator      | string | 是   | 无     | 创建者       | creator                  |
| create_time  | time   | 是   | 无     | 创建时间     | create time              |

- output:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": null,
  	"permission": null,
  	"result": true
  }
  ```



##### 删除云账户 （有发现任务在使用时，不可删除）

- API DELETE  /delete/cloud/account/{accoutID}

- API名称：delete_cloud_account

- 功能说明：

  - 中文：删除云账户
  - English：delete cloud account

- input:

  ```

  ```

- output:

  ```
  {
  	"bk_error_code": 0,
  	"bk_error_msg": "success",
  	"data": null,
  	"permission": null,
  	"result": true
  }
  ```

- output字段说明：
    | 名称  | 类型  | 说明 |Description|
    |---|---|---|---|
    | result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
    | bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
    | bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
    | data | object| 请求返回的数据 |return data|
