<template>
    <div class="topo-wrapper clearfix">
        <div class="topo-tree-ctn fl">
            <div class="biz-selector-ctn">
                <v-application-selector :selected.sync="tree.bkBizId" @on-selected="handleBizSelected" :filterable="true"></v-application-selector>
            </div>
            <div class="topo-options-ctn" hidden>
                <i class="topo-option-del icon-cc-del fr" v-if="isShowOptionDel && Object.keys(tree.treeData).length" @click="deleteNode"></i>
            </div>
            <div class="tree-list-ctn">
                <v-tree ref="topoTree" 
                    :treeData="tree.treeData" 
                    :initNode="tree.initNode" 
                    :model="tree.model" 
                    @nodeClick="handleNodeClick" 
                    @nodeToggle="handleNodeToggle"
                    @addNode="handleAddNode">
                </v-tree>
            </div>
        </div>
        <div class="topo-view-ctn">
            <bk-tab :active-name="view.tab.active" @tab-changed="tabChanged" class="topo-view-tab">
                <bk-tabpanel name="host" :title="$t('BusinessTopology[\'主机调配\']')">
                    <v-hosts ref="hosts"
                        :outerParams="searchParams"
                        :isShowRefresh="true"
                        :outerLoading="tree.loading"
                        :isShowCrossImport="authority['is_host_cross_biz'] && attributeBkObjId === 'module'"
                        :tableVisible="view.tab.active === 'host'"
                        :wrapperMinusHeight="210"
                        @handleCrossImport="handleCrossImport">
                        <div slot="filter"></div>
                    </v-hosts>
                </bk-tabpanel>
                <bk-tabpanel name="attribute" :title="$t('BusinessTopology[\'节点属性\']')" :show="isShowAttribute">
                    <v-attribute ref="topoAttribute"
                        :bkObjId="attributeBkObjId" 
                        :bkBizId="tree.bkBizId" 
                        :activeNode="tree.activeNode"
                        :activeParentNode="tree.activeParentNode"
                        :formValues="view.attribute.formValues"
                        :type="view.attribute.type"
                        :active="view.tab.active === 'attribute'"
                        :isLoading="view.attribute.isLoading"
                        :editable="tree.bkBizName !== '蓝鲸'"
                        @submit="submitNode"
                        @delete="deleteNode"
                        @cancel="cancelCreate"></v-attribute>
                </bk-tabpanel>
                <bk-tabpanel name="process" :title="$t('ProcessManagement[\'进程信息\']')" 
                    :show="attributeBkObjId === 'module' && ![1,2].includes(tree.activeNode.default)">
                    <v-process
                        :isShow="view.tab.active === 'process'"
                        :bizId="tree.bkBizId"
                        :moduleName="attributeBkObjId === 'module' ? tree.activeNode['bk_inst_name'] : ''"
                    >
                    </v-process>
                </bk-tabpanel>
            </bk-tab>
        </div>
        <bk-dialog :is-show.sync="view.crossImport.isShow" :quick-close="false" :has-header="false" :has-footer="false" :width="700" :padding="0">
            <v-cross-import  slot="content"
                :is-show.sync="view.crossImport.isShow"
                :bizId="tree.bkBizId"
                :moduleId="tree.activeNode['bk_inst_id']"
                @handleCrossImportSuccess="setSearchParams">
            </v-cross-import>
        </bk-dialog>
    </div>
