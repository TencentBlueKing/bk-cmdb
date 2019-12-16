<template>
    <div>
        <div class="form-label">
            <span class="label-text">{{$t('最小值')}}</span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('min') }">
                <bk-input type="text" class="cmdb-form-input"
                    v-model="localValue.min"
                    @input="handleInput"
                    v-validate="`number`"
                    maxlength="11"
                    :disabled="isReadOnly"
                    :name="'min'">
                </bk-input>
                <p class="form-error">{{errors.first('min')}}</p>
            </div>
        </div>
        <div class="form-label">
            <span class="label-text">{{$t('最大值')}}</span>
            <div class="cmdb-form-item" :class="{ 'is-error': errors.has('max') }">
                <bk-input type="text" class="cmdb-form-input"
                    v-model="localValue.max"
                    name="max"
                    @input="handleInput"
                    :disabled="isReadOnly"
                    v-validate="`number|isBigger:${localValue.min}`">
                </bk-input>
                <p class="form-error">{{errors.first('max')}}</p>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                type: [Object, String],
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
                    const min = this.localValue.min
                    const max = this.localValue.max
                    this.$emit('input', {
                        min: min ? Number(min) : null,
                        max: max ? Number(max) : null
                    })
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
