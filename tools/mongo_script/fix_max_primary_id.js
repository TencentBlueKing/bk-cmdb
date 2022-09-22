// usage
// `mongo --host 127.0.0.1 --port 27017 --username cc --password cc cmdb fix_max_primary_id.js`

/*
脚本用途：
    CMDB系统有原有版本升级到v3.5.1~v3.5.9间版本可能造成新增服务实例ID与由进程迁移过来的服务实例ID冲突问题
    
副作用：
    无

对应问题修复版本：v3.5.10
 */

function fixMaxIDField () {
    const collection = db.getSiblingDB('cmdb').cc_idgenerator ;
    const processTemplateMax = collection.findOne({ _id: 'cc_ProcessTemplate' }) ;
    const serviceInstanceMax = collection.findOne({ _id: 'cc_ServiceInstance' }) ;
    print("processTemplateMax:") ;
    printjson(processTemplateMax) ;
    print("serviceInstanceMax: ") ;
    printjson(serviceInstanceMax) ;
    let maxID = processTemplateMax.SequenceID ;
    if (serviceInstanceMax !== null) {
        maxID = processTemplateMax.SequenceID > serviceInstanceMax.SequenceID ? processTemplateMax.SequenceID : serviceInstanceMax.SequenceID ;
    }
    if (maxID > processTemplateMax.SequenceID) {
        collection.update({ _id: 'cc_ProcessTemplate' }, { $set: { SequenceID: maxID } }) ;
        print('fix cc_ProcessTemplate from ', processTemplateMax.SequenceID, ' to ', maxID) ;
    }
    if (serviceInstanceMax !== null && maxID > serviceInstanceMax.SequenceID) {
        collection.update({ _id: 'cc_ServiceInstance' }, { $set: { SequenceID: maxID } }) ;
        print('fix cc_ServiceInstance from ', serviceInstanceMax.SequenceID, ' to ', maxID) ;
    }
    print('max primary id doesn\'t need change') ;
}

fixMaxIDField() ;
