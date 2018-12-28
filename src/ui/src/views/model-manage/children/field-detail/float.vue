<template>
    <div>
        <div class="form-label">
            <span class="label-text">{{$t('ModelManagement["最小值"]')}}</span>
            <div class="cmdb-form-item">
                <input type="text" class="cmdb-form-input"
                    v-model="localValue.min"
                    @input="handleInput"
                    v-validate="`float`"
                    :disabled="isReadOnly"
                    :name="'min'">
            </div>
        </div>
        <div class="form-label">
            <span class="label-text">{{$t('ModelManagement["最大值"]')}}</span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('max')}">
                <input type="text" class="cmdb-form-input"
                    v-model="localValue.max"
                    name="max"
                    @input="handleInput"
                    :disabled="isReadOnly"
                    v-validate="`float|isBigger:${localValue.min}`">
                <p class="form-error">{{errors.first('max')}}</p>
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
        &:last-child {
           margin: 0;
        }
    }
</style>
