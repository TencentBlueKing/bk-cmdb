### Function Description

Install a host to the BlueKing business with the following details:

1. Can only operate in the BlueKing business.
2. Cannot transfer hosts to built-in modules such as idle machines and faulty machines.
3. Will not delete existing host modules; it will only add the host to the module.
4. Newly add hosts if they do not exist; the rule is determined by the intranet IP and cloud ID to check if the host exists.
5. No error will be reported if the process does not exist.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field           | Type   | Required | Description                                                  |
| --------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_set_name     | string | Yes      | Name of the cluster where the host is located                |
| bk_module_name  | string | Yes      | Name of the module where the host is located                 |
| bk_host_innerip | string | Yes      | Intranet IP of the host                                      |
| bk_cloud_id     | int    | No       | Control area where the host is located; default is 0         |
| host_info       | object | No       | Host details, corresponding to all fields and values of the host model |
| proc_info       | object | No       | Values of processes in the service instance under the current module, {"process name": {"process property": value}}, refer to the process model |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_set_name": "set1",
    "bk_module_name": "module2",
    "bk_host_innerip": "127.0.0.1",
    "bk_cloud_id": 0,
    "host_info": {
        "bk_comment": "test bk_comment 1",
        "bk_os_type": "1"
    },
    "proc_info": {
        "p1": {"description": "xxx"}
    }
}
```

### Response Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |