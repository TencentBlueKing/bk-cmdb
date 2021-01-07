db.cc_ObjAsst.find({ "bk_asst_id": "bk_mainline", "bk_obj_id": {"$nin": ["set","module","host"]} }).forEach(function (myDoc) {
    var ret = db.cc_ObjAsst.remove({"bk_asst_id": "bk_mainline", "bk_obj_id": myDoc.bk_obj_id})
    print("delete cc_ObjAsst for ", myDoc.bk_obj_id, " result: ", ret)
    ret = db.cc_ObjAttDes.remove({"bk_obj_id": myDoc.bk_obj_id})
    print("delete cc_ObjAttDes for ", myDoc.bk_obj_id, " result: ", ret)
    ret = db.cc_ObjectBase.remove({"bk_obj_id": myDoc.bk_obj_id})
    print("delete cc_ObjectBase for ", myDoc.bk_obj_id, " result: ", ret)
});

var ret = db.cc_ObjAsst.update({"bk_asst_id": "bk_mainline", "bk_obj_id": "set"}, {"$set": {"bk_asst_obj_id": "biz"}});
print("update cc_ObjAsst for ", "set", " result: ", ret)

db.cc_SetBase.find().forEach(function (myDoc) {
    ret = db.cc_SetBase.update({"bk_set_id": myDoc.bk_set_id,"bk_supplier_account": myDoc.bk_supplier_account}, {"$set": {"bk_parent_id": myDoc.bk_biz_id}})
    print("update cc_SetBase for ", myDoc.bk_set_id, " result: ", ret)
});

print("done")
