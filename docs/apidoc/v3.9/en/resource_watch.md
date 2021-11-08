### Functional description

watch and get the change of a kind of resource's event.


### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description                                                    |
| ------------------- | -------------- | ------ | ------------------------------------------------------------ |
| bk_event_types      | array   string | No     | resource's event kind, enum: create, update, delete 。 if empty, it's means all kind of a resource. |
| bk_fields           | array string   | Depends on resource | the resource event fields you want to return, the host resource need to be set.  |
| bk_start_from       | Int64          | No     | watch from a unix seconds time, It is the number of seconds that have elapsed since the Unix epoch, minus leap seconds;  |
| bk_cursor           | string         | No     | a cursor to represent the event you are watched at, you can use a cursor to watch the next event after it. |
| bk_resource         | string         | Yes     | resource you can watch, now is host, host_relation, biz, set, module, process. "host" means a host's detail info, "host_relation" means host's relation with biz, set and module, "biz" means a biz's detail info, "set" means a set's detail info, "module" means a module's detail info, "process" means a process's detail info |
| bk_supplier_account | string         | Yes    | supplier account                                                  |


### Request Parameters Example

```json
{
    "bk_event_types": ["create","update","delete"],
    "bk_fields": ["bk_host_innerip", "bk_mac"],
    "bk_start_from": 12345678999,
    "bk_cursor": "MQ0yDTE1ODkyMDcyODENMQ01ZWI3ZWZjNTBiOTA5ZTYyMGFmYWQzZGY=",
    "bk_resource": "host",
    "bk_supplier_account": "0"
}

```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
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

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | the responses event details |

- data description

| Field                   | Type     | Description                                                                                          |
| ---------- | ---------------- | ------------------------------------------------------------ |
| bk_watched | bool             | have watched event or not. true：watched event；false: not watched event.|
| bk_events  | watched events details | the events list details being watched |

- bk_events description

| Field                   | Type     | Description                                                                                          |
| ------------- | ----------- | ------------------------------------------------------------ |
| bk_cursor     | string      | represent the current event's cursor location, user can use this to get the next event after it. |
| bk_resource   | enum string | the resource type being watched.                                         |
| bk_event_type | enum string | the event type, can be create/update/delete, which means create/update/delete a resource。 |
| bk_detail     | object      | event's details info, different resources event have different details。   |

#### host_relation resource details example：
```json

{
	"bk_biz_id" : 1,
	"bk_host_id" : 2,
	"bk_module_id" : 3,
	"bk_set_id" : 4,
	"bk_supplier_account" : "0"
}
```

#### host resource details example：
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



### Usage description

how to use this api：

1. find a way to tell us where you want to start to watch a resource：

- 1.1 from a time point, then set bk_start_from field, a unix time seconds value.

- 1.2 to watch from now on.

- 1.3 watch with a cursor to get event next it.

  then you fire your request.

2. the api will return your event according to your request, it will be：
   - 2.1: if bk_watched is true, it means you have event(s)，event details is in bk_events
   and then you can use the last cursor to get events next to it with item 1.3.
   - 2.2: if bk_watched is false, it means you have not got events. it will return only one
    event detail, and you use this cursor to get next events.

**Note**：

the event's expired time is 3h at now. if event is expire, then it will be removed from system
, and the cursor target to it will turn to illegal. you can watch resource with "start from now
" or "start from a time" policy to watch events again.




