/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

$(document).ready(function() {
    $('.addedButtom').click(function(){
        adduser();
    });
    $(document).on("click", ".c-userlist-edit", evtEditRow);
    $(document).on("click", ".c-userlist-cancel", evtCancelEdit);
    $(document).on("click", ".c-userlist-delete", evtDelRow);
    $(document).on("click", ".c-userlist-save", evtSaveUser);
    $(document).on("click", ".c-userlist-reset", evtResetPassword);
    refreshGrid();
});

function enableRowEdit(row){
    var id = row.find('.id').text();
    var UserName = row.find('.UserName').text();
    var ChName = row.find('.ChName').text();
    var QQ = row.find('.QQ').text();
    var Tel = row.find('.Tel').text();
    var Email = row.find('.Email').text();
    var Role = row.find('.Role').attr('role');
    var disabledField = (curUserRole == 'admin') ? '' : 'disabled="disabled"';
    row.find('.UserName').html('<input required data-required-msg="请输入用户名" ' +
        'pattern="[A-Za-z0-9]{4,11}" validationMessage="用户名包含数字和字母，长度在4-10个字符" name="UserName" ' +
        'placeholder="请输入用户名" type="text" value="' + UserName + '" class="txt_username" ' + disabledField + '/>');
    row.find('.ChName').html('<input placeholder="请输入姓名" name="ChName" ' +
        'required data-required-msg="请输入姓名！" type="text" value="' + ChName + '" class="txt_chname" />');
    row.find('.QQ').html('<input type="text" value="'+QQ+'" class="txt_qq"' +
        'pattern="[0-9]{4,13}" validationMessage="QQ号只能是数字组合" name="QQ" ' + '/>');
    row.find('.Tel').html('<input placeholder="请输入手机号码" type="tel" name="Tel" ' +
        'pattern="\\d{11}" validationMessage="手机号码只能为11位数字" value="' + Tel + '" class="txt_tel"/>');
    row.find('.Email').html('<input placeholder="请输入常用邮箱" name="Email" ' +
        'validationMessage="邮箱格式不正确" pattern="^([\.a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$" type="email" value="' + Email + '" class="txt_email" />');
    var idHtml = row.find('.id').html();
    row.find('.id').html(idHtml + '<div style="display:none">' +
        '<span class="k-invalid-msg" data-for="UserName"></span>' +
        '<span class="k-invalid-msg" data-for="ChName"></span>' +
        '<span class="k-invalid-msg" data-for="Tel"></span>' +
        '<span class="k-invalid-msg" data-for="Email"></span>' +
        '</div>');
    row.find('.Role').html(
        '<select style="width:100%;" ' + disabledField + '>'+
           '<option value="admin">管理员</option>'+
           '<option value="user">用户</option>'+
       '</select>'
    );
    row.find('.Role select').val(Role);
    row.find('.Role select').select2({
        allowClear: true,
        minimumResultsForSearch: Infinity
    });
}

function evtEditRow(){
    var tr = $(this).closest('tr');
    $(this).hide();
    $(this).siblings().filter('a[name="deletes"]').hide();
    $(this).siblings().filter('a[name="saves"]').show();
    $(this).siblings().filter('a[name="cancels"]').show();
    enableRowEdit(tr);
}

function evtCancelEdit(){
    var tr = $(this).closest('tr');
    disableRowEdit(tr);
    refreshGrid();
    $(this).hide();
    $(this).siblings().filter('a[name="edits"]').show();
    $(this).siblings().filter('a[name="deletes"]').show();
    $(this).siblings().filter('a[name="saves"]').hide();
}

function evtDelRow(){
    var thisObj = this;
    var tr = $(thisObj).closest('tr');
    var UserName = tr.find('.UserName').text();
    if('admin' == UserName) {
        showWindows('admin用户不能删除', 'error');
        return;
    }
    var id = tr.find('.hid_id').val();
    var gridBatDel = dialog({
        title: '确认',
        width: 250,
        content: '是否删除选中用户',
        okValue: '确定',
        cancelValue: '取消',
        ok: function (){
            $.post("/account/delUser", {'id': id}, function(result){
                var resultInfo = $.parseJSON(result);
                if(resultInfo.success == false) {
                    window.location.reload();
                }else{
                    delete userList[id];
                    showWindows('删除成功！', 'success');
                    window.location.reload();
                }
            });
        },
        cancel: function () {
        }
    });
    gridBatDel.showModal();
}

