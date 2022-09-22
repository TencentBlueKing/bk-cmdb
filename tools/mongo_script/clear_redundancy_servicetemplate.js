// usage
// `mongo --host 127.0.0.1 --port 27017 --username cc --password cc cmdb clear_redundancy_servicetemplate.js`

/*
脚本用途：
    CMDB升级到v3.5给没有进程的模块绑定了空服务模板, 本脚本用于删除这些空服务模板，并解除模块对其引用关系
    
清理内容：
- 空服务模板
- 模块的服务模板字段
- 空服务模板创建出来的服务实例

- 删除检查：如果服务模板下新建了进程，则不应该删除

对应问题修复版本：v3.5.14
 */

function clearRedundancyServiceTemplate () {
    const collection = database.cc_ApplicationBase ;
    let cursor = collection.find();
    while ( cursor.hasNext() ) {
        let biz = cursor.next() ;
        print("------------------------ 华丽分割线 ----------------------------")
        print("业务: ", biz.bk_biz_name) ;
        let logPrefix = "    ";
        bizClear(logPrefix, biz) ;
        print("\n\n\n") ;
    }
}

function bizClear(logPrefix, biz) {
    const collection = database.cc_ModuleBase ;
    let filter = {
        bk_biz_id: biz.bk_biz_id,
        default: 0
    } ;
    let cursor = collection.find(filter);
    while ( cursor.hasNext() ) {
        let module = cursor.next() ;
        print(logPrefix, "模块: ", module.bk_module_name) ;
        moduleClear(logPrefix, biz, module) ;
        print("\n") ;
    }
}

function moduleClear(logPrefix, biz, module) {
    const collection = database.cc_Proc2Module ;
    let filter = {
        bk_biz_id: biz.bk_biz_id,
        bk_module_name: module.bk_module_name
    } ;
    let count = collection.find(filter).count() ;
    if (count > 0) {
        print(logPrefix, "模块下有 ", count, " 个进程，不需要清理服务模板") ;
        return
    }
    
    let serviceTemplateID = module.service_template_id ;
    let processTemplateFilter = {
        service_template_id: serviceTemplateID
    };
    let processCount = database.cc_ProcessTemplate.find(processTemplateFilter).count()
    if (processCount > 0) {
        print(logPrefix, "Warning: 模块下有 ", count, " 个进程模板（升级到3.5后添加的进程模板），不需要清理服务模板") ;
        return
    }

    let serviceInstanceFilter = {
        service_template_id: serviceTemplateID
    };
    let serviceInstanceCount = database.cc_ServiceInstance.find(serviceInstanceFilter).count()
    if (serviceInstanceCount > 0) {
        print(logPrefix, "Warning: 服务模板有 ", serviceInstanceCount, " 个实例，不需要清理服务模板");
        return
    }
    
    logPrefix = logPrefix + "    "
    doRealClear(logPrefix, biz, module) ;
}

function doRealClear(logPrefix, biz, module) {
    // step1. clear related service instance
    let serviceInstanceFilter = {
        service_template_id: module.service_template_id
    } ;
    let removeServiceInstanceResult = database.cc_ServiceInstance.remove(serviceInstanceFilter) ;
    print(logPrefix, "删除服务实例返回：", removeServiceInstanceResult) ;
    
    // step2. update module service template id field
    let moduleUpdateFilter = {
        bk_module_id: module.bk_module_id
    } ;
    let moduleUpdateDoc = {
        service_template_id: 0
    } ;
    let moduleUpdateResult = database.cc_ModuleBase.updateOne(moduleUpdateFilter, {$set: moduleUpdateDoc}) ;
    print(logPrefix, "更新模块属性字段返回：", moduleUpdateResult) ;
    
    // step3. delete service template
    let serviceTemplateFilter = {
        id: module.service_template_id
    };
    let serviceTemplateDeleteResult = database.cc_ServiceTemplate.remove(serviceTemplateFilter) ;
    print(logPrefix, "服务模板删除返回：", serviceTemplateDeleteResult) ;
}



// start 
const database = db.getSiblingDB('cmdb') ;

clearRedundancyServiceTemplate() ;
