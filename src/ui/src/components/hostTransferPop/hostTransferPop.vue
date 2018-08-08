<template>
    <div class="transfer-pop" v-show="isShow">
        <div class="transfer-content" ref="drag" v-drag="'#drag'" v-bkloading="{isLoading: loading}">
            <div class="content-title" id="drag">
                <i class="icon icon-cc-shift mr5"></i>
                {{$t('Common[\'主机转移\']')}}
                <template v-if="hosts.length === 1">
                    {{hosts[0].host['bk_host_innerip']}}
                </template>
            </div>
            <div class="content-section clearfix" v-bkloading="{isLoading: $loading('transfer')}">
                <div class="section-left fl">
                    <div class="section-biz">
                        <label class="biz-label">{{$t('Common[\'业务\']')}}</label>
                        <v-application-selector class="biz-selector" :disabled="true" :selected.sync="bkBizId">
                        </v-application-selector>
                    </div>
                    <div class="section-tree">
                        <v-tree ref="topoTree"
                            :hideRoot="true"
                            :treeData="treeData"
                            :initNode="initNode"
                            @nodeClick="handleNodeClick"
                            @nodeToggle="handleNodeToggle"></v-tree>
                    </div>
                </div>
                <div class="section-right fr">
                    <h3 class="section-title">{{$t('Hosts["已选中模块"]')}}</h3>
                    <ul class="selected-list">
                        <li class="selected-item clearfix" v-for="(node, index) in selectedList" :key="node['bk_inst_id']">
                            <div class="module-info fl">
                                <span class="module-info-name">{{node['bk_inst_name']}}</span>
                                <span class="module-info-path">{{getModulePathLabel(node)}}</span>
                            </div>
                            <i class="bk-icon icon-close fr" @click="removeSelected(index)"></i>
                        </li>
                        <li class="selected-item disabled" v-for="(node, index) in otherBizNodes" :key="index">
                            <div class="module-info">
                                <span class="module-info-name">{{node['module']['bk_module_name']}}</span>
                                <span class="module-info-path">{{node['biz']['bk_biz_name']}}-{{node['set']['bk_set_name']}}</span>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="content-footer">
                <div class="button-group clearfix">
                    <template v-if="chooseId.length > 1 && !isNotModule">
                        <label class="transfer-type" for="increment">
                            <input type="radio" id="increment" v-model="isIncrement" :value="true">
                            <span class="transfer-description">{{$t('Hosts["增量更新"]')}}</span>
                        </label>
                        <label class="transfer-type" for="replacement">
                            <input type="radio" id="replacement" v-model="isIncrement" :value="false">
                            <span class="transfer-description">{{$t('Hosts["完全替换"]')}}</span>
                        </label>
                    </template>
                    <div class="fr">
                        <bk-button type="primary" v-show="selectedList.length" :loading="$loading('transfer')" @click="doTransfer">{{$t('Common[\'确认转移\']')}}</bk-button>
                        <button class="bk-button vice-btn" @click="cancel">{{$t('Common[\'取消\']')}}</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script type="text/javascript">
    import vApplicationSelector from '@/components/common/selector/application'
    import vTree from '@/components/tree/tree.v2'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            isShow: Boolean,
            chooseId: Array,
            hosts: Array
        },
        data () {
            return {
                bkBizId: '',
                treeData: {
                    'default': 0,
                    'bk_obj_id': 'root',
                    'bk_obj_name': 'root',
                    'bk_inst_id': 'root',
                    'bk_inst_name': 'root',
                    'isFolder': true,
                    'child': [{
                        'default': 0,
                        'bk_obj_id': 'source',
                        'bk_obj_name': this.$t('HostResourcePool[\'资源池\']'),
                        'bk_inst_id': 'source',
                        'bk_inst_name': this.$t('HostResourcePool[\'资源池\']'),
                        'isFolder': false
                    }]
                },
                initNode: {
                    level: 1,
                    open: true,
                    active: true,
                    'bk_inst_id': 'root'
                },
                activeNode: {},
                activeParentNode: {},
                selectedList: [],
                allowType: ['source', 'module'],
                isNotModule: true,
                otherBizNodes: [],
                isIncrement: true,
                loading: true
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            maxExpandedLevel () {
                return this.getLevel(this.treeData) - 1
            },
            allModulesNodes () {
                return this.getNodes(this.treeData, 'module')
            },
            allSetNodes () {
                return this.getNodes(this.treeData, 'set')
            },
            allBizNodes () {
                return this.getNodes(this.treeData, 'biz')
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow && this.bkBizId) {
                    this.init()
                } else {
                    this.initNode = {
                        level: 1,
                        open: true,
                        active: true,
                        'bk_inst_id': 'root'
                    }
                    this.$refs.drag.style = ''
                    this.treeData.child.splice(1)
                }
            },
            maxExpandedLevel (level) {
                let treeListEl = this.$refs.topoTree.$el
                if (level > 8) {
                    let width = treeListEl.getBoundingClientRect().width
                    treeListEl.style.minWidth = `${width + 40}px`
                } else {
                    treeListEl.style.minWidth = 'auto'
                }
            }
        },
        methods: {
            async init () {
                this.selectedList = []
                this.otherBizNodes = []
                await this.getTopoTree().then(() => {
                    this.loading = false
                })
                if (this.hosts.length === 1) {
                    this.setSelectedList()
                }
            },
            setSelectedList () {
                let selected = []
                let otherBizNodes = []
                const hostData = this.hosts[0]
                const moduleArr = hostData.module
                const setArr = hostData.set
                const bizArr = hostData.biz
                moduleArr.forEach(module => {
                    const targetSet = setArr.find(set => set['bk_biz_id'] === module['bk_biz_id'] && set['bk_set_id'] === module['bk_set_id'])
                    const targetBiz = bizArr.find(biz => biz['bk_biz_id'] === module['bk_biz_id'])
                    if (targetBiz['bk_biz_id'] === this.bkBizId) {
                        const node = this.allModulesNodes.find(moduleNode => moduleNode['bk_inst_id'] === module['bk_module_id'])
                        if (node) {
                            selected.push(node)
                        }
                    } else {
                        otherBizNodes.push({
                            module,
                            set: targetSet,
                            biz: targetBiz
                        })
                    }
                })
                this.selectedList = selected
                this.otherBizNodes = otherBizNodes
            },
            getNodes (node, type) {
                let result = []
                if (node['bk_obj_id'] === type) {
                    result.push(node)
                } else if (node.child && node.child.length) {
                    node.child.forEach(childNode => {
                        result = [...result, ...this.getNodes(childNode, type)]
                    })
                }
                return result
            },
            getModulePath (node) {
                return {
                    biz: this.allBizNodes.find(biz => biz['bk_inst_id'] === this.bkBizId),
                    set: this.allSetNodes.find(set => {
                        if (set.child && set.child.length) {
                            return set.child.find(module => module['bk_inst_id'] === node['bk_inst_id'])
                        }
                        return false
                    }),
                    module: node
                }
            },
            getModulePathLabel (node) {
                if (node['bk_inst_id'] === 'source') {
                    return this.$t('Common["主机资源池"]')
                } else {
                    const path = this.getModulePath(node)
                    if ([1, 2].includes(node.default)) {
                        return `${path['biz']['bk_inst_name']}-${path['module']['bk_inst_name']}`
                    }
                    return `${path['biz']['bk_inst_name']}-${path['set']['bk_inst_name']}`
                }
            },
            getLevel (node) {
                let level = node.level || 1
                if (node.isOpen && node.child && node.child.length) {
                    level = Math.max(level, Math.max.apply(null, node.child.map(childNode => this.getLevel(childNode))))
                }
                return level
            },
            handleNodeClick (activeNode, nodeOptions) {
                this.activeNode = activeNode
                this.activeParentNode = nodeOptions.parent
                this.checkNode(activeNode)
            },
            handleNodeToggle (isOpen, node, nodeOptions) {
                if (!node.child || !node.child.length) {
                    this.$set(node, 'isLoading', true)
                    this.$axios.get(`topo/inst/child/${this.bkSupplierAccount}/${node['bk_obj_id']}/${this.bkBizId}/${node['bk_inst_id']}`).then(res => {
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
            checkNode (node) {
                if (this.allowType.indexOf(node['bk_obj_id']) !== -1) {
                    if (node['default'] || node['bk_inst_id'] === 'source') {
                        if (this.selectedList.length && !this.selectedList[0]['default'] && this.selectedList[0]['bk_inst_id'] !== 'source') {
                            this.$bkInfo({
                                title: this.$t('Common[\'转移确认\']', {target: node['bk_inst_name']}),
                                confirmFn: () => {
                                    this.selectedList = [node]
                                    this.isNotModule = true
                                }
                            })
                        } else {
                            this.selectedList = [node]
                        }
                    } else {
                        if (this.selectedList.length && (this.selectedList[0]['default'] || this.selectedList[0]['bk_inst_id'] === 'source')) {
                            this.selectedList = []
                        }
                        let isExist = this.selectedList.find(selectedNode => {
                            return selectedNode['bk_obj_id'] === node['bk_obj_id'] && selectedNode['bk_inst_id'] === node['bk_inst_id']
                        })
                        if (!isExist) {
                            this.selectedList.push(node)
                        }
                        this.isNotModule = false
                    }
                }
            },
            doTransfer () {
                if (this.selectedList[0]['bk_obj_id'] === 'source') {
                    this.$axios.post('hosts/modules/resource', {
                        'bk_biz_id': this.bkBizId,
                        'bk_host_id': this.chooseId
                    }, {id: 'transfer'}).then(res => {
                        if (res.result) {
                            this.$emit('success', res)
                            this.$alertMsg(this.$t('Common[\'转移成功\']'), 'success')
                            this.cancel()
                        } else {
                            if (res.data && res.data['bk_host_id']) {
                                this.$alertMsg(`${res['bk_error_msg']} : ${res.data['bk_host_id']}`)
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        }
                    })
                } else {
                    let isIncrement
                    if (this.chooseId.length === 1) {
                        isIncrement = false
                    } else {
                        isIncrement = this.isNotModule ? false : this.isIncrement
                    }
                    let modulesDefault = this.selectedList[0]['default']
                    let transferType = {0: '', 1: 'idle', 2: 'fault'}
                    let url = `hosts/modules/${transferType[modulesDefault]}`
                    this.$axios.post(url, {
                        'bk_biz_id': this.bkBizId,
                        'bk_host_id': this.chooseId,
                        'bk_module_id': this.selectedList.map(node => {
                            return node['bk_inst_id']
                        }),
                        'is_increment': isIncrement
                    }, {id: 'transfer'}).then(res => {
                        if (res.result) {
                            this.$emit('success', res)
                            this.$alertMsg(this.$t('Common[\'转移成功\']'), 'success')
                            this.cancel()
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    }).catch(e => {
                        if (e.response && e.response.status === 403) {
                            this.$alertMsg(this.$t('Common[\'您没有主机转移的权限\']'))
                        }
                    })
                }
            },
            removeSelected (index) {
                this.selectedList.splice(index, 1)
            },
            getTopoInst () {
                return this.$axios.get(`topo/inst/${this.bkSupplierAccount}/${this.bkBizId}?level=-1`).then(res => {
                    return res
                })
            },
            getTopoInternal () {
                return this.$axios.get(`topo/internal/${this.bkSupplierAccount}/${this.bkBizId}`).then(res => {
                    return res
                })
            },
            getTopoTree () {
                return this.$Axios.all([this.getTopoInst(), this.getTopoInternal()]).then(this.$Axios.spread((instRes, internalRes) => {
                    if (instRes.result && internalRes.result) {
                        let internalModule = internalRes.data.module.map(module => {
                            return {
                                'default': module['bk_module_name'] === '空闲机' || module['bk_module_name'] === 'idle machine' ? 1 : 2,
                                'bk_obj_id': 'module',
                                'bk_obj_name': '模块',
                                'bk_inst_id': module['bk_module_id'],
                                'bk_inst_name': module['bk_module_name'],
                                'isFolder': false
                            }
                        })
                        instRes.data[0]['child'] = internalModule.concat(instRes.data[0]['child'])
                        this.treeData.child = [...this.treeData.child, ...instRes.data]
                        this.$nextTick(() => {
                            this.initNode = {
                                level: 2,
                                open: true,
                                active: false,
                                'bk_inst_id': instRes.data[0]['bk_inst_id']
                            }
                        })
                    } else {
                        this.$alertMsg(internalRes.result ? instRes.message : internalRes.message)
                    }
                }))
            },
            cancel () {
                this.selectedList = []
                this.otherBizNodes = []
                this.$emit('update:isShow', false)
            }
        },
        components: {
            vApplicationSelector,
            vTree
        }
    }
</script>
<style lang="scss" scoped>
    .transfer-pop{
        position: fixed;
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
        z-index: 2000;
        background: rgba(0, 0, 0, 0.6);
    }
    .transfer-content{
        width: 720px;
        height: 590px;
        background: #fff;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        border-radius: 2px;
    }
    .content-title{
        height: 50px;
        background: #f9f9f9;
        color: #333948;
        font-weight: bold;
        line-height: 50px;
        font-size: 14px;
        padding-left: 30px;
        border-bottom: 1px solid #e7e9ef;
        .icon{
            position: relative;
            top: -1px;
        }
    }
    .content-section{
        height: 478px;
        .section-left{
            height: 100%;
            width: 367px;
            border-right: 1px solid #e7e9ef;
        }
        .section-right {
            height: 100%;
            width: 353px;
            .section-title{
                padding: 0 0 0 25px;
                margin: 0;
                height: 64px;
                line-height: 63px;
                font-size: 14px;
                font-weight: normal;
                border-bottom: 1px solid #e7e9ef;
            }
        }
    }
    .content-footer{
        height: 60px;
        line-height: 60px;
        background: #f9f9f9;
        border-top: 1px solid #e7e9ef;
    }
    .section-biz{
        height: 64px;
        border-bottom: 1px solid #e7e9ef;
        font-size: 14px;
        padding: 14px 0;
        .biz-label{
            display: inline-block;
            vertical-align: middle;
            padding: 0 20px 0 30px;
        }
        .biz-selector{
            display: inline-block;
            vertical-align: center;
            width: 245px;
        }
    }
    .section-tree{
        height: 413px;
        padding: 10px 5px 0 0;
        overflow: auto;
        @include scrollbar;
    }
    .selected-list{
        height: calc(100% - 64px);
        color: #3c96ff;
        font-size: 12px;
        overflow-y: auto;
        @include scrollbar;
        .selected-item{
            height: 44px;
            line-height: 16px;
            padding: 6px 0 6px 25px;
            &:hover:not(.disabled){
                background-color: #e2efff;
                .module-info-name{
                    color: #498fe0;
                }
                .module-info-path{
                    color: #93bff5;
                }
            }
            .module-info{
                width: 250px;
            }
            .module-info-name{
                display: block;
                font-size: 14px;
                color: $textColor;
                @include ellipsis;
            }
            .module-info-path{
                display: block;
                font-size: 12px;
                color: #b9bdc1;
                @include ellipsis;
            }
            .icon-close{
                cursor: pointer;
                color: #9196a1;
                font-size: 15px;
                margin: 10px 25px 0 0;
                &:hover{
                    color: #3c96ff;
                }
            }
        }
    }
    .selected-count{
        padding: 0 0 0 32px;
        .color-info{
            color: #3c96ff;
        }
    }
    .button-group{
        padding: 0 20px 0 0;
        .transfer-type{
            line-height: normal;
            font-size: 14px;
            margin: 0 0 0 20px;
            .transfer-description,
            input[type="radio"]{
                display: inline-block;
                vertical-align: middle;
            }
        }
        .bk-button{
            margin: 0 5px;
        }
    }
</style>