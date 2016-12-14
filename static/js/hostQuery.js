/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

window.CC = window.CC || {};
/**
 * CC.host
 */
!function(){
    $('body').bind("click",function(event){
        var target = $(event.target);
        if (!target.closest('#panelContain_9899').length){
            $('#panelContain_9899').find('[data-toggle="popover"]').popover('hide');
        }
    });


    var ismouseenter=false;
    $('body').on('mouseenter','.copy-menu-box',function(event) {
        if( $('.k-grid-copyIP').attr('disabled')!='disabled' ) {
            $(this).find('.copy-menu').addClass('active').removeClass('none');
        }
    });
    $('body').on('mouseleave','.copy-menu-box',function(event) {
        if(!ismouseenter){
            cpyTime=setTimeout(function (){
                $('.copy-menu').removeClass('active');
            },1000)
        }
    }).on('mouseenter','.global-zeroclipboard-container',function(e){
        ismouseenter = true;
    }).on('mouseleave','.global-zeroclipboard-container',function(e){
        ismouseenter = false;
    })

    /**
    * 拓扑试图kendoTreeView树结构配置对象
    */
    var treeContainerChartObj = {
        dataTextField:['Name'],
        template:function(data){
          var node=data.item;
          return '<span>'+node.Name+'</span></span><span class="host-count label">'+node.number+'</span>';
        },
        loadOnDemand:false,
        select:function(e){
            if(this.emptyView){
                this.emptyView.select($());
            }else{
                // $('#emptyContainer').data('kendoTreeView').select($());
            }

            var node = this.dataItem(e.node);
            var data = {'application' : {'ApplicationID' : node.id}, 'set' : {'ApplicationID' : node.appId, 'SetID' : node.id}, 'module' : {'ApplicationID' : node.appId, 'ModuleID' : node.id}};

            if(node.type=='application'){
                var ModuleID = [];

                if(node.lvl==2){
                    var modules = typeof node['items']=='undefined' ? [] : node['items'];
                    for(var i=0,mlen=modules.length; i<mlen; i++){
                        ModuleID.push(modules[i].id);
                    }
                }else{
                    var sets = typeof node['items']=='undefined' ? [] : node['items'];
                    for(var i=0,slen=sets.length; i<slen; i++){
                        var modules = typeof sets[i]['items']=='undefined' ? [] : sets[i]['items'];
                        for(var j=0,mlen=modules.length; j<mlen; j++){
                            ModuleID.push(modules[j].id);
                        }
                    }

                }

                data.application['ModuleID'] = ModuleID.join(',');
            }


            gridObj.dataSource.transport.read.data = data[node.type];
            gridObj.dataSource.transport.read.url = '/host/getHostById';
            // CC.host.hostlist.init();
            $('#batRes').hide();
            $('#batDel').show();


            window.intval = setInterval(function(){
                $('.k-state-focused', '#treeContainer').removeClass('k-state-focused');
                if($('.k-state-focused', '#treeContainer').length==0){
                    clearInterval(intval);
                }
            }, 20);
      }
    };

    /**
    * 空闲机kendoTreeView树结构配置对象
    */
    var emptyContainerChartObj = {
        dataTextField:['Name'],
        template:function(data){
          var node=data.item;
          return '<span>'+node.Name+'</span></span><span class="host-count label">'+node.number+'</span>';
        },
        select:function(e){
            if(this.treeView){
                this.treeView.select($());
            }else{
                // $('#treeContainer').data('kendoTreeView').select($());
            }
            var node = this.dataItem(e.node);
            var data = {'application' : {'ApplicationID' : node.id}, 'set' : {'ApplicationID' : node.appId, 'SetID' : node.id}, 'module' : {'ApplicationID' : node.appId, 'ModuleID' : node.id}};

            if(node.type=='application'){
                var ModuleID = [];

                if(node.lvl==2){
                    var modules = typeof node['items']=='undefined' ? [] : node['items'];
                    for(var i=0,mlen=modules.length; i<mlen; i++){
                        ModuleID.push(modules[i].id);
                    }
                }else{
                    var sets = typeof node['items']=='undefined' ? [] : node['items'];
                    for(var i=0,slen=sets.length; i<slen; i++){
                        var modules = typeof sets[i]['items']=='undefined' ? [] : sets[i]['items'];
                        for(var j=0,mlen=modules.length; j<mlen; j++){
                            ModuleID.push(modules[j].id);
                        }
                    }
                }

                data.application['ModuleID'] = ModuleID.join(',');
            }

            gridObj.dataSource.transport.read.data = data[node.type];
            gridObj.dataSource.transport.read.url = '/host/getHostById';
            // CC.host.hostlist.init();
            $('#batRes').show();
            $('#batDel').hide();

            window.intval = setInterval(function(){
                $('.k-state-focused', '#treeContainer').removeClass('k-state-focused');
                if($('.k-state-focused', '#treeContainer').length==0){
                    clearInterval(intval);
                }
            }, 20);
        }
    };

    /**
    * 拓扑试图主机数量按钮点击事件处理函数
    */
    $("#treeContainer").on('click','.host-count',function(e){
        $(e.target).prev('span').trigger('click');
    });

    /**
    * 空闲机主机数量按钮点击事件处理函数
    */
    $("#emptyContainer").on('click','.host-count',function(e){
        $(e.target).prev('span').trigger('click');
    });

    /**
    * 查询条件按钮"我"点击事件处理函数
    */
    // $('#filter_module_mine').click(function(e){
    //     e.stopPropagation();
    //     $('#emptyContainer').data('kendoTreeView').select($());
    //     $('#treeContainer').data('kendoTreeView').select($());

    //     //现改为搜索维护人为当前用户的主机
    //     gridObj.dataSource.transport.read.data = {'ApplicationID' : cookie.get('defaultAppId'), 'Operator' : window.currUser.username};
    //     gridObj.dataSource.transport.read.url = '/host/getHostByCondition';
    //     // CC.host.hostlist.init();
    //     $('#batRes').hide();
    //     $('#batDel').show();
    // });

    /**
    * 查询条件按钮"ALL"点击事件处理函数
    */
    // $('#filter_module_all').click(function(e){
    //     e.stopPropagation();
    //     $('#emptyContainer').data('kendoTreeView').select($());
    //     $('#treeContainer').data('kendoTreeView').select($());

    //     //现改为搜索当前业务所有主机
    //     gridObj.dataSource.transport.read.data = {'ApplicationID' : cookie.get('defaultAppId')};
    //     gridObj.dataSource.transport.read.url = '/host/getHostByCondition';
    //     // CC.host.hostlist.init();
    //     $('#batRes').hide();
    //     $('#batDel').show();
    // });

    /**
    * 查询条件按钮"空闲机"点击事件处理函数
    */
    // $('#filter_module_empty').click(function(e){
    //     e.stopPropagation();
    //     $('#emptyContainer').data('kendoTreeView').select($());
    //     $('.host-count', '#emptyContainer').click();
    // });

    /**
    * 查询条件横条点击事件处理函数
    */
    $('#collapseOneBtn').click(function(){
        var isUp = $('#collapseOne').attr('class').indexOf('in') == -1;
        if(isUp){
            $(this).find('span').removeClass('fa-angle-down').addClass('fa-angle-up');
        }else{
            $(this).find('span').removeClass('fa-angle-up').addClass('fa-angle-down');
        }
        $('#collapseOne').collapse('toggle');

        // CC.host.hostlist.init();
    });

    window.CC.host=CC.host||{};
    /**
    * 树结构初始化、刷新函数定义
    */
    CC.host.topology={
        topoView:null,
        emptyView:null,
        init:function(data){
            var topo = data.topo;
            var empty = data.empty;
            // treeContainerChartObj.dataSource = new kendo.data.HierarchicalDataSource({data:topo});
            // $("#treeContainer").kendoTreeView(treeContainerChartObj);
            this.topoView=$("#treeContainer").data('kendoTreeView');
            if(typeof topo[0].items!='undefined' && topo[0].items.length>4){
                this.topoView.collapse('.k-item:not(.k-first)');
            }

            // emptyContainerChartObj.dataSource = new kendo.data.HierarchicalDataSource({data:empty});
            // $("#emptyContainer").kendoTreeView(emptyContainerChartObj);
            // this.emptyView=$("#emptyContainer").data('kendoTreeView');
            /**
            * 拓扑试图拖放主机事件
            */
            // $('#treeContainer,#emptyContainer').kendoDropTarget({
            //     group:"gridGroup",
            //     drop:function(e){
            //         var treeContainer = $(e.target).parents('.k-treeview')
            //         var grid = treeContainer.data('kendoTreeView');
            //         var data = grid.dataItem($(e.target).closest('.k-item'));

            //         if(data.type!='module'){
            //             var d = dialog({
            //                 content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>必需指定模块节点</div>'
            //             });
            //             d.show();
            //             setTimeout(function() {
            //                 d.close().remove();
            //             }, 2500);
            //             return false;
            //         }

            //         var param = {
            //             ApplicationID:cookie.get('defaultAppId'),
            //         };

            //         if(typeof JSON=='undefined'){
            //             $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //         }
            //         var grid = $('#host-list').data('kendoGrid');
            //         var newData = JSON.parse(JSON.stringify(hostlist.dataSource.data()));
            //         var hostId = [];
            //         for(var i=0,len=newData.length; i<len; i++){
            //             var d = newData[i];
            //             if(d.Checked==='checked'){
            //                 hostId.push(d.HostID);
            //             }
            //         }

            //         if(hostId.length==0){
            //             var sourceData = $('#host-list').data('kendoGrid').dataItem($(e.draggable.currentTarget).closest('tr'));
            //             param.HostID = sourceData.HostID;
            //         }else{
            //             param.HostID = hostId;
            //         }

            //         if(treeContainer.attr('id')=='treeContainer'){
            //             param.ModuleID = data.id;
            //         }

            //         var d = dialog({
            //             content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>转移中...</div>'
            //         });
            //         d.showModal();

            //         setTimeout(function(){
            //                 $.ajax({
            //                 url:treeContainer.attr('id')=='emptyContainer' ? '/host/delHostModule' : '/host/modHostModule/',
            //                 dialog:d,
            //                 data:param,
            //                 dataType:'json',
            //                 method:'post',
            //                 //async:false,
            //                 success:function(response){
            //                     this.dialog.close().remove();
            //                     var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
            //                     var d = dialog({
            //                             content: '<div class="c-dialogdiv2">'+content+'</div>'
            //                         });
            //                     d.show();
            //                     setTimeout(function () {
            //                         d.close().remove();
            //                     }, 2500);
            //                     // CC.host.hostlist.init();
            //                     CC.host.topology.refresh();
            //                     return true;
            //                 }
            //             });
            //         },100);
            //     }
            // });

            // 拓扑树根据浏览器调整自身高度
            function treeHeightChange () {
                $('.host-sidebar-left').css('position','fixed');
                setTimeout(function (){
                    $('.c-host-side').css('height',$(window).outerHeight()-70-20);
                    if ($('.c-host-side').height()>820)$('.c-host-side').css('height',820);
                    $(".c-tree-box").css('height',$('.c-host-side').outerHeight()-$('.free-group').outerHeight()-$('.c-host-side>h4').outerHeight()-40);
                },200)

            }
            treeHeightChange();
            $(window).resize(function (){
                treeHeightChange();
            })
            // 拓扑树根据浏览器调整自身高度 end
        },

        refresh:function(clickEmptyContainer){
            var me = this;
            $.ajax({
                url:'/host/getTopoTree4view',
                data:{ApplicationID:cookie.get('defaultAppId')},
                dataType:'json',
                method:'post',
                //async:false,
                success:function(response){
                    // kendo.destroy($("#treeContainer"));
                    $('#treeContainer').empty();
                    // kendo.destroy($("#emptyContainer"));
                    $('#emptyContainer').empty();
                    me.init(response);

                    clickEmptyContainer && $('.host-count', '#emptyContainer').click();
                }
            });

        }
    }

    /**
    * 主机列表kendoGrid相关
    */
    CC.host.hostlist={
        modHostModuleKeeper:$('#modHostModuleKeeper').html(),
        init:function(){
            var me=this;
            $('#host-list').empty();
            me._initColumns();
            // $("#host-list").kendoGrid(gridObj);
            // hostlist = $("#host-list").data('kendoGrid');

            if(typeof JSON=='undefined'){
                $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            }
            var data = JSON.parse(JSON.stringify(hostlist.dataSource.data()));
            for(var i=0,len=data.length; i<len; i++){
                data[i].Checked = '';
            }
            hostlist.dataSource.data(data);
            hostlist.refresh();

            me._init_checkAll();
            me._init_copy('.copy-inner-ip','InnerIP');
            me._init_copy('.copy-outer-ip','OuterIP');
            me._init_copy('.copy-asset-id','AssetID');
            if(typeof this.view=='undefined'){
                this.events();
            }
            this.view=hostlist;

            hostlist.checkedNum = 0;
            me._init_modHostModule();

            me._init_drag_drop();
        },
        _init_drag_drop:function(){
            // this.view.table.kendoDraggable({
            //     filter:".a-innerip",
            //     group:"gridGroup",
            //     cursorOffset:{top:5,left:5},
            //     hint:function(e){
            //         //return e.clone();
            //         var grid = $('#host-list').data('kendoGrid');
            //         if(typeof JSON=='undefined'){
            //             $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //         }

            //         var data = grid.dataSource.data();
            //         var data = JSON.parse(JSON.stringify(data));
            //         var html = '<div><table>';
            //         for(var i=0,len=data.length; i<len; i++){
            //             var d = data[i];
            //             if(d.Checked==='checked'){
            //                 html += '<tr><td>'+d.InnerIP+'</td></tr>';
            //             }
            //         }

            //         if(html==='<div><table>'){
            //             var d = grid.dataItem($(e).closest('tr'));
            //             html += '<tr><td>'+d.InnerIP+'</td></tr>';
            //         }
            //         return html+= '</table></div>';
            //     }
            // });
        },
        events:function(){
            var me = this;
            // hostlist.element.on('change.kendoGrid', 'input[type=checkbox]', function(e){
            //     var target = e.target;

            //     if($(target).parent('th').length>0){//表头
            //         var data = hostlist.dataSource.data();
            //         if(typeof JSON=='undefined'){
            //             $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //         }
            //         var d = JSON.parse(JSON.stringify(data));

            //         var checked = $(target).prop('checked');
            //         for(var i=0,len=d.length; i<len; i++){
            //             d[i].Checked = checked ? 'checked' : '';
            //         }

            //         hostlist.dataSource.data(d);
            //         hostlist.refresh();//kendoGrid初始化时，直接设置复选框的状态，不用手动设置
            //         hostlist.checkedNum = checked ? d.length : 0;
            //     }else{
            //         var checked = $(target).prop('checked');
            //         var dataItem = hostlist.dataItem($(target).closest('tr'));
            //         dataItem.Checked = checked ? 'checked' : '';
            //         if(checked){
            //            hostlist.checkedNum++;
            //         }else{
            //             hostlist.checkedNum--;
            //         }

            //         hostlist.thead.find('input[type=checkbox]').prop('checked', hostlist.checkedNum==hostlist.dataSource._total);
            //     }


            //     if(hostlist.checkedNum>0){
            //         $('#moveIp,#batDel,#batEdit,#batRes').attr('disabled', false);
            //         hostlist.element.find('.k-grid-copyIP').attr('disabled', false);
            //     }else{
            //         $('#moveIp,#batDel,#batEdit,#batRes').attr('disabled', true);
            //         hostlist.element.find('.k-grid-copyIP').attr('disabled', true);
            //     }

            //     me._init_modHostModule();

            //     if(hostlist.checkedNum==hostlist.dataSource._total){
            //         var selectAllDialog = dialog({
            //             align: 'bottom',
            //             content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>全选<i class="redFont">'+hostlist.dataSource._total+'</i>台主机</div>'
            //         });
            //         selectAllDialog.show(document.getElementById('dialogs'));
            //         setTimeout(function(){selectAllDialog.close().remove();}, 2000);
            //     }
            // });

            // hostlist.element.on('click.kendoGrid', '.a-innerip', function(e){
            //     var me = e.target;
            //     var grid = CC.host.hostlist.view;
            //     var data = grid.dataItem($(me).closest('tr'));

            //     var param = {};
            //     param['HostID'] = data.HostID;
            //     param['ApplicationID'] = cookie.get('defaultAppId');
            //     $.ajax({
            //         url:'/host/details',
            //         data:param,
            //         dataType:'html',
            //         method:'post',
            //         success:function(data){
            //             CC.rightPanel.show();
            //             CC.rightPanel.render(data);

            //             $(".sidebar-panel").mCustomScrollbar({
            //                 theme: "minimal-dark" //设置风格
            //             });

            //             $('.show-all-details-info').on('click',function(e){
            //                 $('.show-all-details-info').popover('destroy');
            //                 $(e.target).popover('show');
            //             });

            //             $('.sidebar-panel-container-new').on('click',function(e){
            //                 var className = $(e.target).prop('class');
            //                 if(className != 'show-all-details-info' && className != 'popover-content'){
            //                   $('.show-all-details-info').popover('destroy');
            //                 }
            //             });
            //         }
            //     });
            // });

            // /*导出到excel*/
            // hostlist.element.on('click.kendoGrid','.k-grid-saveAsExcel',function(e){
            //     if($(e.target).attr('disabled')=='disabled'){
            //         return false;
            //     }

            //     if($('#hostExport').length>0){
            //         $('#hostExport').remove();
            //     }

            //     if(typeof JSON=='undefined'){
            //         $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //     }
            //     var grid = $('#host-list').data('kendoGrid');
            //     var d = JSON.parse(JSON.stringify(grid.dataSource.data()));

            //     var hostId = [];
            //     for(var i=0,len=d.length; i<len; i++){
            //         hostId.push(d[i].HostID);
            //     }

            //     if(hostId.length==0){
            //         var d = dialog({
            //                 content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>没有主机可导出</div>'
            //             });
            //         d.show();
            //         setTimeout(function() {
            //             d.close().remove();
            //         }, 2500);
            //         return false;
            //     }
            //     $('body').append('<form id="hostExport" action="/host/hostExport" method="post" style="display:none;" target="_self"><input type="text" name="HostID" value="'+hostId.join(',')+'"><input type="hidden" name="ApplicationID" value="'+cookie.get('defaultAppId')+'"></form>');
            //     var d = dialog({
            //         content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>导出中...</div>'
            //     });
            //     setTimeout(function(){
            //         d.showModal();
            //         window.dintval = setInterval(function(){
            //             var dowcomplete = cookie.get('comdownload');
            //             if(dowcomplete == 1){
            //                 d.close().remove();
            //                 cookie.set('comdownload','');
            //                 clearInterval(dintval);
            //             }
            //         }, 500);
            //         $('#hostExport').submit();
            //     },500);

            // });

            /* 点击删除按钮 */
            // hostlist.element.on('click.kendoGrid', '#batDel', function(e){
            //     if($(e.target).attr('disabled')=='disabled'){
            //         return false;
            //     }
            //     var data = CC.host.hostlist.view.dataSource.data();
            //     if(typeof JSON=='undefined'){
            //         $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //     }
            //     var newData = JSON.parse(JSON.stringify(data));
            //     var hostId = [];
            //     for(var i=0,len=newData.length; i<len; i++){
            //         if(newData[i].Checked=='checked'){
            //             hostId.push(newData[i].HostID);
            //         }
            //     }

            //     if(hostId.length==0){
            //         var noHostSelectDialog = dialog({
            //                 content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            //             });
            //         noHostSelectDialog.show();
            //         setTimeout(function () {
            //             noHostSelectDialog.close().remove();
            //         }, 2000);
            //         return false;
            //     }

            //     var param = {};
            //     param['ApplicationID'] = cookie.get('defaultAppId');
            //     param['HostID'] = hostId.join(',');

            //     var gridBatDel = dialog({
            //         title:'确认',
            //         width:300,
            //         content: '确认是否将已勾选的<i class="redFont">'+hostId.length+'</i>台主机移动至空闲机?',
            //         okValue:"确定",
            //         cancelValue:"取消",
            //         ok:function (){
            //             var d = dialog({
            //                 content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>转移中...</div>'
            //             });
            //             d.showModal();
            //             $.ajax({
            //                 url:'/host/delHostModule/',
            //                 dialog:d,
            //                 data:param,
            //                 method:'post',
            //                 dataType:'json',
            //                 //async:false,
            //                 success:function(response){
            //                     this.dialog.close().remove();
            //                     var content = '<i class="c-dialogimg-'+ (response.success==true ? 'success' : 'prompt') +'"></i>'+ response.message +'</div>';
            //                     var d = dialog({
            //                         content: '<div class="c-dialogdiv2">'+content+'</div>'
            //                     });
            //                     d.show();
            //                     setTimeout(function() {
            //                         d.close().remove();
            //                     }, 2500);
            //                     // CC.host.hostlist.init();
            //                     CC.host.topology.refresh();
            //                     return true;
            //                 }
            //             });
            //         },
            //         cancel: function () {
            //         }
            //     });

            //     gridBatDel.showModal();
            // });

            // /* 点击上交按钮 */
            // hostlist.element.on('click.kendoGrid', '#batRes', function (e){
            //     if($(e.target).attr('disabled')=='disabled'){
            //         return false;
            //     }
            //     var grid = CC.host.hostlist.view;
            //     var data = grid.dataSource.data();
            //     if(typeof JSON=='undefined'){
            //         $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //     }
            //     var newData = JSON.parse(JSON.stringify(data));
            //     var hostId = [];
            //     for(var i=0,len=newData.length; i<len; i++){
            //         if(newData[i].Checked==='checked'){
            //             hostId.push(newData[i].HostID);
            //         }
            //     }

            //     if(hostId.length==0){
            //         var noHostSelectDialog = dialog({
            //                 content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            //             });
            //         noHostSelectDialog.show();
            //         setTimeout(function () {
            //             noHostSelectDialog.close().remove();
            //         }, 2000);
            //         return false;
            //     }


            //     var param = {};
            //     param['ApplicationID'] = cookie.get('defaultAppId');
            //     param['HostID'] = hostId.join(',');

            //     var gridBatRes = dialog({
            //         title:'确认',
            //         width:300,
            //         content: '确认是否将已勾选的<i class="redFont">'+hostId.length+'</i>台主机上交至资源池',
            //         okValue:"确定",
            //         cancelValue:"取消",
            //         ok:function (){
            //             var d = dialog({
            //                 content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>上交中...</div>'
            //             });
            //             d.showModal();
            //             $.ajax({
            //                 url:'/host/resHostModule/',
            //                 dialog:d,
            //                 data:param,
            //                 method:'post',
            //                 dataType:'json',
            //                 //async:false,
            //                 success:function(response){
            //                     this.dialog.close().remove();
            //                     var content = '<i class="c-dialogimg-'+ (response.success==true ? 'success' : 'prompt') +'"></i>'+ response.message +'</div>';
            //                     var d = dialog({
            //                         content: '<div class="c-dialogdiv2">'+content+'</div>'
            //                     });
            //                     d.show();
            //                     setTimeout(function() {
            //                         d.close().remove();
            //                     }, 2500);
            //                     // CC.host.hostlist.init();
            //                     CC.host.topology.refresh(true);

            //                     return true;
            //                 }
            //             });
            //         },
            //         cancel: function () {
            //         }
            //     });

            //     gridBatRes.showModal();
            // });

            /* 点击修改选中按钮 */
            // hostlist.element.on('click.kendoGrid', '#batEdit', function (e){
            //     if($(e.target).attr('disabled')=='disabled'){
            //         return false;
            //     }

            //     var hostId = [];
            //     var appId = cookie.get('defaultAppId');
            //     var grid = $('#host-list').data('kendoGrid');
            //     var data = grid.dataSource.data();
            //     if(typeof JSON=='undefined'){
            //         $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            //     }
            //     var newData = JSON.parse(JSON.stringify(data));
            //     var hostInfo = [];
            //     for(var i=0,len=newData.length; i<len; i++){
            //         if(newData[i].Checked==='checked'){
            //             hostId.push(newData[i].HostID);
            //             hostInfo = newData[i];
            //         }
            //     }

            //     if(hostId.length==0){
            //         var noHostSelectDialog = dialog({
            //                 content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            //             });
            //         noHostSelectDialog.show();
            //         setTimeout(function () {
            //             noHostSelectDialog.close().remove();
            //         }, 2000);
            //         return false;
            //     }

            //     if(hostId.length==1){
            //         $('#moduleHost_HostName').val(hostInfo['HostName']);
            //         $('#moduleHost_Description').val(hostInfo['Description']);
            //         $('#moduleHost_Source').val(hostInfo['Source']);

            //         $("#moduleHost_Operator").select2({
            //             placeholder:'请选择',
            //             data:window.userList,
            //             formatResult:function format(state) {
            //             return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            //         },
            //         formatSelection:function format(state) {
            //             return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            //         }
            //     }).select2('val',hostInfo['Operator']).select2("enable", false);

            //         $("#moduleHost_BakOperator").select2({
            //             placeholder:'请选择',
            //             data:window.userList,
            //             formatResult:function format(state) {
            //                         return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            //                     },
            //             formatSelection:function format(state) {
            //                         return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            //                     }
            //         }).select2('val',hostInfo['BakOperator']).select2("enable", false);

            //         $("#moduleHost_Source").select2({
            //             placeholder:'请选择',
            //             data:window.hostSource,
            //             formatResult:function format(state) {
            //                 return state.text;
            //             },
            //             formatSelection:function format(state) {
            //                 return state.text;
            //             }
            //         }).select2('val',hostInfo['Source']).select2("enable", false);
            //     }else{
            //         $("#moduleHost_Operator,#moduleHost_BakOperator").select2({
            //             placeholder:'请选择',
            //             data:window.userList,
            //             formatResult:function format(state) {
            //                         return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            //                     },
            //             formatSelection:function format(state) {
            //                         return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            //                     }
            //         }).select2('val','').select2("enable", false);

            //         $("#moduleHost_Source").select2({
            //             placeholder:'请选择',
            //             data:window.hostSource,
            //             formatResult:function format(state) {
            //                 return state.text;
            //             },
            //             formatSelection:function format(state) {
            //                 return state.text;
            //             }
            //         }).select2('val',hostInfo['Source']).select2("enable", false);
            //     }

            //     $('#batEdit_window').modal('show');
            // });

            $(document.body).on('change', '.ui-multiselect-checkboxes', function(e){
                var disabled = true;
                $(e.target).closest('ul').find('input[type=checkbox]').each(function(index, el){
                    if($(el).prop('checked')){
                        disabled = false;
                        return false;
                    }
                });

                $('.btn-primary', '#modSelectMenu').prop('disabled', disabled);
            });

            $(document.body).on('mouseleave', '.ui-multiselect-checkboxes', function(e){
                $(this).find('.ui-state-hover').removeClass('ui-state-hover');
            });

            $(document.body).on('click', '#modSelectMenu', function(e){
                if($(e.target).hasClass('btn-default')){
                    $(e.target).parents('.ui-multiselect-menu').find('.ui-icon-circle-close').click();
                }else{

                    if($(e.target).prop('disabled')!==false){
                        return false;
                    }
                    var appId = cookie.get('defaultAppId');
                    var moduleId = [];
                    var hostId = [];
                    $('#modSelectMenu').prev('ul').find('input[type=checkbox]').each(function(index, el){
                        if($(el).attr('aria-selected')=='true'){
                            moduleId.push($(el).val());
                        }
                    });
                    $(e.target).parents('.ui-multiselect-menu').find('.ui-icon-circle-close').click();

                    if(moduleId.length==0){
                        var content = '';
                        var d = dialog({
                                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请至少选择一个模块</div>'
                            });
                        d.show();
                        setTimeout(function () {
                            d.close().remove();
                        }, 2500);

                        return false;
                    }

                    // var grid = $('#host-list').data('kendoGrid');
                    var data = grid.dataSource.data();
                    if(typeof JSON=='undefined'){
                        $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
                    }
                    var d = JSON.parse(JSON.stringify(data));
                    for(var i in d){
                        if(d[i].Checked==='checked'){
                            hostId.push(d[i]['HostID']);
                        }
                    }

                    if(hostId.length==0){
                        var content = '';
                        var d = dialog({
                                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请至少选择一台主机</div>'
                            });
                        d.show();
                        setTimeout(function () {
                            d.close().remove();
                        }, 2500);

                        return false;
                    }

                    var param = {};
                    param['ApplicationID'] = appId;
                    param['ModuleID'] = moduleId.join(',');
                    param['HostID'] = hostId.join(',');

                    var d = dialog({
                        content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>转移中...</div>'
                    });
                    d.showModal();

                    setTimeout(function(){
                            $.ajax({
                            url:'/host/modHostModule/',
                            dialog:d,
                            data:param,
                            dataType:'json',
                            method:'post',
                            success:function(response){
                                this.dialog.close().remove();
                                var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
                                var d = dialog({
                                        content: '<div class="c-dialogdiv2">'+content+'</div>'
                                    });
                                d.show();
                                setTimeout(function () {
                                    d.close().remove();
                                }, 2500);
                                // CC.host.hostlist.init();
                                CC.host.topology.refresh();
                                return true;
                            }
                        });
                    },100);
                }
            });
        },

        /**该属性可覆盖**/
        defaultColumns:[],
        _initColumns:function(){
            var me = this;

            if(!me.defaultColumns || me.defaultColumns.length===0){
                me.defaultColumns = ['InnerIP','OuterIP','SetName','ModuleName','HostName','SN','ApplicationName'];
            }
            for(var i=0,len=allColumns.length; i<len; i++){
                allColumns[i].hidden = false;
                if(allColumns[i].field!='checkbox' && -1==$.inArray(allColumns[i].field,me.defaultColumns)){
                    allColumns[i].hidden=true;
                }
            }
            var disColumn =[];
            var disColumnFalse =[];
            var disColumncp = [];
            for(var j in allColumns){
                if(allColumns[j].hidden == false) {
                    disColumn.push(allColumns[j]);
                }else{
                    disColumnFalse.push(allColumns[j]);
                }
            }
            disColumncp.push(allColumns[0]);
            for(var i in me.defaultColumns){
                for(var j in disColumn){
                    if(me.defaultColumns[i] == disColumn[j].field){
                        disColumncp.push(disColumn[j]);
                    }
                }
            }
            var totalColumn = disColumncp.concat(disColumnFalse);
            gridObj.columns = totalColumn;
        },

        _init_checkAll:function(){
            var checkAll = $('<input type="checkbox"/>');
            hostlist.thead.find('th[data-field=checkbox]').empty().append(checkAll);
        },

        _init_copy:function(btnclass,filedname){
            /**
             * copy
             */
             var me = this;
            var copyIpBtn = hostlist.element.find(btnclass);
            ZeroClipboard.config({moviePath:'/assets/ZeroClipboard/ZeroClipboard.swf'});
            var clip = new ZeroClipboard(copyIpBtn.get(0));
            clip.on('copy',function(e){
                     var clipboarde=e.clipboardData;
                     me._copyIp(copyIpBtn,clipboarde,filedname);
            });
            clip.on('aftercopy',function(e){
                if(e.data['text/plain']){
                    var d = dialog({
                            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>复制成功</div>'
                        });
                    d.show();
                    setTimeout(function() {
                        d.close().remove();
                    }, 2500);
                }
            });
        },

        _copyIp:function(btn,clipboard,filedname){
            var me=this,
                list=me._getCopyList(filedname);
            if(list.length){
                clipboard.setData('text/plain',list.join("\n"));
            }else{
                if($('.k-grid-copyIP').attr('disabled')=='disabled'){
                    return false;
                }
                var noHostSelectDialog = dialog({
                        content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
                    });
                noHostSelectDialog.show();
                setTimeout(function () {
                    noHostSelectDialog.close().remove();
                }, 2000);
                return false;
            }
        },
        _getCopyList:function(key){
            var list=[];
            if(typeof JSON=='undefined'){
                $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
            }
            var newData = JSON.parse(JSON.stringify($('#host-list').data('kendoGrid').dataSource.data()));
            for(var i=0,len=newData.length; i<len; i++){
                if(newData[i].Checked==='checked'){
                    list.push(newData[i][key]);
                }
            }
            return list;
        },
        _init_modHostModule:function(){
            if(hostlist.checkedNum>0){
                var tmp = $(this.modHostModuleKeeper).removeAttr('disabled').get(0).outerHTML;
            }else{
                var tmp = $(this.modHostModuleKeeper).attr('disabled', 'disabled').get(0).outerHTML;
            }

            $('#modHostModuleKeeper').empty().append(tmp);

            $('.k-grid-moveIp').find('button').remove().end().find('select').remove();
            $('.ui-multiselect-menu').remove();
            $('#modHostModule').appendTo($('.k-grid-moveIp')).multiselect({
                noneSelectedText: "转移主机至",
                checkAllText: "全选",
                uncheckAllText: '全不选',
                selectedList:0
            }).multiselectfilter({
                label:'',
                placeholder:'',
                filter: function(event){
                }
            });

            $(".ui-multiselect-checkboxes").mCustomScrollbar({
                theme: "minimal-dark" //设置风格
            });
            $('.ui-multiselect-close').siblings().css('display','none');
            $('.ui-multiselect-menu').append("<div class='text-center' id='modSelectMenu'><button class='btn btn-xs btn-default mr10'>取消</button><button class='btn btn-xs btn-primary' disabled='disabled'>转移</button></div>");
        }
    };

    var wH=$('.content-wrapper').height(),
        hostlist,
        allColumns=tablesFields,
        gridObj={
            change:function(e){
                var selectedRows = this.select();
                var table = $('#host-list').find('.k-selectable tbody')
                var userSelection = [];
                var preParent = 0;
                var rowText =  '';
                var lens=0;

                for (var i= 0,len=selectedRows.length; i<len; i++){
                    var cur = selectedRows.eq(i);
                    if (i>0){
                        preParent = selectedRows.eq(i-1).closest('tr').index();
                        if (cur.closest('tr').index() == preParent){
                            rowText = rowText + '   ' + cur.text()
                        }else{
                            rowText = rowText + '\n' + cur.text()
                        }
                    }else{
                        rowText = cur.text()
                    }
                    lens=(selectedRows.eq(i).width());
                }
                var scrollLeft=$('.k-grid-content').scrollLeft();
                var scrollTop=$('.k-grid-content').scrollTop();
                var grid = $('#host-list').data('kendoGrid');
                var pos=grid.select().last().offset();
                userSelection = rowText;
                if(userSelection!=null){
                    var relativeX = (pos.left-555+lens+scrollLeft);
                    var relativeY = (pos.top-212+scrollTop);
                    var static_src = static_url+"/static/img/copy.png";
                    var srcimg = '<div id="c-hereis"><img src="'+static_src+'" alt="复制文本"></div>' ;
                    $('#host-list').find('.k-selectable tbody').append(srcimg);
                    $('#c-hereis').css({"top":relativeY,"left":relativeX});
                    $('#c-hereis').show();
                }else{
                    $('#c-hereis').hide();
                }
                ZeroClipboard.config({moviePath:'/assets/ZeroClipboard/ZeroClipboard.swf'});
                var clips=new ZeroClipboard($('#c-hereis')[0]);
                clips.on('copy',function(e){
                    var clipboarde=e.clipboardData;
                    clipboarde.setData('text/plain',userSelection);
                });
                clips.on('aftercopy',function(e){
                    if(e.data['text/plain']){
                        var d = dialog({
                            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>复制成功</div>'
                        });
                        d.show();
                        setTimeout(function() {
                            d.close().remove();
                        }, 2000);
                        $('#c-hereis').hide();
                    }
                });
                $(selectedRows).closest('tr').find('input').prop('checked', true);
                for(var i=0,len=selectedRows.length; i<len; i++){
                    var data = grid.dataItem($(selectedRows[i]).closest('tr'));
                    if(data.Checked != 'checked'){
                        data.Checked = 'checked';
                        this.checkedNum++;
                    }
                }
                $('#moveIp,#batDel,#batEdit,#batRes').attr('disabled', this.checkedNum==0);
                CC.host.hostlist._init_modHostModule();
                hostlist.element.find('.k-grid-copyIP').attr('disabled', this.checkedNum==0);
            },
            dataBound:function(){
                this.checkedNum = 0;
            },
            columnHide:function(e){
                var param = {};
                param['ApplicationID'] = cookie.get('defaultAppId');
                param['DefaultColumn'] = [];

                var columns = this.columns;
                for(var i in columns){
                    if(!columns[i].hidden && columns[i].field!='checkbox'){
                        param['DefaultColumn'].push(columns[i].field);
                    }
                }

                $.ajax({
                    url:'/UserCustom/setUserCustom/',
                    dataType:'json',
                    data:param,
                    method:'post',
                    success:function(response){
                    }
                });

                CC.host.hostlist.defaultColumns = param['DefaultColumn'];
            },
            columnShow:function(e){
                var param = {};
                param['ApplicationID'] = cookie.get('defaultAppId');
                param['DefaultColumn'] = [];

                var columns = this.columns;
                for(var i in columns){
                    if(!columns[i].hidden && columns[i].field!='checkbox'){
                        param['DefaultColumn'].push(columns[i].field);
                    }
                }

                $.ajax({
                    url:'/UserCustom/setUserCustom/',
                    dataType:'json',
                    data:param,
                    method:'post',
                    success:function(response){
                    }
                });

                CC.host.hostlist.defaultColumns = param['DefaultColumn'];
            },
            dataSource: {//数据源配置项
                transport: {
                    read: {
                        url: "/host/getHostById",
                        data:{'ApplicationID':cookie.get('defaultAppId')},
                        method:'post',
                        dataType:"json"
                    }
                },
                pageSize:20,
                serverPaging: false,
                schema: {
                    data: function (response) {
                        var param = $('#host-list').data('kendoGrid').dataSource.options.transport.read.data;
                        if(typeof param.SetID=='undefined' && param.ModuleID==''){
                            return [];
                        }
                        return response.data;
                    },
                    total: function (response) {
                        var param = $('#host-list').data('kendoGrid').dataSource.options.transport.read.data;
                        if(typeof param.SetID=='undefined' && param.ModuleID==''){
                            return 0;
                        }
                        return response.total;
                    }
                }
            },
            toolbar: [//头部工具栏kendoToolBar,可以参考ui.toolbar的api
                {
                    text: '<div class="copy-menu-box">'+
                    '<div>复制<span class="caret"></span></div>'+
                    '<ul class="copy-menu none">'+
                    '<li class="copy-inner-ip" data-type="InnerIP">复制内网IP</li>'+
                    '<li class="copy-outer-ip" data-type="HostID">复制外网IP</li>'+
                    '<li class="copy-asset-id" data-type="copy-things-num">复制固资编号</li>'+
                    '</ul>'+
                    '</div>',
                    name:'copyIP',
                    attr:{"href":"javascript:void(0);","disabled":"true"}
                },
                {text:'修改',name:'batEdit',attr:{"id":"batEdit","href":"javascript:void(0);","disabled":"true"}},
                {text:'',name:'moveIp',attr:{"id":"moveIp","href":"javascript:void(0);","disabled":"true"}},
                {text:'导出Excel',name:'saveAsExcel',attr:{"href":"javascript:void(0);"}},
                {text:'上交',name:'batRes',attr:{"id":"batRes","href":"javascript:void(0);","style":"display:none;float:right","disabled":"true"}},
                {text:'移至空闲机',name:'batDel',attr:{"id":"batDel","href":"javascript:void(0);","disabled":"true"}}
            ],
            selectable:'multiple cell',
            allowCopy:{delimiter : ';'},
            height:753,
            resizable:true,
            pageable: true
        };

    /* 点击修改选中弹窗的checkbox*/
    $('#batEdit_window').on('click', '.select-area', function(e){
        if(e.target.nodeName.toLowerCase()=='td'){
            $(e.target).find('input').click();
            return false;
        }
        var target = e.target;
        var checked = $(target).prop('checked');
        if($(target).attr('data')=='selectAll'){//全选checkbox
            $('#batEdit_window input[type=checkbox]').not(target).prop('checked', checked);
            $('#moduleHost_Operator,#moduleHost_BakOperator').select2("enable", checked);
            checked && $('#batEdit_window').find('.edit-area-mask').addClass('none');
            !checked && $('#batEdit_window').find('.edit-area-mask').removeClass('none');
            $('#batEdit_window input[type=text]').prop('disabled', !checked);

        }else{//其他checkbox
            if(!checked){
                $('#batEdit_window input[data=selectAll]').prop('checked', checked);
            }

            $(target).closest('tr').find('input[type=text]').prop('disabled',!checked);
            $(target).closest('tr').find('select').prop('disabled',!checked);

            var id = $(target).closest('tr').find('td:last>div:eq(1)').attr('id');
            if(id){
                checked && $(target).closest('tr').find('.edit-area-mask').addClass('none');
                !checked && $(target).closest('tr').find('.edit-area-mask').removeClass('none');
                $('#'+id).select2("enable", checked);
            }

            var allChecked = true;
            $('#batEdit_window input[type=checkbox]').each(function(i, el){
                if($(el).attr('data')=='selectAll'){
                    return true;
                }
                if($(el).prop('checked')==false){
                    allChecked = false;
                    return false;
                }
            });

            $('#batEdit_window input[data=selectAll]').prop('checked', allChecked);
        }
    });

    /* 点击修改选中弹窗的取消按钮*/
    $('#batEdit_window').on('click', '#moduleHostHide,.close', function(e){
        $('#batEdit_window .edit-area-mask').removeClass('none');
        $("#s2id_moduleHost_Operator,#s2id_moduleHost_BakOperator").select2('val', '').select2('enable', false);
        $('#batEdit_window').find('input[type=checkbox]').prop('checked', false).end().find('input[type=text]').val('').attr('disabled', true).end().modal('hide');
    });
    /* 点击修改选中弹窗的保存按钮*/
    // $('#batEdit_window').on('click', '#moduleHostSubmit', function(e){
    //     var hostInfo = {};
    //     var stdProperty = {};
    //     var cusProperty = {};

    //     if(!$('#moduleHost_HostName').prop('disabled')){
    //         stdProperty['HostName'] = $('#moduleHost_HostName').val();
    //         if(stdProperty['HostName']==null || stdProperty['HostName']==''){
    //             var d = dialog({
    //                     content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>主机名不能为空</div>',
    //                     zIndex:1051
    //                 });
    //             d.show();
    //             setTimeout(function () {
    //                 d.close().remove();
    //                 $('#moduleHost_HostName').focus();
    //             }, 1500);
    //             $('#moduleHost_HostName').focus();
    //             return false;
    //         }
    //     }

    //     var empty_field = [];
    //     var empty_field_disable = [];
    //     if(!$('#moduleHost_Operator').prop('disabled')){
    //         stdProperty['Operator'] = $('#moduleHost_Operator').select2('val');
    //         stdProperty['Operator']=='' && empty_field.push('负责人');
    //     }

    //     if(!$('#moduleHost_BakOperator').prop('disabled')){
    //         stdProperty['BakOperator'] = $('#moduleHost_BakOperator').select2('val');
    //         stdProperty['BakOperator']=='' && empty_field.push('备份负责人');
    //     }

    //     if(!$('#moduleHost_Source').prop('disabled')){
    //         stdProperty['Source'] = $('#moduleHost_Source').val();
    //         stdProperty['Source']=='' && empty_field_disable.push('云供应商');
    //     }

    //     if(!$('#moduleHost_Description').prop('disabled')){
    //         stdProperty['Description'] = $('#moduleHost_Description').val();
    //         stdProperty['Description']=='' && empty_field.push('备注信息');
    //     }

    //     var data = CC.host.hostlist.view.dataSource.data();
    //     if(typeof JSON=='undefined'){
    //         $('head').append('<script type="text/javascript" src="/static/js/json2.js"></script>');
    //     }
    //     var newData = JSON.parse(JSON.stringify(data));
    //     var hostId = [];
    //     for(var i in newData){
    //         if(newData[i].Checked==='checked'){
    //             hostId.push(newData[i].HostID);
    //         }
    //     }
    //     hostInfo['HostID'] = hostId.join(',');
    //     hostInfo['ApplicationID'] = cookie.get('defaultAppId');
    //     hostInfo['stdProperty'] = stdProperty;
    //     hostInfo['cusProperty'] = cusProperty;

    //     if(empty_field_disable.length > 0){
    //         var promptDialog = dialog({
    //             title:'提示',
    //             width:300,
    //             okValue:"确定",
    //             zIndex:1051,
    //             ok:function(){},
    //             content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>' + empty_field_disable.join('、') + '不能设置为空</div>'
    //         });

    //         promptDialog.showModal();
    //         return true;
    //     }

    //     if(empty_field.length > 0){
    //         var confirmDialog = dialog({
    //             title:'确认',
    //             width:300,
    //             zIndex:1051,
    //             content: empty_field.join('、')+'将设置为空，确认继续？',
    //             okValue:"继续",
    //             cancelValue:"取消",
    //             ok:function (){
    //                 if(!$.isEmptyObject(hostInfo['stdProperty']) && hostInfo['ApplicationID']!='' && hostInfo['HostID']!=''){
    //                     $.ajax({
    //                         url:'/host/updateHostInfo/',
    //                         data:hostInfo,
    //                         dataType:'json',
    //                         method:'post',
    //                         success:function(response){
    //                             var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
    //                             var d = dialog({
    //                                     content: '<div class="c-dialogdiv2">'+content+'</div>'
    //                                 });
    //                             d.show();
    //                             setTimeout(function () {
    //                                 d.close().remove();
    //                             }, 2500);
    //                             // CC.host.hostlist.init();
    //                             return true;
    //                         }
    //                     });

    //                     $("#s2id_moduleHost_Operator,#s2id_moduleHost_BakOperator").select2('val', '');
    //                     $('#batEdit_window').find('input[type=checkbox]').prop('checked', false).end().find('input[type=text]').val('').attr('disabled', true).end().modal('hide');
    //                 }
    //             },
    //             cancel: function () {
    //             }
    //         });

    //         confirmDialog.showModal();
    //         return true;
    //     }

    //     if(!$.isEmptyObject(hostInfo['stdProperty']) && hostInfo['ApplicationID']!='' && hostInfo['HostID']!=''){
    //         $.ajax({
    //             url:'/host/updateHostInfo/',
    //             data:hostInfo,
    //             dataType:'json',
    //             method:'post',
    //             success:function(response){
    //                 var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
    //                 var d = dialog({
    //                         content: '<div class="c-dialogdiv2">'+content+'</div>'
    //                     });
    //                 d.show();
    //                 setTimeout(function () {
    //                     d.close().remove();
    //                 }, 2500);
    //                 // CC.host.hostlist.init();
    //                 return true;
    //             }
    //         });
    //     }

    //     $("#s2id_moduleHost_Operator,#s2id_moduleHost_BakOperator").select2('val', '');
    //     $('#batEdit_window').find('input[type=checkbox]').prop('checked', false).end().find('input[type=text]').val('').attr('disabled', true).end().modal('hide');
    // });

    /* 点击修改选中弹窗的输入框*/
    $('#batEdit_window').on('click', '.edit-area', function(e){

        var checkbox = $(e.target).closest('tr').find('input').eq(0);
        if(!checkbox.prop('checked')){
            checkbox.click();
            if(e.target.nodeName.toLowerCase()=='div'){
                $(e.target).addClass('none');
                $(e.target).prev('div').select2('enable', true);
            }else{
                $(e.target).focus();
            }
        }
    });


    /**
    * 查询条件“重置”按钮点击事件处理函数
    */
    $('#host_query_reset').click(function(){
        $('#InnerIP').val('');
        $('#OuterIP').val('');
        $('#set_select-selectDialogWarp').html('<div class="btn-toolbar"><span style="color:#999999;">(不选为全部)</span></div>');
        $('#set_select-content').find('button').removeClass('btn-primary');
        $('#module_select-selectDialogWarp').html('<div class="btn-toolbar"><span style="color:#999999;">(不选为全部)</span></div>');
        $('#module_select-content button').find('button').removeClass('btn-primary');
        $("[name='IfInnerIPexact']").prop('checked', false);
        $("[name='IfOuterexact']").prop('checked', false);
        $('[filter-name="CreateTime"]').val('');
        $.each(customerQueryFields, function(index, v) {
            if($.inArray(index, ['Operator', 'BakOperator', 'OSName']) > -1) {
                $('#' + index).select2('val','');
            }else {
                $('#' + index).val('');
            }
        });
    });

    /**
    * 查询条件“查询”按钮点击事件处理函数
    */
    $('#host_query_submit').click(function(){
        var param = {};
        var InnerIP = $('#InnerIP').val();
        if(InnerIP){
            param['InnerIP'] = InnerIP.replace(/\s+/g, ',').replace(/[\s,]+$/g, '');
        }

        var OuterIP = $('#OuterIP').val();
        if(OuterIP){
            param['OuterIP'] = OuterIP.replace(/\s+/g, ',').replace(/[\s,]+$/g, '');
        }

        var ModuleID = [];
        $('#module_select-selectDialogWarp a').each(function(i, el){
            ModuleID.push($(el).attr('data'));
        });
        if(ModuleID.length>0){
            param['ModuleID']=ModuleID.join(',');
        }

        var SetID = [];
        $('#set_select-selectDialogWarp a').each(function(i, el){
            SetID.push($(el).attr('data'));
        });
        if(SetID.length>0){
            param['SetID']=SetID.join(',');
        }

        var IfInnerIPexact = $("[name='IfInnerIPexact']").prop('checked');
        if(IfInnerIPexact){
            param['IfInnerIPexact'] = IfInnerIPexact;
        }

        var IfOuterexact = $("[name='IfOuterexact']").prop('checked');
        if(IfOuterexact){
            param['IfOuterexact'] = IfOuterexact;
        }

        $.each(customerQueryFields, function(index, v) {
            var stdProperty = $('#' + index).val();
            if(stdProperty){
                param[index] = stdProperty.replace(/\s+/g, ',').replace(/[\s,]+$/g, '');
            }
        });

        param['ApplicationID'] = cookie.get('defaultAppId');

        gridObj.dataSource.transport.read.data = param;
        gridObj.dataSource.transport.read.url = '/host/getHostByCondition/';
        // CC.host.hostlist.init();

        var emptyModule = $('#module_select-selectDialogWarp a');
        if(emptyModule.length===1 && emptyModule.html().indexOf('空闲机')>-1){
            $('#batDel').hide();
            $('#batRes').show();
        }
    });
}();

