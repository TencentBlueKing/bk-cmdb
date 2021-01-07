<template>
    <div class="form-bool-layout">
        <span class="default">{{$t('默认值')}}</span>
        <bk-switcher
            size="small"
            theme="primary"
            :value="localValue"
            :disabled="isReadonly"
            @change="handleChange">
        </bk-switcher>
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                type: [String, Boolean],
                default: false
            },
            isReadonly: Boolean
        },
        data () {
            return {
                localValue: false
            }
        },
        watch: {
            value: {
                immediate: true,
                handler (value) {
                    this.localValue = typeof value === 'boolean' ? value : false
                    // 将空字符转为false
                    this.handleChange(this.localValue)
                }
            }
        },
        methods: {
            handleChange (selected) {
                this.$emit('input', selected)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-bool-layout {
        .default {
            display: block;
            line-height: 36px;
            font-size: 14px;
        }
    }
</style>
