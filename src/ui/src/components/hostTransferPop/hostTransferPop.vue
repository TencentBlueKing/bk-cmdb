<template>
    <div class="transfer-pop" v-show="isShow">
        <div class="transfer-content">
            <div class="content-title">
                <i class="icon icon-cc-shift mr5"></i>
                {{$t('Common[\'主机转移\']')}}
            </div>
            <div class="content-section clearfix">
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
                    <ul class="selected-list">
                        <li class="selected-item clearfix" v-for="(node, index) in selectedList" :key="index">
                            <span>{{node['bk_inst_name']}}</span>
                            <i class="bk-icon icon-close fr" @click="removeSelected(index)"></i>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="content-footer clearfix">
                <i18n path="Common['已选N项']" tag="div" class="selected-count fl">
                    <span class="color-info" place="N">{{selectedList.length}}</span>
                </i18n>
                <div class="button-group fr">
                    <bk-button type="primary" v-if="isNotModule" v-show="selectedList.length" @click="doTransfer(false)">{{$t('Common[\'确认转移\']')}}</bk-button>
                    <template v-else v-show="selectedList.length">
                        <bk-button type="primary" @click="doTransfer(false)">{{$t('Common[\'覆盖\']')}}</bk-button>
                        <bk-button type="primary" @click="doTransfer(true)">{{$t('Common[\'更新\']')}}</bk-button>
                    </template>
                    <button class="bk-button vice-btn" @click="cancel">{{$t('Common[\'取消\']')}}</button>
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
            chooseId: Array
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
                isNotModule: true
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            maxExpandedLevel () {
                return this.getLevel(this.treeData) - 1
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
            init () {
                this.getTopoTree()
                this.selectedList = []
            },
            getLevel (node) {
                let level = node.level || 1
                if (node.isOpen && node.child && node.child.length) {
                    level = Math.max(level, Math.max.apply(null, node.child.map(childNode => this.getLevel(childNode))))
                }
                return level
            },
            handleNodeClick (activeNode, activeParentNode, activeRootNode) {
                this.activeNode = activeNode
                this.activeParentNode = activeParentNode
                this.checkNode(activeNode)
            },
            handleNodeToggle (isOpen, node, parentNode, rootNode, level, nodeId) {
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
                                }
                            })
                        } else {
                            this.selectedList = [node]
                        }
                        this.isNotModule = true
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
            doTransfer (type) {
                if (this.selectedList[0]['bk_obj_id'] === 'source') {
                    this.$axios.post('hosts/modules/resource', {
                        'bk_biz_id': this.bkBizId,
                        'bk_host_id': this.chooseId
                    }).then(res => {
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
                    let modulesDefault = this.selectedList[0]['default']
                    let transferType = {0: '', 1: 'idle', 2: 'fault'}
                    let url = `hosts/modules/${transferType[modulesDefault]}`
                    this.$axios.post(url, {
                        'bk_biz_id': this.bkBizId,
                        'bk_host_id': this.chooseId,
                        'bk_module_id': this.selectedList.map(node => {
                            return node['bk_inst_id']
                        }),
                        'is_increment': type
                    }).then(res => {
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
                return this.$axios.get(`topo/inst/${this.bkSupplierAccount}/${this.bkBizId}`).then(res => {
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
        width: 643px;
        height: 588px;
        background: #fff;
        position: absolute;
        top: 50%;
        left: 50%;
        margin-left: -321px;
        border-radius: 2px;
        margin-top: -294px;
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
            width: 273px;
            padding: 20px 20px 20px 30px;
            overflow: auto;
            @include scrollbar;
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
        color: #3c96ff;
        font-weight: bold;
        font-size: 12px;
        .selected-item{
            padding: 5px 0;
            .icon-close{
                cursor: pointer;
                color: #9196a1;
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
        .bk-button{
            margin: 0 5px;
        }
    }
</style>