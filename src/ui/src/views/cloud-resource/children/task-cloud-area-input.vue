<template>
    <div class="cloud-area-input-wrapper">
        <input class="cloud-area-input"
            v-model="localValue"
            :disabled="disabled"
            :readonly="readonly"
            :class="{
                'has-tips': hasTips,
                'has-error': error
            }">
        <i class="tips-icon icon icon-cc-tips"
            v-if="!!cloudId"
            v-bk-tooltips="{
                content: $t('VPC已绑定云区域')
            }">
        </i>
        <i class="tips-icon bk-icon icon-exclamation-circle-shape"
            v-else-if="error"
            v-bk-tooltips="{
                content: errorMessage || $t('请填写云区域')
            }">
        </i>
    </div>
</template>

<script>
    import CacheLoader from '@/utils/cache-loader'
    const LOADER_ID = Symbol('cloudArea')
    export default {
        props: {
            value: String,
            cloudId: Number,
            disabled: Boolean,
            readonly: Boolean,
            errorMessage: String,
            mode: {
                type: String,
                default: 'create',
                validator (val) {
                    return ['create', 'read'].includes(val)
                }
            }
        },
        data () {
            return {
                error: true,
                cloudAreaList: []
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                }
            },
            hasTips () {
                return this.error || !!this.cloudId
            }
        },
        async created () {
            try {
                const { info } = await CacheLoader.use(LOADER_ID, this.dataRequest)
                this.cloudAreaList = info
            } catch (e) {
                console.error(e)
            }
        },
        methods: {
            dataRequest () {
                return this.$store.dispatch('cloud/area/findMany', {
                    params: {
                        page: {}
                    }
                })
            },
            validate () {
                if (this.mode === 'read') {
                    return true
                }
                return !!this.value.length
            },
            getMessage (val) {
                if (!val.length) {
                    return this.$t('请填写云区域')
                }
                return this.errorMessage
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cloud-area-input-wrapper {
        position: relative;
        .cloud-area-input {
            height: 26px;
            line-height: normal;
            color: $textColor;
            background-color: #fff;
            border-radius: 2px;
            width: 100%;
            font-size: 12px;
            border: 1px solid #c4c6cc;
            padding: 0 10px;
            text-align: left;
            outline: none;
            resize: none;
            transition: border .2s linear;
            &:focus {
                border-color: $primaryColor;
            }
            &:disabled {
                cursor: not-allowed;
                background-color: #fafbfd;
                border-color: #dcdee5;
            }
            &[readonly] {
                cursor: default;
                background-color: #fafbfd;
                border-color: #dcdee5;
            }
            &.has-tips {
                padding-right: 20px;
            }
            &.has-error {
                border-color: $dangerColor;
                & ~ .tips-icon {
                    color: $dangerColor;
                }
            }
        }
        .tips-icon {
            position: absolute;
            font-size: 12px;
            right: 7px;
            top: 7px;
        }
    }
</style>
