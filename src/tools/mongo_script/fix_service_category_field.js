// usage:
// `mongo --host 127.0.0.1 --port 27017 --username cc --password cc cmdb fix_service_category_field.js`

/*
脚本用途： 修复模块的服务分类为0导致无法编辑问题
 */

const cmdb = db.getSiblingDB('cmdb') ;

function getDefaultServiceCategory() {
	let ret = db.getCollection('cc_ServiceCategory').findOne({is_built_in: true, name: "Default", bk_parent_id: {$ne: 0}}) ;
	return ret.id;
}

function fixModuleServiceCategoryField() {
	let defaultServiceCategoryID = getDefaultServiceCategory();
	print("defaultServiceCategoryID: ", defaultServiceCategoryID);

	const collection = cmdb.cc_ModuleBase;
	const filter = {service_category_id: 0};
	let count = collection.count(filter);
	print("service_category_id:0 count: ", count);

	let cursor = collection.find(filter);
	while ( cursor.hasNext() ) {
		let module = cursor.next();
		print("module to be fix: biz", JSON.stringify(module));
	}
	print("\n\n");

	let doc = {$set: {service_category_id: defaultServiceCategoryID}};
	let update_result = collection.update(filter, doc, {multi:true});
	print("update result: ", update_result);

	count = collection.count(filter);
	print("service_category_id:0 count: ", count);
}

// run fix
fixModuleServiceCategoryField() ;
