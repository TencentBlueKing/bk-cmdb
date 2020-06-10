<template>
    <div class="config-layout">
        <div class="config-form">
            <div class="inputs">
                <bk-input
                    :class="{ 'has-error': hasError }"
                    :input-style="{ height: $APP.height - 160 + 'px' }"
                    :type="'textarea'"
                    @input="handleInput"
                    @blur="handleBlur"
                    v-model="configValue">
                </bk-input>
            </div>
            <div class="buttons">
                <bk-button size="large" class="mr10" @click="handleReset">{{$t('取消')}}</bk-button>
                <bk-button theme="primary" size="large"
                    :disabled="hasError"
                    :loading="loading"
                    @click="handleSave">
                    {{$t('保存')}}
                </bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { Base64 } from 'js-base64'
    import { updateValidator } from '@/setup/validate'
    export default {
        data () {
            return {
                configValue: '',
                disabled: false,
                hasError: false,
                request: {
                    update: Symbol('updateConfig')
                },
                loading: false
            }
        },
        computed: {
            ...mapGetters(['config'])
        },
        created () {
            this.configValue = this.getJSONString(this.config)
        },
        methods: {
            handleInput (value) {
                try {
                    JSON.parse(value)
                    this.hasError = false
                } catch (e) {
                    this.hasError = true
                }
            },
            handleBlur (value) {
                if (!this.hasError) {
                    const config = JSON.parse(value)
                    this.configValue = this.getJSONString(config)
                }
            },
            async handleSave () {
                this.loading = true
                try {
                    const configData = JSON.parse(this.configValue)
                    const { validationRules } = configData
                    for (const rule of Object.values(validationRules)) {
                        rule.value = Base64.encode(rule.value)
                    }
                    await this.$store.dispatch('updateConfig', {
                        params: configData,
                        config: { requestId: this.request.update }
                    })

                    updateValidator()

                    this.$success(this.$t('保存成功'))
                } catch (e) {
                    console.error(e)
                } finally {
                    this.loading = false
                }
            },
            getJSONString (data) {
                return JSON.stringify(data, null, 4)
            },
            handleReset () {
                this.configValue = this.getJSONString(this.config)
                this.hasError = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .config-layout {
        padding: 15px 20px 0;

        .config-form {
            width: 60%;
            margin: 0 auto;

            /deep/ .inputs {
                .bk-form-textarea {
                    font-size: 16px;
                    font-family: Consolas, 'Courier New', monospace;
                    line-height: 2;
                    white-space: nowrap;
                }
                .has-error {
                    .bk-textarea-wrapper {
                        border: 1px solid #f00 !important;
                        &:focus-within {
                            border: 1px solid #f00!important;
                        }
                    }
                }
            }
            .buttons {
                text-align: right;
                margin-top: 12px;
            }
        }
    }
</style>
