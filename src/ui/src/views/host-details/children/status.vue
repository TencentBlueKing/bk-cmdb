<template>
    <div :class="['status', { 'is-offline': !snapshot }]" v-bkloading="{ isLoading: $loading('getHostSnapshot') }">
        <div class="status-info" v-if="snapshot && !$loading('getHostSnapshot')">
            <h2 class="info-title">
                <span>{{$t('基本值')}}</span>
                <span class="update-time fr">
                    <span class="time-label">{{$t('最近更新时间')}}：</span>
                    <span class="time-value">{{snapshot.upTime}}</span>
                </span>
            </h2>
            <ul class="info-list clearfix">
                <li class="info-item fl">
                    <span class="item-label">{{$t('总流入量')}}：</span>
                    <span class="item-value">{{(snapshot.rcvRate / 100).toFixed(2)}}Mb/s</span>
                </li>
                <li class="info-item fl">
                    <span class="item-label">{{$t('启动时间')}}：</span>
                    <span class="item-value">{{$tools.formatTime(snapshot.bootTime * 1000)}}</span>
                </li>
                <li class="info-item fl">
                    <span class="item-label">{{$t('总流出量')}}：</span>
                    <span class="item-value">{{(snapshot.sendRate / 100).toFixed(2)}}Mb/s</span>
                </li>
                <li class="info-item fl">
                    <span class="item-label">{{$t('磁盘总量')}}：</span>
                    <span class="item-value">{{snapshot.Disk}}GB</span>
                </li>
                <li class="info-item fl">
                    <span class="item-label">{{$t('内存总量')}}：</span>
                    <span class="item-value">{{(snapshot.Mem / 1024).toFixed(2)}}GB</span>
                </li>
                <li class="info-item fl" v-if="!isWindows">
                    <span class="item-label">loadavg：</span>
                    <span class="item-value">{{snapshot.loadavg}}</span>
                </li>
            </ul>
            <div class="info-gauge">
                <div class="gauge-chart" ref="cpuChart"></div>
                <div class="gauge-chart" ref="memoryChart"></div>
                <div class="gauge-chart" ref="diskChart"></div>
            </div>
        </div>
        <div class="status-offline" v-else-if="!$loading('getHostSnapshot')">
            <div class="offline-image"></div>
            <p class="offline-text">
                {{$t('主机快照提示语')}}
                <a href="javascript:void(0)" @click="openAgentApp">
                    <i class="icon-cc-skip"></i>
                    {{$t('点此进入节点管理')}}
                </a>
            </p>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-host-status',
        data () {
            return {
                snapshot: null,
                Echarts: null
            }
        },
        computed: {
            info () {
                return this.$store.state.hostDetails.info || {}
            },
            isWindows () {
                return this.info.host.bk_os_type === 'windows'
            },
            id () {
                return this.$route.params.id
            }
        },
        async mounted () {
            try {
                const [Echarts, snapshot] = await Promise.all([
                    import(/* webpackChunkName: "echart" */ 'echarts'),
                    this.getHostSnapshot()
                ])
                this.Echarts = Echarts
                this.snapshot = snapshot
                this.$nextTick(() => {
                    this.initCharts()
                })
            } catch (e) {
                console.log(e)
                this.snapshot = null
            }
        },
        methods: {
            async getHostSnapshot () {
                const snapshot = await this.$store.dispatch('hostSearch/getHostSnapshot', {
                    hostId: this.id,
                    config: {
                        cancelPrevious: true,
                        requestId: 'getHostSnapshot'
                    }
                })
                if (snapshot) {
                    return Promise.resolve(snapshot)
                }
                return Promise.reject(new Error('Get host snapshot failed.'))
            },
            initCharts () {
                this.initCpuChart()
                this.initMemoryChart()
                this.initDiskChart()
            },
            initCpuChart () {
                const cpuChart = this.Echarts.init(this.$refs.cpuChart)
                cpuChart.setOption({
                    title: {
                        text: this.$t('总CPU使用率'),
                        textStyle: {
                            color: '#333948'
                        },
                        x: 'center',
                        y: 'bottom'
                    },
                    series: [{
                        name: '',
                        type: 'gauge',
                        radius: '80%',
                        axisLine: {
                            lineStyle: {
                                color: [[0.2, '#30d878'], [0.8, '#3C96FF'], [1, '#FF5656']]
                            }
                        },
                        detail: {
                            formatter: '{value}%',
                            fontSize: 26,
                            offsetCenter: [0, '60%']
                        },
                        data: [{
                            value: this.snapshot.cpuUsage / 100,
                            name: ''
                        }]
                    }]
                })
            },
            initMemoryChart () {
                const memoryChart = this.Echarts.init(this.$refs.memoryChart)
                memoryChart.setOption({
                    title: {
                        text: this.$t('总内存使用率'),
                        textStyle: {
                            color: '#333948'
                        },
                        x: 'center',
                        y: 'bottom'
                    },
                    series: [{
                        name: '',
                        type: 'gauge',
                        radius: '80%',
                        axisLine: {
                            lineStyle: {
                                color: [[0.2, '#30d878'], [0.8, '#3C96FF'], [1, '#FF5656']]
                            }
                        },
                        detail: {
                            formatter: '{value}%',
                            fontSize: 26,
                            offsetCenter: [0, '60%']
                        },
                        data: [{
                            value: this.snapshot.memUsage / 100,
                            name: ''
                        }]
                    }]
                })
            },
            initDiskChart () {
                const diskChart = this.Echarts.init(this.$refs.diskChart)
                diskChart.setOption({
                    title: {
                        text: this.$t('磁盘使用情况'),
                        textStyle: {
                            color: '#333948'
                        },
                        x: 'center',
                        y: 'bottom'
                    },
                    series: [{
                        name: '',
                        type: 'gauge',
                        radius: '80%',
                        axisLine: {
                            lineStyle: {
                                color: [[0.2, '#30d878'], [0.8, '#3C96FF'], [1, '#FF5656']]
                            }
                        },
                        detail: {
                            formatter: '{value}%',
                            fontSize: 26,
                            offsetCenter: [0, '60%']
                        },
                        data: [{
                            value: this.snapshot.diskUsage / 100,
                            name: ''
                        }]
                    }]
                })
            },
            openAgentApp () {
                const topWindow = window.top
                const isPaasConsole = topWindow !== window
                const [cloud = {}] = this.info.host.bk_cloud_id || []
                const urlSuffix = `#/plugin-manager/list?cloud_id=${cloud.bk_inst_id}&ip=${this.info.host.bk_host_innerip}`
                if (isPaasConsole) {
                    topWindow.postMessage(JSON.stringify({
                        action: 'open_other_app',
                        app_code: 'bk_nodeman',
                        app_url: urlSuffix
                    }), '*')
                } else {
                    const agentAppUrl = window.CMDB_CONFIG.site.agent
                    if (agentAppUrl) {
                        window.open(agentAppUrl + urlSuffix)
                    } else {
                        this.$warn(this.$t('未配置节点管理地址'))
                    }
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .status {
        height: 100%;
        &.is-offline {
            text-align: center;
            &:before {
                content: "";
                display: inline-block;
                vertical-align: middle;
                height: 100%;
            }
        }
    }
    .status-info {
        width: 720px;
        padding: 22px 0 0 0;
        text-align: left;
        .info-title {
            color: #333948;
            font-size: 14px;
            font-weight: bold;
            line-height: 1;
        }
        .update-time {
            font-weight: normal;
        }
    }
    .info-list {
        font-size: 14px;
        .info-item {
            width: 50%;
            margin: 20px 0 0 0;
            .item-label {
                color: #6b7baa;
            }
            .item-value {
                color: #4d597d;
            }
        }
    }
    .info-gauge {
        display: flex;
        .gauge-chart {
            flex: 1;
            height: 250px;
        }
    }
    .status-offline {
        display: inline-block;
        vertical-align: middle;
        .offline-image {
            height: 140px;
            background-image: url("../../../assets/images/box-circle.png"), url("../../../assets/images/box.png"), url("../../../assets/images/box-light.png");
            background-repeat: no-repeat;
            background-position: center -10px, center, center;
            background-size: 75px 75px, 127px 100%, 127px 100%;
        }
        .offline-text {
            margin: 40px 0 0 0;
            font-size: 14px;
            a {
                color: #3c96ff;
            }
        }
    }
</style>
