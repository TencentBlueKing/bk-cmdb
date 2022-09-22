// usage:
// `mongo --host 127.0.0.1 --port 27017 --username cc --password cc cmdb fix_service_category_field.js`

/*
脚本用途： 修复模块的服务分类为0导致无法编辑问题
 */

const cmdb = db.getSiblingDB('cmdb');

function getNextSequence (name) {
    const ret = cmdb.cc_idgenerator.findAndModify(
        {
            query: { _id: name },
            update: { $inc: { SequenceID: 1 } },
            new: true
        }
    );
    return ret.SequenceID;
}

function getDefaultServiceCategory () {
    const filter = { is_built_in: true, name: 'Default', bk_parent_id: { $ne: 0 } };
    const ret = db.getCollection('cc_ServiceCategory').findOne(filter);
    return ret.id;
}

function fixModuleServiceCategoryField () {
    const defaultServiceCategoryID = getDefaultServiceCategory();
    print('defaultServiceCategoryID: ', defaultServiceCategoryID);

    const collection = cmdb.cc_ModuleBase;
    const filter = { service_category_id: 0 };
    let count = collection.count(filter);
    print('service_category_id:0 count: ', count);

    let cursor = collection.find(filter);
    while (cursor.hasNext()) {
        const module = cursor.next();
        print('module to be fix: biz', JSON.stringify(module));
    }
    print('\n\n');

    const doc = { $set: { service_category_id: defaultServiceCategoryID } };
    const updateResult = collection.update(filter, doc, { multi: true });
    print('update result: ', updateResult);

    count = collection.count(filter);
    print('service_category_id:0 count: ', count);
}

function changeProcessName () {
    const mapping = {
        bk_process_name: '进程别名',
        bk_func_name: '进程名称',
        description: '备注'
    };

    for (const key in mapping) {
        const cond = {
            bk_obj_id: 'process',
            bk_property_id: key
        };
        const doc = { bk_property_name: mapping[key] };
        const updateResult = cmdb.cc_ObjAttDes.update(cond, doc);
        print('update process attribute ', key, 'result: ', updateResult);
    }

    const cond = { bk_obj_id: 'process', 'bk_property_id': 'bk_process_name' };
    const updateResult = cmdb.cc_ObjAttDes.update(cond, { 'isrequired': false });
    print('updateResult: ', updateResult);

    const procNameAttr = cmdb.cc_ObjAttDes.findOne({ bk_obj_id: 'process', bk_property_id: 'bk_process_name' });
    const bizIDAttr = cmdb.cc_ObjAttDes.findOne({ bk_obj_id: 'process', bk_property_id: 'bk_biz_id' });
    const deleteFilter = {
        'bk_obj_id': 'process',
        'keys': {
            '$all': [
                {
                    'key_kind': 'property',
                    'key_id': procNameAttr.id
                },
                {
                    'key_kind': 'property',
                    'key_id': bizIDAttr.id
                }
            ]
        }
    };
    const deleteResult = db.cc_ObjectUnique.delete(deleteFilter);
    print('deleteResult: ', deleteResult);
}

// run fix
fixModuleServiceCategoryField();
changeProcessName();