/**
/* 弹窗选择插件
 * auth：v_weilli
 * 依赖组件 ：jquery artDialog bootstrap
 * 使用方式 ：$('#id').selectDialog({配置项});
 * 取值方式 ：$('#id').val();
 * 可用配置项 ：弹窗的宽：width, 弹窗的高：height, 弹窗的标题：title, 确定按钮的文字:okVal,关闭按钮的文字：cancelVal, 无选择时显示的文字:emptyText
 */

(function($) {
    $.fn.selectDialog=function(options) {
        var defaults= {
            width: 800, height: 'auto', title: '选择', okValue: '保存', cancelValue: '取消', emptyText: ''
        }
        var options=$.extend(defaults, options);
        var values=[];
        var optionText=[];
        var select=this;
        var module='';
        if(select.attr('multiple')==='multiple') {
            //多选
            module='multiple';
        }
        else {
            module='single';
        }
        options.selected=[]; //选中的行
        options.emptyText=options.emptyText?options.emptyText: select.attr('placeholder'); //无选择时显示的文字
        select.find('option').each(function(index, value) {
            var v=$(value).attr('value');
            if($(value).attr('selected')==="selected") {
                options.selected.push(v);
            }
            values.push(v);
            optionText.push($(value).text());
        }
        );
        var id=select.attr('id');
        var warpId=id + '-selectDialogWarp';
        var contentId=id + '-content';
        var dialogId=id + '-dialogId';
        var open_dialog=function() {
            var d=dialog( {
                okValue: options.okValue,
                cancelValue: options.cancelValue,
                title: options.title, width: options.width, height: options.height, lock: true, content: creat_content(optionText), ok: function() {
                    $(document).unbind('keydown');
                    creat_warp_content();
                    //点保存以后，设置select值
                    select.find('option').each(function(index, value) {
                        var option=$(value), v=option.attr('value');
                        if($.inArray(v, options.selected)===-1) {
                            option.removeAttr('selected');
                        }
                        else {
                            option.attr('selected', 'selected')
                        }
                    }
                    );
                    select.trigger('change');
                }
                , onshow: function() {
                    //对话框弹出时执行的函数
                    $(".ui-dialog-body").css('vertical-align', 'top');
                    $(".ui-dialog-body").css('text-align', 'left');
                    $(".ui-dialog-content .btn-toolbar .btn").click(function() {
                        //绑定按钮点击事件
                        btnClick(this);
                    }
                    );
                    $(".ui-dialog-content").find('button[data-id="check_all"]').click(function() {
                        check_all();
                    }
                    );
                    $(".ui-dialog-content").find('button[data-id="inverse"]').click(function() {
                        inverse();
                    }
                    );
                    $(document).keyup(function(event) {
                        if(event.keyCode==27) {
                            d.close();
                        }
                    }
                    );
                    //过滤
                    $('input[data-id="keyword"]').bind('keyup', filter);
                    //如果是单选，隐藏全选和反选
                    if(module==='single') {
                        $('#'+dialogId).find('div[data-id="btn-group"]').hide();
                    }
                }
                , okVal: options.okVal, cancel: function() {
                    $(".ui-dialog-body").css('vertical-align', 'middle');
                    $(".ui-dialog-body").css('text-align', 'center');
                    $(".ui-dialog-close").trigger('mouseout');
                    $(document).unbind('keydown');
                }
                , cancelVal: options.cancelVal
            });
            d.showModal();
        }
        var check_all=function() {
            $(".ui-dialog-content .btn-toolbar .btn").each(function() {
                var v=values[parseInt($(this).attr('data-index'))];
                if(v && !$(this).is(':hidden')) {
                    add(v);
                    $(this).addClass('btn-primary');
                }
            }
            );
        }
        var inverse=function() {
            $(".ui-dialog-content .btn-toolbar .btn").each(function() {
                if(!$(this).is(':hidden'))$(this).trigger('click');
            }
            );
            //反选后执行一次过滤操作
            $('input[data-id="keyword"]').trigger('keyup');
        }
        var btnClick=function (btn) {
            var btn=$(btn);
            var value=values[parseInt(btn.attr('data-index'))];
            if(btn.hasClass('btn-primary')) {
                btn.removeClass('btn-primary');
                remove(value);
            }
            else {
                //单选模式下，取消其他选中的
                if(module==='single') {
                    $('#'+ contentId + ' .btn-primary').each(function(index, v) {
                        $(v).trigger('click');
                    }
                    );
                }
                btn.addClass('btn-primary');
                add(value);
            }
        }
        var add=function(v) {
            if($.inArray(v, options.selected)===-1) {
                options.selected.push(v);
            }
        }
        var remove=function(v) {
            var key=$.inArray(v, options.selected);
            options.selected.splice(key, 1);
        }
        var creat_warp_content=function() {
            var content='<div class="btn-toolbar">';
            if(options.selected.length===0 && options.emptyText) {
                content +='<span style="color:#999999;">'+ options.emptyText +'</span>'
            }
            else {
                $.each(options.selected, function(index, v) {
                    var key=$.inArray(v, values);
                    var value=optionText[key];
                    if(value) {
                        content +='<div class="btn-group"><a class="btn btn-primary btn-xs" data="'+v+'">'+ value +'</a><button type="button" class="u-btn-close btn btn-primary btn-xs dropdown-toggle"><span></span></button></div>';
                    }
                }
                );
            }
            content +='</div>';
            $('#'+warpId).html(content);
            $('.u-btn-close').on('click',function (){
                var closeIndex=optionText.indexOf($(this).siblings('a').text());
                $(this).parent().remove();
                select.find('option').eq(closeIndex).removeAttr('selected');
                var _closeIndex=options.selected.indexOf(select.find('option').eq(closeIndex).val());
                options.selected.splice(_closeIndex,1);
                return false;
            })
        }
        ;
        var creat_content=function(data) {
            var toolbar='<div class="container-fluid"><div class="row" style="height:44px;" id="'+ dialogId +'">'+ '<div class="col-xs-8" style="padding:0;">'+ '<div class="btn-group btn-group-sm" data-id="btn-group" style="margin-bottom: 10px;">'+ '<button class="btn btn-primary" data-id="check_all">全选</button>'+ '<button class="btn btn-success" data-id="inverse">反选</button>'+ '</div>'+ '</div>'+ '<div class="col-xs-4" style="padding:0;">'+ '<input type="text" data-id="keyword" placeholder="搜索..." style="float: right;margin-top:  7px;margin-right: 6px;width:200px;">'+ '</div>'+ '</div></div>';
            var content=toolbar+'<div class="btn-toolbar" id='+ contentId +' style="width:700px;">';
            $.each(data, function(index, v) {
                var default_class='btn';
                if($.inArray(values[index], options.selected) > -1) {
                    default_class +=' btn-primary'
                }
                content +='<div class="btn-group btn-group-sm"><button class="'+ default_class +'" data-index="'+index+'" data-text="'+ v.toLowerCase() +'">'+ v +'</button></div>';
            }
            );
            content +='</div>';
            return content;
        }
        var filter=function(event) {
            var input=$(event.target), key=input.val().toLowerCase(), content=$('#'+contentId);
            if(key=='') {
                content.find('button').show();
            }
            else {
                var buttons=content.find('button:not([data-text*="'+ key +'"]):not(.btn-primary)');
                content.find('button').show(); //避免输入中文时的bug
                buttons.hide();
            }
        }
        this.each(function() {
            $(this).hide();
            if($('#'+warpId).length>0) {
                //如果存在，清空内容
                $('#'+warpId).html('').off('click');
            }
            else {
                $(this).after('<div class="form-control" id="'+warpId+'" style="min-height:34px;height:auto;cursor: pointer;"></div>');
            }
            creat_warp_content();
            $('#'+warpId).click(function() {
                open_dialog(options);
            }
            );
        }
        );
        return this;
    }
})(jQuery);
