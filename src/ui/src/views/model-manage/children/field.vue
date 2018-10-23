<template>
    <div class="model-field-wrapper">
        <bk-button class="create-btn" type="primary" @click="showSlider">
            {{$t('ModelManagement["新建字段"]')}}
        </bk-button>
        <cmdb-table
            class="field-table"
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :wrapperMinusHeight="220"
            @handlePageChange="handlePageChange"
            @handleSizeChange="handleSizeChange"
            @handleSortChange="handleSortChange">
        </cmdb-table>
        <cmdb-slider
            :width="514"
            :isShow="slider.isShow"
        >
            <div slot="content" class="slider-content">
                <label class="form-label">
                    <span class="label-text">
                        {{$t('ModelManagement["唯一标识"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                    <i class="bk-icon icon-info-circle"></i>
                </label>
                <label class="form-label">
                    <span class="label-text">
                        {{$t('ModelManagement["名称"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                    <i class="bk-icon icon-info-circle"></i>
                </label>
                <div class="form-label">
                    <span class="label-text">
                        {{$t('ModelManagement["字段类型"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <bk-selector
                        :list="fieldTypeList"
                        :selected.sync="fieldInfo['bk_property_type']"
                    ></bk-selector>
                    <i class="bk-icon icon-info-circle"></i>
                </div>
                <div class="field-detail">
                    <div class="form-label">
                        <span class="label-text">{{$t('ModelManagement["字段设置"]')}}</span>
                        <label class="cmdb-form-checkbox cmdb-checkbox-small">
                            <input type="checkbox">
                            <span class="cmdb-form-text">{{$t('ModelManagement["可编辑"]')}}</span>
                        </label>
                        <label class="cmdb-form-checkbox cmdb-checkbox-small">
                            <input type="checkbox">
                            <span class="cmdb-form-text">{{$t('ModelManagement["必填"]')}}</span>
                        </label>
                    </div>
                    <div class="form-label">
                        <span class="label-text">{{$t('ModelManagement["正则校验"]')}}</span>
                        <textarea name="" id="" cols="30" rows="10"></textarea>
                    </div>
                </div>
                <label class="form-label">
                    <span class="label-text">
                        {{$t('ModelManagement["单位"]')}}
                    </span>
                    <input type="text" class="cmdb-form-input">
                    <i class="bk-icon icon-info-circle"></i>
                </label>
                <div class="form-label">
                    <span class="label-text">{{$t('ModelManagement["用户提示"]')}}</span>
                    <textarea name="" id="" cols="30" rows="10"></textarea>
                </div>
            </div>
        </cmdb-slider>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                slider: {
                    isShow: false
                },
                table: {
                    header: [{
                        id: 'name',
                        name: this.$t('ModelManagement["字段名称"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["类型"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["唯一"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["必填"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["字段描述"]')
                    }, {
                        id: 'last_time',
                        name: this.$t('ModelManagement["创建时间"]')
                    }, {
                        id: 'operation',
                        name: this.$t('ModelManagement["操作"]')
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time'
                },
                fieldTypeList: [{
                    id: 'singlechar',
                    name: this.$t('ModelManagement["短字符"]')
                }, {
                    id: 'int',
                    name: this.$t('ModelManagement["数字"]')
                }, {
                    id: 'enum',
                    name: this.$t('ModelManagement["枚举"]')
                }, {
                    id: 'date',
                    name: this.$t('ModelManagement["日期"]')
                }, {
                    id: 'time',
                    name: this.$t('ModelManagement["时间"]')
                }, {
                    id: 'longchar',
                    name: this.$t('ModelManagement["长字符"]')
                }, {
                    id: 'objuser',
                    name: this.$t('ModelManagement["用户"]')
                }, {
                    id: 'timezone',
                    name: this.$t('ModelManagement["时区"]')
                }, {
                    id: 'bool',
                    name: 'bool'
                }],
                fieldInfo: {
                    bk_property_name: '',
                    bk_property_id: '',
                    unit: '',
                    placeholder: '',
                    bk_property_type: 'singlechar',
                    editable: true,
                    isrequired: false,
                    isonly: false,
                    option: '',
                    bk_asst_obj_id: ''
                }
            }
        },
        methods: {
            showSlider () {
                this.slider.isShow = true
            },
            handlePageChange (current) {
                this.pagination.current = current
                this.refresh()
            },
            handleSizeChange (size) {
                this.pagination.size = size
                this.handlePageChange(1)
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
    .slider-content {
        padding: 20px;
        .form-label {
            .cmdb-form-input {
                width: calc(100% - 145px);
            }
            textarea {
                padding: 10px;
                width: 329px;
                height: 84px;
                font-size: 14px;
                border: 1px solid $cmdbBorderColor;
                border-radius: 2px;
                outline: none;
                resize: none;
            }
        }
        .label-text {
            width: 110px;
        }
        .icon-info-circle {
            font-size: 18px;
            color: $cmdbBorderColor;
            padding-left: 5px;
        }
        .field-detail {
            margin-bottom: 20px;
            padding: 20px 0;
            background: #f3f8ff;
            .form-label:last-child {
                margin: 0;
            }
            .label-text {
                vertical-align: top;
            }
            .cmdb-form-checkbox {
                width: 90px;
                vertical-align: middle;
            }
        }
    }
</style>
