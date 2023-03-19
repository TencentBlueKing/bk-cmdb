### 功能描述

监听系统资源变化产生的事件(v3.8以上)


**该watch功能的主要特性包括：**


* 在有限的时间内（目前为3小时,可能会调整，请勿依赖此时间）为用户提供高可用的数据变更watch服务。


* 在有限时间内，用户可以根据自己上一次事件的cursor(游标)进行事件回溯或者追数据，适用于异常数据回溯，或者系统变更进行数据补录。


* 支持根据时间点进行变更数据回溯，支持根据游标进行变更数据回溯，支持从当前时间点进行数据变更watch。


* 支持根据事件类型进行watch的能力，包括增、删、改。事件中包含全量的数据。


* 支持主机与主机关系数据变化的事件watch能力。


* 采用短长链的设计，当用户通过游标进行事件watch时，如果没有事件，则会保持会话连接，在20s内有事件变更则直接直接将事件推回。避免用户不断请求，同时保证用户能及时的拿到变更的数据。


* 支持批量事件watch能力，提升系统吞吐能力。


* 支持定制关注的事件数据字段，满足用户轻量级的watch需求。


### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                | 类型           | 必选   | 描述                                                         |
| ------------------- | -------------- | ------ | ------------------------------------------------------------ |
| bk_event_types      | array   string | 否     | 事件类型，如果填了，即为只关注此类型的事件。可选的值为：create(新增)/update(更新)/delete(删除)。如，若使用create，则只关注该资源的新增事件。不填为空，则为关注所有。 |
| bk_fields           | array string   | 看情况 | 返回的事件中需要返回的字段列表，目前监听主机资源该字段为必填字段，不能置空，主机关系可以置空。置空则默认为返回所有字段。 |
| bk_start_from       | Int64          | 否     | 监听事件的起始时间，该值为unix time的秒数，即为从UTC1970年1月1日0时0分0秒起至你要watch的时间点的总秒数。 |
| bk_cursor           | string         | 否     | 监听事件的游标，代表了要开始或者继续watch(监听)的事件地址，系统会返回这个游标的下一个、或一批事件。 |
| bk_resource         | string         | 是     | 要监听的资源类型，枚举值为：host, host_relation, biz, set, module, process, object_instance, mainline_instance, biz_set, biz_set_relation, plat, project。其中host代表主机详情事件，host_relation代表主机的关系事件，biz代表业务详情事件，set代表集群详情事件，module代表模块详情事件，process代表进程详情事件，object_instance代表通用模型实例事件，mainline_instance代表主线模型实例事件，biz_set代表业务集事件，biz_set_relation代表业务集和业务的关系事件, plat代表云区域事件, project代表项目事件。 |
| bk_supplier_account | string         | 是     | 开发商账号                                                   |
| bk_filter           | object         | 否     | 过滤条件                                                     |

**注: biz_set_relation事件会在业务集的新增、删除和更新"bk_scope"字段时和业务的新增、删除、更新涉及到业务集关系变更时触发。所有业务集关系事件的事件类型(bk_event_type)均为update类型，事件详情中会返回关系发生了变更的业务集的ID和该业务集所包含的所有业务ID列表。当事件是由业务集删除事件触发时，返回的事件详情中的业务ID列表为空**

#### bk_filter

| 字段                 | 类型           | 必选   | 描述                                                         |
| ------------------- | -------------- | ------ | ------------------------------------------------------------ |
| bk_sub_resource     | string         | 否     | 要监听的下级资源类型，仅支持bk_resource为object_instance或mainline_instance时使用，代表需要监听的模型的bk_obj_id |


### 请求参数示例

主机：

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_event_types": ["create","update","delete"],
    "bk_fields": ["bk_host_innerip", "bk_mac"],
    "bk_start_from": 12345678999,
    "bk_cursor": "MQ0yDTE1ODkyMDcyODENMQ01ZWI3ZWZjNTBiOTA5ZTYyMGFmYWQzZGY=",
    "bk_resource": "host"
}

