<template>
    <section class="export-content">
        <div class="header">
            <h1 class="title">{{title}}</h1>
            <div class="subtitle-wrapper">
                <i18n class="subtitle" tag="h2" path="分批下载副标题">
                    <strong class="count" place="count">{{count}}</strong>
                    <span place="limit">{{limit}}</span>
                </i18n>
                <span class="process-counter">
                    <span class="finished">{{finishedCount}}</span>
                    <i>&nbsp;/&nbsp;&nbsp;</i>
                    <span class="total">{{all.length}}</span>
                </span>
                <bk-button class="restart" v-if="hasError" text @click="processSchedule">{{$t('重试失败项')}}</bk-button>
            </div>
            <i class="close-trigger bk-icon icon-close"
                v-if="isFinished"
                @click="$emit('close')">
            </i>
            <bk-popconfirm class="close" v-else
                placement="top"
                trigger="click"
                :title="$t('关闭批量导出提示')"
                @confirm="$emit('close')">
                <i class="close-confirm-trigger bk-icon icon-close"></i>
            </bk-popconfirm>
        </div>
        <ul class="list" ref="list">
            <li class="list-item"
                v-for="(task, index) in all"
                ref="listItem"
                :key="index">
                <span :class="['state', task.state]">
                    <i :class="['state-icon', iconMapping[task.state]]" v-if="task.state !== 'waiting'"></i>
                    {{textMapdding[task.state]}}
                </span>
                <span class="info">
                    <span class="info-name">{{`${name}_download_${index + 1}.${ext}`}}</span>
                    <span class="info-error"
                        v-if="task === current && current.state === 'error'">
                        {{message}}
                    </span>
                </span>
            </li>
        </ul>
    </section>
</template>

<script>
    export default {
        props: {
            name: {
                type: String,
                default: 'host'
            },
            count: {
                type: Number,
                required: true
            },
            options: {
                type: Function,
                required: true
            },
            limit: {
                type: Number,
                default: 10000
            },
            ext: {
                type: String,
                default: 'xlsx'
            }
        },
        data () {
            return {
                current: null,
                message: null,
                queue: [],
                all: [],
                request: {
                    id: Symbol('id')
                },
                iconMapping: {
                    error: 'bk-icon icon-close-circle-shape',
                    finished: 'bk-icon icon-check-circle-shape',
                    pending: 'loading'
                },
                textMapdding: {
                    error: this.$t('失败'),
                    finished: this.$t('已完成'),
                    pending: this.$t('下载中'),
                    waiting: this.$t('等待中')
                }
            }
        },
        computed: {
            isFinished () {
                return this.all.every(task => task.state === 'finished')
            },
            finishedCount () {
                return this.all.filter(task => task.state === 'finished').length
            },
            hasError () {
                return this.all.some(task => task.state === 'error')
            },
            title () {
                if (this.hasError) {
                    return this.$t('下载失败')
                }
                if (this.isFinished) {
                    return this.$t('全部下载完成')
                }
                return this.$t('分批下载标题')
            }
        },
        mounted () {
            this.setupSchedule()
        },
        beforeDestroy () {
            this.$http.cancel(this.request.id)
        },
        methods: {
            setupSchedule () {
                const queue = new Array(Math.ceil(this.count / this.limit)).fill(null).map((_, index) => ({
                    state: 'waiting',
                    page: {
                        start: index * this.limit,
                        limit: this.limit
                    }
                }))
                this.queue = queue
                this.all = [...queue]
                this.processSchedule()
            },
            async processSchedule () {
                if (!this.queue.length) return
                try {
                    this.message = null
                    const [current] = this.queue.splice(0, 1)
                    this.current = current
                    this.current.state = 'pending'
                    this.$nextTick(this.syncScrollbar)
                    const options = this.options(current.page)
                    options.config = options.config || {}
                    options.config.requestId = this.request.id
                    await this.$http.download(options)
                    this.current.state = 'finished'
                    this.processSchedule()
                } catch (error) {
                    this.queue.unshift(this.current)
                    this.current.state = 'error'
                    this.message = error.message
                    console.error(error)
                }
            },
            syncScrollbar () {
                const index = this.all.indexOf(this.current)
                const item = this.$refs.listItem[index]
                const top = item.offsetTop + item.offsetHeight + 10 // margin
                const listViewport = this.$refs.list.offsetHeight
                const scrollHeight = top - listViewport
                if (scrollHeight > 0) {
                    this.$refs.list.scrollTop = scrollHeight
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .export-content {
        margin: -3px -24px 0;
        background-color: #fff;
        .header {
            padding: 0 60px;
            position: relative;
            .title {
                text-align: center;
                font-size: 22px;
                font-weight: normal;
                color: #313238;
            }
            .subtitle-wrapper {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-top: 32px;
                line-height: 20px;
                .subtitle {
                    flex: 1;
                    font-size: 14px;
                    font-weight: normal;
                    color: $textColor;
                    .count {
                        font-weight: bold;
                    }
                }
                .process-counter {
                    font-size: 14px;
                    color: $textColor;
                    font-weight: bold;
                    padding-right: 8px;
                    display: inline-flex;
                    .finished {
                        color: $successColor;
                    }
                }
                .restart {
                    font-size: 12px;
                    height: 20px;
                    margin-left: 4px;
                }
            }
            .close,
            .close-trigger{
                position: absolute;
                top: -25px;
                right: 5px;
            }
            .close-trigger,
            .close-confirm-trigger {
                color: #979ba5;
                width: 26px;
                height: 26px;
                line-height: 26px;
                text-align: center;
                border-radius: 50%;
                font-weight: 700;
                font-size: 22px;
                cursor: pointer;
                &:hover {
                    background-color: #f0f1f5;
                }
            }
        }
    }
    .list {
        position: relative;
        margin-top: 12px;
        padding: 0 60px;
        max-height: 300px;
        @include scrollbar-y;
        .list-item {
            display: flex;
            align-items: center;
            margin-bottom: 10px;
            height: 52px;
            border: 1px solid #dcdee5;
            border-radius: 2px;
            box-shadow: 0px 2px 4px 0px rgba(0,0,0,0.1);
        }
    }
    .state {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 60px;
        height: 100%;
        font-size: 12px;
        flex-direction: column;
        &.pending {
            background-color: #e1ecff;
        }
        &.finished {
            background-color: #e4faf0;
            .state-icon {
                color: $successColor;
            }
        }
        &.waiting {
            background-color: #f0f1f5;
        }
        &.error {
            background-color: #fedddc;
            .state-icon {
                color: #ea3636;
            }
        }
        .state-icon {
            font-size: 20px;
            &.loading {
                display: inline-block;
                width: 16px;
                height: 16px;
                background-color: transparent;
                background-image: url("../../assets/images/icon/loading.svg");
                background-position: center center;
                background-size: 16px;
                background-repeat: no-repeat;
            }
        }
    }
    .info {
        flex: 1;
        display: flex;
        justify-content: start;
        flex-direction: column;
        padding: 0 17px;
        .info-name {
            font-size: 14px;
            font-weight: bold;
            color: $textColor;
        }
        .info-error {
            font-size: 12px;
            color: $dangerColor;
        }
    }
</style>
