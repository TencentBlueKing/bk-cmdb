<template>
    <div class="model-field-wrapper">
        <div>
            <bk-button class="create-btn" type="primary"
                :disabled="isReadOnly || !authority.includes('update')"
                @click="createField">
                {{$t('ModelManagement["新建字段"]')}}
            </bk-button>
        </div>
        <cmdb-table
            class="field-table"
            :loading="$loading('initFieldList')"
            :header="table.header"
            :has-footer="false"
            :list="table.list"
            :wrapperMinusHeight="300"
            @handleSortChange="handleSortChange">
            <template v-for="(header, index) in table.header" :slot="header.id" slot-scope="{ item }">
                <div :key="index" :class="{'disabled': item.ispre || isReadOnly}">
                    <template v-if="header.id==='bk_property_type'">
                        {{fieldTypeMap[item['bk_property_type']]}}
                    </template>
                    <template v-else-if="header.id==='isrequired'">
                        <i class="bk-icon icon-check-1" v-if="item.isrequired"></i>
                    </template>
                    <template v-else-if="header.id==='create_time'">
                        {{$tools.formatTime(item['create_time'])}}
                    </template>
                    <template v-else-if="header.id==='operation'">
                        <span class="text-primary disabled mr10" v-if="isReadOnly">
                            {{$t('Common["编辑"]')}}
                        </span>
                        <span class="text-primary mr10"
                            v-else
                            @click.stop="editField(item)">
                            {{$t('Common["编辑"]')}}
                        </span>
                        <span class="text-primary" v-if="!item.ispre && !isReadOnly" @click.stop="deleteField(item)">
                            {{$t('Common["删除"]')}}
                        </span>
                        <span class="text-primary disabled" v-else>
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
            :width="450"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <the-field-detail
                class="slider-content"
                slot="content"
                :isReadOnly="isReadOnly"
                :isEditField="slider.isEditField"
                :field="slider.curField"
                @save="saveField"
                @cancel="slider.isShow = false">
            </the-field-detail>
        </cmdb-slider>
    </div>
</template>

<script>
    import theFieldDetail from './field-detail'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            theFieldDetail
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    isEditField: false,
                    title: this.$t('ModelManagement["新建字段"]'),
                    curField: {}
                },
                fieldTypeMap: {
                    'singlechar': this.$t('ModelManagement["短字符"]'),
                    'int': this.$t('ModelManagement["数字"]'),
                    'float': this.$t('ModelManagement["浮点"]'),
                    'enum': this.$t('ModelManagement["枚举"]'),
                    'date': this.$t('ModelManagement["日期"]'),
                    'time': this.$t('ModelManagement["时间"]'),
                    'longchar': this.$t('ModelManagement["长字符"]'),
                    'objuser': this.$t('ModelManagement["用户"]'),
                    'timezone': this.$t('ModelManagement["时区"]'),
                    'bool': 'bool'
                },
                table: {
                    header: [{
                        id: 'bk_property_id',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'bk_property_type',
                        name: this.$t('ModelManagement["字段类型"]')
                    }, {
                        id: 'isrequired',
                        name: this.$t('ModelManagement["必填"]')
                    }, {
                        id: 'bk_property_name',
                        name: this.$t('ModelManagement["名称"]')
                    }, {
                        id: 'create_time',
                        name: this.$t('ModelManagement["创建时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Common["操作"]'),
                        sortable: false
                    }],
                    list: [],
                    defaultSort: '-create_time',
                    sort: '-create_time'
                }
            }
        },
        computed: {
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            objId () {
                return this.$route.params.modelId
            },
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            authority () {
                const cantEdit = ['process', 'plat']
                if (cantEdit.includes(this.objId)) {
                    return []
                }
                return this.$store.getters.admin ? ['search', 'update', 'delete'] : []
            }
        },
        watch: {
            '$route.params.modelId' () {
                this.initFieldList()
            }
        },
        created () {
            this.initFieldList()
            if (!this.authority.includes('update')) {
                this.table.header.pop()
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute',
                'deleteObjectAttribute'
            ]),
            createField () {
                this.slider.isEditField = false
                this.slider.isReadOnly = false
                this.slider.title = this.$t('ModelManagement["新建字段"]')
                this.slider.curField = {}
                this.slider.isShow = true
            },
            editField (item) {
                this.slider.isEditField = true
                this.slider.isReadOnly = this.isReadOnly
                this.slider.title = this.$t('ModelManagement["编辑字段"]')
                this.slider.curField = item
                this.slider.isShow = true
            },
            deleteField (field) {
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除字段？"]', field['bk_property_name'], {name: field['bk_property_name']}),
                    confirmFn: async () => {
                        await this.deleteObjectAttribute({
                            id: field.id,
                            config: {
                                requestId: 'deleteObjectAttribute'
                            }
                        }).then(() => {
                            this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                        })
                        this.initFieldList()
                    }
                })
            },
            async initFieldList () {
                const res = await this.searchObjectAttribute({
                    params: {
                        bk_obj_id: this.objId
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_${this.objId}`
                    }
                })
                this.table.list = res
            },
            saveField () {
                this.slider.isShow = false
                this.initFieldList()
            },
            handleSortChange (sort) {
                if (!sort.length) {
                    sort = this.table.defaultSort
                }
                let key = sort
                if (sort[0] === '-') {
                    key = sort.substr(1, sort.length - 1)
                }
                this.table.list.sort((itemA, itemB) => {
                    if (key === 'isrequired') {
                        return itemA[key] > itemB[key]
                    }
                    return itemA[key].localeCompare(itemB[key])
                })
                if (sort[0] === '-') {
                    this.table.list.reverse()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-btn {
        margin: 10px 0;
    }
    .field-table {
        .disabled {
            color: #bfc7d2;
        }
    }
</style>