```

通用模型实例：

```json
{
    "bk_event_types": [],
    "bk_fields": ["bk_inst_id", "bk_inst_name"],
    "bk_start_from": 12345678999,
    "bk_cursor": "MQ0yDTE1ODkyMDcyODENMQ01ZWI3ZWZjNTBiOTA5ZTYyMGFmYWQzZGY=",
    "bk_resource": "object_instance",
    "bk_filter": {
        "bk_sub_resource": "xxx"
    },
    "bk_supplier_account": "0"
}

```

## 返回参数

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "bk_watched": true,
        "bk_events": [
            {
                "bk_cursor": "MQ0yDTE1ODkyMDcyODENMQ01ZWI3ZWZjNTBiOTA5ZTYyMGFmYWQzZGY=",
                "bk_resource": "host",
                "bk_event_type": "update",
                "bk_detail": {
                    "bk_cpu": 2
                }
            },
            {
                "bk_cursor": "MQ0yDTE1ODkzNDExMDcNMQ01ZWI3ZWZjNTBiOTA5ZTYyMGFmYWQzZGY=",
                "bk_resource": "host",
                "bk_event_type": "update",
                "bk_detail": {
                    "bk_cpu": 2
                }
            }
        ]
    }
}

```

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | Array | 事件数据详情，是一个有序的数组，数组尾部的事件为新的事件。 |

- data 数据描述

| 名称       | 类型             | 描述                                                         |
| ---------- | ---------------- | ------------------------------------------------------------ |
| bk_watched | bool             | 是否监听到了事件，true：监听到了事件；false:未监听到事件     |
| bk_events  | 监听到的事件详情 | 监听到的事件详情列表，最大长度为200，后续可能会调，请勿依赖此长度。 |

- bk_events 数据描述

| 名称          | 类型        | 描述                                                         |
| ------------- | ----------- | ------------------------------------------------------------ |
| bk_cursor     | string      | 代表当前资源事件的游标值，调用方可以用该游标获取该事件后的下一个事件 |
| bk_resource   | enum string | 该事件对应的资源类型                                         |
| bk_event_type | enum string | 该事件对应的事件类型，枚举值为：create(新增)/update(更新)/delete(删除)。 |
| bk_detail     | object      | 该事件的对应的资源的详情数据，不同的资源，对应的详情不同。   |

#### host_relation资源 bk_detail字段数据示例：
```json

{
	"bk_biz_id" : 1,
	"bk_host_id" : 2,
	"bk_module_id" : 3,
	"bk_set_id" : 4,
	"bk_supplier_account" : "0"
}
```

#### host资源 bk_detail字段数据示例：
```json
{
	"bk_host_name" : "hostname",
	"bk_mem" : null,
	"bk_cloud_id" : 0,
	"operator" : "user",
	"bk_cpu" : null,
	"bk_mac" : "",
	"bk_host_innerip" : "192.168.1.1",	
        "bk_supplier_account" : "0",
	....
}
```

#### biz_set_relation资源 bk_detail字段数据示例：
```json
{
	"bk_biz_set_id": 1,
	"bk_biz_ids": [1 ,2, 3]
}
```

- biz_set_relation资源 bk_detail数据描述

| 名称          | 类型      | 描述                                   |
| ------------- | --------- | -------------------------------------- |
| bk_biz_set_id | int       | 业务集和业务的关系发生了变化的业务集ID |
| bk_biz_ids    | int array | 该业务集所包含的所有业务的ID列表      |


### 使用说明

该接口的使用的流程：

1. 确定起始的监听方式，分别可为：

- 1.1 指定从某一时间开始监听事件

- 1.2 指定当前时间开始监听事件

- 1.3 指定用游标bk_cursor位置开始监听事件

  确定后，发起请求。

2. cmdb返回符合你预期的事件列表：
   - 2.1: 如果bk_watched为ture, 表明已经监听到了事件，bk_events中为事件详情列表，调用方可拿数组最后一个事件的cursor，用1.3步骤中的方式继续进行下一次的watch。
   - 2.2: 如果bk_watched为false，表明未监听到事件，bk_events中只有一个事件，拿该事件的cursor进行下一次的watch。

**注意**：

该事件是有过期时间的，目前为3小时，如果过期，事件将会释放，且过期的游标（cursor）也会失效。可以用以下两种方式重新watch：

1. 指定从某一时间开始监听事件
2. 指定当前时间开始监听事件