function evtSaveUser(){
    var tr = $(this).closest('tr');
    var id = tr.find('.hid_id').val();
    var UserName = tr.find('.txt_username').val();
    var ChName = tr.find('.txt_chname').val();
    var QQ = tr.find('.txt_qq').val();
    var Tel = tr.find('.txt_tel').val();
    var Email = tr.find('.txt_email').val();
    var Role = tr.find('.Role select').select2('val');

    var form_validate = validate_save($(this));
    // alert(form_validate);
    if (form_validate) {
        var form_data = {};
        form_data['id'] = id;
        form_data['UserName'] = UserName;
        form_data['ChName'] = ChName;
        form_data['QQ'] = QQ;
        form_data['Tel'] = Tel;
        form_data['Email'] = Email;
        form_data['Role'] = Role;
        saveUserAjax(form_data);
        disableRowEdit(tr);
        $(this).hide();
        $(this).siblings().filter('a[name="cancels"]').hide();
        $(this).siblings().filter('a[name="edits"]').show();
        $(this).siblings().filter('a[name="deletes"]').show();
        $(this).siblings().filter('a[name="resets"]').show();
    } else {
        //表单验证未通过
        // var errors = form_validate.errors();
        // showWindows(errors.join('<br/>'), 'error');
        alert("表单验证未通过");
        return false;
    }
}

function disableRowEdit(row){
    var UserName = row.find('.txt_username').val();
    var ChName = row.find('.txt_chname').val();
    var QQ = row.find('.txt_qq').val();
    var Tel = row.find('.txt_tel').val();
    var Email = row.find('.txt_email').val();
    var Role = row.find('.Role select').select2('data').text;
    row.find('.UserName').html(UserName);
    row.find('.ChName').html(ChName);
    row.find('.QQ').html(QQ);
    row.find('.Tel').html(Tel);
    row.find('.Email').html(Email);
    row.find('.Role').html(Role);
}

function refreshGrid(){
    if(!userList){
        userList = {};
    }
    var tbl_html = '';
    //var grid = $('#grid_userManager').data('kendoGrid');
    // if(grid){
    //     grid.destroy();
    //     $('#grid_userManager').remove();
    // }
    // var gridLength=$('#grid_userManager').length;
    // if(gridLength!=0){
    //     //grid.destroy();
    //     //$('#grid_userManager').remove();
    // }
    $("#ctn_userlist").html($("#tpl_grid").html());
    for(userId in userList){
        var userInfo = userList[userId];
        if(userInfo['Role'] == 'admin'){
            userInfo['RoleName'] = '管理员';
        }else{
            userInfo['RoleName'] = '用户';
        }
        // tbl_html += kendo.template($("#tpl_user_tr").html())(userInfo);
    }
    // $("#tbl_userlist_body").html(tbl_html);
    // $("#grid_userManager").kendoGrid({
    //     pageable: {
    //         pageSize: 10,
    //         buttonCount: 3,
    //         refresh: false
    //     },
    //     scrollable: false,
    //     dataBound: function(e) {
    //         $("#grid_userManager").css("opacity","1")
    //     }
    // });
    var language = {
        search: '搜索：',
        lengthMenu: "每页显示 _MENU_ 记录",
        zeroRecords: "没找到相应的数据！",
        info: "分页 _PAGE_ / _PAGES_",
        infoEmpty: "暂无数据！",
        infoFiltered: "(从 _MAX_ 条数据中搜索)",
        paginate: {
            first: '首页',
            last: '尾页',
            previous: '上一页',
            next: '下一页',
        }
    };
    $('#grid_userManager').dataTable({
        retrieve: true,
        paging: true, //隐藏分页
        ordering: false, //关闭排序
        info: false, //隐藏左下角分页信息
        searching: false, //关闭搜索
        language:language,
        lengthChange: false //不允许用户改变表格每页显示的记录数
    });
}

