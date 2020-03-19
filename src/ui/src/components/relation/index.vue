<template>
    <div class="relation-layout">
        <div class="relation-options clearfix">
            <div class="fl">
                <cmdb-auth class="inline-block-middle" :auth="authResources">
                    <bk-button slot-scope="{ disabled }"
                        class="options-button options-button-update"
                        theme="primary"
                        :disabled="!hasRelation || disabled"
                        :class="{ active: activeComponent === 'cmdbRelationUpdate' }"
                        @click="handleShowUpdate">
                        {{$t('关联管理')}}
                        <i class="bk-icon icon-angle-down"></i>
                    </bk-button>
                </cmdb-auth>
            </div>
            <div class="fr">
                <bk-button theme="default" class="options-full-screen"
                    v-show="activeComponent === 'cmdbRelationTopology'"
                    v-bk-tooltips="$t('全屏')"
                    @click="handleFullScreen">
                    <i class="icon-cc-resize-full"></i>
                </bk-button>
            </div>
        </div>
        <div class="relation-component">
            <component ref="dynamicComponent"
                :is="activeComponent">
            </component>
        </div>
    </div>
</template>

<script>
    // import cmdbRelationTopology from './_topology.vue'
    import cmdbRelationTopology from './_topology.new.vue'
    import cmdbRelationUpdate from './_update.vue'
    export default {
        components: {
            cmdbRelationTopology,
            cmdbRelationUpdate
        },
        props: {
            objId: {
                type: String,
                required: true
            },
            inst: {
                type: Object,
                required: true
            },
            auth: {
                type: [String, Array],
                default: ''
            }
        },
        data () {
            return {
                hasRelation: false,
                fullScreen: false,
                activeComponent: 'cmdbRelationTopology',
                previousComponent: 'cmdbRelationTopology',
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
                const auth = this.auth
                if (!auth) return {}
                if (Array.isArray(auth) && !auth.length) return {}
                return this.$authResources({ type: auth })
            }
        },
        created () {
            this.getRelation()
        },
        methods: {
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
                return this.$store.dispatch('objectAssociation/searchObjectAssociation', {
                    params: this.$injectMetadata({ condition }),
                    config
                })
            },
            handleShowUpdate () {
                if (this.activeComponent === 'cmdbRelationUpdate') {
                    this.activeComponent = this.previousComponent
                } else {
                    this.previousComponent = this.activeComponent
                    this.activeComponent = 'cmdbRelationUpdate'
                }
            },
            handleFullScreen () {
                this.$refs.dynamicComponent.toggleFullScreen(true)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-layout {
        height: 100%;
        .relation-options {
            padding: 20px 0 10px;
            font-size: 0;
        }
    }
    .relation-options {
        .options-full-screen {
            width: 32px;
            height: 32px;
            padding: 0;
            text-align: center;
            margin-right: 10px;
        }
        .icon-angle-down {
            font-size: 20px;
            vertical-align: baseline;
            transition: transform .2s linear;
        }
        .options-button {
            border-radius: 0;
            margin: 0 0 0 -1px;
        }
        .options-button-update {
            position: relative;
            margin: 0;
            &.active {
                background-color: #fff;
                color: $cmdbTextColor;
                border-color: $cmdbBorderColor;
                .icon-angle-down {
                    transform: rotate(-180deg);
                }
                &:after {
                    position: absolute;
                    top: 100%;
                    left: 0;
                    width: 100%;
                    height: 17px;
                    margin: -1px -1px 0;
                    border: 1px solid $cmdbBorderColor;
                    border-top: none;
                    border-bottom: none;
                    content: "";
                    background-color: #fff;
                    z-index: 1;
                }
            }
        }
    }
    .relation-component {
        height: calc(100% - 80px);
    }
</style>
