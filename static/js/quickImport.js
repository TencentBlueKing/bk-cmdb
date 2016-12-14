/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

$(document).ready(function() {
    var params={ApplicationID:cookie.get('defaultAppId'),IsDistributed:false,Source:"3"};
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
    $("#private").dataTable({
        paging: true, //隐藏分页
        ordering: false, //关闭排序
        //info: false, //隐藏左下角分页信息
        //searching: false, //关闭搜索
        lengthChange: false, //不允许用户改变表格每页显示的记录数
        language: language, //汉化
        autoWidth: false,
        destroy: true,
        ajax: {
            url: "/host/getHost4QuickImport",
            type:"post",
            data: function ( d ) {
               return d = $.extend({},d, params);
            }
        },
        initComplete: function(){
            showSwitchNum();
        },
        columns : [
            {data:null,title:'<input type="checkbox" class="selectAll">',
             width:"15px",defaultContent:"<input type='checkbox' class='selectTr c-grid-checkbox'>"},
            {data:'InnerIP',title:"内网IP",width:"120px"},
            {data:"OuterIP", title: "外网IP",width:"120px"},
            {data:"AssetID",title:'固资编号',width:"120px"},
            {data:"ApplicationName",title:'所属业务',width:"120px"},
            {data:"SetName",title:'所属集群',width:"120px"},
            {data:"ModuleName",title:'所属模块',width:"120px"},
            {data:"HostName",title:'主机名称',width:"120px"},
            {data:"OSName",title:'操作系统',width:"120px"}
        ]
    });
    function showSwitchNum(){
        var total=$('#private').DataTable().data().length;
        if($('.cc_switch_btn').attr('data-fp')=='1'){
            $('.cc_switch_btn').removeClass('cc_switch_btn_left').addClass('cc_switch_btn_right');
            var text = '未分配('+total+')';
            $('.switch').siblings('.num').text(text);
        }else if($('.cc_switch_btn').attr('data-fp')=='0'){
            $('.cc_switch_btn').removeClass('cc_switch_btn_right').addClass('cc_switch_btn_left');
            var text = '已分配('+total+')';
            $('.switch').siblings(".num").text(text);
        }
    }
    $('.cc_switch_btn .switch').click(function(){
        if($('.cc_switch_btn').attr('data-fp')=='1'){
            $('.cc_switch_btn').attr('data-fp',0);
            //$('.cc_switch_btn').removeClass('cc_switch_btn_right').addClass('cc_switch_btn_left');
            var state = true;
            $('.k-grid-delete,.k-grid-quickDistribute', '#private').attr('disabled', true);
            $('.k-grid-delete', '#private').hide();
            btnHide();
        }else if($('.cc_switch_btn').attr('data-fp')=='0'){
            $('.cc_switch_btn').attr('data-fp',1);
            //$('.cc_switch_btn').removeClass('cc_switch_btn_left').addClass('cc_switch_btn_right');
            $('.k-grid-delete', '#private').show();
            var state = false;
            btnHide();
        }
        var grid = $('#private').dataTable();
        grid.find('thead input').attr('checked', false);
        var table = $('#private').DataTable();
        params={ApplicationID:cookie.get('defaultAppId'),IsDistributed:state,Source:"3"};
        table.ajax.reload();
        setTimeout(function () {
            showSwitchNum();
        }, 500);
    });

    (function(){

        /*标签点击事件*/
        $('.nav-tabs').on('click', function(e){
            if($(e.target).attr('href')){
                cookie.set('quick_destribute_current_tab', $(e.target).attr('href').replace('#',''));
            }
        });
    })();

    /**
    * 表格工具栏的切换按钮
    * size：按钮大小
    * labelWidth：切换按钮的label宽度
    * onText：切换按钮on状态的文字
    * offText：切换按钮off状态的文字
    * onSwitchChange：状态切换时间处理函数
    */
    $(".host-state-switcher").bootstrapSwitch({
        size: 'small',
        labelWidth:'60px',
        onText: '已分配',
        offText: '未分配',
        onSwitchChange: function(e, state){
            var id = 'private';
            var grid = $('#'+id).data('kendoGrid');
            grid.thead.find('input').attr('checked', false);
            grid.destroy();
            var gridObj = eval(id+'KendoGridObj');
            gridObj.dataSource.transport.read.data.IsDistributed = state;
            $('#'+id).kendoGrid(gridObj);
            grid = $('#'+id).data('kendoGrid');
            grid.refresh();
            $('#filter-'+id).data('data',{});
            if(state == true) {
                $('.k-grid-quickDistribute', '#'+id).attr('title', '配置平台禁止跨业务分配主机，如果实在要用，请联系原主机业务的运维同学上交后再分配');
            }
            else{
                $('.k-grid-quickDistribute', '#'+id).removeAttr('title');
            }

            $('.k-grid-delete,.k-grid-quickDistribute', '#'+id).attr('disabled', true);
            if(state){
                $('.k-grid-delete', '#'+id).hide();
            }else{
                $('.k-grid-delete', '#'+id).show();
            }

            $('#filter-'+id).val('');
        }
    });

    /**
    * 表头搜索框输入时间处理函数
    * 支持所有字段的搜索，字段之间逻辑关系为or，搜索方式为contains，即包含
    */
    $('#filter-private').on('keyup', function(e){
        var type = $(e.target).attr('id').split('-').pop();
        var grid = $('#'+type).data('kendoGrid');
        if(typeof JSON=='undefined'){
            $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
        }

        if($.isEmptyObject($(e.target).data('data'))){
            var data = grid.dataSource.data();
            $(e.target).data('data', data);
        }else{
            var data = $(e.target).data('data');
        }
        var d = JSON.parse(JSON.stringify(data));
        for(var i in d){
            d[i].Checked = '';
        }
        grid.dataSource.data(d);
        grid.refresh();
        grid.thead.find('input[type=checkbox]').prop('checked', false);

        filter = {logic: "or", filters: []};
        $searchValue = $(e.target).val();
        if ($searchValue) {
            $.each(grid.columns, function (key, column) {
                if (column.filterable) {
                    filter.filters.push({field: column.field, operator: "contains", value: $searchValue});
                }
            });
        }

        grid.dataSource.options.serverFiltering = false;
        grid.dataSource.filter(filter);
        grid.selectNum = 0;
        var query = new kendo.data.Query(grid.dataSource.data());

        if($('.cc_switch_btn').attr('data-fp')=='1'){
            var text = '未分配('+query.filter(filter).data.length+')';
            $('.switch').siblings('.num').text(text);
        }else if($('.cc_switch_btn').attr('data-fp')=='0'){
            $('.cc_switch_btn').removeClass('cc_switch_btn_right').addClass('cc_switch_btn_left');
            var text = '已分配('+query.filter(filter).data.length+')';
            $('.switch').siblings(".num").text(text);
        }
    });

    /**
    * 表头工具栏“分配至”点击事件处理函数
    */
    $('.c-grid-quickDistribute').on('click', function(e){
        if($(e.target).attr('disabled')=='disabled'){
            return false;
        }
        var grid = $('#private').dataTable();
        var data = $('#private').DataTable().data();
        var d = data;
        var appId = [];
        var hostId = [];
        var IsDistributed = false
        for(var i=0,len=d.length; i<len; i++){
            if(d[i].Checked==='checked'){
                if($.inArray(d[i].ApplicationID, appId)==-1){
                    appId.push(d[i].ApplicationID);
                }

                hostId.push(d[i].HostID);

                if(d[i].ApplicationName.indexOf('资源池')>-1){
                    IsDistributed = true;
                }
            }
        }

        if(hostId.length===0){
            var noHostSelectDialog = dialog({
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            });
            noHostSelectDialog.show();
            setTimeout(function () {
                noHostSelectDialog.close().remove();
            }, 2000);
            return false;
        }

        var param = {};
        param['HostID'] = hostId.join(',');
        /*重构完成后，需改动*/
        param['ApplicationID'] = appId.join(',');//标签上放appId
        param['ToApplicationID'] = $('#appId').attr('ApplicationID') ? $('#appId').attr('ApplicationID') : cookie.get('defaultAppId');//标签上放appId

        var options = {
            title:'确认',
            width:300,
            content: '当前选择的主机已经被分配到其他业务使用，确认继续分配至<i class="redFont">'+ cookie.get('defaultAppName')+'</i>',
            okValue:"继续",
            cancelValue:"我再想想",
            ok:function (){
                var d = dialog({
                    content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>正在分配...</div>'
                });
                d.showModal();

                $.ajax({
                    dialog:d,
                    url:'/host/quickDistribute',
                    data:param,
                    method:'post',
                    dataType:'json',
                    success:function(response){
                        this.dialog.close().remove();
                        if(response.success && cookie.get('play_user_guide_quickimport')==1){
                            cookie.set('quickimport_user_guide_step', 5);
                            var distributeHostDialog = dialog({
                                title:'提示',
                                width:300,
                                height:50,
                                content: '<div class="c-dialogdiv"><i class="c-dialogimg-success"></i>恭喜，您离成功越来越近，<a href="/topology/index">点此</a>将您的业务拓扑绘制到配置平台</div>',
                                okValue:"确定",
                                ok:function (){
                                    location.href = '/topology/index';
                                }
                            });

                            distributeHostDialog.showModal();
                            return true;
                        }else{
                            var content = response.success==true ? '<i class="c-dialogimg-success"></i>'+response.message : '<i class="c-dialogimg-prompt"></i>'+response.errInfo;
                            var d = dialog({
                                    content: '<div class="c-dialogdiv2">'+content+'</div>'
                                });
                            d.showModal();
                            setTimeout(function() {
                                d.close().remove();
                                window.location.reload();
                            }, 2500);
                        }
                        return true;
                    }
                });
            },
            cancel:function(){}
        };

        if(IsDistributed){
            options.content = '当前操作会将已勾选的<i class="redFont">'+ hostId.length +'</i>台主机分配至<i class="redFont">'+ cookie.get('defaultAppName') +'</i>的空闲机池，确认继续？';
        }

        var quickDistributeDialog = dialog(options);
        quickDistributeDialog.showModal();
        $('.redFont').tooltip();
        return true;
    });
    function getChecked(){
        var checked =0;
        var cluster_select_s = $('#private').find('.selectTr');
        cluster_select_s.each(function(i,v){
            if($(v).is(':checked') == true) {
                checked ++;
            }
        });
        return checked;
    };
    var quickDistribute = $('.c-grid-quickDistribute');
    var deleteBtn=$('.c-grid-delete');
    function btnShow(){
        var flag = $('.cc_switch_btn').attr('data-fp');
        if(flag > 0) {
            quickDistribute.attr('disabled', false);
            deleteBtn.attr('disabled', false);
            $(".c-grid-quickDistribute,.c-grid-delete").css('color', "#555");
        }
    };
    function btnHide(){
        quickDistribute.attr('disabled', true);
        deleteBtn.attr('disabled', true);
        $(".c-grid-quickDistribute,.c-grid-delete").css('color', "#a1a1a1");
    }
    $('#private').on('change','input[type=checkbox]',function(e){
        if($(e.target).attr('class')==='selectAll'){
            var isChecked = e.target.checked;
            var data=$('#private').DataTable().data();
            if(isChecked){
                $('#private').find('.selectTr').prop('checked',true);
                btnShow();
                $.each(data,function(i,m){
                    m.Checked = 'checked';
                });
            }else{
                $('#private').find('.selectTr').prop('checked',false);
                btnHide();
                $.each(data,function(i,m){
                    m.Checked = '';
                });
            }
        }else{
            //被勾选个数
            var checked = getChecked();
            //勾选框总个数
            var trLength=$('#private').find('.selectTr').length;
            //当前是否勾选
            var isChecked = e.target.checked;
           //当前change行的IP及数据
            var tdIp = $(e.target).closest('tr').find('td').eq(1).text();
            var trdata=getData(tdIp);
            if(checked ==0 ){
                $(".c-grid-quickDistribute").attr('disabled', true);
                btnHide();
            }else if(checked==trLength){
                $('#private').find('.selectAll').prop('checked',true);
                btnShow();
            }else{
                btnShow();
                $('#private').find('.selectAll').prop('checked',false);
            }
            if(isChecked){
                trdata.Checked = 'checked';
            }else{
               trdata.Checked = '';
            }
        }
    })

    function getData(tdIp){
        var trdata='';
        var data=$('#private').DataTable().data();
        $.each(data,function(i,m){
            if(m.InnerIP==tdIp){
                trdata=m;
            }
        })
        return trdata;
    }

    /**
    * 表头工具栏“删除”按钮点击事件处理函数
    * 将从其它云，利用excel导入的主机从配置平台彻底删除
    */
    $(".c-grid-delete").on('click', function(e){
        if($(e.target).attr('disabled')=='disabled'){
            return false;
        }

        var param = {};
        //var type = $(e.target).parents('.tab-pane').attr('id');
        // var grid = $('#'+type).data('kendoGrid');
        // var data = grid.dataSource.data();

        var grid = $('#private').dataTable();
        var data = $('#private').DataTable().data();

        // if(typeof JSON=='undefined'){
        //     $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
        // }
        //var d = JSON.parse(JSON.stringify(data));
        var d = data;

        var hostId = [];
        var appId = [];
        for(var i=0,len=d.length; i<len; i++){
            if(d[i].ApplicationName!=='资源池'){
                var notAllowToDeleteDialog = dialog({
                    content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>只能删除未分配机器</div>'
                });
                notAllowToDeleteDialog.show();
                setTimeout(function () {
                    notAllowToDeleteDialog.close().remove();
                }, 2000);
                return false;
            }

            if(d[i].Checked==='checked'){
                var tmp = d[i].ApplicationID.split(',');
                if(tmp.length==1){
                    if($.inArray(d[i].ApplicationID, appId)==-1){
                        appId.push(d[i].ApplicationID);
                    }
                }else{
                    var tmp = d[i].ApplicationID.split(',');
                    for(var j=0,jlen=tmp.length; j<jlen; j++){
                       if($.inArray(tmp[j], appId)==-1){
                            appId.push(tmp[j]);
                        }
                    }
                }

                hostId.push(d[i].HostID);
            }
        }

        if(hostId.length==0){
            var noHostSelectDialog = dialog({
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            });
            noHostSelectDialog.show();
            setTimeout(function () {
                noHostSelectDialog.close().remove();
            }, 2000);
            return false;
        }

        param['HostID'] = hostId.join(',');
        /*重构完成后，需改动*/
        param['ApplicationID'] = appId.join(',');//标签上放appId
        var options = {
            title:'确认',
            width:300,
            content: '您勾选的<i class="redFont">' + hostId.length + '</i>台主机即将离开配置平台，确认是否继续？',
            okValue:"继续",
            cancelValue:"我再想想",
            ok:function (){
                var d = dialog({
                    content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>正在删除...</div>'
                });
                d.showModal();

                $.ajax({
                    dialog:d,
                    url:'/host/delPrivateDefaultApplicationHost',
                    data:param,
                    method:'post',
                    dataType:'json',
                    success:function(response){
                        this.dialog.close().remove();
                        if(response.success){
                            var distributeHostDialog = dialog({
                                title:'提示',
                                width:300,
                                height:50,
                                content: '<div class="c-dialogdiv"><i class="c-dialogimg-success"></i>删除成功!</div>',
                                okValue:"确定",
                                ok:function (){
                                    location.href = '/host/quickImport';
                                }
                            });

                            distributeHostDialog.showModal();
                            return true;
                        }else{
                            var content = response.success==true ? '<i class="c-dialogimg-success"></i>'+response.message : '<i class="c-dialogimg-prompt"></i>'+response.message;
                            var d = dialog({
                                width: 150,
                                content: '<div class="c-dialogdiv2">'+content+'</div>'
                            });
                            d.showModal();
                            setTimeout(function() {
                                d.close().remove();
                                window.location.reload();
                            }, 2500);
                        }
                        return true;
                    }
                });
            },
            cancel:function(){}
        };

        var quickDistributeDialog = dialog(options);
        quickDistributeDialog.showModal();

        return true;
    });

    $(document.body).on('click', '#user_guide_import_private', function(e){
        d.close().remove();step1();
    });


    $('.import-page-mask').click(function(e){
        $(this).hide();
    });

    $('#importOtherHost').click(function(){
        cookie.set('quick_destribute_current_tab', 'private');
        cookie.set('quickimport_user_guide_step', 1);
        step1();
    });

    // 导入私有云机器
    $('#importPrivateHostByExcel').on('click',function (e){
        var importPrivateHost = dialog({
                title:'导入主机',
                width:530,
                content: '<div class="pt10">'+
                         '<form action="/host/getImportPrivateHostTableFieldsByExcel" id="upload_form" enctype="multipart/form-data" method="post" target="upload_proxy" style="display:inline-block;">'+
                         '<lable><span class="c-gridinputmust pr10">*</span>请选择导入文件：</lable>'+
                         '<a class="k-button king-btn-mini king-file-btn filebox">选择文件'+
                         '<input type="file" id="importPrivateHost" name="importPrivateHost">'+
                         '</a>'+
                         '<span class="import-file-name ml15"></span>'+
                         '<p style="color:#666;padding:10px 0 0 5px;"></p>'+
                         '<p class="">温馨提示：<br>1.文件类型支持xls、xlsx、csv;  <br>2.格式如下示例，其中<lable class="redFont">内网IP</lable>是必填项,其它的均非必填且可以自定义;<br><br><img src="/static/img/import.jpg"/></p>'+
                         '</form>'+
                         '</div>',
                okValue:"导入",
                cancelValue:"关闭",
                skin:'dia-grid-batDel',
                ok:function (){

                    var file = $('#importPrivateHost').val();
                    if(file){
                        $("#upload_form").submit();
                    }else{
                        var noFileSelectDialog = dialog({
                            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择文件</div>'
                        });
                        noFileSelectDialog.show();
                        setTimeout(function () {
                            noFileSelectDialog.close().remove();
                        }, 2000);
                        return false;
                    }

                    uploadDialog = dialog({
                        content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>正在导入...</div>'
                    });
                    uploadDialog.showModal();
                    setTimeout(function(){
                        //clearUpload($(e.target).attr('id'));
                    },500);
                }
            });
        importPrivateHost.showModal();
        $('#importPrivateHost').on('change', function(){
          if (!$('.import-file-name').text($('#importPrivateHost').val().split('\\')[$('#importPrivateHost').val().split('\\').length-1])) {
               $('.import-file-name').text($('#importPrivateHost').val().split('/')[$('#importPrivateHost').val().split('/').length-1])
          };
        });
    })

});


