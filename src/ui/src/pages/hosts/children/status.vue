/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div>
        <div class="tab-content" v-if="isLoaded && hostSnapshot">
            <div class="attribute-list clearfix">
                <div class="title clearfix">
                    <h3 class="fl">{{$t('HostResourcePool[\'基本值\']')}}</h3>
                    <div class="content fr clearfix">
                        <div class="info-title">{{$t('HostResourcePool[\'最近更新时间\']')}}：</div>
                        <div class="info-detail">{{hostSnapshot.upTime}}</div>
                    </div>
                </div>
                <ul class="info-list">
                    <li class="attr-item clearfix">
                        <div class="item-title">{{$t('HostResourcePool[\'总流入量\']')}}：</div>
                        <div class="item-detail">{{(hostSnapshot.rcvRate / 100).toFixed(2)}}Mb/s</div>
                    </li>
                    <li class="attr-item clearfix">
                        <div class="item-title">{{$t('HostResourcePool[\'启动时间\']')}}：</div>
                        <div class="item-detail">{{bootTime}}</div>
                    </li>
                    <li class="attr-item clearfix">
                        <div class="item-title">{{$t('HostResourcePool[\'总流出量\']')}}：</div>
                        <div class="item-detail">{{(hostSnapshot.sendRate / 100).toFixed(2)}}Mb/s</div>
                    </li>
                    <li class="attr-item clearfix">
                        <div class="item-title">{{$t('HostResourcePool[\'磁盘总量\']')}}：</div>
                        <div class="item-detail">{{hostSnapshot.Disk}}GB</div>
                    </li>
                    <li class="attr-item clearfix">
                        <div class="item-title">{{$t('HostResourcePool[\'内存总量\']')}}：</div>
                        <div class="item-detail">{{(hostSnapshot.Mem / 1024).toFixed(2)}}GB</div>
                    </li>
                    <li class="attr-item clearfix" v-if="!isWindowsOSType">
                        <div class="item-title">loadavg：</div>
                        <div class="item-detail">{{hostSnapshot.loadavg}}</div>
                    </li>
                </ul>
            </div>
            <div class="chart-wrapper">
                <div class="chart-box">
                    <div ref="chart1" class="chart-item"></div>
                </div>
                <div class="chart-box">
                    <div ref="chart2" class="chart-item"></div>
                </div>
                <div class="chart-box">
                    <div ref="chart3" class="chart-item"></div>
                </div>
            </div>
        </div>
        <div class="tab-content" v-else>
            <div class="box-content">
                <div class="box">
                    <div class="box-light"></div>
                    <div class="box-circle">
                        <div class="circle">
                            <img src="../../../common/images/box-circle.png" alt="">
                        </div>
                    </div>
                </div>
            </div>
            <p class="box-text">{{$t('HostResourcePool[\'当前主机没有安装 Agent 或者 Agent 已经离线\']')}}</p>
        </div>
    </div>
</template>

