<template>
    <div class="module-selector-layout"
        v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <div class="wrapper">
            <div class="wrapper-column wrapper-left">
                <h2 class="title">{{$t('转移主机到其他业务')}}</h2>
                <bk-select class="business-selector"
                    v-model="targetBizId"
                    :clearable="false"
                    :searchable="true"
                    :placeholder="$t('请选择xx', { name: $t('业务') })"
                    @change="getModules">
                    <cmdb-auth-option v-for="business in businessList"
                        :key="business.bk_biz_id"
                        :id="business.bk_biz_id"
                        :name="`[${business.bk_biz_id}] ${business.bk_biz_name}`"
                        :auth="{ type: $OPERATION.HOST_TRANSFER_ACROSS_BIZ, relation: [bizId, business.bk_biz_id] }">
                    </cmdb-auth-option>
                </bk-select>
                <template v-if="targetBizId">
                    <bk-input class="tree-filter" clearable right-icon="icon-search" v-model="filter" :placeholder="$t('请输入关键词')"></bk-input>
                    <bk-big-tree ref="tree" class="topology-tree"
                        display-matched-node-descendants
                        :options="{
                            idKey: getNodeId,
                            nameKey: 'bk_inst_name',
                            childrenKey: 'child'
                        }"
                        :height="278"
                        :node-height="36"
                        :show-checkbox="isShowCheckbox"
                        @node-click="handleNodeClick"
                        @check-change="handleNodeCheck">
                        <template slot-scope="{ node, data }">
                            <i class="internal-node-icon fl"
                                v-if="data.default !== 0"
                                :class="getInternalNodeClass(node, data)">
                            </i>
                            <i v-else :class="['node-icon fl', { 'is-template': isTemplate(data) }]">{{data.bk_obj_name[0]}}</i>
                            <span class="node-name" :title="node.name">{{node.name}}</span>
                        </template>
                    </bk-big-tree>
                </template>
                <template v-else>
                    <bk-exception class="business-tips" type="empty" scene="part">
                        <span>{{$t('请先选择业务')}}</span>
                    </bk-exception>
                </template>
            </div>
            <div class="wrapper-column wrapper-right">
                <div class="selected-info clearfix">
                    <i18n class="selected-count fl" path="已选择N个模块">
                        <span class="count" place="count">{{checked.length}}</span>
                    </i18n>
                    <bk-button class="fr" text theme="primary"
                        v-show="checked.length"
                        @click="handleClearModule">
                        {{$t('清空')}}
                    </bk-button>
                </div>
                <ul class="module-list">
                    <li class="module-item" v-for="node in checked"
                        :key="node.id">
                        <div class="module-info clearfix">
                            <span class="info-icon fl">{{node.data.bk_obj_name[0]}}</span>
                            <span class="info-name" :title="node.data.bk_inst_name">
                                {{node.data.bk_inst_name}}
                            </span>
                        </div>
                        <div class="module-topology" :title="getNodePath(node)">{{getNodePath(node)}}</div>
                        <i class="bk-icon icon-close" @click="handleDeleteModule(node)"></i>
                    </li>
                </ul>
            </div>
        </div>
        <div class="layout-footer">
            <bk-button class="mr10" theme="primary"
                :disabled="!checked.length"
                @click="handleNextStep">
                {{$t('下一步')}}
            </bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import debounce from 'lodash.debounce'
    export default {
        name: 'cmdb-across-business-module-selector',
        props: {
            defaultChecked: {
                type: Array,
                default: () => ([])
            },
            defaultBizId: {
                type: Number,
                default: null
            }
        },
        data () {
            return {
                filter: '',
                businessList: [],
                targetBizId: '',
                handlerFilter: null,
                checked: [],
                request: {
                    internal: Symbol('internal'),
                    business: Symbol('business')
                },
                nodeIconMap: {
                    1: 'icon-cc-host-free-pool',
                    2: 'icon-cc-host-breakdown',
                    default: 'icon-cc-host-free-pool'
                }
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', ['getModelById']),
            ...mapGetters('objectBiz', ['bizId'])
        },
        watch: {
            filter () {
                this.handlerFilter()
            }
        },
        async created () {
            this.handlerFilter = debounce(() => {
                this.$refs.tree.filter(this.filter)
            }, 300)
            this.getFullAmountBusiness()
        },
        methods: {
            async getFullAmountBusiness () {
                try {
                    const data = await this.$http.get('biz/simplify?sort=bk_biz_id')
                    const availableBusiness = (data.info || []).filter(business => business.bk_biz_id !== this.bizId)
                    this.businessList = Object.freeze(availableBusiness)
                    this.targetBizId = this.defaultBizId || ''
                } catch (e) {
                    console.error(e)
                    this.businessList = []
                }
            },
            async getModules () {
                try {
                    this.checked = []
                    const [internalTopo, businessTopo] = await Promise.all([
                        this.getInternalModules(),
                        this.getBusinessModules()
                    ])
                    const [businessNode] = businessTopo
                    businessNode.child = [...internalTopo, ...(businessNode.child || [])]
                    this.$refs.tree.setData(Object.freeze(businessTopo))
                    this.$refs.tree.setExpanded(this.getNodeId(businessTopo[0]))
                    this.setDefaultChecked()
                } catch (e) {
                    this.$refs.tree.setData([])
                    console.error(e)
                }
            },
            setDefaultChecked () {
                this.$nextTick(() => {
                    this.defaultChecked.forEach(id => {
                        const node = this.$refs.tree.getNodeById(this.getNodeId({
                            bk_obj_id: 'module',
                            bk_inst_id: id
                        }))
                        if (node) {
                            this.checked.push(node)
                            this.$refs.tree.setChecked(node.id)
                        }
                    })
                })
            },
            getInternalModules () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.targetBizId,
                    config: {
                        requestId: this.request.internal
                    }
                }).then(data => {
                    return [{
                        bk_inst_id: data.bk_set_id,
                        bk_inst_name: data.bk_set_name,
                        bk_obj_id: 'set',
                        bk_obj_name: this.getModelById('set').bk_obj_name,
                        default: 0,
                        child: (data.module || []).map(module => ({
                            bk_inst_id: module.bk_module_id,
                            bk_inst_name: module.bk_module_name,
                            bk_obj_id: 'module',
                            bk_obj_name: this.getModelById('module').bk_obj_name,
                            default: module.default
                        }))
                    }]
                })
            },
            getBusinessModules () {
                return this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                    bizId: this.targetBizId,
                    config: {
                        requestId: this.request.instance
                    }
                })
            },
            getNodeId (data) {
                return `${data.bk_obj_id}-${data.bk_inst_id}`
            },
            getInternalNodeClass (node, data) {
                const clazz = []
                clazz.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)
                return clazz
            },
            getNodePath (node) {
                const parents = node.parents
                return parents.map(parent => parent.data.bk_inst_name).join(' / ')
            },
            handleNodeClick (node) {
                if (node.data.bk_obj_id !== 'module') {
                    return false
                }
                this.$refs.tree.setChecked(node.id, { checked: !node.checked, emitEvent: true })
            },
            handleNodeCheck (checked, currentNode) {
                const currentChecked = []
                const removeChecked = []
                checked.forEach(id => {
                    const node = this.$refs.tree.getNodeById(id)
                    if (node.data.default === currentNode.data.default) {
                        currentChecked.push(id)
                    } else {
                        removeChecked.push(id)
                    }
                })
                this.checked = currentChecked.map(id => this.$refs.tree.getNodeById(id))
                this.$refs.tree.setChecked(removeChecked, { checked: false })
            },
            handleDeleteModule (node) {
                this.checked = this.checked.filter(exist => exist !== node)
                this.$refs.tree.setChecked(node.id, { checked: false })
            },
            handleClearModule () {
                this.$refs.tree.setChecked(this.checked.map(node => node.id), { checked: false })
                this.checked = []
            },
            handleCancel () {
                this.$emit('cancel')
            },
            handleNextStep () {
                if (!this.checked.length) {
                    this.$warn('请选择模块')
                    return false
                }
                this.$emit('confirm', this.checked, this.targetBizId)
            },
            isTemplate (data) {
                return data.service_template_id || data.set_template_id
            },
            isShowCheckbox (data) {
                return data.bk_obj_id === 'module'
            }
        }
    }
