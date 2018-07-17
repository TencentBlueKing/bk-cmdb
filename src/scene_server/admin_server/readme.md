## admin_server


### Usage of cmdb_adminserver bkbiz:
```
      --config="conf/api.conf": The config path. e.g conf/api.conf
      --dryrun[=false]: dryrun flag, if this flag seted, we will just print what we will do but not execute to db
      --export[=false]: export flag
      --file="": export or import filepath
      --import[=false]: import flag
      --mini[=false]: mini flag, only export required fields
      --scope="all": export model, could be [biz] or [process], default all
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
