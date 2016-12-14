/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

$(document).ready(function(){
    var setDis = false;
    if(2 == alevel) {
        setDis = true;
    }
    var clusterList = {
        init:function(){
            this.initClusterTable($('#cluster_table'));     //初始化集合表格
            this.initModuleTable($('#module_table'));       //初始化模块表格
            this.initToolbarEvent();                        //工具栏初始化
            this.initBatchCheckbox($('.cluster-edit-cluster'));
            this.initBatchCheckbox($('.cluster-edit-module'));
        },
        CloneSetID:0,
        EnviType:[{ id: 1, text: '测试' }, { id: 2, text: '体验' }, { id: 3, text: '正式' } ],
        environmentType:{ 1:'测试', 2: '体验',3:'正式' } ,
        ServiceStatus:[{ id: 0, text: '关闭' }, { id: 1, text: '开放' }],
        serverStatus: {0:'关闭',1:'开放'},
        Operator:userlist,
        BakOperator:userlist,
        clusterTableColumns:function(){
            var topThis = this;
            return [{
                    title:'<input type="checkbox" class="checked_all">',
                    width:'15px',
                    defaultContent:"<input type='checkbox' class='tr_select'>"
                },
                {
                    title:'集群ID',data:null,field:"SetID",
                    render: function ( data,type,row ){
                        return '<div data-SetID="'+data.SetID+'">'+data.SetID+'</div>';
                    }
                },
                {   title : '集群名称',data:"SetName",field:"SetName",width:'200px'},
                {   title : '环境类型', data:null,field:"EnviType",width:'52px',
                    render: function ( data, type, row ) {
                        var val = data.EnviType;
                        return '<div data-selid="'+val+'">'+topThis.environmentType[val]+'</div>';
                    }
                },
                {   title : '服务状态',data:null, field:"ServiceStatus",width:'52px',
                    render: function ( data, type, row ) {
                        var val = data.ServiceStatus;
                        return '<div data-selid="'+val+'">'+topThis.serverStatus[val]+'</div>';
                    }
                },
                {   title : '设计容量',data:null, field:"Capacity",width:'52px',
                    render: function ( data, type, row ) {
                        if(0 == data.Capacity) {
                            return '';
                        } else {
                            return data.Capacity;
                        }
                    }
                },
                {title : 'Openstatus', data:"Openstatus",field:"Openstatus",width:'50px'},
                {title : '包含模块数', data:"ModuleNum"},
                {
                    title:'操作',
                    defaultContent: '<a name="edits_set" class="btn btn-success  btn-sm ml10">修改</a>'+
                                    '<a name="deletes_set" class="btn btn-danger btn-sm ml10">删除</a>'+
                                    '<a name="cancels_set" class="btn btn-default btn-sm ml10 none">取消</a>'+
                                    '<a name="saves_set" class="btn btn-primary btn-sm ml10  none">保存</a>'
                }]
        },
        moduleTableColumns:function(){
           var topThis = this;
           return [{
                    title:'<input type="checkbox" class="checked_all">',
                    width:'15px',
                    defaultContent:"<input type='checkbox' class='tr_select'>"
                },
                {
                    title:'模块ID',data:null,field:"ModuleID",
                    render: function ( data,type,row ){
                        return '<div data-ModuleID="'+data.ModuleID+'">'+data.ModuleID+'</div>';
                    }
                },
                {   title : '大区ID', data:null,field:"SetID",
                    render: function ( data, type, row ) {
                        return '<div data-SetID="'+data.SetID+'">'+data.SetID+'</div>';
                    }
                },
                {   title : '模块名称',data:"ModuleName",field:"ModuleName",inputType:'text' },
                {   title : '所属集群',data:"SetName", field:"SetName"},
                {   title : '维护人',data:null, field:"Operator",
                    render: function ( data, type, row ) {
                        var val = data.Operator;
                        var text = userkv[val];
                        return '<div data-selid="'+val+'">'+text+'</div>';
                    }
                },
                {   title : '备份维护人',data:null, field:"BakOperator",
                    render: function ( data, type, row ) {
                        var val = data.BakOperator;
                        var text = userkv[val];
                        if(text === 'undefined' || typeof(text) == 'undefined' ) {
                            text = val;
                        }
                        return '<div data-selid="'+val+'">'+text+'</div>';
                    }
                },
                {title : '主机数', data:"HostCount",field:"HostCount"},
                {
                    title:'操作',
                    defaultContent:
                        '<a name="edits_module" class="btn btn-success  btn-sm ml10">修改</a>'+
                        '<a name="deletes_module" class="btn btn-danger btn-sm ml10">删除</a>'+
                        '<a name="cancels_module" class="btn btn-default btn-sm ml10 none">取消</a>'+
                        '<a name="saves_module" class="btn btn-primary btn-sm ml10  none">保存</a>'
                }]
        },
        createSelect2:function(eln,type,val){
            var topThis = this;
            eln.select2({
                width:'100%',
                data:topThis[type]
            });
            if(val){
                eln.select2('val',val);
            }
            return eln;
        },
        initBatchCheckbox:function(tableEln){
            var topThis = this;
            tableEln.on('click','.selectAll',function(e){
                var isChecked = e.currentTarget.checked;
                if(isChecked){
                    tableEln.find('.tr_select').prop('checked',false);
                    tableEln.find('.tr_select').trigger('click');
                }else{
                    tableEln.find('.tr_select').prop('checked',true);
                    tableEln.find('.tr_select').trigger('click');
                }
            });
            tableEln.on('click','.tr_select',function(e){
                var isChecked = e.currentTarget.checked;
                var checked = topThis.getChecked(tableEln);
                if(checked == tableEln.find('.tr_select').length){
                    tableEln.find('.selectAll').prop('checked',true);
                }else{
                    tableEln.find('.selectAll').prop('checked',false);
                }
                if(isChecked){
                    var dataId = $(this).attr('data-id');
                    var dataType = $(this).attr('data-type');
                    if("select" == dataType){
                        $('#'+dataId).select2('enable',true);
                    }else{
                        $('#'+dataId).prop('disabled','');
                    }
                    $('[data-for="'+dataId+'"]').addClass('none');
                }else{
                    var dataId = $(this).attr('data-id');
                    var dataType = $(this).attr('data-type');
                    if("select" == dataType){
                        $('#'+dataId).select2('disable',false);
                    }else{
                        $('#'+dataId).prop('disabled','disabled');
                    }
                    $('[data-for="'+dataId+'"]').removeClass('none');
                }
            });

            tableEln.on('click','.edit-area-mask',function(){
                var dataFor = $(this).attr('data-for');
                $('[data-id="'+dataFor+'"]').trigger('click');
                $(this).addClass('none');
            });

        },
        initCheckBoxEvent:function(tableEln){
            console.log(tableEln)
            var topThis = this;
            //tableEln.find('[role="columnheader"]:eq(0)').html('<input type="checkbox" class="checked_all" />');
            tableEln.on('click','.checked_all',function(e){
                var isChecked = e.currentTarget.checked;
                tableEln.find('.tr_select').prop('checked',isChecked);
                checkedBtnDisabled();
            });
            tableEln.on('click','.tr_select',function(e){
                var checked = topThis.getChecked(tableEln);
                if(checked == tableEln.find('.tr_select').length){
                    tableEln.find('.checked_all').prop('checked',true);
                }else{
                    tableEln.find('.checked_all').prop('checked',false);
                }
                checkedBtnDisabled();
            })

            function checkedBtnDisabled(){
                var checked = topThis.getChecked(tableEln);
                if(1 == checked){
                    tableEln.closest('.tab-pane').find('.clone').removeAttr('disabled');
                }else{
                    tableEln.closest('.tab-pane').find('.clone').attr('disabled',"");
                }
                if( 1 <= checked){
                    $(".b_edit,.clone").css("color","#555");
                    tableEln.closest('.tab-pane').find('.b_edit').removeAttr('disabled');
                }else{
                    $(".b_edit,.clone").css("color","#a1a1a1");
                    tableEln.closest('.tab-pane').find('.b_edit').attr('disabled',"");
                }
            }
        },
        getChecked:function(tableEln){  //获取当前页已经勾选的checkbox数量;
            var checked =0;
            var cluster_select_s = tableEln.find('.tr_select');
            cluster_select_s.each(function(i,v){
                if($(v).is(':checked') == true) {
                    checked ++;
                }
            });
            return checked;
        },
        initToolbarEvent:function() {    //工具栏上的初始化和事件绑定
            var topThis = this;
            topThis.initCloneModal();   //初始化克隆modal
            topThis.initClusterBatchEditEvent();    //初始化集群批量修改
            topThis.initModalBatchEditEvent();      //初始化模块批量修改

            $('#cluster').on('click','.clone',function(){
                if(1 == topThis.getChecked($('#cluster_table'))){
                    var SetName = $('#cluster_table').find('.tr_select').parent().parent().children().eq(2).text();
                    var SetID = $('#cluster_table').find('.tr_select').parent().parent().children().eq(1).text();
                    topThis.CloneSetID = SetID;
                    var title = '克隆集群['+SetName+']';
                    $("#myModalLabel").text(title);
                    $('#cloneModal').modal();
                }
            })
            .on('click','.b_edit',function(){
                if(topThis.getChecked($('#cluster_table')) >= 1){
                    console.log(topThis.getChecked($('#cluster_table')))
                    $('.cluster-edit-module').addClass('none');
                    $('.cluster-edit-cluster').removeClass('none');
                    //切换为集群编辑
                    $('.btn-save').removeClass('Module').addClass('Set');
                    $('#b_edit_modal').modal();

                }

            })

            $('#module').on('click','.b_edit',function(){
                if(topThis.getChecked($('#module_table')) >= 1){
                    $('.cluster-edit-module').removeClass('none');
                    $('.cluster-edit-cluster').addClass('none');
                    //切换为集群编辑
                    $('.btn-save').removeClass('Set').addClass('Module');
                    $('#b_edit_modal').modal();
                }
            });

            var cluster_grid = $('#cluster_table').data('kendoGrid');
            $('#filter-cluster').on('keyup',function(){
                topThis.filterGrid($(this),cluster_grid);
            });

            var module_grid = $('#module_table').data('kendoGrid');
            $('#filter-module').on('keyup',function(){
               topThis.filterGrid($(this),module_grid);
            });

            $('body').on('keyup','.regNumber',function(){
                var val = $(this).val().trim();
                if(val.substr(-1) == '.' && val.split('.').length >2){
                    $(this).val(val.substring(0,val.length-1));
                }else{
                    $(this).val(val.replace(/[^0-9|.]/gi,''));
                }
            });
        },
        filterGrid:function(eln,grid){
            var keyword = eln.val();
            var filter = {logic: "or", filters: []};
            if (keyword) {
                console.log(keyword);
                if(keyword=='体验') {
                    keyword=1;
                    console.log(keyword);
                }
                $.each(grid.columns, function (key, column) {
                    if (column.filterable) {
                        filter.filters.push({field: column.field, operator: "contains", value: keyword});
                    }
                });
            }
            grid.dataSource.options.serverFiltering = false;
            grid.dataSource.filter(filter);
        },
        initOperateEvent:function(tableEln,columns){
            var trData = [];
            var topThis = this;
            tableEln.on('click','[name="edits_set"]',function(){
                var tr = $(this).parent().parent().find('td');
                $(this).parent().find('[name="edits_set"]').addClass('none');
                $(this).parent().find('[name="deletes_set"]').addClass('none');
                $(this).parent().find('[name="saves_set"]').removeClass('none');
                $(this).parent().find('[name="cancels_set"]').removeClass('none');
                for(var i=0;i<tr.length;i++){
                    var v = tr[i];
                    var tdInput = {};
                    if(columns[i].field == 'EnviType' || columns[i].field== 'ServiceStatus'){
                        var indexval = $(v).find('div').attr('data-selid');
                        var indexhtml = $(v).find('div').text();
                        tdInput = $('<input type="hidden" indeval="'+indexval+'" indexhtml="'+indexhtml+'" class="'+columns[i].field+'">');
                        trData[i] = $(v).find('div').attr('data-selid');
                        $(v).find('div').html(tdInput);
                        topThis.createSelect2(tdInput,columns[i].field,$(v).find('div').attr('data-selid'));
                    }else if(columns[i].field == 'SetName' || columns[i].field == 'Capacity'){
                        var className = columns[i].field;
                        if('Capacity' == columns[i].field) {
                            className="regNumber"
                        }
                        tdInput = $('<input type="text" style="width:100%" maxlength="10" class="form-control cluster-edit-input  '+className+' '+columns[i].field+'">');
                        trData[i] = $(v).html();
                        tdInput.val($(v).html());
                        $(v).html(tdInput);
                    }else if(columns[i].field == 'Openstatus'){
                        tdInput = $('<input type="text" style="width:100%" maxlength="16" class="form-control cluster-edit-input '+columns[i].field+'">');
                        trData[i] = $(v).html();
                        tdInput.val($(v).html());
                        $(v).html(tdInput);
                    }
                }
            }).on('click','[name="edits_module"]',function(){
                var tr = $(this).parent().parent().find('td');
                $(this).parent().find('[name="edits_module"]').addClass('none');
                $(this).parent().find('[name="deletes_module"]').addClass('none');
                $(this).parent().find('[name="saves_module"]').removeClass('none');
                $(this).parent().find('[name="cancels_module"]').removeClass('none');
                for(var i=0;i<tr.length;i++){
                    var v = tr[i];
                    var tdInput = {};
                    if(columns[i].field == 'Operator' || columns[i].field== 'BakOperator'){
                        var indexval = $(v).find('div').attr('data-selid');
                        var indexhtml = $(v).find('div').text();
                        tdInput = $('<input type="hidden" indeval="'+indexval+'" indexhtml="'+indexhtml+'" class="'+columns[i].field+'">');
                        trData[i] = $(v).find('div').attr('data-selid');
                        $(v).find('div').html(tdInput);
                        topThis.createSelect2(tdInput,columns[i].field,$(v).find('div').attr('data-selid'));
                    }else if(columns[i].field == 'ModuleName'){
                        var className = "";
                        if('capacity' == columns[i].field) {
                            className="regNumber"
                        }
                        tdInput = $('<input type="text" maxlength="10" class="form-control cluster-edit-input  '+className+' '+columns[i].field+'">');
                        trData[i] = $(v).html();
                        tdInput.val($(v).html());
                        $(v).html(tdInput);
                    }
                }
            }).on('click','[name="saves_set"]',function(){
                var tr = $(this).parent().parent();
                var ApplicationID = cookie.get('defaultAppId');
                var SetID = tr.children().eq(1).text();
                var SetName = tr.children().children('.cluster-edit-input.SetName').val();
                var ServiceStatus = tr.find('.ServiceStatus').select2('val');
                var EnviType = tr.find('.EnviType').select2('val');
                var Capacity = tr.children().children('.cluster-edit-input.Capacity').val();
                var Openstatus = tr.children().children('.cluster-edit-input.Openstatus').val();
                if(SetName.length==0 ){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群名不能为空</span>'
                    });
                    diaCopyMsg.show(tr.children().children('.cluster-edit-input.SetName').get(0));
                    return ;
                }
                if(SetName.length > 64){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群名过长</span>'
                    });
                    diaCopyMsg.show(tr.children().children('.cluster-edit-input.SetName').get(0));
                    return ;
                }
                if(Capacity >9990000) {
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">请输入合理容量</span>'
                    });
                    diaCopyMsg.show(tr.children().children('.cluster-edit-input.Capacity').get(0));
                    return ;
                }
                $.post("/Set/editSet",
                    {ApplicationID:ApplicationID,SetName:SetName,SetEnviType:EnviType,ServiceStatus:ServiceStatus,Capacity:Capacity,SetID:SetID,Openstatus:Openstatus}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            showWindows('修改集群成功！', 'success');
                            window.location.reload();
                            return;
                        }
                    });
            }).on('click','[name="saves_module"]',function(){
                var newData = [];
                var tr = $(this).parent().parent();
                var tds = tr.find('td');
                var ApplicationID = cookie.get('defaultAppId');;
                var SetID = tr.children().eq(2).text();
                var ModuleID = tr.children().eq(1).text();
                var ModuleName = tr.children().eq(3).children('.cluster-edit-input.ModuleName').val();
                var Operator = tr.find('.Operator').select2('val');
                var BakOperator = tr.find('.BakOperator').select2('val');
                if(ModuleName.length==0 ){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">模块名不能为空</span>'
                    });
                    diaCopyMsg.show(tr.children().eq(3).children('.cluster-edit-input.ModuleName').get(0));
                    return ;
                }
                if(ModuleName.length>60){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">模块名过长</span>'
                    });
                    diaCopyMsg.show(tr.children().eq(3).children('.cluster-edit-input.ModuleName').get(0));
                    return ;
                }
                $.post("/Module/editModule",
                    {ApplicationID:ApplicationID,ModuleID:ModuleID,ModuleName:ModuleName,SetID:SetID,Operator:Operator,BakOperator:BakOperator}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            showWindows('更新模块成功！', 'success');
                            cookie.set('level',2,5);
                            window.location.reload();
                            return;
                        }
                    });
            }).on('click','[name="cancels_set"]',function(){
                $(this).parent().find('[name="edits_set"]').removeClass('none');
                $(this).parent().find('[name="deletes_set"]').removeClass('none');
                $(this).parent().find('[name="saves_set"]').addClass('none');
                $(this).parent().find('[name="cancels_set"]').addClass('none');
                var tr = $(this).parent().parent().find('td');
                for(var i=0;i<tr.length;i++){
                    var v = tr[i];
                    if(columns[i].field == 'EnviType' || columns[i].field== 'ServiceStatus'){
                        var index = $(v).find('div').find('[type="hidden"]').attr('indexval');
                        var html = $(v).find('div').find('[type="hidden"]').attr('indexhtml');
                        $(v).find('div').html(html);
                        $(v).find('div').attr('data-selid',index);
                    }else if(columns[i].field == 'SetName' || columns[i].field== 'Capacity' || columns[i].field== 'Openstatus'){
                        $(v).html(trData[i]);
                    }
                }
            }).on('click','[name="cancels_module"]',function(){
                $(this).parent().find('[name="edits_module"]').removeClass('none');
                $(this).parent().find('[name="deletes_module"]').removeClass('none');
                $(this).parent().find('[name="saves_module"]').addClass('none');
                $(this).parent().find('[name="cancels_module"]').addClass('none');
                var tr = $(this).parent().parent().find('td');
                for(var i=0;i<tr.length;i++){
                    var v = tr[i];
                    if(columns[i].field == 'Operator' || columns[i].field== 'BakOperator'){
                        var index = $(v).find('div').find('[type="hidden"]').attr('indexval');
                        var html = $(v).find('div').find('[type="hidden"]').attr('indexhtml');
                        $(v).find('div').html(html);
                        $(v).find('div').attr('data-selid',index);
                    }else if(columns[i].field == 'ModuleName'){
                        $(v).html(trData[i]);
                    }
                }
            }).on('click','[name="deletes_set"]',function(){
                var tr = $(this).parent().parent().find('td');
                var SetID = tr.children().eq(1).text();
                var gridBatDel = dialog({
                    title:'确认',
                    width:250,
                    content: '是否确认删除集群？',
                    okValue:"确定",
                    cancelValue:"取消",
                    ok:function () {
                        var ApplicationID = cookie.get('defaultAppId');
                        $.post("/Set/delSet",
                            {ApplicationID:ApplicationID,SetID:SetID}
                            ,function(result) {
                                rere = $.parseJSON(result);
                                if (rere.success == false) {
                                    showWindows(rere.errInfo, 'notice');
                                    return;
                                }
                                else {
                                    showWindows('删除集群成功！', 'success');
                                    window.location.reload();
                                }
                            });
                    },
                    cancel: function () {
                    }
                });
                gridBatDel.showModal();
            }).on('click','[name="deletes_module"]',function(){
                var tr = $(this).parent().parent().find('td');
                var ModuleID = tr.children().eq(1).text();
                var SetID = tr.children().eq(2).text();
                var gridBatDel = dialog({
                    title:'确认',
                    width:250,
                    content: '是否删除选中模块',
                    okValue:"确定",
                    cancelValue:"取消",
                    ok:function (){
                        var appId = cookie.get('defaultAppId');
                        $.post("/Module/delModule",
                            {ApplicationID:appId ,ModuleID:ModuleID, SetID:SetID}
                            ,function(result) {
                                re = $.parseJSON(result);
                                if (re.success == false) {
                                    showWindows(re.errInfo, 'notice');
                                    return;
                                }
                                else {
                                    showWindows('删除模块成功！', 'success');
                                    setTimeout(function(){cookie.set('level',2,5);window.location.reload();},1000);
                                }
                            });
                    },
                    cancel: function () {
                    }
                });
                gridBatDel.showModal();
            });
        },
        initCloneModal:function(){       //克隆modal初始化和事件绑定
            var topThis = this;
            $('#cloneModal').on('hidden.bs.modal', function (e) {
                $('#cloneTextarea').val(null);
            }).on('click','.btn-save',function(){
                var SetID = topThis.CloneSetID;
                var SetName = $("#cloneTextarea").val();
                if($.trim(SetName)==''){
                    $("#cloneerrtips").text('集群名不能为空');
                    return;
                }
                var appId = cookie.get('defaultAppId');
                $.post("/set/cloneset",
                    {ApplicationID:appId, SetID:SetID, SetName:SetName}
                    ,function(result) {
                        re = $.parseJSON(result);
                        if (re.success == false) {
                            $('#cloneModal').modal('hide');
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            $('#cloneModal').modal('hide');
                            showWindows('克隆集群成功！', 'success');
                            window.location.reload();
                            return;
                        }
                    });
            });
        },
        initClusterBatchEditEvent:function(){
            var topThis = this, eln=$('#b_edit_modal');
            eln.on('hidden.bs.modal', function (e){
                $('#e_EnviType').select2('disable',true);
                $('#e_ServiceStatus').select2('disable',true);
                $('#e_Capacity').prop('disabled','disabled');
                $('#e_Openstatus').prop('disabled','disabled');
                $('#c_EnviType').prop('checked',false);
                $('#c_ServiceStatus').prop('checked',false);
                $('#c_Capacity').prop('checked',false);
                $('.selectAll').prop('checked',false);
            }).on('show.bs.modal', function (e){
                $('#e_EnviType').select2('disable',true);
                $('#e_ServiceStatus').select2('disable',true);
                $('#e_Capacity').prop('disabled','disabled');
                $('#e_Openstatus').prop('disabled','disabled');
                $('.edit-area-mask').removeClass('none');
            }).on('click','.btn-save.Set',function(){
                var ApplicationID = cookie.get('defaultAppId');
                var SetID = [];
                $("#cluster_table tr td:first-child input:checked").parent().parent().each(function() {
                    SetID.push($(this).children().eq(1).text());
                })
                var Capacity = $("#e_Capacity").val();
                var Openstatus = $("#e_Openstatus").val();
                var EnviType = $("#e_EnviType").select2('val');
                var ServiceStatus = $("#e_ServiceStatus").select2('val');

                if($("#c_ServiceStatus").prop("checked") && ServiceStatus == "") {
                    $("#errtips").text('服务状态不能为空');
                    return ;
                }

                if($("#c_EnviType").prop("checked") && EnviType == "") {
                    $("#errtips").text('环境类型不能为空');
                    return ;
                }

                if($("#c_Capacity").prop("checked") && Capacity == "") {
                    $("#errtips").text('容量不能为空');
                    return ;
                }

                if($("#c_Openstatus").prop("checked") && Openstatus == "") {
                    $("#errtips").text('Openstatus不能为空');
                    return ;
                }

                if(Capacity >9990000) {
                    $("#errtips").text('输入合理容量');
                    return ;
                }
                $("#errtips").text('');
                $.post("/Set/editSet",
                    {ApplicationID:ApplicationID,SetEnviType:EnviType,ServiceStatus:ServiceStatus,Capacity:Capacity,SetID:SetID,Openstatus:Openstatus}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        } else {
                            showWindows('修改集群成功！', 'success');
                            window.location.reload();
                            return;
                        }
                    });
            });

            topThis.createSelect2($('#e_ServiceStatus'),'ServiceStatus').select2('disable',true);
            topThis.createSelect2($('#e_EnviType'),'EnviType').select2('disable',true);
        },
        initModalBatchEditEvent:function(){
            var topThis = this, eln=$('#b_edit_modal');
            eln.on('hidden.bs.modal', function (e){
                $('#e_Operator').prop('disabled','disabled');
                $('#e_BakOperator').prop('disabled','disabled');
                $('#c_BakOperator').prop('checked',false);
                $('#c_Operator').prop('checked',false);
                $('.selectAll').prop('checked',false);
            }).on('show.bs.modal', function (e){
                if(!$("#c_Operator").prop("checked")) {
                    $('#e_Operator').select2('disable',true);
                }
                if(!$("#c_BakOperator").prop("checked")) {
                    $('#e_BakOperator').select2('disable',true);
                }
                $('.edit-area-mask').removeClass('none');
            }).on('click','.btn-save.Module',function(){
                var ApplicationID = cookie.get('defaultAppId');
                var ModuleID = [];
                $("#module_table tr td:first-child input:checked").parent().parent().each(function(){
                // $("#module_table").find("table tr td:first-child input:checked").parent().parent().each(function(){
                    ModuleID.push($(this).children().eq(1).text());
                })
                var Operator = $("#e_Operator").select2('val');
                var BakOperator = $("#e_BakOperator").select2('val');
                if($("#c_Operator").prop("checked") && Operator == "") {
                    $("#errtips").text('维护人不能为空');
                    return ;
                }
                if($("#c_BakOperator").prop("checked") && BakOperator== "") {
                    $("#errtips").text('备份维护人不能为空');
                    return ;
                }
                $("#errtips").text();
                $.post("/Module/editModule",
                    {ApplicationID:ApplicationID,ModuleID:ModuleID,ModuleName:'',SetID:'',Operator:Operator,BakOperator:BakOperator}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            showWindows('更新模块成功！', 'success');
                            setTimeout(function(){cookie.set('level',2,5);window.location.reload();},1000);
                            return;
                        }
                    });
                eln.modal('hide');
            });

            topThis.createSelect2($('#e_Operator'),'Operator').select2('disable',true);
            topThis.createSelect2($('#e_BakOperator'),'BakOperator').select2('disable',true);
        },
        initClusterTable:function(tableEln){
            var topThis = this;
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
            }
            var params={'ApplicationID':cookie.get('defaultAppId')};
            tableEln.dataTable({
                paging: true, //隐藏分页
                ordering: false, //关闭排序
                //info: false, //隐藏左下角分页信息
                //searching: false, //关闭搜索
                lengthChange: false, //不允许用户改变表格每页显示的记录数
                language: language, //汉化
                autoWidth: false,
                destroy: true,
                "ajax": {
                    "url": "/topology/setlist",
                    type:"post",
                    "data": function ( d ) {
                       return d = $.extend({},d, params);
                    }
                },
                columns: topThis.clusterTableColumns()
            });
                topThis.initCheckBoxEvent(tableEln);
                topThis.initOperateEvent(tableEln,topThis.clusterTableColumns());
        },
        initModuleTable:function(tableEln){
            var topThis = this;
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
            }
            var params={'ApplicationID':cookie.get('defaultAppId')};
            tableEln.dataTable({
                paging: true, //隐藏分页
                ordering: false, //关闭排序
                //info: false, //隐藏左下角分页信息
                //searching: false, //关闭搜索
                lengthChange: false, //不允许用户改变表格每页显示的记录数
                language: language, //汉化
                autoWidth: false,
                destroy: true,
                "ajax": {
                    "url": "/topology/modulelist",
                    type:"post",
                    "data": function ( d ) {
                       return d = $.extend({},d, params);
                    }
                },
                columns: topThis.moduleTableColumns()
            });
            topThis.initCheckBoxEvent(tableEln);
            topThis.initOperateEvent(tableEln,topThis.moduleTableColumns());
        }
    }

    clusterList.init(); //初始化
    level = cookie.get('level');
    if(level==2 && alevel==3)
    {
        $(".nav.nav-tabs li").eq(0).removeClass('active');
        $(".nav.nav-tabs li").eq(1).addClass('active');
        $("#cluster").removeClass('in').removeClass('active');
        $("#module").addClass('in').addClass('active');
    }

});
