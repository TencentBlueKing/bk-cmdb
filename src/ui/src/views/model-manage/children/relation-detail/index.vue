<template>
    <div>
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
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading(['updateObjectAttribute', 'createObjectAttribute'])" @click="saveRelation">
                {{$t('ModelManagement["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('ModelManagement["取消"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            relation: {
                type: Object
            },
            isReadOnly: {
                type: Boolean,
                default: false
            },
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                modelList: [],
                selected: ''
            }
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createObjectAssociation',
                'updateObjectAssociation'
            ]),
            saveRelation () {

            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-relation-wrapper {
        padding: 20px;
    }
</style>
