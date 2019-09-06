<template>
    <bk-sideslider
        v-transfer-dom
        :is-show.sync="isShow"
        :width="600"
        :title="$t('拓扑显示设置')"
        @close="handleClose">
        <div class="config-layout" slot="content" v-if="isShow">
            <div class="config-section">
                <h2 class="config-type">{{$t('模型显示设置')}}</h2>
                <div class="config-name">
                    <cmdb-form-bool class="config-name-checkbox"
                        v-model="config.label.node"
                        :size="14">
                        {{$t('显示模型名称')}}
                    </cmdb-form-bool>
                    <cmdb-form-bool class="config-name-checkbox"
                        v-model="config.label.edge"
                        :size="14">
                        {{$t('显示关联名称')}}
                    </cmdb-form-bool>
                </div>
            </div>
            <div class="config-section">
                <h2 class="config-type">{{$t('关系显示设置')}}</h2>
                <ul class="association-group">
                    <li class="group-item"
                        v-for="(group, index) in associationGroups"
                        :key="index">
                        <p class="group-name">{{group.name}}</p>
                        <ul class="model-list clearfix">
                            <li class="model-item fl"
                                v-for="(model, index) in group.models"
                                :key="index">
                                <cmdb-form-bool class="model-checkbox"
                                    :size="14"
                                    :checked="isModelAllChecked(model)"
                                    :indeterminate="isModelPartialChecked(model)"
                                    @click.native.stop
                                    @click="handleCheckModelAssociation(model)">
                                </cmdb-form-bool>
                                <span class="model-collapse-trigger" @click.stop="toggleCollapseStatus(model)">
                                    <span class="model-name">{{model.name}}</span>
                                    <span class="model-association-count"
                                        :class="{
                                            'has-angle': collapseStatus[model.id]
                                        }">
                                        ({{model.count}})
                                    </span>
                                    <i class="bk-icon icon-angle-down"></i>
                                </span>
                                <cmdb-collapse-transition>
                                    <div class="model-association"
                                        v-if="collapseStatus[model.id]"
                                        v-click-outside="hideAssociation">
                                        <cmdb-form-bool class="association-checkbox"
                                            v-for="(association, index) in associations[model.id]"
                                            :checked="edgeConfig[association['bk_inst_id']]"
                                            :key="association['bk_inst_id']"
                                            :size="14"
                                            @click="handleCheckAssociiation(association)">
                                            <span class="association-desc">{{getAssociationDesc(association)}}</span>
                                        </cmdb-form-bool>
                                    </div>
                                </cmdb-collapse-transition>
                            </li>
                        </ul>
                    </li>
                </ul>
            </div>
            <div class="config-section button-section">
                <bk-button theme="primary" @click="handleConfirm">
                    {{$t('确定')}}
                </bk-button>
                <bk-button theme="default" @click="handleReset">
                    {{$t('重置')}}
                </bk-button>
            </div>
        </div>
    </bk-sideslider>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { color } from './graphics.js'
    export default {
        name: 'cmdb-graphics-config',
        inject: ['parentRefs'],
        data () {
            return {
                isShow: false,
                isQuickClose: true,
                config: {
                    label: {
                        node: true,
                        edge: true
                    },
                    edge: {}
                },
                backupConfig: null,
                collapseStatus: {},
                associationStatus: {}
            }
        },
        computed: {
            ...mapGetters('globalModels', ['topologyData']),
            ...mapGetters('objectAssociation', ['associationList']),
            ...mapGetters('objectModelClassify', ['classifications', 'models']),
            associationGroups () {
                const groups = []
                const modelIdKey = 'bk_obj_id'
                const modelNameKey = 'bk_obj_name'
                const classificationIdKey = 'bk_classification_id'
                const classificationNameKey = 'bk_classification_name'
                this.topologyData.forEach(data => {
                    if ((data.assts || []).length) {
                        const model = this.models.find(model => model[modelIdKey] === data[modelIdKey])
                        const classification = this.classifications.find(classification => classification[classificationIdKey] === model[classificationIdKey])
                        const group = groups.find(group => group.id === classification[classificationIdKey])
                        const modelInfo = {
                            id: model[modelIdKey],
                            name: model[modelNameKey],
                            count: data.assts.length
                        }
                        if (group) {
                            group.models.push(modelInfo)
                        } else {
                            groups.push({
                                id: classification[classificationIdKey],
                                name: classification[classificationNameKey],
                                models: [modelInfo]
                            })
                        }
                    }
                })
                return groups
            },
            associations () {
                const associations = {}
                this.topologyData.forEach(data => {
                    associations[data['bk_obj_id']] = data.assts || []
                })
                return associations
            },
            edgeConfig () {
                const edgeConfig = {}
                const localConfig = this.config.edge
                const idKey = 'bk_inst_id'
                this.topologyData.forEach(data => {
                    (data.assts || []).forEach(association => {
                        const id = association[idKey]
                        edgeConfig[id] = typeof localConfig[id] === 'undefined' ? true : localConfig[id]
                    })
                })
                return edgeConfig
            },
            options () {
                const { node, edge } = this.config.label
                return {
                    nodes: {
                        font: {
                            color: node ? color.node.label : 'transparent'
                        }
                    },
                    edges: {
                        font: {
                            color: edge ? color.edge.label : 'transparent'
                        }
                    }
                }
            },
            edgeOptions () {
                return Object.keys(this.edgeConfig).map(id => {
                    return {
                        id,
                        hidden: !this.edgeConfig[id]
                    }
                })
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.isQuickClose = true
                    this.backupConfig = this.$tools.clone(this.config)
                }
            }
        },
        methods: {
            toggleSlider () {
                this.isShow = !this.isShow
            },
            handleClose () {
                if (this.isQuickClose) {
                    this.config = this.backupConfig
                }
            },
            toggleCollapseStatus ({ id }) {
                const previousState = this.collapseStatus[id]
                this.hideAssociation()
                this.$set(this.collapseStatus, id, !previousState)
            },
            getAssociationDesc (association) {
                const associationData = this.associationList.find(data => data.id === association['bk_asst_inst_id'])
                const associationModel = this.models.find(data => data['bk_obj_id'] === association['bk_obj_id'])
                if (associationData['bk_asst_name']) {
                    return `${associationData['bk_asst_name']} -> ${associationModel['bk_obj_name']}`
                }
                return `${associationData['bk_asst_id']} -> ${associationModel['bk_obj_name']}`
            },
            hideAssociation () {
                Object.keys(this.collapseStatus).forEach(modelId => {
                    this.collapseStatus[modelId] = false
                })
            },
            handleCheckModelAssociation (model) {
                const id = model.id
                const isAllChecked = this.isModelAllChecked(model)
                const associations = this.associations[id]
                associations.forEach(association => {
                    this.$set(this.config.edge, association['bk_inst_id'], !isAllChecked)
                })
            },
            handleCheckAssociiation (association) {
                const id = association['bk_inst_id']
                const previousState = this.config.edge[id]
                this.$set(this.config.edge, id, typeof previousState === 'undefined' ? false : !previousState)
            },
            isModelAllChecked ({ id }) {
                const associations = this.associations[id]
                const checked = associations.filter(association => this.edgeConfig[association['bk_inst_id']])
                return checked.length === associations.length
            },
            isModelPartialChecked ({ id }) {
                const associations = this.associations[id]
                const checked = associations.filter(association => this.edgeConfig[association['bk_inst_id']])
                return checked.length !== 0 && checked.length !== associations.length
            },
            handleConfirm () {
                this.$store.commit('globalModels/setOptions', this.options)
                this.$store.commit('globalModels/setEdgeOptions', this.edgeOptions)
                this.isQuickClose = false
                this.toggleSlider()
            },
            handleReset () {
                this.config.label = {
                    node: true,
                    edge: true
                }
                this.config.edge = {}
            }
        }
    }
