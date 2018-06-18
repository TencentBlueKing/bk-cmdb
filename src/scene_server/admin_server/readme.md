admin_server


Usage ofcmdb_adminserver bkbiz:
    --import[=false] import flag
    --export[=false] export flag
    --config="" config file path, normaly migrate.conf
    --file="" the export or export file path
    --dryrun[=false] dryrun print what we will do but not execute to db
example usage:
cmdb_adminserver bkbiz --export --config /data/cmdb/cmdb_adminserver/configures/migrate.conf --file export_2018_06_18_14_59_00.json --dryrun
cmdb_adminserver bkbiz --import --config /data/cmdb/cmdb_adminserver/configures/migrate.conf --file export_2018_06_18_14_59_00.json --dryrun