/**
* 重置上传表单，防止相同文件上传没反应
*/
function clearUpload(id){
    $("#"+id).parents('form').submit().end().remove();//移除原来的
    $("<input/>").attr("name",id).attr("id",id).attr("type","file").appendTo(".filebox");//添加新的
}

/**
* 导入主机回调函数
* 负责页面显示成功or失败的提示
*/
function uploadCallback(data){
    uploadDialog.close().remove();
    if(data.success){
        var d = dialog({
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>'+ data.errInfo +'</div>'
        });
        d.showModal();
        setTimeout(function(){
            //window.location.replace(location.href);
            window.location.reload();
        }, 2000);
    }else{
        var d = dialog({
            title:'确认',
            width:300,
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>'+ data.errInfo +'</div>',
            okValue:"确定",
            ok:function(){
                window.location.reload();
            }
        });

        d.showModal();
        $(".import-error-list").mCustomScrollbar({
            theme: "minimal-dark" //设置风格
        });
    }
}

/**
* 导入主机回调函数
* 负责页面显示成功or失败的提示
*/
function uploadCallbackToHostField(data){
    uploadDialog.close().remove();
    if(data.success){
        var titles = data.keys;
        var select2Arr = data.fields;
        var select2data = data.fields;
        var readTilte_html='';
        $.each(titles,function(i,e) {
            readTilte_html += '<div class="row import_readdialoga">'+
                            '   <div class="col-4">'+
                            '       <input type="text" class="form-control tableHeader" style="width:100%;" readOnly="true" value='+e.name+' >'+
                            '   </div>'+
                            '   <div class="col-4 user-radio">'+
                            '       <label><input type="radio" name="'+e.name+'" class="user-filter" value="select" checked="checked">映射已有字段</label>'+
                            '       <label><input type="radio" name="'+e.name+'" class="user-defined" value="customer">自定义</label>'+
                            '   </div>'+
                            '   <div class="col-4">'+
                            '       <input type="text" name="'+e.name+'_select" class="select2_box user-filter-input" style="width:135px;display:block;">'+
                            '       <input type="text"  name="'+e.name+'_customer" class="form-control user-defined-input" style="width:135px;display:none;">'+
                            '   </div>'+
                            '</div>';
        });
        var importPrivateHost = dialog({
            title:'导入主机字段映射',
            width:520,
            content: '<div class="pt10">'+
                     '<form action="/host/importPrivateHostByExcel" id="upload_form" method="post" target="upload_proxy" style="display:inline-block;">'+
                     '<input type="hidden" name="filename" value="'+ data.filename +'">'+
                     '<div class="pt10">'+readTilte_html+'</div>'+
                     '</form>'+
                     '</div>',
            okValue:"确认",
            cancelValue:"关闭",
            skin:'dia-grid-batDel',
            ok:function (){
                $("#upload_form").submit();
                uploadDialog = dialog({
                    content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>正在导入...</div>'
                });
                uploadDialog.showModal();
                setTimeout(function(){
                    //clearUpload($(e.target).attr('id'));
                },500);
            }
        });
        importPrivateHost.showModal();
        $(".select2_box").select2({ data: select2data });
        var currentVal="";
        var selectVal="";
        //下拉框 筛选时获取当前值
        $("input.select2_box").on("select2-open", function(e) {
            currentVal=$(this).select2("val");
        }).on("change", function(e) {
            selectVal=$(this).select2("val");
            var num="";
            //选中值在数组中的位置
            $.each(select2data,function(n,m) {
                if(m.id==selectVal){
                    num=n;
                }
            })
            //删除该位置的节点
            select2data.splice(num,1);
            $.each(select2Arr,function(n,m) {
                if(m.id==currentVal){
                    select2data.push(m);
                }
            })
        })
        $('.import_readdialoga .user-defined').on('click',function(){
            var customerValue = $(this).closest('.import_readdialoga').find('.tableHeader').val();
            $(this).closest('.import_readdialoga').find('.user-defined-input').show().val(customerValue).siblings().hide();
        });
        $('.import_readdialoga .user-filter').on('click',function(){
            $(this).closest('.import_readdialoga').find('div.select2_box').show().siblings().hide();
        });
    }else{
        var d = dialog({
            title:'确认',
            width:300,
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>'+ data.errInfo +'</div>',
            okValue:"确定",
            ok:function(){}
        });

        d.showModal();
        $(".import-error-list").mCustomScrollbar({
            theme: "minimal-dark" //设置风格
        });
    }
}
