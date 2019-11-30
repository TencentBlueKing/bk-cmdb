<template>
    <cmdb-dialog v-model="visible" :width="400" @close="handleEvent('return')">
        <div class="status status-loading" v-if="loading">
            <img class="status-icon" src="../../assets//images/icon/loading.svg">
            <p class="status-title">{{$t('正在应用中')}}</p>
        </div>
        <div class="status status-loading" v-else-if="error">
            <span class="status-icon bk-icon icon-close"></span>
            <p class="status-title">{{$t('应用异常')}}</p>
        </div>
        <div class="status status-result" v-else-if="success.length && !fail.length">
            <p class="result-title">{{$t('应用成功')}}</p>
            <i18n class="result-subtitle" tag="p" path="应用结果">
                <span class="result-count" place="success">{{success.length}}</span>
                <span class="result-count" place="fail">0</span>
            </i18n>
            <div class="result-options">
                <bk-button class="mr10" theme="primary" @click="handleEvent('view-host')">{{$t('立即查看主机')}}</bk-button>
                <bk-button @click="handleEvent('return')">{{$t('返回')}}</bk-button>
            </div>
        </div>
        <div class="status status-result" v-else>
            <p class="result-title">{{$t('应用完成')}}</p>
            <i18n class="result-subtitle" tag="p" path="应用结果">
                <span class="result-count" place="success">{{success.length}}</span>
                <span class="result-count fail" place="fail">{{fail.length}}</span>
            </i18n>
            <div class="result-options">
                <bk-button class="mr10" theme="primary" @click="handleEvent('view-failed')">{{$t('查看失败')}}</bk-button>
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
        height: 220px;
    }
    .status-loading {
        padding: 40px 0;
        .status-icon {
            display: block;
            width: 58px;
            height: 58px;
            margin: 0 auto;
            &:not(img) {
                line-height: 58px;
                border-radius: 50%;
                color: #FFF;
                font-size: 30px;
                background-color: $dangerColor;
            }
        }

        .status-title {
            margin: 30px 0 0;
            font-size: 24px;
            color: #313238;
        }
    }

    .status-result {
        .result-title {
            margin: 50px auto 0;
            line-height: 30px;
            font-size: 24px;
            color: #313238;
        }
        .result-subtitle {
            margin: 30px auto 0;
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
            margin-top: 28px;
        }
    }
</style>
