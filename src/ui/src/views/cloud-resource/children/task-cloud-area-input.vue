<template>
    <div class="cloud-area-input-wrapper" v-if="display === 'input'">
        <cmdb-auth :auth="{ type: $OPERATION.C_CLOUD_AREA }">
            <input class="cloud-area-input"
                slot-scope="{ disabled }"
                v-model.trim="localValue"
                :readonly="readonly"
                :disabled="disabled"
                :class="{
                    'has-tips': hasTips,
                    'has-error': error
                }"
                :placeholder="$t('请输入xx', { name: $t('云区域') })">
        </cmdb-auth>
        <i class="tips-icon icon icon-cc-tips"
            v-if="readonly"
            v-bk-tooltips="{
                content: $t('VPC已绑定云区域')
            }">
        </i>
        <i class="tips-icon bk-icon icon-exclamation-circle-shape"
            v-else-if="error"
            v-bk-tooltips="{
                content: error
            }">
        </i>
    </div>
    <span class="cloud-area-info" v-bk-overflow-tips v-else>{{localValue}}</span>
</template>

<script>
    import symbols from '../common/symbol'
    import Bus from '@/utils/bus'
    export default {
        props: {
            id: Number,
            value: String,
            error: [Boolean, String],
            mode: {
                type: String,
                default: 'create',
                validator (val) {
                    return ['create', 'read'].includes(val)
                }
            },
            display: {
                type: String,
                default: 'input'
            }
        },
        data () {
            return {
                list: []
            }
        },
        computed: {
            readonly () {
                return this.id !== -1
            },
            localValue: {
                get () {
                    const area = this.list.find(area => area.bk_cloud_id === this.id)
                    return area ? area.bk_cloud_name : this.value
                },
                set (value) {
                    this.$emit('input', value)
                }
            },
            hasTips () {
                return this.error || this.readonly
            }
        },
        async created () {
            Bus.$on('refresh-cloud-area', this.refresh)
            try {
                const { info } = await this.getList()
                this.list = info
            } catch (e) {
                console.error(e)
            }
        },
        beforeDestroy () {
            Bus.$off('refresh-cloud-area', this.refresh)
        },
        methods: {
            refresh () {
                this.$http.cancelCache(symbols.get('cloudArea'))
                this.$nextTick(this.getList)
            },
            getList () {
                return this.$store.dispatch('cloud/area/findMany', {
                    params: {
                        page: {}
                    },
                    config: {
                        requestId: symbols.get('cloudArea'),
                        fromCache: true
                    }
                })
            },
            validate () {
                if (this.mode === 'read') {
                    return true
                }
                return !!this.localValue.length
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
    .cloud-area-info {
        display: block;
        @include ellipsis;
    }
</style>
