### Functional description

Install the host to the blue whale business as follows:
1.  Can only operate blue whale business
2. Host can not be transferred to built-in modules such as idle and failed machines
3. Host modules that already exist in the host will not be deleted, only new hosts and modules will be added. 
4. The non-existent host will be added, and the rule will judge whether the host exists through the intranet IP and cloud id
5. No error reported if process does not exist

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_set_name  | string     | yes  | The set name where the host resides |
| bk_module_name | string  |yes   | Module name where the host resides|
| bk_host_innerip | string  |yes   | Host intranet IP|
| bk_cloud_id | int  |no   | Cloud area where the host is located, default 0 |
| host_info | object  |no   | Host details, all fields and values of host model correspond|
| proc_info | object |no| The value of the process in the service instance of the host under the current module, {"process name":{"Process properties": value}}, referencing the process model|




### Request Parameters Example

```python

{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_set_name":"set1",
    "bk_module_name":"module2",
    "bk_host_innerip":"127.0.0.1",
    "bk_cloud_id":0,
    "host_info":{
            "bk_comment":"test bk_comment 1",
            "bk_os_type":"1"
    },
    "proc_info":{
            "p1":{"description":"xxx"}
    }
}

```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
}
```

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |


