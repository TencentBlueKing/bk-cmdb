## admin_server


### Usage of cmdb_adminserver bkbiz:
```
    --import[=false] import flag
    --export[=false] export flag
    --config="" config file path, normaly migrate.conf
    --file="" the export or export file path
    --dryrun[=false] dryrun print what we will do but not execute to db
```
#### example usage:

- export:
```sh
cmdb_adminserver bkbiz --export --config /data/cmdb/cmdb_adminserver/configures/migrate.conf --file bkbiz_export_2018_06_18_14_59_00.json
```

- dryrun import:
```sh
cmdb_adminserver bkbiz --import --config /data/cmdb/cmdb_adminserver/configures/migrate.conf --file bkbiz_export_2018_06_18_14_59_00.json --dryrun
```

- import:
```sh
cmdb_adminserver bkbiz --import --config /data/cmdb/cmdb_adminserver/configures/migrate.conf --file bkbiz_export_2018_06_18_14_59_00.json
```
