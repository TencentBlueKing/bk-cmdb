<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<!-- 主面板Content -->
<link href="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.css" rel="stylesheet"/>
<link href="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.css" rel="stylesheet">
<style>
    td{position: relative;vertical-align: middle!important;}
    .copy-btn{position: absolute;top: 13px;right: 8px;visibility: hidden;z-index: 8888;}
    td:hover .copy-btn,.copy-btn.hover{display:inline;visibility: visible;}
</style>
<div class="content-wrapper">
    <!-- 主面板 Header  -->
    <section class="content-header">
        <h1>业务查询</h1>
        <a id="newapp" class="btn btn-success" href="/app/newapp" style="position:absolute;right:20px;top:20px;"><span class="fa fa-plus pr10"></span>新增业务</a>
    </section>
    <!-- 主面板 Main-->
    <section class="content" id="searcBusiness">
        <table id="table_demo" class="table table-bordered table-striped mt30">
            <thead>
                <tr>
                    <th><span id="appname">业务名称</span></th>
                    <th><span id='maintainers'>运维人员</span></th>
                    <th>创建人</th>
                    <th>创建时间</th>
                    <th>主机个数</th>
                    <th><div class='limit-width'>操作</div></th>
                </tr>
            </thead>
            <tbody>
            <?php foreach($app as $a): ?>
                <tr >
                    <td><div class="business_name" title="<?php echo $a['ApplicationName'] ?>"><?php echo $a['ApplicationName'] ?></div>
                        <input type="hidden" class="ApplicationID" value="<?php echo $a['ApplicationID'] ?>">
                        <input type="hidden" class="ApplicationName" value="<?php echo $a['ApplicationName'] ?>">
                    </td>
                    <td><div class="operaer" title="<?php echo $a['Maintainers']?>"><?php echo $a['Maintainers']?></div>
                        <input type="hidden" class="Maintainers" value="<?php echo $a['Maintainers'] ?>">
                        <span class='btn btn-default btn-xs copy-btn'>复制</span>
                    </td>
                    <td><?php echo $a['Creator']?></td>
                    <td><?php echo $a['CreateTime']?></td>
                    <td><?php echo $a['HostNum']?></td>
                    <td style="text-align: center;" role="gridcell">
                            <a name="edits" class="btn btn-success  btn-sm c-searchyw-edit" >编辑</a>
                            <a name="cancels"  class="btn btn-default btn-sm c-searchyw-cancel" href="javascript:void(0)">取消</a>
                            <a name="saves" style="display:none" class="btn btn-primary btn-sm c-searchyw-save"  >保存</a>
                            <a name="deletes"  class="btn btn-danger btn-sm c-searchyw-delete" href="javascript:void(0)">删除</a>
                    </td>
                </tr>
            <?php endforeach;?>

            </tbody>
        </table>
    </section>
</div>

<div class="control-sidebar-bg"></div>
</div>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/jquery.dataTables.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/jquery.zeroclipboard/jquery.zeroclipboard.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/js/app.index.js?version=<?php echo $version;?>"></script>
<script>
    $(document).ready(function() {
         //表格(DataTables)
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
        $('#table_demo').dataTable({
          paging: false, //隐藏分页
          ordering: false, //关闭排序
          info: false, //隐藏左下角分页信息
          searching: false, //关闭搜索
          lengthChange: false, //不允许用户改变表格每页显示的记录数
          pageLength : 5, //每页显示几条数据
          language: language //汉化
        });
        $("body")
            .on("copy", ".copy-btn", function(/* ClipboardEvent */ e) {
                e.clipboardData.clearData();
                e.clipboardData.setData("text/plain", $(this).prev().val());
                e.preventDefault();
            })
            .on('aftercopy','.copy-btn',function (e){
                var setDefault=$(this).get(0);
                if(e.success['text/plain']===true){
                    var d = dialog({
                        quickClose: true,/*点击空白处快速关闭*/
                        width: 150,
                        align:"top",
                        padding:6,
                        content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>复制成功</div>'
                    });
                }else{
                    var d = dialog({
                        quickClose: true,/* 点击空白处快速关闭*/
                        width: 150,
                        align:"top",
                        padding:6,
                        content: '<div class="c-dialogdiv2"><i class="c-dialogimg-failure"></i>复制失败</div>'
                    });
                }
                d.show(setDefault);
                setTimeout(function() {
                    d.close().remove();
                }, 1500);

            })
    });
</script>