</script>

<style lang="scss" scoped>
    $leftPadding: 0 12px 0 23px;
    .module-selector-layout {
        height: 460px;
        min-height: 300px;
        padding: 0 0 50px;
        position: relative;
        .layout-footer {
            position: sticky;
            bottom: 0;
            left: 0;
            width: 100%;
            height: 50px;
            padding: 8px 20px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
            text-align: right;
            font-size: 0;
            z-index: 100;
            .footer-tips {
                display: inline-block;
                vertical-align: middle;
            }
        }
    }
    .wrapper {
        display: flex;
        height: 100%;
        .wrapper-column {
            flex: 1;
        }
        .wrapper-left {
            max-width: 380px;
            border-right: 1px solid $borderColor;
            .title {
                margin-top: 15px;
                padding: $leftPadding;
                font-size: 24px;
                font-weight: normal;
                color: #444444;
                line-height:32px;
            }
            .business-selector {
                display: block;
                margin: 10px 12px 0 23px;
            }
            .business-tips {
                height: calc(100% - 90px);
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 12px;
                color: $textColor;
            }
            .tree-filter {
                display: block;
                width: auto;
                margin: 10px 12px 0 23px;
            }
        }
        .wrapper-right {
            padding: 57px 23px 0;
            .selected-info {
                font-size: 14px;
                line-height: 20px;
                color: $textColor;
                .count {
                    padding: 0 4px;
                    font-weight: bold;
                    color: #2DCB56;
                }
            }
        }
    }
    .topology-tree {
        width: 100%;
        padding: 0 0 0 10px;
        @include scrollbar;
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
            &.is-template {
                background-color: #97AED6;
            }
            &.is-selected {
                background-color: #3A84FF;
            }
        }
        .node-name {
            height: 36px;
            line-height: 36px;
            overflow: hidden;
            @include ellipsis;
        }
        .node-checkbox {
            width: 16px;
            height: 16px;
            margin: 10px 17px 0 10px;
            background: #FFF;
            border-radius: 50%;
            border: 1px solid #979BA5;
            &.is-checked {
                padding: 3px;
                border-color: $primaryColor;
                background-color: $primaryColor;
                background-clip: content-box;
            }
        }
        .internal-node-icon{
            width: 20px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            margin: 8px 4px 8px 0;
        }
    }

    .module-list {
        height: calc(100% - 35px);
        margin-top: 12px;
        @include scrollbar-y;
        .module-item {
            position: relative;
            margin-top: 12px;
            .icon-close {
                position: absolute;
                top: 6px;
                right: 0px;
                width: 28px;
                height: 28px;
                line-height: 28px;
                text-align: center;
                color: #D8D8D8;
                cursor: pointer;
                &:hover {
                    color: #979BA5;
                }
            }
        }
    }
    .module-info {
        .info-icon {
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
        }
        .info-name {
            display: block;
            width: 250px;
            padding-left: 10px;
            font-size:14px;
            color: $textColor;
            line-height:20px;
            @include ellipsis;
        }
    }
    .module-topology {
        width: 250px;
        padding-left: 30px;
        margin-top: 3px;
        font-size: 12px;
        color: #C4C6CC;
        @include ellipsis;
    }
</style>
