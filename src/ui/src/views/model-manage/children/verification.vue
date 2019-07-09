<template>
    <div class="verification-layout">
        <span class="inline-block-middle"
            v-if="!isTopoModel"
            v-cursor="{
                active: !$isAuthorized($OPERATION.U_MODEL),
                auth: [$OPERATION.U_MODEL]
            }">
            <bk-button class="create-btn" theme="primary"
                :disabled="isReadOnly || !updateAuth"
                @click="createVerification">
                {{$t('ModelManagement["新建校验"]')}}
            </bk-button>
        </span>
        <bk-table
            class="verification-table"
            v-bkloading="{
                isLoading: $loading(['searchObjectUniqueConstraints', 'deleteObjectUniqueConstraints'])
            }"
            :data="table.list"
            :max-height="$APP.height - 220">
            <bk-table-column :label="$t('ModelManagement[\'校验规则\']')">
                <template slot-scope="{ row }">
                    {{getRuleName(row.keys)}}
                </template>
            </bk-table-column>
            <bk-table-column :label="$t('ModelManagement[\'是否为必须校验\']')">
                <template slot-scope="{ row }">
                    {{row.must_check ? $t('ModelManagement["是"]') : $t('ModelManagement["否"]')}}
                </template>
            </bk-table-column>
            <bk-table-column
                v-if="updateAuth && !isTopoModel"
                :label="$t('Common[\'操作\']')">
                <template slot-scope="{ row }">
                    <button class="text-primary mr10"
                        :disabled="!isEditable(row)"
                        @click.stop="editVerification(row)">
                        {{$t('Common["编辑"]')}}
                    </button>
                    <button class="text-primary"
                        :disabled="!isEditable(row)"
                        @click.stop="deleteVerification(row)">
                        {{$t('Common["删除"]')}}
                    </button>
                </template>
            </bk-table-column>
        </bk-table>
        <bk-sideslider
            :width="450"
            :title="slider.title"
            :is-show.sync="slider.isShow">
            <the-verification-detail
                class="slider-content"
                slot="content"
                :is-read-only="isReadOnly"
                :is-edit="slider.isEdit"
                :verification="slider.verification"
                :attribute-list="attributeList"
                @save="saveVerification"
                @cancel="slider.isShow = false">
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
        data () {
            return {
                slider: {
                    isShow: false,
                    isEdit: false,
                    verification: {}
                },
                table: {
                    list: []
                },
                attributeList: []
            }
        },
        computed: {
            ...mapGetters(['isAdminView', 'isBusinessSelected']),
            ...mapGetters('objectModel', [
                'activeModel',
                'isInjectable'
            ]),
            isTopoModel () {
                return this.activeModel.bk_classification_id === 'bk_biz_topo'
            },
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            updateAuth () {
                const cantEdit = ['process', 'plat']
                if (cantEdit.includes(this.$route.params.modelId)) {
                    return false
                }
                const editable = this.isAdminView || (this.isBusinessSelected && this.isInjectable)
                return editable && this.$isAuthorized(this.$OPERATION.U_MODEL)
            }
        },
        async created () {
            this.initAttrList()
            this.searchVerification()
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
                if (!this.isAdminView) {
                    return !!this.$tools.getMetadataBiz(item)
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
                    params: this.$injectMetadata({
                        bk_obj_id: this.activeModel['bk_obj_id']
                    }, {
                        inject: this.isInjectable
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`
                    }
                })
            },
            createVerification () {
                this.slider.title = this.$t('ModelManagement["新建校验"]')
                this.slider.isEdit = false
                this.slider.isShow = true
            },
            editVerification (verification) {
                this.slider.title = this.$t('ModelManagement["编辑校验"]')
                this.slider.verification = verification
                this.slider.isEdit = true
                this.slider.isShow = true
            },
            saveVerification () {
                this.slider.isShow = false
                this.searchVerification()
            },
            deleteVerification (verification) {
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除唯一校验？"]', this.getRuleName(verification.keys), { name: this.getRuleName(verification.keys) }),
                    confirmFn: async () => {
                        await this.deleteObjectUniqueConstraints({
                            objId: verification['bk_obj_id'],
                            id: verification.id,
                            params: this.$injectMetadata({}, {
                                inject: !!this.$tools.getMetadataBiz(verification)
                            }),
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
                    params: this.$injectMetadata({}, { inject: this.isInjectable }),
                    config: {
                        requestId: 'searchObjectUniqueConstraints'
                    }
                })
                this.table.list = res
            }
        }
    }
</script>

<style lang="scss" scoped>
    .verification-layout {
        padding: 10px 0;
    }
    .verification-table {
        margin: 10px 0 0 0;
    }
</style>
