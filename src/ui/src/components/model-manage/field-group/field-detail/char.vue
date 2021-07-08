<template>
    <div class="form-label">
        <span class="label-text">{{$t('正则校验')}}</span>
        <textarea
            v-model="localValue"
            :disabled="isReadOnly"
            v-validate="'remoteRegular'"
            data-vv-validate-on="blur"
            data-vv-name="regular"
            @input="handleInput"
        ></textarea>
        <p class="form-error" v-if="errors.has('regular')">{{errors.first('regular')}}</p>
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            },
            isReadOnly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localValue: ''
            }
        },
        watch: {
            value () {
                this.localValue = this.value === '' ? '' : this.value
            }
        },
        created () {
            this.localValue = this.value === '' ? '' : this.value
        },
        methods: {
            handleInput () {
                this.$emit('input', this.localValue)
            },
            validate () {
                return this.$validator.validateAll()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-label {
        position: relative;
        .form-error {
            font-size: 12px;
            color: $cmdbDangerColor;
        }
    }
</style>
