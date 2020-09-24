<template>
    <div class="verification-layout">
        <div class="options">
            <cmdb-auth class="inline-block-middle"
                v-if="!isTopoModel"
                :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
                @update-auth="handleReceiveAuth">
                <bk-button slot-scope="{ disabled }"
                    class="create-btn"
                    theme="primary"
                    :disabled="isReadOnly || disabled"
                    @click="createVerification">
                    {{$t('新建校验')}}
                </bk-button>
            </cmdb-auth>
        </div>
        <bk-table
            class="verification-table"
            v-bkloading="{
                isLoading: $loading(['searchObjectUniqueConstraints', 'deleteObjectUniqueConstraints'])
            }"
            :data="table.list"
            :max-height="$APP.height - 320"
            :row-style="{
                cursor: 'pointer'
            }"
            @cell-click="handleShowDetails">
            <bk-table-column :label="$t('校验规则')" class-name="is-highlight" show-overflow-tooltip>
                <template slot-scope="{ row }">
                    {{getRuleName(row.keys)}}
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('属性为空值是否校验')">
                <template slot-scope="{ row }">
                    {{row.must_check ? $t('是') : $t('否')}}
                </template>
            </bk-table-column>
            <bk-table-column prop="operation"
                v-if="updateAuth && !isTopoModel"
                :label="$t('操作')">
                <template slot-scope="{ row }">
                    <button class="text-primary mr10 operation-btn"
                        :disabled="!isEditable(row)"
                        @click.stop="editVerification(row)">
                        {{$t('编辑')}}
                    </button>
                    <button class="text-primary operation-btn"
                        :disabled="!isEditable(row) || row.must_check"
                        @click.stop="deleteVerification(row)">
                        {{$t('删除')}}
                    </button>
                </template>
            </bk-table-column>
            <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :width="450"
            :title="slider.title"
            :is-show.sync="slider.isShow"
            :before-close="handleSliderBeforeClose">
            <the-verification-detail
                ref="verificationForm"
                slot="content"
                v-if="slider.isShow"
                :is-read-only="isReadOnly || slider.isReadOnly"
                :is-edit="slider.isEdit"
                :verification="slider.verification"
                :attribute-list="attributeList"
                @save="saveVerification"
                @cancel="handleSliderBeforeClose">
            </the-verification-detail>
        </bk-sideslider>
    </div>
</template>

<script>
    import theVerificationDetail from './verification-detail'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            theVerificationDetail
        },
        props: {
            modelId: {
                type: Number,
                default: null
            }
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    isEdit: false,
                    verification: {}
                },
                table: {
                    list: [],
                    stuff: {
                        type: 'default',
                        payload: {
                            emptyText: this.$t('bk.table.emptyText')
                        }
                    }
                },
                attributeList: [],
                updateAuth: false
            }
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isTopoModel () {
                return ['bk_biz_topo', 'bk_organization'].includes(this.activeModel.bk_classification_id)
            },
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            }
        },
        watch: {
            activeModel: {
                immediate: true,
                handler (activeModel) {
                    if (activeModel.bk_obj_id) {
                        this.initAttrList()
                        this.searchVerification()
                    }
                }
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            ...mapActions('objectUnique', [
                'searchObjectUniqueConstraints',
                'deleteObjectUniqueConstraints'
            ]),
            isEditable (item) {
                if (item.ispre || this.isReadOnly) {
                    return false
                }
                return true
            },
            getRuleName (keys) {
                const name = []
                keys.forEach(key => {
                    if (key['key_kind'] === 'property') {
                        const attr = this.attributeList.find(({ id }) => id === key['key_id'])
                        if (attr) {
                            name.push(attr['bk_property_name'])
                        }
                    }
                })
                return name.join('+')
            },
            async initAttrList () {
                this.attributeList = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: this.activeModel['bk_obj_id']
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`
                    }
                })
            },
            createVerification () {
                this.slider.title = this.$t('新建校验')
                this.slider.isEdit = false
                this.slider.isReadOnly = false
                this.slider.isShow = true
            },
            editVerification (verification) {
                this.slider.title = this.$t('编辑校验')
                this.slider.verification = verification
                this.slider.isEdit = true
                this.slider.isReadOnly = false
                this.slider.isShow = true
            },
            saveVerification () {
                this.slider.isShow = false
                this.searchVerification()
            },
            deleteVerification (verification) {
                this.$bkInfo({
                    title: this.$tc('确定删除唯一校验', this.getRuleName(verification.keys), { name: this.getRuleName(verification.keys) }),
                    confirmFn: async () => {
                        await this.deleteObjectUniqueConstraints({
                            objId: verification['bk_obj_id'],
                            id: verification.id,
                            params: {},
                            config: {
                                requestId: 'deleteObjectUniqueConstraints'
                            }
                        })
                        this.searchVerification()
                    }
                })
            },
            async searchVerification () {
                const res = await this.searchObjectUniqueConstraints({
                    objId: this.activeModel['bk_obj_id'],
                    params: {},
                    config: {
                        requestId: 'searchObjectUniqueConstraints'
                    }
                })
                this.table.list = res
            },
            handleShowDetails (row, column, cell) {
                if (column.property === 'operation') return
                this.slider.title = this.$t('查看校验')
                this.slider.verification = row
                this.slider.isEdit = true
                this.slider.isReadOnly = true
                this.slider.isShow = true
            },
            handleReceiveAuth (auth) {
                this.updateAuth = auth
            },
            handleSliderBeforeClose () {
                const hasChanged = Object.keys(this.$refs.verificationForm.changedValues).length
                if (hasChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.slider.isShow = false
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.slider.isShow = false
                return true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .verification-layout {
        padding: 20px;
    }
    .verification-table {
        margin: 14px 0 0 0;
    }
    .operation-btn[disabled] {
        color: #dcdee5 !important;
        opacity: 1 !important;
    }
</style>