</script>

<style lang="scss" scoped>
    .config-layout {
        padding: 10px 0 0;
        max-height: 100%;
        @include scrollbar-y;
        .config-section {
            padding: 0 40px;
            background-color: #fff;
            &.button-section {
                position: sticky;
                bottom: 0;
                left: 0;
                font-size: 0;
                .bk-button {
                    margin: 0 10px 0 0;
                }
            }
            .config-type {
                margin: 20px 0 0 0;
                line-height: 18px;
                font-size: 14px;
                color: #333948;
            }
        }
    }
    .config-name {
        margin: 10px 0 0 0;
        font-size: 0;
        .config-name-checkbox {
            width: 155px;
            margin: 0 30px 0 0;
        }
    }
    .association-group {
        .group-item {
            .group-name {
                position: relative;
                padding: 0 0 0 12px;
                margin: 10px 0 0 0;
                line-height: 18px;
                &:before {
                    position: absolute;
                    left: 0;
                    top: 2px;
                    width: 4px;
                    height: 14px;
                    background-color: $cmdbBorderColor;
                    content: "";
                }
            }
        }
    }
    .model-list {
        position: relative;
        margin: 10px 0 25px 0;
        .model-item {
            width: 155px;
            padding: 7px 0;
            margin: 0 30px 0 0;
            font-size: 0;
            .model-collapse-trigger {
                display: inline-block;
                vertical-align: middle;
                line-height: 18px;
                font-size: 0px;
                cursor: pointer;
            }
        }
    }
    .model-collapse-trigger {
        .model-name,
        .model-association-count {
            display: inline-block;
            vertical-align: middle;
            font-size: 14px;
        }
        .model-association-count {
            position: relative;
            color: $cmdbBorderFocusColor;
            &.has-angle {
                &:before,
                &:after {
                    position: absolute;
                    content: '';
                    border: 6px solid transparent;
                    border-bottom-color: $cmdbTableBorderColor;
                    top: 14px;
                    left: 3px;
                    z-index: 101;
                }
                &:after {
                    top: 15px;
                    border-bottom-color: #fff;
                }
            }
        }
        .bk-icon {
            font-size: 12px;
            margin: 0 0 0 4px;
        }
    }
    .model-association {
        position: absolute;
        left: 0;
        width: 100%;
        margin: 7px 0 0 0;
        padding: 10px 18px;
        background-color: #fff;
        border: 1px solid $cmdbTableBorderColor;
        z-index: 100;
        .association-checkbox {
            width: 175px;
            margin: 0 9px 0 0;
            line-height: 18px;
            .association-desc {
                display: block;
                max-width: 160px;
                @include ellipsis;
            }
        }
    }
</style>
