### Functional description

Install the host under the blueking business, the details are as follows:
Can only operate blueking business
2. Cannot transfer host to built-in modules such as idle machine and failure machine
3. The existing host module of the host will not be deleted, only the host and module will be added.
4. Non-existing hosts will be added. The rule will judge whether the host exists by the internal IP and cloud id.
5. Process does not exist without error

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description       | 
|----------------------|------------|--------|-------------|
| bk_set_name | string | Yes | The name of the cluster where the host is located |
| bk_module_name | string | Yes | Module name where host is |
| bk_host_innerip | string | Yes | Host Intranet IP |
| bk_cloud_id | int | No | The cloud area where the host is located. The default value is 0 |
host_info | object | No | host details, all fields of the host model and the corresponding |
| proc_info | object | No | The value of the process in the service instance of the host under the current module, {"process name": {"process attribute": value}}, reference process model |


### Request Parameters Example

``` python

{
        "bk_set_name":"set1",
        "bk_module_name":"module2",
        "bk_host_innerip":"127.127.0.1",
        "bk_cloud_id":0,
        "host_info":{
                "bk_comment":"test bk_comment 1",
                "bk_os_type":"1"
        },
        "proc_info":{
                "p1":{"bind_ip":"127.127.0.1"}
        }
}

```

### Return Result Example


```python

{
    "result": true,
    "code": 0,
    "message": "",

}

```
### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