<script>
    import echarts from 'echarts'
    import {mapGetters} from 'vuex'
    export default {
        props: {
            isShow: {
                default: false,
                type: Boolean
            },
            isSidesliderShow: {
                default: false,
                type: Boolean
            },
            isLoaded: {
                type: Boolean,
                required: true
            },
            isWindowsOSType: {
                type: Boolean,
                required: true
            }
        },
        computed: {
            ...mapGetters({
                hostSnapshot: 'getHostSnapshot'
            }),
            bootTime () {
                return this.$formatTime(this.hostSnapshot.bootTime * 1000)
            }
        },
        watch: {
            isSidesliderShow (val) {
                if (!val) {
                    this.$emit('update:isLoaded', false)
                }
            },
            isShow (val) {
                if (this.isLoaded && val) {
                    if (this.hostSnapshot !== '') {
                        this.initChart()
                    }
                }
            },
            isLoaded (isLoaded) {
                if (isLoaded && this.isShow) {
                    this.initChart()
                }
            }
        },
        methods: {
            /*
                初始化图表
            */
            initChart () {
                let chart1 = echarts.init(this.$refs['chart1'])
                let chart2 = echarts.init(this.$refs['chart2'])
                let chart3 = echarts.init(this.$refs['chart3'])
                chart1.setOption({
                    title: {
                        text: this.$t('Hosts["总CPU使用率"]'),
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
                            value: this.hostSnapshot.cpuUsage / 100,
                            name: ''
                        }]
                    }]
                })
                chart2.setOption({
                    title: {
                        text: this.$t('Hosts["总内存使用率"]'),
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
                            value: this.hostSnapshot.memUsage / 100,
                            name: ''
                        }]
                    }]
                })
                chart3.setOption({
                    title: {
                        text: this.$t('Hosts["磁盘使用情况"]'),
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
                            value: this.hostSnapshot.diskUsage / 100,
                            name: ''
                        }]
                    }]
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tab-content{
        padding-top: 62px!important;
        font-size: 14px;
        padding-left:40px;
        height: 100%;
        overflow: auto;
        &::-webkit-scrollbar{
            width: 6px;
            height: 5px;
        }
        &::-webkit-scrollbar-thumb{
            border-radius: 20px;
            background: #a5a5a5;
        }
        .attribute-list{
            margin-bottom:40px;
            .title{
                h3{
                    color:#333948;
                    font-size:14px;
                    font-weight: bold;
                    line-height:1;
                    margin: 0;
                    margin-bottom: 20px;
                }
                .content{
                    font-size: 12px;
                    padding-right: 30px;
                    .info-title{
                        float: left;
                        color: #333948;
                    }
                    .info-detail{
                        float: left;
                        color: #6b7baa;
                    }
                }
            }
            ul{
                li.attr-item{
                    font-size: 14px;
                    width: 50%;
                    float: left;
                    line-height:1;
                    color:#6b7baa;
                    margin-bottom: 20px;
                    padding-right: 20px;
                    .item-title{
                        float: left;
                    }
                    .item-detail{
                        float: left;
                        color: #4d597d;
                    }
                }
            }
        }
        .chart-wrapper{
            display: flex;
            padding-right: 20px;
            .chart-box{
                flex: 1;
                .chart-item{
                    // width: 100%;
                    height: 250px;
                }
            }
        }
        .box-content{
            width: 127px;
            height: 153px;
            margin: 40px auto;
            text-align: center;
            .box{
                width: 127px;
                height: 153px;
                background: url(../../../common/images/box.png) no-repeat;
                background-size: 100%;
                position: relative;
                .box-light{
                    position: absolute;
                    width: 127px;
                    left: 0;
                    background: url(../../../common/images/box-light.png) no-repeat;
                    -webkit-transform-origin: center bottom;
                    -moz-transform-origin: center bottom;
                    -o-transform-origin: center bottom;
                    -ms-transform-origin: center bottom;
                    transform-origin: center bottom;
                    -webkit-animation: light 2s ease-in-out .5s both;
                    -moz-animation: light 2s ease-in-out .5s both;
                    -ms-animation: light 2s ease-in-out .5s both;
                    -o-animation: light 2s ease-in-out .5s both;
                    animation: light 2s ease-in-out .5s both;
                    background-size: 100%;
                }
                .box-circle{
                    perspective: 1200px;
                    position: relative;
                    .circle{
                        position: absolute;
                        top: -10px;
                        left: 26px;
                        width: 75px;
                        height: 75px;
                        transform-style: preserve-3d;
                        img{
                            position: absolute;
                            left: 0;
                            top: 0;
                            width: 100%;
                            height: 100%;
                        }
                    }
                }
            }
        }
        .box-text{
            width: 100%;
            text-align: center;
        }
    }
    @-webkit-keyframes light {
        from {
            height: 0;
            bottom: 70px;
        }
        to {
            height: 153px;
            top: 0;
            bottom: auto;
        }
    }
    @-moz-keyframes light {
        from {
            height: 0;
            bottom: 70px;
        }
        to {
            height: 153px;
            top: 0;
            bottom: auto;
        }
    }
    @-ms-keyframes light {
        from {
            height: 0;
            bottom: 70px;
        }
        to {
            height: 153px;
            top: 0;
            bottom: auto;
        }
    }
    @-o-keyframes light {
        from {
            height: 0;
            bottom: 70px;
        }
        to {
            height: 153px;
            top: 0;
            bottom: auto;
        }
    }
    @keyframes light {
        from {
            height: 0;
            bottom: 70px;
        }
        to {
            height: 153px;
            top: 0;
            bottom: auto;
        }
    }
</style>