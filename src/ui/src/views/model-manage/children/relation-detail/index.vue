<template>
    <div class="model-relation-wrapper">
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["唯一标识"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input" >
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["别名"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input" >
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["模型源"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="selected"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="selected"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="selected"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源-目标约束"]')}}
                <span class="color-danger">*</span>
            </span>
            <bk-selector
                :list="modelList"
                :selected.sync="selected"
            ></bk-selector>
            <i class="bk-icon icon-info-circle"></i>
        </label>
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
</template>

<script>
    export default {
        data () {
            return {
                modelList: [],
                selected: ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-relation-wrapper {
        padding: 20px;
    }
    .radio-box {
        margin: 0 -20px;
        font-size: 0;
        .label-text {
            display: inline-block;
            padding-right: 10px;
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
</style>
