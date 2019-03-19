<template>
    <div>
        <bk-button class="create-btn" type="primary"
            :disabled="isReadOnly || !authority.includes('update')"
            @click="createVerification">
            {{$t('ModelManagement["新建校验"]')}}
        </bk-button>
        <cmdb-table
            class="relation-table"
            :loading="$loading(['searchObjectUniqueConstraints', 'deleteObjectUniqueConstraints'])"
            :sortable="false"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <div :key="index" :class="{'disabled': isReadOnly}">
                    <template v-if="header.id==='keys'">
                        {{getRuleName(item.keys)}}
                    </template>
                    <template v-else-if="header.id==='must_check'">
                        {{item['must_check'] ? $t('ModelManagement["是"]') : $t('ModelManagement["否"]')}}
                    </template>
                    <template v-else-if="header.id==='operation'">
                        <template v-if="isReadOnly || item.ispre">
                            <span class="text-primary disabled mr10">
                                {{$t('Common["编辑"]')}}
                            </span>
                            <span class="text-primary disabled">
                                {{$t('Common["删除"]')}}
                            </span>
                        </template>
                        <template v-else>
                            <span class="text-primary mr10" @click.stop="editVerification(item)">
                                {{$t('Common["编辑"]')}}
                            </span>
                            <span class="text-primary" @click.stop="deleteVerification(item)">
                                {{$t('Common["删除"]')}}
                            </span>
                        </template>
                    </template>
                </div>
            </template>
        </cmdb-table>
        <cmdb-slider
            :width="450"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <the-verification-detail
                class="slider-content"
                slot="content"
                :isReadOnly="isReadOnly"
                :isEdit="slider.isEdit"
                :verification="slider.verification"
                :attributeList="attributeList"
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
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            authority () {
                const cantEdit = ['process', 'plat']
                if (cantEdit.includes(this.$route.params.modelId)) {
                    return []
                }
                return this.$store.getters.admin ? ['search', 'update', 'delete'] : []
            }
        },
        async created () {
            if (!this.authority.includes('update')) {
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
            getRuleName (keys) {
                let name = []
                keys.forEach(key => {
                    if (key['key_kind'] === 'property') {
                        let attr = this.attributeList.find(({id}) => id === key['key_id'])
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
                    title: this.$tc('ModelManagement["确定删除唯一校验？"]', this.getRuleName(verification.keys), {name: this.getRuleName(verification.keys)}),
                    confirmFn: async () => {
                        await this.deleteObjectUniqueConstraints({
                            objId: verification['bk_obj_id'],
                            id: verification.id,
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
                    requestId: 'searchObjectUniqueConstraints'
                })
                this.table.list = res
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
</style>

