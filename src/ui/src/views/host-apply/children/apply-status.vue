<template>
    <cmdb-dialog v-model="visible" :width="490" @close="handleEvent('return')">
        <div class="status status-loading" v-if="loading">
            <div class="status-inner">
                <img class="status-icon" src="../../../assets/images/icon/loading.svg">
                <p class="status-title">{{$t('正在应用中')}}</p>
            </div>
        </div>
        <div class="status status-loading" v-else-if="error">
            <div class="status-inner">
                <span class="status-icon bk-icon icon-close icon-error"></span>
                <p class="status-title">{{$t('应用异常')}}</p>
            </div>
        </div>
        <div class="status status-result" v-else>
            <div class="result-icon">
                <span class="status-icon bk-icon icon-check-1 icon-success" v-if="fail.length === 0"></span>
                <span class="status-icon bk-icon icon-exclamation icon-abnormal" v-else-if="fail.length > 0"></span>
                <span class="status-icon bk-icon icon-close icon-fail" v-else-if="success.length === 0"></span>
            </div>
            <p class="result-title">{{$t('应用成功')}}</p>
            <p class="result-subtitle" v-if="fail.length === 0">{{$t('保存策略并应用到当前模块下主机成功')}}</p>
            <i18n class="result-stat" tag="p" path="应用结果">
                <span class="result-count" place="success">{{success.length}}</span>
                <span :class="['result-count', { fail: fail.length > 0 }]" place="fail">{{fail.length}}</span>
            </i18n>
            <div class="result-options">
                <bk-button class="mr10" theme="primary" @click="handleEvent('view-host')" v-if="fail.length === 0">{{$t('立即查看主机')}}</bk-button>
                <bk-button class="mr10" theme="primary" @click="handleEvent('view-failed')" v-else>{{$t('查看失败')}}</bk-button>
                <bk-button @click="handleEvent('return')">{{$t('返回')}}</bk-button>
            </div>
        </div>
    </cmdb-dialog>
</template>

<script>
    export default {
        props: {
            request: {
                validator (request) {
                    return request instanceof Promise
                }
            }
        },
        data () {
            return {
                visible: false,
                loading: false,
                error: false,
                success: [],
                fail: []
            }
        },
        watch: {
            request (request) {
                this.initStatus()
            }
        },
        methods: {
            show () {
                this.visible = true
            },
            hide () {
                this.visible = false
            },
            async initStatus () {
                try {
                    this.loading = true
                    this.error = false
                    const results = await this.request
                    const success = []
                    const fail = []
                    results.forEach(result => {
                        result.error_code ? fail.push(result) : success.push(result)
                    })
                    this.success = success
                    this.fail = fail
                    this.loading = false
                } catch (e) {
                    this.loading = false
                    this.error = true
                    console.error(e)
                }
            },
            handleEvent (event) {
                if (event) {
                    this.$emit(event)
                }
                this.hide()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .status {
        overflow: hidden;
        text-align: center;
        height: 300px;

        .status-icon {
            display: block;
            width: 54px;
            height: 54px;
            margin: 0 auto;
            &:not(img) {
                line-height: 54px;
                border-radius: 50%;
                color: #FFF;
                font-size: 30px;
            }

            &.icon-error {
                background-color: $dangerColor;
            }
            &.icon-success {
                background-color: #2dcb56;
            }
            &.icon-fail {
                background-color: #ff5656;
            }
            &.icon-abnormal {
                background-color: #ffb848;
            }
        }
    }

    .status-loading {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 40px 0;

        .status-title {
            margin: 30px 0 0;
            font-size: 24px;
            color: #313238;
        }
    }

    .status-result {
        .result-icon {
            margin: 30px auto 0;
        }
        .result-title {
            margin: 20px auto 0;
            line-height: 30px;
            font-size: 24px;
            color: #313238;
        }
        .result-subtitle {
            font-size: 14px;
            color: $textColor;
            margin: 12px auto 0;
            & + .result-stat {
                margin-top: 4px;
            }
        }
        .result-stat {
            margin: 28px auto 0;
            font-size: 14px;
            color: $textColor;
            .result-count {
                padding: 0 2px;
                font-weight: bold;
                &.fail {
                    color: $dangerColor;
                }
            }
        }
        .result-options {
            font-size: 0;
            margin-top: 24px;

            .bk-button {
                min-width: 90px;
            }
        }
    }
</style>
