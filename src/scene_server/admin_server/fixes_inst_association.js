db.cc_ObjAsst.find().forEach(function (objasst) {
    var ret = db.cc_InstAsst.update(
        {
            "bk_obj_id": objasst.bk_obj_id,
            "bk_asst_obj_id": objasst.bk_asst_obj_id,
            "bk_supplier_account": objasst.bk_supplier_account
        },
        {
            "$set": {
                "bk_obj_asst_id": objasst.bk_obj_asst_id,
                "bk_asst_id": objasst.bk_asst_id,
                "last_time": new Date(),
            }
        })
    print(ret)
});
