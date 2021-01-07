<template>
    <div class="relation">
        <div class="options clearfix">
            <div class="fl" v-show="activeView === viewName.list">
                <cmdb-auth class="inline-block-middle mr10"
                    v-if="hasRelation"
                    :auth="authResources">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        class="options-button"
                        :disabled="disabled"
                        @click="showCreate = true">
                        {{$t('新增关联')}}
                    </bk-button>
                </cmdb-auth>
                <span class="inline-block-middle mr10" v-else v-bk-tooltips="$t('当前模型暂未定义可用关联')">
                    <bk-button theme="primary" class="options-button" disabled>
                        {{$t('新增关联')}}
                    </bk-button>
                </span>
                <bk-checkbox v-if="hasRelation"
                    :size="16" class="options-checkbox"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('全部展开')}}</span>
                </bk-checkbox>
                <bk-button theme="default" class="options-button" v-show="false">{{$t('批量取消')}}</bk-button>
            </div>
            <div class="fr">
                <bk-button class="options-button options-button-view"
                    :theme="activeView === viewName.list ? 'primary' : 'default'"
                    @click="handleToggleView(viewName.list)">
                    {{$t('列表')}}
                </bk-button>
                <bk-button class="options-button options-button-view"
                    :theme="activeView === viewName.graphics ? 'primary' : 'default'"
                    @click="handleToggleView(viewName.graphics)">
                    {{$t('拓扑')}}
                </bk-button>
                <bk-button class="options-full-screen"
                    v-show="activeView === viewName.graphics"
                    v-bk-tooltips="$t('全屏')"
                    @click="handleFullScreen">
                    <i class="icon-cc-resize-full"></i>
                </bk-button>
            </div>
        </div>
        <div class="relation-view">
            <component
                ref="dynamicComponent"
                :association-types="associationTypes"
                :association-object="associationObject"
                :is="activeView">
            </component>
        </div>
        <bk-sideslider v-transfer-dom :is-show.sync="showCreate" :width="800" :title="$t('新增关联')">
            <cmdb-relation-create
                slot="content"
                :obj-id="objId"
                :inst="formatedInst"
                :association-types="associationTypes"
                :association-object="associationObject"
                v-if="showCreate">
            </cmdb-relation-create>
        </bk-sideslider>
    </div>
</template>

<script>
    import bus from '@/utils/bus.js'
    import { mapActions } from 'vuex'
    import cmdbRelationList from './list.vue'
    import cmdbRelationGraphics from './graphics.vue'
    import cmdbRelationCreate from './create.vue'
    import authMixin from '../mixin-auth'
    export default {
        components: {
            cmdbRelationList,
            cmdbRelationGraphics,
            cmdbRelationCreate
        },
        mixins: [authMixin],
        props: {
            objId: {
                type: String,
                required: true
            },
            inst: {
                type: Object,
                required: true
            },
            resourceType: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                associationObject: [],
                associationTypes: [],
                hasRelation: false,
                fullScreen: false,
                viewName: {
                    'list': cmdbRelationList.name,
                    'graphics': cmdbRelationGraphics.name
                },
                activeView: cmdbRelationList.name,
                showCreate: false,
                idKeyMap: {
                    host: 'bk_host_id',
                    biz: 'bk_biz_id'
                },
                nameKeyMap: {
                    host: 'bk_host_innerip',
                    biz: 'bk_biz_name'
                }
            }
        },
        computed: {
            formatedInst () {
                const idKey = this.idKeyMap[this.objId] || 'bk_inst_id'
                const nameKey = this.nameKeyMap[this.objId] || 'bk_inst_name'
                return {
                    ...this.inst,
                    'bk_inst_id': this.inst[idKey],
                    'bk_inst_name': this.inst[nameKey]
                }
            },
            authResources () {
                if (this.resourceType === 'business') {
                    return this.INST_AUTH.U_BUSINESS
                }
                return this.INST_AUTH.U_INST
            }
        },
        created () {
            this.getRelation()
            this.getAssociationType()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType',
                'searchObjectAssociation'
            ]),
            async getRelation () {
                try {
                    let [dataAsSource, dataAsTarget, mainLineModels] = await Promise.all([
                        this.getObjectAssociation({ 'bk_obj_id': this.objId }, { requestId: 'getSourceAssocaition' }),
                        this.getObjectAssociation({ 'bk_asst_obj_id': this.objId }, { requestId: 'getTargetAssocaition' }),
                        this.$store.dispatch('objectMainLineModule/searchMainlineObject', {
                            config: {
                                requestId: 'getMainLineModels'
                            }
                        })
                    ])
                    mainLineModels = mainLineModels.filter(model => !['biz', 'host'].includes(model['bk_obj_id']))
                    dataAsSource = this.getAvailableRelation(dataAsSource, mainLineModels)
                    dataAsTarget = this.getAvailableRelation(dataAsTarget, mainLineModels)
                    this.associationObject = [...dataAsSource, ...dataAsTarget]
                    if (dataAsSource.length || dataAsTarget.length) {
                        this.hasRelation = true
                    }
                } catch (e) {
                    this.hasRelation = false
                }
            },
            getAvailableRelation (data, mainLine) {
                return data.filter(relation => {
                    return !mainLine.some(model => [relation['bk_obj_id'], relation['bk_asst_obj_id']].includes(model['bk_obj_id']))
                })
            },
            getObjectAssociation (condition, config) {
                return this.searchObjectAssociation({
                    params: { condition },
                    config
                })
            },
            getAssociationType () {
                return this.searchAssociationType({}, {
                    config: {
                        requestId: 'getAssociationType'
                    }
                }).then(data => {
                    this.associationTypes = data.info
                    return data
                })
            },
            handleFullScreen () {
                this.$refs.dynamicComponent.toggleFullScreen(true)
            },
            handleToggleView (view) {
                this.activeView = view
            },
            handleExpandAll (expandAll) {
                bus.$emit('expand-all-relation-table', expandAll)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation {
        height: 100%;
        .relation-view {
            height: calc(100% - 62px);
            @include scrollbar-y;
        }
    }
    .options {
        padding: 15px 0;
        font-size: 0;
        .options-button {
            height: 32px;
            line-height: 30px;
            font-size: 14px;
            &.options-button-view {
                margin: 0 0 0 -1px;
                border-radius: 0;
            }
        }
        .options-checkbox {
            margin: 0 0 0 25px;
            line-height: 32px;
            .checkbox-label {
                padding-left: 4px;
            }
        }
        .options-full-screen {
            width: 32px;
            height: 32px;
            padding: 0;
            text-align: center;
            margin-left: 10px;
        }
    }
</style>
