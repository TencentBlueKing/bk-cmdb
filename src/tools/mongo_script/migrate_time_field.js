// usage
// mongo migrate_time_field.js

/*
脚本用途：修复cmdb业务列表接口返回的时间字符串不统一问题
默认情况返回的格式串为: '2019-09-03T03:41:40.909Z'
异常情况返回的格式串为: '2019-09-03T03:41:40Z'

修复原理：
如果检测到 create_time/last_time 字段没有毫秒，则将号码段设置成001，比如上述异常返回结果修正为: '2019-09-03T03:41:40.001Z'
 */

let repeatCount = 40;
let startStr = '<'.repeat(repeatCount);
let endStr = '>'.repeat(repeatCount);

function fixTimeField(collection, fieldName, primaryKey) {
    let cursor = collection.find();
    while ( cursor.hasNext() ) {
        let item = cursor.next();
        print(startStr);
        printjson(item);
        for (let idx=0;idx<fieldName.length;idx++) {
            let field = fieldName[idx];
            print("field: ", field)
            if (item[field].getMilliseconds() == 0) {
                let before = item[field];
                let uniformTime = new Date(before.toISOString());
                uniformTime.setMilliseconds(1);
                print(before.toLocaleString(), " ===> ", uniformTime.toISOString(), "\n");
                let updateData = {};
                updateData[field] = uniformTime;
                let filter = {};
                filter[primaryKey] = item[primaryKey];
                collection.update(filter, {$set: updateData});
            }
        }
        print(endStr, "\n".repeat(2));
    };
}

// fix biz table time field
let collection = db.getSiblingDB('cmdb').cc_ApplicationBase;
const primaryKey = "bk_biz_id";
const fieldName = ["create_time", "last_time"];
fixTimeField(collection, fieldName, primaryKey);
