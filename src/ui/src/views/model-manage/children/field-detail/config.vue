<template>
    <div class="form-label">
        <span class="label-text">{{$t('ModelManagement["字段设置"]')}}</span>
        <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="isEditableShow">
            <input type="checkbox" v-model="localValue.editable" :disabled="isReadOnly">
            <span class="cmdb-checkbox-text">{{$t('ModelManagement["可编辑"]')}}</span>
        </label>
        <label class="cmdb-form-checkbox cmdb-checkbox-small" v-if="isRequiredShow">
            <input type="checkbox" v-model="localValue.isrequired" :disabled="isReadOnly">
            <span class="cmdb-checkbox-text">{{$t('ModelManagement["必填"]')}}</span>
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
            }
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
                    'bool'
                ],
                isrequiredMap: [
                    'singlechar',
                    'int',
                    'float',
                    'date',
                    'time',
                    'longchar',
                    'objuser',
                    'timezone'
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
