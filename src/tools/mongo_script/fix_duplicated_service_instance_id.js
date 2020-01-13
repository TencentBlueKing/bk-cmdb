// usage:
// `mongo --host 127.0.0.1 --port 27017 --username cc --password cc cmdb fix_duplicated_service_instance_id.js`

/*
脚本用途： 修复服务实例ID冲突问题
实现逻辑：对出现问题的服务实例重新分配ID，并修复cc_ProcessInstanceRelation表数据
 */



// 出现问题的最小和最大服务实例ID
MinServiceInstanceID = 48;
MaxServiceInstanceID = 70;
// 出现问题的升级时间点
const MigrationTime = ISODate("2019-09-25T07:41:57.570Z");

function getTopoPath(bk_module_id) {
	let m = db.cc_ModuleBase.findOne({bk_module_id: bk_module_id});
	let s = db.cc_SetBase.findOne({bk_set_id: m.bk_set_id});
	let b = db.cc_ApplicationBase.findOne({bk_biz_id: m.bk_biz_id});
	return b.bk_biz_name + ":" + b.bk_biz_id + "->" + s.bk_set_name + ":" + s.bk_set_id + "->" + m.bk_module_name + ":" + m.bk_module_id;
}

function getNextSequence(name) {
    let ret = db.cc_idgenerator.findAndModify(
        {
            query: { _id: name },
            update: { $inc: { SequenceID: 1 } },
            new: true
        }
    ) ;
    return ret.SequenceID;
}


let cursor = db.cc_ServiceInstance.find({id: {$gte: MinServiceInstanceID, $lte: MaxServiceInstanceID}, create_time: {$gt: MigrationTime}});
while ( cursor.hasNext() ) {
	let item = cursor.next();
	const topoPath = getTopoPath(item.bk_module_id);
	print(topoPath + "==> ServiceInstanceID:" + item.id);
	const newID = getNextSequence("cc_ServiceInstance");
	const instanceFilter = {
		_id: item._id,
		id: item.id,
		create_time: {$gt: MigrationTime}
	};
	const r1 = db.cc_ServiceInstance.update(instanceFilter, {$set: {id: newID}}, { $multi: true });
	const relationFilter = {
		bk_host_id: item.bk_host_id,
		service_instance_id: item.id
	};
	const r2 = db.cc_ProcessInstanceRelation.update(relationFilter, {$set: {service_instance_id: newID}}, { $multi: true });
	print("update instance result:", r1, "update relation result:", r2, "\n");
}
