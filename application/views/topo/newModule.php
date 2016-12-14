<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<div class="creat-container creat-module-container">
    <h4 class="c-conf-title pl15 pr15">新建模块</h4>
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
                <td id="newmodulebe">所属集群</td>
                <td id="newmodulegroupname" style="text-align: left;">
                </td>
            </tr>
            <tr>
                <td>名称</td>
                <td class="tl">
                    <input type="text" class="form-control c-gridinput" value="" id="newmoduleModuleName" placeholder="创建的模块名，必填" maxlength="10" style="width:90%">
                    <span class="c-gridinputmust">*</span>
                </td>
            </tr>
            <tr>
                <td>维护人</td>
                <td class="tl">
                    <select id="Operator" class="form-control c-gridinput">
                    </select>
                    <span class="c-gridinputmust">*</span>
                </td>
            </tr>
            <tr>
                <td>备份维护人</td>
                <td class="tl">
                    <select id="BakOperator" class="form-control c-gridinput">
                    </select>
                    <span class="c-gridinputmust">*</span>
                </td>
            </tr>
            </tbody>
        </table>
        <div class="text-center">
            <button class="btn btn-primary" id="newsavemodule">保存</button>
        </div>
    </div>
</div>