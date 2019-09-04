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
                    {{$t('新建字段')}}
                </bk-button>
            </span>
        </div>
        <bk-table
            class="field-table"
            v-bkloading="{ isLoading: $loading(`post_searchObjectAttribute_${objId}`) }"
            :data="table.list"
            :max-height="$APP.height - 300"
            :row-style="{
                cursor: 'pointer'
            }"
            @sort-change="handleSortChange"
            @cell-click="handleShowDetails">
            <bk-table-column prop="isrequired" :label="$t('必填')" width="60" sortable="custom">
                <template slot-scope="{ row }">
                    <i class="field-required-icon bk-icon icon-check-1" v-if="row.isrequired"></i>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_property_id"
                min-width="110"
                sortable="custom"
                class-name="is-highlight"
                :label="$t('唯一标识')">
                <template slot-scope="{ row }">
                    <span
                        v-if="row['ispre']"
                        :class="['field-pre', $i18n.locale]">
                        {{$t('内置')}}
                    </span>
                    <span class="field-id">{{row['bk_property_id']}}</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="bk_property_name" :label="$t('名称')" sortable="custom"></bk-table-column>
            <bk-table-column prop="bk_property_type" :label="$t('字段类型')" sortable="custom">
                <template slot-scope="{ row }">
                    <span>{{fieldTypeMap[row['bk_property_type']]}}</span>
                </template>
            </bk-table-column>
            <bk-table-column prop="create_time" :label="$t('创建时间')" sortable="custom">
                <template slot-scope="{ row }">
                    {{$tools.formatTime(row['create_time'])}}
                </template>
            </bk-table-column>
            <bk-table-column prop="operation" :label="$t('操作')" v-if="updateAuth">
                <template slot-scope="{ row }">
                    <button class="text-primary mr10 operation-btn"
                        :disabled="!isFieldEditable(row)"
                        @click.stop="editField(row)">
                        {{$t('编辑')}}
                    </button>
                    <button class="text-primary operation-btn"
                        :disabled="!isFieldEditable(row)"
                        @click.stop="deleteField(row)">
                        {{$t('删除')}}
                    </button>
                </template>
            </bk-table-column>
        </bk-table>
        <bk-sideslider
            v-transfer-dom
            :width="450"
            :title="slider.title"
            :is-show.sync="slider.isShow"
            :before-close="handleSliderBeforeClose">
            <the-field-detail
                ref="fieldForm"
                slot="content"
                v-if="slider.isShow"
                :is-read-only="isReadOnly"
                :is-edit-field="slider.isEditField"
                :field="slider.curField"
                :only-read-of-type="slider.type"
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
                    title: this.$t('新建字段'),
                    curField: {},
                    type: ''
                },
                fieldTypeMap: {
                    'singlechar': this.$t('短字符'),
                    'int': this.$t('数字'),
                    'float': this.$t('浮点'),
                    'enum': this.$t('枚举'),
                    'date': this.$t('日期'),
                    'time': this.$t('时间'),
                    'longchar': this.$t('长字符'),
                    'objuser': this.$t('用户'),
                    'timezone': this.$t('时区'),
                    'bool': 'bool'
                },
                table: {
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
                this.slider.title = this.$t('新建字段')
                this.slider.curField = {}
                this.slider.type = false
                this.slider.isShow = true
            },
            editField (item) {
                this.slider.isEditField = true
                this.slider.isReadOnly = this.isReadOnly
                this.slider.title = this.$t('编辑字段')
                this.slider.curField = item
                this.slider.type = false
                this.slider.isShow = true
            },
            deleteField (field) {
                this.$bkInfo({
                    title: this.$tc('确定删除字段？', field['bk_property_name'], { name: field['bk_property_name'] }),
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
                sort = this.$tools.getSort(sort)
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
            },
            handleShowDetails (row, column, cell) {
                if (column.property === 'operation') return
                this.slider.isEditField = true
                this.slider.isReadOnly = true
                this.slider.title = this.$t('字段详情')
                this.slider.curField = row
                this.slider.type = true
                this.slider.isShow = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .options {
        padding: 20px 0 14px;
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
    .operation-btn[disabled] {
        color: #dcdee5 !important;
        opacity: 1 !important;
    }
</style>
