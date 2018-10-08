<template>
    <ul class="label-wrapper clearfix">
        <li v-if="isOnlyShow">
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" v-model="localValue.isonly" :disabled="isReadOnly">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["唯一"]')}}</span>
            </label>
        </li>
        <li v-if="isRequiredShow">
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" v-model="localValue.isrequired" :disabled="isReadOnly">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["必填"]')}}</span>
            </label>
        </li>
        <li v-if="isEditableShow">
            <label class="cmdb-form-checkbox cmdb-checkbox-small">
                <input type="checkbox" v-model="localValue.editable" :disabled="isReadOnly">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["可编辑"]')}}</span>
            </label>
        </li>
    </ul>
</template>

<script>
    export default {
        props: {
            isReadOnly: {
                type: Boolean,
                default: false
            },
            type: {
                type: String,
                required: true
            },
            editable: {
                type: Boolean,
                default: true
            },
            isrequired: {
                type: Boolean,
                default: false
            },
            isonly: {
                type: Boolean,
                default: false
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
                ],
                localValue: {
                    editable: this.editable,
                    isrequired: this.isrequired,
                    isonly: this.isonly
                }
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
        },
        watch: {
            editable (editable) {
                this.localValue.editable = editable
            },
            isrequired (isrequired) {
                this.localValue.isrequired = isrequired
            },
            isonly (isonly) {
                this.localValue.isonly = isonly
            },
            'localValue.editable' (editable) {
                this.$emit('update:editable', editable)
            },
            'localValue.isrequired' (isrequired) {
                if (!isrequired && this.isOnlyShow) {
                    this.localValue.isonly = false
                }
                this.$emit('update:isrequired', isrequired)
            },
            'localValue.isonly' (isonly) {
                if (isonly && this.isRequiredShow) {
                    this.localValue.isrequired = true
                }
                this.$emit('update:isonly', isonly)
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
                input {
                    margin-right: 5px;
                }
            }
        }
    }
</style>
