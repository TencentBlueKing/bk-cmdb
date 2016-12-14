<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<div class="creat-container edit-group-container">
    <h4 class="c-conf-title pl15 pr15">集群(<span id="editset_property"></span>)属性</h4>
    <div class="c-attr-box c-conf-inner">
        <table class="table table-bordered">
            <thead>
            <tr class="active">
                <th>属性分组</th>
                <th>属性名</th>
                <th>属性值</th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td rowspan="4">基本属性</td>
                <td>名称</td>
                <td class="tl">
                    <input type="text" class="form-control c-gridinput" value="" id="editSetSetName" maxlength="10">
                    <span class="c-gridinputmust">*</span>
                </td>
            </tr>
            <tr>
                <td>环境类型</td>
                <td class="live-time" id ="edit_set_setenctype">
                    <div class="btn-group" data-toggle="buttons">
                        <label class="btn btn-radio ">
                            <input type="radio" name="options" id="edit_option1" autocomplete="off" value="1">测试
                        </label>
                        <label class="btn btn-radio ">
                            <input type="radio" name="options" id="edit_option2" autocomplete="off" value="2">体验
                        </label>
                        <label class="btn btn-radio">
                            <input type="radio" name="options" id="edit_option3" autocomplete="off" value="3">正式
                        </label>
                    </div>
                </td>
            </tr>
            <tr>
                <td>服务状态</td>
                <td class="live-time" id ="edit_set_sersta">
                    <input type="checkbox" name="servsta-checkbox" id="servstacheck" data-on-color="success" data-off-color="warning" checked>
<!--                    <div class="btn-group" data-toggle="buttons">-->
<!--                        <label class="btn btn-primary ">-->
<!--                            <input type="radio" name="options" id="edit_option4" autocomplete="off" value="开放">开放-->
<!--                        </label>-->
<!--                        <label class="btn btn-primary">-->
<!--                            <input type="radio" name="options" id="edit_option5" autocomplete="off" value="关闭">关闭-->
<!--                        </label>-->
<!--                    </div>-->
                </td>
            </tr>
            <tr>
                <td>中文名称</td>
                <td>
                    <input type="text" class="form-control" value="" id="editSetChnName" maxlength="32" style="width:90%">
                </td>
            </tr>
            <tr>
                <td rowspan="3">扩展属性</td>
                <td>设计容量</td>
                <td>
                    <input type="text" class="form-control" value="" id="editSetCapacity" maxlength="8" style="width:90%">
                </td>
            </tr>
            <tr>
                <td>描述</td>
                <td>
                    <input type="text" class="form-control" value="" id="editSetDes" maxlength="250" style="width:90%">
                </td>
            </tr>
            <tr>
                <td>Openstatus</td>
                <td>
                    <input type="text" class="form-control" value="" id="editOpenstatus" maxlength="16" style="width:90%">
                </td>
            </tr>
            </tbody>
        </table>
        <div class="text-center">
            <button class="btn btn-danger" id="editsetdelete">删除</button>
            <button class="btn btn-primary" id="editsetsave">保存</button>
        </div>
    </div>
</div>