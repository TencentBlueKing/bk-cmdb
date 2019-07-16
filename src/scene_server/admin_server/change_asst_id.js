// params:
var bk_obj_asst_id = "a_default_b"
var newAsstID = "run"
var bk_supplier_account = "0"

// execute
var objAsst = db.cc_ObjAsst.findOne({ "bk_obj_asst_id": bk_obj_asst_id, "bk_supplier_account": bk_supplier_account })
if (objAsst !=null){

    var new_bk_obj_asst_id = objAsst.bk_obj_id + "_" + newAsstID + "_" + objAsst.bk_asst_obj_id
    db.cc_ObjAsst.update({ "bk_obj_asst_id": bk_obj_asst_id, "bk_supplier_account": bk_supplier_account }, { "$set": { "bk_asst_id": newAsstID, "bk_obj_asst_id": new_bk_obj_asst_id } })
    
    db.cc_InstAsst.update({ "bk_obj_asst_id": bk_obj_asst_id, "bk_supplier_account": bk_supplier_account }, { "$set": { "bk_asst_id": newAsstID, "bk_obj_asst_id": new_bk_obj_asst_id } }, { "multi": 1 })
}
