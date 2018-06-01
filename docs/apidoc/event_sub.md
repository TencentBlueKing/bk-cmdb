
### 订阅事件

- API: POST /api/{version}/event/subscribe/{supplier_account}/{bk_biz_id}
- API 名称: subscribe_event
- 功能说明：
	- 中文：事件订阅
	- English：subscribe event

- input body:

``` json
{
  "subscription_name":"mysubscribe",
  "system_name":"SystemName",
  "callback_url":"http://127.0.0.1:8080/callback",
  "confirm_mode":"httpstatus",
  "confirm_pattern":"200",
  "subscription_form":"hostcreate",
  "timeout":10
}
```

- input 字段说明

|字段|类型|是否必须|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id|business id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|subscription_name|string|是|无|订阅的名字|the subscription name|
|system_name|string|是|无|订阅事件的系统的名字|the subscriber's name|
|callback_url|string|是|无|回调函数|the callbacks of the subscribers|
|confirm_mode|string|是|无|事件发送成功校验模式,可选 1-httpstatus,2-regular|confirm success mode of send to callback success, could be 1-httpstatus,2-regular |
|confirm_pattern|string|是|无|callback的httpstatus或正则|the correct return httpstatus or regular|
|subscription_form|string|是|无|订阅的事件,以逗号分隔|subcription event names, should split by comma|
|timeout|int|是|无|发送事件超时时间|time out when send event message to callback|


- output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "subscription_id": 1
    }
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|


data 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|subscription_id|int|新增订阅的订阅ID|the id of the new subscription |



### 退订事件

- API: DELETE /api/{version}/event/subscribe/{supplier_account}/{bk_biz_id}/{subscription_id}
- API 名称：unsubcribe_event
- 功能说明：
	- 中文：退订
	- English：event unsubscribe

- input body

``` json
{}
```

- input 字段说明

|名称|类型|默认值|说明|Description|
|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id| business id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|subscription_id|int|订阅ID|无|subscription_id|


- output
``` json
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```
- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|



### 修改订阅

- API: PUT  /api/{version}/event/subscribe/{supplier_account}/{biz_id}/{subscription_id}
- API 名称：update_event_subscribe
	- 中文：修改订阅
	- English：update the event subscription

- input body

``` json
{
  "subscription_name":"mysubscribe",
  "system_name":"SystemName",
  "callback_url":"http://127.0.0.1:8080/callback",
  "confirm_mode":"httpstatus",
  "confirm_pattern":"200",
  "subscription_form":"hostcreate",
  "timeout":10
}
```

- input 字段说明

|字段|类型|是否必须|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id| business id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|subscription_id|int|是|无|订阅ID|subscription_id|
|subscription_name|string|是|无|订阅的名字|the subscription name|
|system_name|string|是|无|订阅事件的系统的名字|the subscriber's name|
|callback_url|string|是|无|回调函数|the callbacks of the subscribers|
|confirm_mode|string|是|无|事件发送成功校验模式,可选 1-httpstatus,2-regular|confirm success mode of send to callback success, could be 1-httpstatus,2-regular |
|confirm_pattern|string|是|无|callback的httpstatus或正则|the correct return httpstatus or regular|
|subscription_form|string|是|无|订阅的事件,以逗号分隔|subcription event names, should split by comma|
|timeout|int|是|无|发送事件超时时间|time out when send event message to callback|



- output

``` json
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|数据对象，在操作成功后如果有返回值数据会在此字段设置|The result, it will include the data ,only the error code is zero.|


### 查询订阅

- API: POST /api/{version}/event/subscribe/search/{supplier_account}/{biz_id}
- API 名称：search_subscription
	- 中文：查询订阅
	- English：search subscriptions

- input body

``` json
{
    "page":{
        "start":0,
        "limit":10,
        "sort":"HostName"
    }
}
```

- input 输入字段说明
无

- output 

``` json
{
	"result":true,
	"bk_error_code":0,
	"bk_error_msg":"",
	"data":[
		{
			"subscription_id":1,
			"subscription_name":"mysubscribe",
			"system_name":"SystemName",
			"callback_url":"http://127.0.0.1:8080/callback",
			"confirm_mode":"httpstatus",
			"confirm_pattern":"200",
			"subscription_form":"hostcreate",
			"timeout":10,
			"last_time": "2017-09-19 16:57:07",
			"operator": "user",
			"statistics": {
				"total": 30,
				"failure": 2
			}
		}
	]
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data 字段说明

| 名称  | 类型     | 说明   |Description|
| --- | ---|--- |---|
| subscription_id   | int    |订阅ID |the subscription id|
| subscription_name | string |订阅名|the subscription|
| system_name       | string |系统名称|the subscriber's name |
| callback_url      | string |回调地址|the callback of the subscription|
| confirm_mode      | string |回调成功确认模式，可选:httpstatus，regular |the http status|
| confirm_pattern   | string |回调成功标志|the http result pattern|
| subscription_form | string |订阅单，用","分隔|subscribed events,split by comma|
| timeout          | int    |超时时间，单位：秒|time out|
| operator       | int    |本条数据的最后更新人员|updator of this subscription|
| last_time         | int    |更新时间|update time of this subscription|
| statistics.total  | int    |推送总数|the total count one push|
| statistics.failure| int    |推送失败数|the failure total count |

### 测试推送

- API: POST /api/{version}/event/subscribe/ping
- API 名称：ping_subscription
	- 中文：推送测试
	- English：push test

- input body

``` json
{
	"callback_url": "127.0.0.1:8080/callback",
	"data": {}	
}
```

- input 字段说明

|字段|类型|说明|Description|
|---|---|---|---|
|callback_url|string|回调方法|the callback URL|
|data|string|回调方法|data that would send to callback url|

- output

``` json
{
	"result":true,
	"bk_error_code":0,
	"bk_error_msg":"",
	"data":[
		{
			"http_status": 200,
			"response_body": "xxxxx"
		}
	]
}
```

- output 字段说明

| 名称               | 类型     | 说明                             |Description|
| --- | --- | --- |---|
| http_status| int  | 返回的HTTP STATUS |the http status|
| response_body| string | 订阅者的callback返回体| the response data from subcription callback|

### 测试推送（只测试连通性）

- API: POST /api/{version}/event/subscribe/telnet
- API 名称： testing_connection
	- 中文：连通性测试
	- English：connectivity testing

-  input body

``` json
{
	"callback_url": "127.0.0.1:8080/callback"
}
```

- input 字段说明

|字段|类型|说明|Description|
|---|---|---|---|
|callback_url|string|回调方法|the callback URL|

- output 

``` json
{
	"result":true,
	"bk_error_code":0,
	"bk_error_msg":"",
	"data":  "success"
}
```

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|


