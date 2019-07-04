<template>
    <div class="model-field-wrapper">
        <div class="options">
            <span class="inline-block-middle" v-cursor="{
                active: !$isAuthorized($OPERATION.U_MODEL),
                auth: [$OPERATION.U_MODEL]
            }">
                <bk-button class="create-btn" theme="primary"
                    :disabled="createDisabled || isReadOnly || !updateAuth"
                    @click="createField">
                    {{$t('ModelManagement["新建字段"]')}}
                </bk-button>
            </span>
        </div>
        <cmdb-table
            class="field-table"
            :loading="$loading(`post_searchObjectAttribute_${objId}`)"
            :header="table.header"
            :has-footer="false"
            :list="table.list"
            :wrapper-minus-height="300"
            @handleSortChange="handleSortChange">
            <template slot="bk_property_id" slot-scope="{ item }">
                <span
                    v-if="item['ispre']"
                    :class="['field-pre', $i18n.locale]">
                    {{$t('ModelManagement["内置"]')}}
                </span>
                <span class="field-id">{{item['bk_property_id']}}</span>
            </template>
            <template slot="bk_property_type" slot-scope="{ item }">
                <span>{{fieldTypeMap[item['bk_property_type']]}}</span>
            </template>
            <template slot="isrequired" slot-scope="{ item }">
                <i class="field-required-icon bk-icon icon-check-1" v-if="item.isrequired"></i>
                <i class="field-required-icon bk-icon icon-close" v-else></i>
            </template>
            <template slot="create_time" slot-scope="{ item }">
                {{$tools.formatTime(item['create_time'])}}
            </template>
            <template slot="operation" slot-scope="{ item }">
                <button class="text-primary mr10"
                    :disabled="!isFieldEditable(item)"
                    @click.stop="editField(item)">
                    {{$t('Common["编辑"]')}}
                </button>
                <button class="text-primary"
                    :disabled="!isFieldEditable(item)"
                    @click.stop="deleteField(item)">
                    {{$t('Common["删除"]')}}
                </button>
            </template>
        </cmdb-table>
        <bk-sideslider
            :width="450"
            :title="slider.title"
            :is-show.sync="slider.isShow"
            :before-close="handleSliderBeforeClose">
            <the-field-detail
                ref="fieldForm"
                class="slider-content"
                slot="content"
                :is-read-only="isReadOnly"
                :is-edit-field="slider.isEditField"
                :field="slider.curField"
                @save="saveField"
                @cancel="handleSliderBeforeClose">
            </the-field-detail>
        </bk-sideslider>
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
                        id: 'isrequired',
                        name: this.$t('ModelManagement["必填"]')
                    }, {
                        id: 'bk_property_id',
                        name: this.$t('ModelManagement["唯一标识"]'),
                        minWidth: 110
                    }, {
                        id: 'bk_property_name',
                        name: this.$t('ModelManagement["名称"]')
                    }, {
                        id: 'bk_property_type',
                        name: this.$t('ModelManagement["字段类型"]')
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
            ...mapGetters(['isAdminView', 'isBusinessSelected']),
            ...mapGetters('objectModel', [
                'activeModel',
                'isPublicModel',
                'isMainLine',
                'isInjectable'
            ]),
            objId () {
                return this.$route.params.modelId
            },
            createDisabled () {
                return !this.isAdminView && this.isPublicModel
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
        watch: {
            '$route.params.modelId' () {
                this.initFieldList()
            }
        },
        created () {
            this.initFieldList()
            if (!this.updateAuth) {
                this.table.header.pop()
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute',
                'deleteObjectAttribute'
            ]),
            isFieldEditable (item) {
                if (item.ispre || this.isReadOnly || !this.updateAuth) {
                    return false
                }
                if (!this.isAdminView) {
                    return !!this.$tools.getMetadataBiz(item)
                }
                return true
            },
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
                    title: this.$tc('ModelManagement["确定删除字段？"]', field['bk_property_name'], { name: field['bk_property_name'] }),
                    confirmFn: async () => {
                        await this.deleteObjectAttribute({
                            id: field.id,
                            config: {
                                data: this.$injectMetadata({}, {
                                    inject: this.isInjectable
                                }),
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
                    params: this.$injectMetadata({ bk_obj_id: this.objId }, { inject: this.isInjectable }),
                    config: { requestId: `post_searchObjectAttribute_${this.objId}` }
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
            },
            handleSliderBeforeClose () {
                const hasChanged = Object.keys(this.$refs.fieldForm.changedValues).length
                if (hasChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
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
    .options {
        padding: 10px 0;
    }
    .field-pre {
        display: inline-block;
        margin-right: -26px;
        padding: 0 6px;
        vertical-align: middle;
        line-height: 32px;
        border-radius: 4px;
        background-color: #a4aab3;
        color: #fff;
        font-size: 20px;
        transform: scale(0.5);
        transform-origin: left center;
        opacity: 0.4;
        &.en {
            margin-right: -40px;
        }
    }
    .field-id {
        display: inline-block;
        vertical-align: middle;
    }
    .field-required-icon {
        font-size: 20px;
        transform: scale(.5);
        transform-origin: left center;
    }
    .text-primary {
        cursor: pointer;
    }
</style>
