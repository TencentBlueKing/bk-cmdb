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
        <cmdb-table
            class="verification-table"
            :loading="$loading(['searchObjectUniqueConstraints', 'deleteObjectUniqueConstraints'])"
            :sortable="false"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapper-minus-height="220">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <div :key="index" :class="{ 'disabled': isReadOnly }">
                    <template v-if="header.id === 'keys'">
                        {{getRuleName(item.keys)}}
                    </template>
                    <template v-else-if="header.id === 'must_check'">
                        {{item['must_check'] ? $t('ModelManagement["是"]') : $t('ModelManagement["否"]')}}
                    </template>
                    <template v-else-if="header.id === 'operation'">
                        <button class="text-primary mr10"
                            :disabled="!isEditable(item)"
                            @click.stop="editVerification(item)">
                            {{$t('Common["编辑"]')}}
                        </button>
                        <button class="text-primary"
                            :disabled="!isEditable(item)"
                            @click.stop="deleteVerification(item)">
                            {{$t('Common["删除"]')}}
                        </button>
                    </template>
                </div>
            </template>
        </cmdb-table>
        <cmdb-slider
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
        </cmdb-slider>
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
                    header: [{
                        id: 'keys',
                        name: this.$t('ModelManagement["校验规则"]')
                    }, {
                        id: 'must_check',
                        name: this.$t('ModelManagement["是否为必须校验"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Common["操作"]')
                    }],
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
            if (!this.updateAuth || this.isTopoModel) {
                this.table.header.pop()
            }
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