function evtResetPassword() {
    var id = $(this).closest('tr').find('.hid_id').val();
    var userName = userList[id].UserName;
    // var form_content = kendo.template($("#tpl_reset_password").html())({});
    var form_content = '<div class="pt10" style="overflow:hidden"><form class="form-horizontal validate-form" action="index.html"><div class="form-group"><label class="control-label col-sm-3">新密码</label><div class="col-sm-7"><input class="form-control password_text" placeholder="请输入密码" type="password"></div></div><div class="form-group"><label class="control-label col-sm-3">确认密码</label><div class="col-sm-7"><input class="form-control password_text2" placeholder="请输入密码" type="password"><p class="userManager_tips"></p></div></div></form></div>';
    var newdialoga = dialog({
        title:'密码重置',
        width:500,
        content:form_content,
        okValue:"重置",
        cancelValue:"取消",
        skin:'dia-grid-batDel',
        ok:function (){
            var password_text=$('.password_text').val();
            var password_text2=$('.password_text2').val();
            if(password_text==''||password_text2==''){
                $('.userManager_tips').html('*输入的密码不能为空');
                return false;
            }else if(password_text!=password_text2){
                $('.userManager_tips').html('*密码不一致，请重新输入');
                return false;
            }else if(password_text==password_text2){
                $('.userManager_tips').html('');
            }
            var postData = {
                'UserName': userName,
                'Password': password_text
            };
            $.post("/account/changePassword", postData, function(result){
                var resultInfo = $.parseJSON(result);
                if(resultInfo.success == false) {
                    showWindows('重置密码失败！' + resultInfo.message, 'error');
                }else{
                    showWindows('重置密码成功！', 'success');
                }
            });
        },
        cancel: function () {
        }
    });
    //弹出框初始化
    newdialoga.showModal();
}

function adduser() {
    var form_content = '<div class="pt10" style="overflow:hidden">'+
        '<form class="form-horizontal validate-form" action="/account/saveUser" id="frm_new_user">'+
            '<input type="hidden" value="" name="id">'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">用户名</label>'+
                '<div class="col-sm-7">'+
                    '<input class="form-control" required="" data-required-msg="请输入用户名！" pattern="[A-Za-z0-9]{4,11}" validationmessage="用户名只能包含数字和字母，长度在4-10个字符" name="UserName" placeholder="请输入用户名" type="text">'+
                    '<span class="k-invalid-msg" data-for="UserName"></span>'+
                '</div>'+
            '</div>'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">姓名</label>'+
                '<div class="col-sm-7">'+
                    '<input class="form-control" placeholder="请输入姓名" name="ChName" required="" data-required-msg="请输入姓名！" type="text">'+
                    '<span class="k-invalid-msg" data-for="ChName"></span>'+
                '</div>'+
            '</div>'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">QQ</label>'+
                '<div class="col-sm-7">'+
                    '<input class="form-control" placeholder="请输入QQ" name="QQ" type="text" pattern="[0-9]{4,13}" validationmessage="QQ号只能是数字组合">'+
                    '<span class="k-invalid-msg" data-for="ChQQ"></span>'+
                '</div>'+
            '</div>'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">手机号码</label>'+
                '<div class="col-sm-7">'+
                    '<input class="form-control" placeholder="请输入手机号码" type="tel" name="Tel" pattern="\d{11}" validationmessage="手机号码只能为11位数字">'+
                    '<span class="k-invalid-msg" data-for="Tel"></span>'+
                '</div>'+
            '</div>'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">常用邮箱</label>'+
                '<div class="col-sm-7">'+
                    '<input class="form-control" placeholder="请输入常用邮箱" name="Email" validationmessage="邮箱格式不正确" type="email" pattern="^([\.a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((.[a-zA-Z0-9_-]{2,3}){1,2})$">'+
                    '<span class="k-invalid-msg" data-for="Email"></span>'+
                '</div>'+
            '</div>'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">角色</label>'+
                    '<div class="col-sm-7" id="sel_add_user">'+
                             '<label class="radio-inline"><input value="user" type="radio" name="Role" checked="checked">普通用户</label>'+
                             '<label class="radio-inline"><input value="admin" type="radio" name="Role">管理员</label>'+
                '</div>'+
            '</div>'+
            '<div class="form-group">'+
                '<label class="control-label col-sm-3">默认密码</label>'+
                '<label class="control-label col-sm-7 tl text-primary">blueking</label>'+
            '</div>'+
        '</form>'+
    '</div>';
    var newdialoga = dialog({
        title:'新增用户',
        width:500,
        content: form_content,
        okValue:"新增",
        cancelValue:"取消",
        skin:'dia-grid-batDel',
        ok:function (){
            // var form_validate = $("#frm_new_user").kendoValidator().data("kendoValidator");
            // 表单验证
            var bz_x=1;
            $(".form-group input.form-control").each(function(){
                var x_index=$(".form-group input.form-control").index(this);
                if(x_index == 0 || x_index == 1 ){
                    if($(this).val() == ""){
                        bz_x=0;
                        x_index == 0?$(this).siblings("span.k-invalid-msg").html("请输入用户名！"):$(this).siblings("span.k-invalid-msg").html("请输入姓名！");
                    }else{
                        if(/[A-Za-z0-9]{4,11}/.test($(this).val())){
                            //
                        }else{
                            bz_x=0;
                            $(this).siblings("span.k-invalid-msg").html($(this).attr("validationmessage"));
                        }
                    }
                }else{
                    switch(x_index){
                        case 2:// QQ验证
                            if($(this).val()){
                                if(!/[0-9]{4,13}/.test($(this).val())){
                                    bz_x=0;
                                    $(this).siblings("span.k-invalid-msg").html($(this).attr("validationmessage"));
                                }
                            }
                            break;
                        case 3:// 手机号码
                            if($(this).val()){
                                if(!/\d{11}/.test($(this).val())){
                                    bz_x=0;
                                    $(this).siblings("span.k-invalid-msg").html($(this).attr("validationmessage"));
                                }
                            }
                            break;
                        case 4:// 常用邮箱
                            if($(this).val()){
                                if(!/^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((.[a-zA-Z0-9_-]{2,3}){1,2})$/.test($(this).val())){
                                    bz_x=0;
                                    $(this).siblings("span.k-invalid-msg").html($(this).attr("validationmessage"));
                                }
                            }
                            break;
                        default:
                            break;
                    }
                }
            })

            if (bz_x==1) {
                //表单验证通过
                var form_data = $("#frm_new_user").serialize();
                saveUserAjax(form_data);
            } else {
                //表单验证未通过
                return false;
            }
        },
        cancel: function () {
        }
    });
    //弹出框初始化
    newdialoga.showModal();
}

