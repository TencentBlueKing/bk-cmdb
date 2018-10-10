<template>
    <div>
        <div class="form-item has-right-content">
            <label class="form-label">{{$t('ModelManagement["最小值"]')}}</label>
            <div class="input-box">
                <input type="text" class="cmdb-form-input" placeholder=""
                    v-model="localValue.min"
                    @input="handleInput"
                    v-validate="`number`"
                    maxlength="11"
                    :disabled="isReadOnly"
                    :name="'min'">
                <span v-show="errors.has('min')" class="error-msg color-danger">{{ errors.first('min') }}</span>
            </div>
        </div>
        <div class="form-item">
            <label class="form-label">{{$t('ModelManagement["最大值"]')}}</label>
            <div class="input-box">
                <input type="text" class="cmdb-form-input" placeholder="" v-model="localValue.max"
                    name="max"
                    @input="handleInput"
                    :disabled="isReadOnly"
                    v-validate="`number|isBigger:${localValue.min}`">
                <span v-show="errors.has('max')" class="error-msg color-danger">{{ errors.first('max') }}</span>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                default: {
                    min: '',
                    max: ''
                }
            },
            isReadOnly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localValue: {
                    min: '',
                    max: ''
                }
            }
        },
        watch: {
            value: {
                handler () {
                    this.initValue()
                },
                deep: true
            }
        },
        created () {
            this.initValue()
        },
        methods: {
            initValue () {
                if (this.value === '' || this.value === null) {
                    this.localValue = {
                        min: '',
                        max: ''
                    }
                } else {
                    this.localValue = this.value
                }
            },
            async handleInput () {
                const res = await this.$validator.validateAll()
                if (res) {
                    this.$emit('input', this.localValue)
                }
            },
            validate () {
                return this.$validator.validateAll()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-label {
        vertical-align: top;
        line-height: 30px;
    }
</style>
