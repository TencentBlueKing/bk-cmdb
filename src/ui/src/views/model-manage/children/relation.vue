<template>
    <div class="model-relation-wrapper">
        <bk-button class="create-btn" type="primary" @click="createRelation">
            {{$t('ModelManagement["新建关联关系"]')}}
        </bk-button>
        <cmdb-table
            class="relation-table"
            :loading="$loading('getOperationLog')"
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
                    <template v-else-if="header.id==='operation'">
                        <span class="text-primary mr10" @click.stop="editRelation(item)">
                            {{$t('Common["编辑"]')}}
                        </span>
                        <span class="text-primary" v-if="!item.ispre && !isReadOnly" @click.stop="deleteRelation(item)">
                            {{$t('Common["删除"]')}}
                        </span>
                    </template>
                    <template v-else>
                        {{item[header.id]}}
                    </template>
                </div>
            </template>
        </cmdb-table>
        <cmdb-slider
            :width="514"
            :isShow.sync="slider.isShow">
            <the-relation-detail
                class="slider-content"
                slot="content"
                :isReadOnly="isReadOnly"
                :isEditField="slider.isEditField"
                :field="slider.curField"
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
                    relation: {}
                },
                table: {
                    header: [{
                        id: 'bk_asst_id',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'bk_asst_name',
                        name: this.$t('ModelManagement["关联类型"]')
                    }, {
                        id: 'mapping',
                        name: this.$t('ModelManagement["约束"]')
                    }, {
                        id: 'bk_obj_name',
                        name: this.$t('ModelManagement["源模型"]')
                    }, {
                        id: 'bk_asst_obj_name',
                        name: this.$t('ModelManagement["目标模型"]')
                    }, {
                        id: 'operation',
                        name: this.$t('ModelManagement["操作"]'),
                        sortable: false
                    }],
                    list: [{
                        'id': 1,
                        'bk_obj_asst_id': 'bk_switch_belong_bk_host',
                        'bk_obj_asst_name': '',
                        'bk_asst_id': 'belong',
                        'bk_obj_id': 'bk_switch',
                        'bk_obj_name': 'bk_switch',
                        'bk_asst_obj_id': 'bk_host',
                        'mapping': '1:n',
                        'on_delete': 'none'
                    }],
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
            }
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchObjectAssociation'
            ]),
            createRelation () {
                this.slider.isEdit = false
                this.slider.isReadOnly = false
                this.slider.relation = {}
                this.slider.title = this.$t('ModelManagement["新建关联关系"]')
                this.slider.isShow = true
            },
            editRelation () {

            },
            deleteRelation () {

            },
            saveRelation () {

            },
            handleSortChange (sort) {
                this.sort = sort
                this.refresh()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
</style>
