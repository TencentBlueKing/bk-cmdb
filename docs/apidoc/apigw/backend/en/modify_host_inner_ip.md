### Description

Modify the inner IPv4/IPv6 addresses of a static-addressing host. This is a dangerous API and is disabled by default.
The administrator must enable `hostServer.enableModifyHostInnerIP` in `common.yaml` or `extra.yaml`.

Permission: host edit permission + global configuration management permission.

Limitations:

1. The host must use static addressing (`bk_addressing = static`).
2. Only one host can be modified per request.
3. `src_*` confirms the current IP list. Its content must match the current list; the order does not matter.
4. `dest_*` is written in the given order, and the first IP is used as the primary IP.
5. Abbreviated IPv6 addresses are accepted and normalized to the full format before being saved.
6. If the stored IPv4/IPv6 value is invalid, `1110072` is returned and the data must be repaired first.

### Parameters

| Name       | Type  | Required | Description |
|------------|-------|----------|-------------|
| bk_host_id | int   | Yes      | Host ID |
| src_ipv4   | array | No       | Current IPv4 list, must be provided together with `dest_ipv4` |
| src_ipv6   | array | No       | Current IPv6 list, must be provided together with `dest_ipv6` |
| dest_ipv4  | array | No       | Target IPv4 list, written in the given order |
| dest_ipv6  | array | No       | Target IPv6 list, written in the given order |

#### src/dest rules

- If a field is omitted or set to `null`, that family is not modified.
- `src` and `dest` of the same family must be provided together or omitted together; providing only one returns `1199010`.
- An empty `dest` array (`[]`) clears that family, but `src` must also be provided and match the current value.
- An empty `src` array (`[]`) confirms that the host currently has no IP of that family; `dest` must also be provided.
- If all four IP fields are omitted or `null`, `1199010` is returned.

### Request Example

```json
{
    "bk_host_id": 111,
    "src_ipv4": ["127.0.0.1"],
    "src_ipv6": ["::1"],
    "dest_ipv4": ["192.168.11.1"],
    "dest_ipv6": ["2001:db8:85a3::8a2e:370:7334"]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### Response Parameters

| Name       | Type   | Description |
|------------|--------|-------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error |
| message    | string | Error message returned in case of request failure |
| permission | object | Permission information |
| data       | object | Data returned in the request |

### Common Errors

| Error Code | Description |
|------------|-------------|
| `1110068` | The host inner IP modification feature is not enabled; contact the administrator |
| `1110008` | Host not found |
| `1110070` | The host is not using static addressing |
| `1110071` | `src` does not match the current host value |
| `1110072` | The stored IPv4/IPv6 value of the host is invalid |
| `1199010` | `src` / `dest` are not paired, or all four IP fields are empty |
| `1199014` | The destination IP is already used by another host in the same cloud area; the duplicated field is returned |
| `1199070` | Invalid IP format |
| `1199047` / `1199048` | Permission-related errors |
