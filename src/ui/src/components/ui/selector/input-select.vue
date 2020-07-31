<template>
    <div class="input-select">
        <bk-select class="select-box"
            :clearable="clearable"
            :searchable="searchable"
            :disabled="disabled"
            :value="localValue"
            v-bind="$attrs"
            @change="handleSelected">
            <div class="input-box" slot="trigger">
                <input :class="['input-text', { 'custom-error': errors.has(name) }]"
                    autocomplete="off"
                    :name="name"
                    :placeholder="placeholder || $t('请选择或输入内容')"
                    :disabled="disabled"
                    v-validate="validate"
                    v-model="localValue">
            </div>
            <bk-option v-for="(option, index) in options"
                :key="index"
                :id="option[settingKey]"
                :name="option[settingKey]">
            </bk-option>
        </bk-select>
        <span class="custom-form-error"
            :title="errors.first(name)">
            {{errors.first(name)}}
        </span>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-input-select',
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            name: {
                type: [String, Number],
                default: ''
            },
            validate: {
                type: Object,
                default: () => {}
            },
            options: {
                type: Array,
                default: () => []
            },
            settingKey: {
                type: String,
                default: 'id'
            },
            displayKey: {
                type: String,
                default: 'name'
            },
            placeholder: {
                type: String,
                default: ''
            },
            searchable: {
                type: Boolean,
                default: false
            },
            clearable: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                localValue: this.value
            }
        },
        watch: {
            value (value) {
                this.localValue = value
            },
            localValue (value) {
                this.$emit('input', value)
                this.$emit('on-change', value)
            }
        },
        methods: {
            handleSelected (value) {
                this.localValue = value
            }
        }
    }
</script>

<style lang="scss" scoped>
    .input-select {
        position: relative;
        .select-box {
            width: 100%;
            border: none !important;

            &.bk-select-small {
                .input-box {
                    .input-text {
                        height: 26px;
                        font-size: 12px;
                        line-height: 26px;
                    }
                }
            }
        }
        .input-box {
            position: relative;
            z-index: 2;
            .input-text {
                width: 100%;
                height: 32px;
                line-height: 30px;
                padding: 0 10px;
                font-size: 14px;
                border: 1px solid #c4c6cc;
                border-radius: 2px;
                outline: none;
                &::placeholder {
                    color: #c4c6cc;
                }
                &[disabled] {
                    color: #c4c6cc;
                    background-color: #fafbfd!important;
                    cursor: not-allowed;
                    border-color: #dcdee5!important;
                }
            }
        }
        .custom-form-error {
            position: absolute;
            top: 100%;
            left: 0;
            line-height: 14px;
            font-size: 12px;
            color: #ff5656;
            max-width: 100%;
            @include ellipsis;
        }
    }
</style>
