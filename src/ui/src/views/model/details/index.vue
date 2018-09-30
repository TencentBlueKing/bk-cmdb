<template>
    <div class="details-wrapper">
        <bk-tab :active-name.sync="curTabName">
            <template slot="setting" v-if="isEdit">
                <div class="btn-group">
                    <span class="btn-box" v-if="!(ispre || isMainLine)" v-tooltip="$t('ModelManagement[\'保留模型和相应实例，隐藏关联关系\']')">
                        <span v-if="ispaused" @click="dialogConfirm('restart')">
                            <i class="icon bk-icon icon-minus-circle-shape"></i>
                            <span class="text">{{$t('ModelManagement["启用模型"]')}}</span>
                        </span>
                        <span v-else @click="dialogConfirm('stop')">
                            <i class="icon bk-icon icon-minus-circle-shape"></i>
                            <span class="text">{{$t('ModelManagement["停用模型"]')}}</span>
                        </span>
                    </span>
                    <span class="btn-box"
                        v-if="!ispre"
                        v-tooltip="$t('ModelManagement[\'删除模型和其下所有实例，此动作不可逆，请谨慎操作\']')"
                        @click="dialogConfirm('delete')" >
                        <i class="icon icon-cc-del"></i>
                        <span class="text">{{$t('ModelManagement["删除模型"]')}}</span>
                    </span>
                </div>
            </template>
            <bk-tabpanel name="host" :title="$t('ModelManagement[\'模型配置\']')">
                <v-field 
                    v-if="curTabName === 'host'"
                    :isEdit.sync="isEdit"
                    @createObject="createObject"
                    @updateTopo="updateModel"
                    @confirm="updateModel"
                    @cancel="cancel"
                    ref="tab"
                ></v-field>
            </bk-tabpanel>
            <bk-tabpanel name="layout" :title="$t('ModelManagement[\'字段分组\']')" :show="isEdit">
                <v-layout v-if="curTabName === 'layout'"
                    @cancel="cancel"
                    ref="tab"
                ></v-layout>
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import vField from './field'
    import vLayout from './layout'
    import vOther from './other'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            vField,
            vLayout,
            vOther
        },
        props: {
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                curTabName: 'host'
            }
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            ispre () {
                return this.activeModel['ispre']
            },
            ispaused () {
                return this.activeModel['bk_ispaused']
            },
            isMainLine () {
                return this.activeModel['bk_classification_id'] === 'bk_biz_topo'
            }
        },
        methods: {
            ...mapActions('objectModel', [
                'updateObject',
                'deleteObject'
            ]),
            ...mapActions('objectMainLineModule', [
                'deleteMainlineObject'
            ]),
            dialogConfirm (type) {
                switch (type) {
                    case 'restart':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要启用该模型？"]'),
                            confirmFn: () => {
                                this.updateModelObject(false)
                            }
                        })
                        break
                    case 'stop':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要停用该模型？"]'),
                            confirmFn: () => {
                                this.updateModelObject(true)
                            }
                        })
                        break
                    case 'delete':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要删除该模型？"]'),
                            confirmFn: () => {
                                this.deleteModel()
                            }
                        })
                        break
                    default:
                }
            },
            async updateModelObject (ispaused) {
                await this.updateObject({
                    id: this.activeModel['id'],
                    params: {
                        bk_ispaused: ispaused
                    },
                    config: {
                        requestId: 'updateModel'
                    }
                }).then(() => {
                    this.$http.cancel('post_searchClassificationsObjects')
                })
                this.updateModel(false)
            },
            async deleteModel () {
                if (this.isMainLine) {
                    await this.deleteMainlineObject({
                        bkObjId: this.activeModel['bk_obj_id'],
                        config: {
                            requestId: 'deleteModel'
                        }
                    })
                } else {
                    await this.deleteObject({
                        id: this.activeModel['id'],
                        config: {
                            requestId: 'deleteModel'
                        }
                    })
                }
                this.$http.cancel('post_searchClassificationsObjects')
                this.updateModel(false)
            },
            cancel () {
                this.$emit('cancel')
            },
            createObject () {
                this.$emit('createModel')
                this.$emit('update:isEdit', true)
            },
            updateModel (isShow) {
                this.$emit('updateModel', isShow)
            },
            isCloseConfirmShow () {
                if (this.curTabName === 'other') {
                    return false
                }
                return this.$refs.tab.isCloseConfirmShow()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-wrapper{
        padding: 8px 0 0;
        height: 100%;
        .btn-group {
            font-size: 0;
            .btn-box {
                cursor: pointer;
            }
            >span {
                font-size: 14px;
                &:last-child {
                    margin-left: 30px;
                }
            }
            .icon,
            .text {
                vertical-align: middle;
                font-size: 14px;
            }
            .icon {
                margin-right: 5px;
            }
        }
    }
</style>
