<template>
    <div class="relation-wrapper">
        <p class="operation-box">
            <bk-button type="primary" @click="slider.isShow = true">
                {{$t('ModelManagement["新增关联类型"]')}}
            </bk-button>
            <label class="search-input">
                <i class="bk-icon icon-search"></i>
                <input type="text" class="cmdb-form-input" :placeholder="$t('ModelManagement[\'请输入关联类型名称\']')">
            </label>
        </p>
        <cmdb-table
            :loading="$loading('getOperationLog')"
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination">
        </cmdb-table>
        <cmdb-slider
            class="relation-slider"
            :width="514"
            :title="slider.title"
            :isShow.sync="slider.isShow">
            <div slot="content" class="content">
                <label class="form-label">
                    <span class="label-text">
                        {{$t('ModelManagement["唯一标识"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                </label>
                <label class="form-label">
                    <span class="label-text">
                        {{$t('Hosts["名称"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                </label>
                <label class="form-label">
                    <span class="label-text">
                        {{$t('ModelManagement["源->目标描述"]')}}
                        <span class="color-danger">*</span>
                    </span>
                    <input type="text" class="cmdb-form-input">
                </label>
                <label class="form-label">
                    <span class="label-text">{{$t('ModelManagement["目标描述->源"]')}}</span>
                    <span class="color-danger">*</span>
                    <input type="text" class="cmdb-form-input">
                </label>
                <div class="radio-box overflow">
                    <label class="label-text">
                        {{$t('ModelManagement["是否有方向"]')}}<span class="text-desc">({{$t('ModelManagement["仅视图"]')}})</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["是，源指向目标"]')}}</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["否"]')}}</span>
                    </label>
                </div>
                <div class="radio-box">
                    <label class="label-text">
                        {{$t('ModelManagement["联动删除"]')}}
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["源不存在联动删除目标"]')}}</span>
                    </label>
                    <label class="cmdb-form-radio cmdb-radio-small">
                        <input type="radio">
                        <span class="cmdb-radio-text">{{$t('ModelManagement["目标不存在联动删除源"]')}}</span>
                    </label>
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
                    isShow: false,
                    title: this.$t('ModelManagement["新增关联类型"]')
                },
                table: {
                    header: [{
                        id: 'name',
                        name: this.$t('Hosts["名称"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["唯一标识"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["源->目标描述"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["目标描述->源"]')
                    }, {
                        id: 'name',
                        name: this.$t('ModelManagement["使用数"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Common["操作"]')
                    }],
                    list: [],
                    pagination: {
                        count: 0,
                        current: 1,
                        size: 10
                    },
                    defaultSort: '-op_time',
                    sort: '-op_time'
                }
            }
        }
    }
</script>


<style lang="scss" scoped>
    .operation-box {
        margin: 20px 0;
        font-size: 0;
        .search-input {
            position: relative;
            display: inline-block;
            margin-left: 10px;
            width: 300px;
            .icon-search {
                position: absolute;
                top: 9px;
                right: 10px;
                font-size: 18px;
                color: $cmdbBorderColor;
            }
            .cmdb-form-input {
                vertical-align: middle;
                padding-right: 36px;
            }
        }
    }
    .relation-slider {
        .content {
            padding: 20px 30px;
        }
        .radio-box {
            margin: 0 -30px;
            font-size: 0;
            .label-text {
                display: inline-block;
                padding-right: 2px;
                width: 130px;
                font-size: 14px;
                line-height: 32px;
                text-align: right;
                .text-desc {
                    color: $cmdbBorderColor;
                }
            }
            .cmdb-form-radio {
                margin-right: 10px;
                width: 170px;
                height: 32px;
            }
        }
    }
</style>