</template>
<script>
    import vApplicationSelector from '@/components/common/selector/application'
    import vTree from '@/components/tree/tree.v2'
    import vHosts from '@/pages/hosts/hosts'
    import vAttribute from './children/attribute'
    import vProcess from './children/process'
    import vCrossImport from './children/crossImport'
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                tree: {
                    bkBizId: -1,
                    bkBizName: '',
                    treeData: {},
                    model: [],
                    activeNode: {},
                    activeNodeOptions: {},
                    activeParentNode: {},
                    initNode: {},
                    loading: false
                },
                view: {
                    tab: {
                        active: 'host'
                    },
                    attribute: {
                        type: 'update',
                        formValues: {},
                        isLoading: true
                    },
                    crossImport: {
                        isShow: false
                    }
                },
                nodeToggleRecord: {},
                searchParams: null
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            ...mapGetters('navigation', ['authority']),
            /* 获取当前属性表单对应的属性obj_id */
            attributeBkObjId () {
                let bkObjId
                if (this.view.attribute.type === 'create') {
                    bkObjId = this.tree.model.find(model => {
                        return model['bk_obj_id'] === this.tree.activeNode['bk_obj_id']
                    })['bk_next_obj']
                } else {
                    bkObjId = this.tree.activeNode['bk_obj_id']
                }
                return bkObjId
            },
            /* 计算是否显示属性修改tab选项卡 */
            isShowAttribute () {
                let isShow = this.tree.activeNode['bk_obj_id'] !== 'biz' && !this.tree.activeNode['default']
                return (isShow || this.view.attribute.type === 'create') && !!Object.keys(this.tree.treeData).length
            },
            /* 计算是否显示添加按钮 */
            isShowOptionAdd () {
                let activeNode = this.tree.activeNode
                return this.tree.model.length && Object.keys(activeNode).length && activeNode['bk_obj_id'] !== 'module' && !activeNode['default']
            },
            /* 查找当前节点对应的主线拓扑模型 */
            optionModel () {
                return this.tree.model.find(model => {
                    return model['bk_obj_id'] === this.tree.activeNode['bk_obj_id']
                })
            },
            /* 计算是否显示删除按钮 */
            isShowOptionDel () {
                return this.tree.activeNode['bk_obj_id'] !== 'biz' && !this.tree.activeNode['default'] && this.view.tab.active === 'host'
            },
            /* 计算当前树展开的最大层次 */
            maxExpandedLevel () {
                return this.getLevel(this.tree.treeData) - 1
            }
        },
        watch: {
            /* 业务切换，初始化拓扑树 */
            'tree.bkBizId' (bkBizId) {
                this.getTopoTree().then(() => {
                    this.tree.initNode = {
                        level: 1,
                        open: true,
                        active: true,
                        bk_inst_id: this.tree.treeData['bk_inst_id']
                    }
                })
            },
            /* 当前节点发生变化且属性修改面板激活时，加载当前节点的具体属性 */
            'tree.activeNode' () {
                if (!this.isShowAttribute || (this.attributeBkObjId !== 'module' && this.view.tab.active === 'process')) {
                    this.tabChanged('host')
                }
                if (this.view.tab.active === 'attribute') {
                    this.getNodeDetails()
                }
            },
            /* tab选项卡处于切换到属性面板时，加载节点具体属性 */
            'view.tab.active' (activeTab) {
                if (activeTab === 'attribute' && this.view.attribute.type === 'update') {
                    this.getNodeDetails()
                }
            },
            /* 根据当前树节点展开的层级设置横向宽度 */
            maxExpandedLevel (level) {
                let extendLength = 40
                let treeListEl = this.$refs.topoTree.$el
                if (level > 4) {
                    let width = treeListEl.getBoundingClientRect().width
                    treeListEl.style.minWidth = `${width + extendLength}px`
                } else {
                    treeListEl.style.minWidth = 'auto'
                }
            }
        },
        methods: {
            handleBizSelected (data) {
                this.tree.bkBizName = data.label
            },
            /* 获取最大展开层级 */
            getLevel (node) {
                let level = node.level
                if (node.isOpen && node.child && node.child.length) {
                    level = Math.max(level, Math.max.apply(null, node.child.map(childNode => this.getLevel(childNode))))
                }
                return level
            },
            /* 获取业务拓扑实例 */
            getTopoInst () {
                return this.$axios.get(`topo/inst/${this.bkSupplierAccount}/${this.tree.bkBizId}`).then(res => {
                    return res
                }).catch(e => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'您没有当前业务的权限\']'))
                    }
                })
            },
            /* 获取内置业务拓扑 */
            getTopoInternal () {
                return this.$axios.get(`topo/internal/${this.bkSupplierAccount}/${this.tree.bkBizId}`).then(res => {
                    return res
                }).catch(e => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'您没有当前业务的权限\']'))
                    }
                })
            },
            /* 获取主线拓扑模型 */
            getTopoModel () {
                this.$axios.get(`topo/model/${this.bkSupplierAccount}`).then(res => {
                    if (res.result) {
                        this.tree.model = res.data
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /* 初始化拓扑树 */
            getTopoTree () {
                this.tree.loading = true
                return this.$Axios.all([this.getTopoInst(), this.getTopoInternal()]).then(this.$Axios.spread((instRes, internalRes) => {
                    if (instRes.result && internalRes.result) {
                        let internalModule = internalRes.data.module.map(module => {
                            return {
                                'default': module['bk_module_name'] === '空闲机' || module['bk_module_name'] === 'idle machine' ? 1 : 2,
                                'bk_obj_id': 'module',
                                'bk_obj_name': this.$t('Hosts[\'模块\']'),
                                'bk_inst_id': module['bk_module_id'],
                                'bk_inst_name': module['bk_module_name'],
                                'isFolder': false
                            }
                        })
                        instRes.data[0]['child'] = internalModule.concat(instRes.data[0]['child'])
                        this.tree.treeData = instRes.data[0]
                    } else {
                        this.$alertMsg(internalRes.result ? instRes.message : internalRes.message)
                    }
                })).then(() => {
                    this.tree.loading = false
                }).catch(() => {
                    this.tree.loading = false
                })
            },
            /* 获取当前节点的具体属性 */
            getNodeDetails () {
                let {
                    bk_inst_id: bkInstId,
                    bk_inst_name: bkInstName,
                    bk_obj_id: bkObjId
                } = this.tree.activeNode
                let url
                let params = {
                    page: {sort: 'id'},
                    fields: [],
                    condition: {}
                }
                if (bkObjId === 'set') {
                    url = `set/search/${this.bkSupplierAccount}/${this.tree.bkBizId}`
                    params['condition']['bk_set_id'] = bkInstId
                } else if (bkObjId === 'module') {
                    url = `module/search/${this.bkSupplierAccount}/${this.tree.bkBizId}/${this.tree.activeParentNode['bk_inst_id']}`
                    params['condition']['bk_module_id'] = bkInstId
                    params['condition']['bk_supplier_account'] = this.bkSupplierAccount
                } else {
                    url = `inst/search/${this.bkSupplierAccount}/${bkObjId}/${bkInstId}`
                }
                this.view.attribute.isLoading = true
                this.$axios.post(url, params).then(res => {
                    if (res.result) {
                        this.view.attribute.formValues = res.data.info[0]
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.view.attribute.isLoading = false
                }).catch(() => {
                    this.view.attribute.isLoading = false
                })
            },
            /* 新增拓扑，切换到属性表单 */
            handleAddNode () {
                this.view.attribute.formValues = {}
                this.view.attribute.isLoading = false
                this.view.attribute.type = 'create'
                this.view.tab.active = 'attribute'
            },
            /* 新增拓扑节点/修改拓扑节点 */
            submitNode (formData, originalData) {
                let url
                let method
                let submitType = this.view.attribute.type
                let {
                    bk_inst_id: bkInstId,
                    bk_obj_id: bkObjId
                } = this.tree.activeNode
                if (submitType === 'create') {
                    method = 'post'
                    formData['bk_parent_id'] = bkInstId
                    if (this.attributeBkObjId === 'set') {
                        url = `set/${this.tree.bkBizId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else if (this.attributeBkObjId === 'module') {
                        url = `module/${this.tree.bkBizId}/${bkInstId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else {
                        url = `inst/${this.bkSupplierAccount}/${this.attributeBkObjId}`
                        formData['bk_biz_id'] = this.tree.bkBizId
                    }
                } else if (submitType === 'update') {
                    method = 'put'
                    if (bkObjId === 'set') {
                        url = `set/${this.tree.bkBizId}/${bkInstId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else if (bkObjId === 'module') {
                        url = `module/${this.tree.bkBizId}/${this.tree.activeParentNode['bk_inst_id']}/${bkInstId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else {
                        url = `inst/${this.bkSupplierAccount}/${bkObjId}/${bkInstId}`
                    }
                }
                this.$axios({
                    url: url,
                    method: method,
                    data: formData,
                    id: 'editAttr'
                }).then(res => {
                    if (res.result) {
                        this.updateTopoTree(this.view.attribute.type, res.data, formData)
                        this.$alertMsg(submitType === 'create' ? this.$t('Common[\'新建成功\']') : this.$t('Common[\'修改成功\']'), 'success')
                        if (this.view.attribute.type === 'create') {
                            this.view.tab.active = 'host'
                        } else {
                            this.getNodeDetails()
                        }
                        this.view.attribute.type = 'update'
                        this.$refs.topoAttribute.displayType = 'list'
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /* 新增、修改拓扑节点成功后更新拓扑树 */
            updateTopoTree (type, response, formData) {
                let node = this.tree.activeNode
                let {
                    bk_next_obj: bkNextObj,
                    bk_next_name: bkNextName,
                    bk_obj_id: bkObjId,
                    bk_obj_name: bkObjName
                } = this.optionModel
                if (type === 'create') {
                    if (node.hasOwnProperty('isFolder')) {
                        node['isFolder'] = true
                    } else {
                        this.$set(node, 'isFolder', true)
                    }
                    node.child = node.child || []
                    node.child.push({
                        'default': 0,
                        'bk_inst_id': bkNextObj === 'set' ? response['bk_set_id'] : bkNextObj === 'module' ? response['bk_module_id'] : response['bk_inst_id'],
                        'bk_inst_name': bkNextObj === 'set' ? formData['bk_set_name'] : bkNextObj === 'module' ? formData['bk_module_name'] : formData['bk_inst_name'],
                        'bk_obj_id': bkNextObj,
                        'bk_obj_name': bkNextName,
                        'child': [],
                        'isFolder': false
                    })
                } else if (type === 'update') {
                    node['bk_inst_name'] = bkObjId === 'set' ? formData['bk_set_name'] : bkObjId === 'module' ? formData['bk_module_name'] : formData['bk_inst_name']
                }
            },
            /* 删除拓扑节点 */
            deleteNode () {
                this.$bkInfo({
                    title: `${this.$t('Common[\'确定删除\']')} ${this.tree.activeNode['bk_inst_name']}?`,
                    content: this.tree.activeNode['bk_obj_id'] === 'module'
                        ? this.$t('Common["请先转移其下所有的主机"]')
                        : this.$t('Common[\'下属层级都会被删除，请先转移其下所有的主机\']'),
                    confirmFn: () => {
                        let url
                        let {
                            bk_obj_id: bkObjId,
                            bk_inst_id: bkInstId
                        } = this.tree.activeNode
                        if (bkObjId === 'set') {
                            url = `set/${this.tree.bkBizId}/${bkInstId}`
                        } else if (bkObjId === 'module') {
                            url = `module/${this.tree.bkBizId}/${this.tree.activeParentNode['bk_inst_id']}/${bkInstId}`
                        } else {
                            url = `inst/${this.bkSupplierAccount}/${bkObjId}/${bkInstId}`
                        }
                        this.$axios.delete(url, {id: 'deleteAttr'}).then(res => {
                            if (res.result) {
                                this.view.tab.active = 'host'
                                this.tree.activeParentNode.child.splice(this.tree.activeNodeOptions.index, 1)
                                this.tree.initNode = {
                                    level: 1,
                                    open: true,
                                    active: true,
                                    bk_inst_id: this.tree.treeData['bk_inst_id']
                                }
                                this.$alertMsg(this.$t('Common[\'删除成功\']'), 'success')
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            /* 点击节点，设置查询参数 */
            handleNodeClick (activeNode, nodeOptions) {
                this.$refs.hosts.clearChooseId()
                this.tree.activeNode = activeNode
                this.tree.activeNodeOptions = nodeOptions
                this.tree.activeParentNode = nodeOptions.parent
                this.view.attribute.type = 'update'
                this.setSearchParams()
            },
            /* node节点展开时，判断是否加载下级节点 */
            handleNodeToggle (isOpen, node, nodeOptions) {
                if (!node.child || !node.child.length) {
                    this.$set(node, 'isLoading', true)
                    this.$axios.get(`topo/inst/child/${this.bkSupplierAccount}/${node['bk_obj_id']}/${this.tree.bkBizId}/${node['bk_inst_id']}`).then(res => {
                        if (res.result) {
                            let child = res['data'][0]['child']
                            if (Array.isArray(child) && child.length) {
                                node.child = child
                            } else {
                                this.$set(node, 'isFolder', false)
                            }
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                        node.isLoading = false
                    })
                }
            },
            /* 设置主机查询参数 */
            setSearchParams () {
                let params = {
                    'bk_biz_id': this.tree.bkBizId,
                    condition: []
                }
                let activeNodeObjId = this.tree.activeNode['bk_obj_id']
                if (activeNodeObjId === 'module' || activeNodeObjId === 'set') {
                    params.condition.push({
                        'bk_obj_id': activeNodeObjId,
                        fields: [],
                        condition: [{
                            field: activeNodeObjId === 'module' ? 'bk_module_id' : 'bk_set_id',
                            operator: '$eq',
                            value: this.tree.activeNode['bk_inst_id']
                        }]
                    })
                } else if (activeNodeObjId !== 'biz') {
                    params.condition.push({
                        'bk_obj_id': 'object',
                        fields: [],
                        condition: [{
                            field: 'bk_inst_id',
                            operator: '$eq',
                            value: this.tree.activeNode['bk_inst_id']
                        }]
                    })
                }
                let defaultObj = ['host', 'module', 'set', 'biz']
                defaultObj.forEach(id => {
                    if (!params.condition.some(({bk_obj_id: bkObjId}) => bkObjId === id)) {
                        params.condition.push({
                            'bk_obj_id': id,
                            fields: [],
                            condition: []
                        })
                    }
                })
                this.searchParams = params
            },
            handleCrossImport () {
                this.view.crossImport.isShow = true
            },
            tabChanged (active) {
                this.view.tab.active = active
                this.view.attribute.formValues = {}
                if (active === 'host') {
                    this.view.attribute.type = 'update'
                }
            },
            cancelCreate () {
                this.tabChanged('host')
            }
        },
        created () {
            this.getTopoModel()
        },
        components: {
            vApplicationSelector,
            vTree,
            vCrossImport,
            vHosts,
            vAttribute,
            vProcess
        }
    }
</script>
<style lang="scss" scoped>
    .topo-wrapper{
        height: 100%;
        .topo-tree-ctn{
            width: 280px;
            height: 100%;
            border-right: 1px solid #e7e9ef;
            background: #fafbfd;
        }
        .topo-view-ctn{
            height: 100%;
            overflow: hidden;
            padding: 0 20px;
        }
    }
    .biz-selector-ctn{
        padding: 20px;
    }
    .topo-options-ctn{
        height: 44px;
        line-height: 44px;
        background: #f9f9f9;
        padding: 0 10px;
        .topo-option-add{
            width: 90px;
            height: 24px;
            font-size: 12px;
            line-height: 22px;
            border-radius: 2px;
            background: #fff;
            padding: 0 5px;
            margin: 0 5px;
            @include ellipsis;
        }
        .topo-option-del{
            font-size: 12px;
            margin: 16px 0 0 0;
            cursor: pointer;
        }
    }
    .topo-view-tab{
        height: 100%;
        border: none;
    }
    .tree-list-ctn{
        padding: 0 0 0 20px;
        max-height: calc(100% - 120px);
        overflow: auto;
        @include scrollbar;
    }
</style>
<style lang="scss">
    .bk-tab2.topo-view-tab{
        .bk-tab2-head{
            height: 80px;
            .tab2-nav-item{
                height: 79px;
                line-height: 79px;
            }
        }
        .bk-tab2-content{
            height: calc(100% - 80px);
            section.active{
                height: 100%;
            }
        }
    }
    .topo-wrapper .hosts-wrapper{
        .table-container{
            padding: 20px 0;
        }
        .breadcrumbs{
            display: none;
        }
    }
</style>
