<template>
    <div class="model-relation-wrapper">
        <bk-button class="create-btn" type="primary"
            :disabled="isReadOnly || !authority.includes('update')"
            @click="createRelation">
            {{$t('ModelManagement["新建关联"]')}}
        </bk-button>
        <cmdb-table
            class="relation-table"
            :loading="$loading()"
            :sortable="false"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handleSortChange="handleSortChange">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <div :key="index" :class="{'disabled': isReadOnly}">
                    <template v-if="header.id==='mapping'">
                        {{mappingMap[item.mapping]}}
                    </template>
                    <template v-else-if="header.id==='bk_asst_name'">
                        {{getRelationName(item['bk_asst_id'])}}
                    </template>
                    <template v-else-if="header.id==='bk_obj_name'">
                        {{getModelName(item['bk_obj_id'])}}
                    </template>
                    <template v-else-if="header.id==='bk_asst_obj_name'">
                        {{getModelName(item['bk_asst_obj_id'])}}
                    </template>
                    <template v-else-if="header.id==='operation'">
                        <template v-if="isReadOnly || item.ispre || item['bk_asst_id'] === 'bk_mainline'">
                            <span class="text-primary mr10 disabled">
                            {{$t('Common["编辑"]')}}
                            </span>
                            <span class="text-primary disabled">
                                {{$t('Common["删除"]')}}
                            </span>
                        </template>
                        <template v-else>
                            <span class="text-primary mr10" @click.stop="editRelation(item)">
                                {{$t('Common["编辑"]')}}
                            </span>
                            <span class="text-primary" @click.stop="deleteRelation(item, index)">
                                {{$t('Common["删除"]')}}
                            </span>
                        </template>
                    </template>
                    <template v-else>
                        {{item[header.id]}}
                    </template>
                </div>
            </template>
        </cmdb-table>
        <cmdb-slider
            :width="450"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <the-relation-detail
                class="slider-content"
                slot="content"
                :isReadOnly="isReadOnly"
                :isEdit="slider.isEdit"
                :relation="slider.relation"
                :relationList="relationList"
                @save="saveRelation"
                @cancel="slider.isShow = false">
            </the-relation-detail>
        </cmdb-slider>
    </div>
</template>

<script>
    import theRelationDetail from './relation-detail'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            theRelationDetail
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    isEdit: false,
                    title: this.$t('ModelManagement["新建关联"]'),
                    relation: {}
                },
                relationList: [],
                table: {
                    header: [{
                        id: 'bk_obj_asst_id',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'bk_asst_name',
                        name: this.$t('ModelManagement["关联类型"]')
                    }, {
                        id: 'mapping',
                        name: this.$t('ModelManagement["源-目标约束"]')
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('ModelManagement["源模型"]')
                    }, {
                        id: 'bk_asst_obj_name',
                        name: this.$t('ModelManagement["目标模型"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Common["操作"]')
                    }],
                    list: [],
                    defaultSort: '-op_time',
                    sort: '-op_time'
                },
                mappingMap: {
                    '1:1': '1-1',
                    '1:n': '1-N',
                    'n:n': 'N-N'
                }
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
        created () {
            if (!this.authority.includes('update')) {
                this.table.header.pop()
            }
            this.searchRelationList()
            this.initRelationList()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchObjectAssociation',
                'deleteObjectAssociation',
                'searchAssociationType'
            ]),
            getRelationName (id) {
                let relation = this.relationList.find(item => item.id === id)
                if (relation) {
                    return relation.name
                }
            },
            async initRelationList () {
                const data = await this.searchAssociationType({
                    params: {},
                    config: {
                        requestId: 'post_searchAssociationType',
                        fromCache: true
                    }
                })
                this.relationList = data.info.map(({bk_asst_id: asstId, bk_asst_name: asstName}) => {
                    if (asstName.length) {
                        return {
                            id: asstId,
                            name: `${asstId}(${asstName})`
                        }
                    }
                    return {
                        id: asstId,
                        name: asstId
                    }
                })
            },
            getModelName (objId) {
                let model = this.$allModels.find(model => model['bk_obj_id'] === objId)
                if (model) {
                    return model['bk_obj_name']
                }
                return ''
            },
            createRelation () {
                this.slider.isEdit = false
                this.slider.isReadOnly = false
                this.slider.relation = {}
                this.slider.title = this.$t('ModelManagement["新建关联"]')
                this.slider.isShow = true
            },
            editRelation (item) {
                this.slider.isEdit = true
                this.slider.isReadOnly = false
                this.slider.relation = item
                this.slider.title = this.$t('ModelManagement["编辑关联"]')
                this.slider.isShow = true
            },
            deleteRelation (relation, index) {
                this.$bkInfo({
                    title: this.$t('ModelManagement["确定删除关联关系?"]'),
                    confirmFn: async () => {
                        await this.deleteObjectAssociation({
                            id: relation.id,
                            config: {
                                requestId: 'deleteObjectAssociation'
                            }
                        }).then(() => {
                            this.$http.cancel(`post_searchObjectAssociation_${this.activeModel['bk_obj_id']}`)
                        })
                        this.searchRelationList()
                    }
                })
            },
            async searchRelationList () {
                const [source, dest] = await Promise.all([this.searchAsSource(), this.searchAsDest()])
                this.table.list = [...source, ...dest]
            },
            searchAsSource () {
                return this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_obj_id': this.activeModel['bk_obj_id']
                        }
                    }
                })
            },
            searchAsDest () {
                return this.searchObjectAssociation({
                    params: {
                        condition: {
                            'bk_asst_obj_id': this.activeModel['bk_obj_id']
                        }
                    }
                })
            },
            saveRelation () {
                this.slider.isShow = false
                this.searchRelationList()
            },
            handleSortChange (sort) {
                this.table.sort = sort
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
</style>
