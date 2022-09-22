/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// 结合 v-test-id 指令使用，辅助获取 data-test-id 的值
// 使用路由名时约定 key 为 menu-symbol 中定义的去掉 menu_ 部分的名称

/**
 * Usage:
 * v-test-id.index = "'ipSearch'"
 * Result:
 * data-test-id="index_button_ipSearch"
 */

export default {
  global: {
    header: 'header_top',
    headerNav: 'nav_headerNav',
    asideNav: 'nav_asideNav'
  },
  index: {
    ipSearch: 'button_ipSearch'
  },
  businessHostAndService: {
    hostList: 'list_hostTable',
    svrInstList: 'list_svrInstTable',
    processList: 'list_processTable',
    addProcess: 'button_addProcess',
    cloneProcess: 'button_cloneProcess',
    delProcess: 'button_delProcess'
  },
  businessServiceCategory: {
    btnConfirm: 'button_confirm',
    btnCancel: 'button_cancel'
  },
  businessServiceTemplate: {
    addForm: 'form_addForm',
    editForm: 'form_editForm',
    processForm: 'form_processForm',
    confirmSaveName: 'button_confirmSaveName',
    confirmSaveCategory: 'button_confirmSaveCategory',
    templateList: 'list_templateTable'
  },
  businessSetTemplate: {
    addForm: 'form_addForm',
    editForm: 'form_editForm',
    instanceTable: 'list_instanceTable',
    batchSync: 'button_batchSync',
    templateList: 'list_templateTable'
  }
}
