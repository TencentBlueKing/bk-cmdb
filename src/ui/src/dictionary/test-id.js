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
