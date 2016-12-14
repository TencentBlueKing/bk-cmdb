<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<?php
    $this->load->library('login');
    $curUser = $this->login->getCurrentUser();
?>
<!-- 主面板Content -->
<link href="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.css" rel="stylesheet"/>
<div class="content-wrapper">
    <!-- 主面板 Header  -->
    <section class="content-header">
        <h1>用户管理</h1>
        <?php if($curUser['role'] == 'admin'){?>
        <div style="position:absolute;right:20px;top:20px;">
            <a class="btn btn-success addedButtom"><span class="fa fa-plus pr10"></span>新增用户</a>
        </div>
        <?php }?>
    </section>
    <!-- 主面板 Main-->
    <style>
        #grid_userManager tr th:nth-child(1){display: none;}
        #grid_userManager tr td:nth-child(1){display: none;}
        #grid_userManager thead tr th:nth-child(2) {background:none!important;}
        #grid_userManager thead tr th {border-bottom: 1px solid #ccc!important;}
        td{position: relative;}
        .copy-btn{position: absolute;top: 13px;right: 8px;visibility: hidden;z-index: 8888;}
        td:hover .copy-btn,.copy-btn.hover{display:inline;visibility: visible;}
        .c-userlist-cancel{display: none}
        .limit-width{min-width: 150px;}
    </style>
    <section class="content" id="ctn_userlist">
        <div id="grid_userManager_wrapper" class="dataTables_wrapper form-inline dt-bootstrap no-footer"><div class="row"><div class="col-sm-6"></div><div class="col-sm-6"></div></div><div class="row"><div class="col-sm-12"><table id="grid_userManager" class="table table-bordered table-striped dataTable no-footer" role="grid">
            <thead>
                <tr role="row"><th data-field="id" class="sorting_disabled" rowspan="1" colspan="1" style="width: 41px;">用户ID</th><th data-field="UserName" class="sorting_disabled" rowspan="1" colspan="1" style="width: 131px;">用户名</th><th data-field="ChName" class="sorting_disabled" rowspan="1" colspan="1" style="width: 140px;">姓名</th><th data-field="QQ" class="sorting_disabled" rowspan="1" colspan="1" style="width: 149px;">QQ</th><th data-field="Tel" class="sorting_disabled" rowspan="1" colspan="1" style="width: 182px;">联系电话</th><th data-field="Email" class="sorting_disabled" rowspan="1" colspan="1" style="width: 278px;">常用邮箱</th><th data-field="Role" class="sorting_disabled" rowspan="1" colspan="1" style="width: 103px;">角色</th><th data-field="Op" class="sorting_disabled" rowspan="1" colspan="1" style="width: 374px;"><div class="limit-width">操作</div></th></tr>
            </thead>
                 <tbody id="tbl_userlist_body">
                 <?php foreach ($users as $key=>$item) : ?>
                    <?php if($key%2){$class = 'even';}else{$class = 'odd';} ?>
                    <tr class="tr_lineedit odd" role="row">
                            <td><div class="id"><?php echo $item['id'];?><input type="hidden" value="<?php echo $item['id'];?>" class="hid_id" name="id"></div></td>
                            <td><div class="UserName"><?php echo $item['UserName'];?></div></td>
                            <td><div class="ChName"><?php echo $item['ChName'];?></div></td>
                            <td><div class="QQ"><?php echo $item['QQ'];?></div></td>
                            <td><div class="Tel"><?php echo $item['Tel'];?></div></td>
                            <td><div class="Email"><?php echo $item['Email'];?></div></td>
                            <td><div class="Role" role="<?php echo $item['Role'];?>"><?php if($item['Role'] == 'admin'){echo '管理员';}else{echo '用户';}?></div></td>
                            <td>
                                <a name="edits" class="btn btn-success btn-sm c-userlist-edit">编辑</a>
                                <a name="saves" style="display:none" class="btn btn-primary btn-sm c-userlist-save">保存</a>
                                <a name="deletes" class="btn btn-danger btn-sm c-userlist-delete" href="javascript:void(0);">删除</a>
                                <a name="cancels" class="btn btn-default btn-sm c-userlist-cancel" href="javascript:void(0);">取消</a>
                                <a name="resets" class="btn btn-default btn-sm c-userlist-reset" href="javascript:void(0);">重置</a>
                            </td>
                        </tr>
                 <?php endforeach;?>
                    </tbody>
        </table></div></div></div>
    </section>
</div>
<div class="control-sidebar-bg"></div>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.js"></script>

<script>
    var userList = {
        <?php
            $cnt = 0;
            foreach($users as $userInfo){
                echo $userInfo['id'] . ': {';
                $columnCnt = 0;
                foreach($userInfo as $columnKey => $value){
                    echo $columnKey . ": '" . $value . "'";
                    if(++$columnCnt < count($userInfo)){
                        echo ',';
                    }
                }
                echo '}';
                if(++$cnt < count($users)){
                    echo ',';
                }
            }
        ?>
    };
    var curUserRole = '<?php echo $curUser['role'];?>';
</script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/jquery.dataTables.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.min.js"></script>
<script src="<?php echo STATIC_URL; ?>/static/js/user.list.js?version=<?php echo $version;?>"></script>


