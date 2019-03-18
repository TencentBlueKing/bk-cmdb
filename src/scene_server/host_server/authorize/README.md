# Design of host authorization base on iam

## Resource_id format
- common format: `business/{business}:set/{set}:module/{module}:host/{host_id}`


#### Authorization for add host
- resource_id: `business/{business}:set/{set}:module/{module}:host`
- action: transferhost

####  Authorization for update/read/delete host
- resource_id: `business/{business}:set/{set}:module/{module}:host/{host_id}`
- action: update/read/delete


#### Deal with one host belong to multiple modules
- respect as multiple iam resource
- ex: host_id belong to module1 and module 2, then there will be two resource
    + resource_id: `business/{business}:set/{set}:module/{module1}:host/{host_id}`
    + resource_id: `business/{business}:set/{set}:module/{module2}:host/{host_id}`