function saveUserAjax(user_info){
    $.post("/account/saveUser", user_info, function(result){
        var resultInfo = $.parseJSON(result);
        if(resultInfo.success == false) {
            showWindows('保存用户失败！' + resultInfo.message, 'error');
            window.location.reload();
            return false;
        }else{
            showWindows('保存用户成功！', 'success');
            userList[resultInfo.user.id] = resultInfo.user;
            window.location.reload();
            return true;
        }
    });
}

function validate_save(Obj){
            var bz_x=1;
            Obj.parent().siblings().each(function(){
                var x_index=$(this).index();
                if(x_index == 1 || x_index == 2 ){
                    if($(this).find("input").val() == ""){
                        bz_x=0;
                    }else{
                        if(/[A-Za-z0-9]{4,11}/.test($(this).find("input").val())){
                            //
                        }else{
                            if(x_index==1){
                                bz_x=0;
                            }
                        }
                    }
                }else{
                    switch(x_index){
                        case 3:// QQ验证
                            if($(this).find("input").val()){
                                if(!/[0-9]{4,13}/.test($(this).find("input").val())){
                                    bz_x=0;
                                }
                            }
                            break;
                        case 4:// 手机号码
                            if($(this).find("input").val()){
                                if(!/\d{11}/.test($(this).find("input").val())){
                                    bz_x=0;
                                }
                            }
                            break;
                        case 5:// 常用邮箱
                            if($(this).find("input").val()){
                                if(!/^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((.[a-zA-Z0-9_-]{2,3}){1,2})$/.test($(this).find("input").val())){
                                    bz_x=0;
                                }
                            }
                            break;
                        default:
                            break;
                    }
                }

            });
            if(bz_x==0){
                    return false;
                }else{
                    return true;
                }
}
