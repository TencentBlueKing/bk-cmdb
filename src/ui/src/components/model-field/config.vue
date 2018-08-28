<template>
    <ul class="label-wrapper clearfix">
        <li v-if="isEditableShow">
            <label for="editable" class="cmdb-form-checkbox cmdb-checkbox-small">
                <span class="cmdb-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</span>
                <input id="editable" type="checkbox">
            </label>
        </li>
        <li v-if="isRequiredShow">
            <label for="isrequired" class="cmdb-form-checkbox cmdb-checkbox-small">
                <span class="cmdb-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</span>
                <input id="isrequired" type="checkbox">
            </label>
        </li>
        <li v-if="isOnlyShow">
            <label for="isonly" class="cmdb-form-checkbox cmdb-checkbox-small">
                <span class="cmdb-checkbox-text mr5">{{$t('ModelManagement["是否唯一"]')}}</span>
                <input id="isonly" type="checkbox">
            </label>
        </li>
    </ul>
</template>

<script>
    export default {
        props: {
            type: {
                type: String,
                required: true
            }
        },
        data () {
            return {
                editableMap: [
                    'singlechar',
                    'int',
                    'enum',
                    'date',
                    'time',
                    'longchar',
                    'singleasst',
                    'multiasst',
                    'objuser',
                    'timezone',
                    'bool'
                ],
                isrequiredMap: [
                    'singlechar',
                    'int',
                    'date',
                    'time',
                    'longchar',
                    'objuser',
                    'timezone'
                ],
                isonlyMap: [
                    'singlechar',
                    'int',
                    'longchar'
                ]
            }
        },
        computed: {
            isEditableShow () {
                return this.editableMap.indexOf(this.type) !== -1
            },
            isRequiredShow () {
                return this.isrequiredMap.indexOf(this.type) !== -1
            },
            isOnlyShow () {
                return this.isonlyMap.indexOf(this.type) !== -1
            }
        }
    }
</script>

<style lang="scss" scoped>
    .label-wrapper {
        float: left;
        vertical-align: middle;
        li {
            float: left;
            height: 30px;
            label {
                line-height: 1;
            }
        }
    }
</style>
