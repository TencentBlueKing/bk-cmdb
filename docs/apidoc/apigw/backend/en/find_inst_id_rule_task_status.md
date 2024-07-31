### Description

Query the status of synchronization instance id rule field (Version：v3.14.1，Permission：None)

### Parameters

| Name      | Type   | Required | Description |
|-----------|--------|----------|-------------|
| task_id | string | Yes      | task id     |

### Request Example

```json
{
  "task_id": "111"
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": {
    "status": "finished"
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Data returned by the request                                        |

#### data

| Name       | Type   | Description |
|------------|--------|------------|
| status     | string | Task execution status: "new" (new), "waiting" (waiting for execution), "executing" (executing), "finished" (completed), "failure" (failed)    |