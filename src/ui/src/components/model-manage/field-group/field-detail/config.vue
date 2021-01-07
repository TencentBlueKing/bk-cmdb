<template>
    <div class="form-label">
        <span class="label-text">{{$t('字段设置')}}</span>
        <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="isEditableShow">
            <input type="checkbox" tabindex="-1" v-model="localValue.editable" :disabled="isReadOnly || ispre">
            <span class="cmdb-checkbox-text">{{$t('可编辑')}}</span>
        </label>
        <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="isRequiredShow && !isMainLineModel">
            <input type="checkbox" tabindex="-1" v-model="localValue.isrequired" :disabled="isReadOnly || ispre">
            <span class="cmdb-checkbox-text">{{$t('必填')}}</span>
        </label>
    </div>
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
            isMainLineModel: {
                type: Boolean,
                default: false
            },
            ispre: Boolean
        },
        data () {
            return {
                editableMap: [
                    'singlechar',
                    'int',
                    'float',
                    'enum',
                    'date',
                    'time',
                    'longchar',
                    'objuser',
                    'timezone',
                    'bool',
                    'list',
                    'organization'
                ],
                isrequiredMap: [
                    'singlechar',
                    'int',
                    'float',
                    'date',
                    'time',
                    'longchar',
                    'objuser',
                    'timezone',
                    'list',
                    'organization'
                ],
                localValue: {
                    editable: this.editable,
                    isrequired: this.isrequired
                }
            }
        },
        computed: {
            isEditableShow () {
                return this.editableMap.indexOf(this.type) !== -1
            },
            isRequiredShow () {
                return this.isrequiredMap.indexOf(this.type) !== -1
            }
        },
        watch: {
            editable (editable) {
                this.localValue.editable = editable
            },
            isrequired (isrequired) {
                this.localValue.isrequired = isrequired
            },
            'localValue.editable' (editable) {
                this.$emit('update:editable', editable)
            },
            'localValue.isrequired' (isrequired) {
                if (!isrequired && this.isOnlyShow) {
                    this.localValue.isonly = false
                }
                this.$emit('update:isrequired', isrequired)
            }
        }
    }
</script>
