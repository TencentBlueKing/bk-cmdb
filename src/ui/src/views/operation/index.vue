<template>
    <div>
        <div class="operation-top">
            <span class="operation-edit" v-if="!editModel" @click="editModel = true">
                {{$t('Operation["编辑模式"]')}}</span>
            <span class="operation-edit" v-if="editModel" @click="editCancel()">
                {{$t('Operation["退出编辑模式"]')}}</span>
            <span class="operation-title">{{$t('Operation["主机统计"]')}}</span>
        </div>
        <div class="operate-menus">
            <div class="menu-items" @click="goRouter('business')">
                <div class="item-left">
                    <i class="icon icon-cc-business"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Common["业务"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：{{ navData.biz }}</span>
                </div>
            </div>
            <div class="menu-items" @click="goRouter('topology')">
                <div class="item-left">
                    <i class="icon icon-cc-host"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Hosts["模块"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：{{ navData.module }}</span>
                </div>
            </div>
            <div class="menu-items" @click="goRouter('resource')">
                <div class="item-left">
                    <i class="icon icon-cc-host"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Nav["主机"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：{{ navData.host }}</span>
                </div>
            </div>
        </div>
        <div v-for="(item, key) in hostData.disList"
            :key="item.report_type + item.config_id"
            :style="{ width: item.width + '%' }"
            class="operation-layout">
            <div class="operation-charts" :id="item.report_type + item.config_id"></div>
            <div v-if="item.noData" class="null-data">
                <span>该图表暂无数据</span>
            </div>
            <div class="charts-options" v-if="editModel">
                <i class="bk-icon icon-arrows-down-shape"
                    @click="moveChart('host', 'down', key, hostData.disList)" disabled></i>
                <i class="bk-icon icon-arrows-up-shape"
                    @click="moveChart('host', 'up', key, hostData.disList)"></i>
                <i class="bk-icon icon-edit"
                    @click="openNew('edit', 'host', item, key)"></i>
                <i class="icon icon-cc-tips-close"
                    @click="deleteChart('host', key, hostData.disList, item.config_id)"></i>
            </div>
        </div>
        <div class="operation-add" v-if="editModel" @click="openNew('add', 'host')">
            <i class="bk-icon icon-plus"></i>
        </div>
        <div class="operation-top">
            <span class="operation-title">{{$t('Operation["实例统计"]')}}</span>
        </div>
        <div class="operate-menus operate-menus-bottom">
            <div class="menu-items" @click="goRouter('model')">
                <div class="item-left">
                    <i class="menu-icon icon-cc-nav-model"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Nav["模型"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：{{ navData.model }}</span>
                </div>
            </div>
            <div class="menu-items" @click="goRouter('index')">
                <div class="item-left">
                    <i class="bk-cc-icon icon-cc-business"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Operation["实例"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：{{ navData.inst }}</span>
                </div>
            </div>
        </div>
        <div v-for="(item, key) in instData.disList"
            :key="item.report_type + item.config_id"
            :style="{ width: item.width + '%' }"
            class="operation-layout">
            <div class="operation-charts" :id="item.report_type + item.config_id"></div>
            <div v-if="item.noData" class="null-data">
                <span>该图表暂无数据</span>
            </div>
            <div class="charts-options" v-if="editModel">
                <i class="bk-icon icon-arrows-down-shape"
                    @click="moveChart('inst', 'down', key, instData.disList)" disabled></i>
                <i class="bk-icon icon-arrows-up-shape"
                    @click="moveChart('inst', 'up', key, instData.disList)"></i>
                <i class="bk-icon icon-edit" @click="openNew('edit', 'inst', item, key)"></i>
                <i class="icon icon-cc-tips-close"
                    @click="deleteChart('inst', key, instData.disList, item.config_id)"></i>
            </div>
        </div>
        <div class="operation-add" v-if="editModel" @click="openNew('add', 'inst')">
            <i class="bk-icon icon-plus"></i>
        </div>
        <v-detail v-if="isShow"
            :open-type="editType.openType"
            :chart-data="newChart"
            @transData="saveData"
            @closeChart="cancelData">
        </v-detail>
    </div>
</template>

<script>
    import Plotly from 'plotly.js'
    import { mapActions } from 'vuex'
    import vDetail from './chart-detail'

    export default {
        name: 'index',
        components: {
            vDetail
        },
        data () {
            return {
                editModel: false,
                isShow: false,
                newChart: {},
                editType: {
                    openType: 'add',
                    hostType: 'host',
                    index: ''
                },
                hostData: {
                    disList: []
                },
                instData: {
                    disList: []
                },
                navData: {
                    biz: '',
                    host: '',
                    module: '',
                    inst: '',
                    model: ''
                }
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Operation["统计报表"]'))
            this.getChartList()
        },
        methods: {
            ...mapActions('operationChart', [
                'getCountedCharts',
                'getCountedChartsData',
                'deleteOperationChart',
                'updateChartPosition'
            ]),
            async getChartList () {
                const res = await this.getCountedCharts({})
                this.hostData.disList = res.info.host
                this.instData.disList = res.info.inst
                res.info.nav.forEach(item => {
                    this.getNavData(item, 'nav')
                })
                this.hostData.disList.forEach((item, key) => {
                    this.getNavData(item, 'host', key, this.hostData.disList.length)
                })
                this.instData.disList.forEach((item, key) => {
                    this.getNavData(item, 'inst', key, this.instData.disList.length)
                })
            },
            async getNavData (item, type, key, length) {
                const res = await this.getCountedChartsData({
                    params: {
                        config_id: item.config_id
                    },
                    config: {
                        globalError: false
                    }
                })
                if (type === 'nav') {
                    res.forEach(items => {
                        this.navData[items.id] = items.count
                    })
                } else {
                    if (JSON.stringify(res) === '{}' || JSON.stringify(res) === '[]') {
                        item.hasData = false
                        item.data = {
                            data: [{
                                labels: [],
                                type: item.chart_type,
                                values: [],
                                x: [],
                                y: []
                            }]
                        }
                    } else {
                        item.hasData = true
                        item.data = await this.dataDeal(item, res)
                    }
                    if (type === 'host') {
                        this.hostData.disList[key] = item
                        if (length === (key + 1)) {
                            setTimeout(() => {
                                this.drawCharts(this.hostData.disList)
                            }, 100)
                        }
                    } else {
                        this.instData.disList[key] = item
                        if (length === (key + 1)) {
                            setTimeout(() => {
                                this.drawCharts(this.instData.disList)
                            }, 100)
                        }
                    }
                }
            },
            dataDeal (data, res) {
                const returnData = {
                    data: []
                }
                if (res && Array.isArray(res)) {
                    const content = {
                        labels: [],
                        values: [],
                        x: [],
                        y: [],
                        type: data.chart_type
                    }
                    res.forEach(item => {
                        if (data.chart_type === 'pie') {
                            content.labels.push(item.id)
                            content.values.push(item.count)
                        } else {
                            content.x.push(item.id)
                            content.y.push(item.count)
                        }
                    })
                    returnData.data.push(content)
                } else if (res && !Array.isArray(res)) {
                    const barModel = {
                        chart: {
                            delete: {
                                name: 'delete',
                                type: 'bar',
                                y: []
                            },
                            update: {
                                name: 'update',
                                type: 'bar',
                                y: []
                            },
                            create: {
                                name: 'create',
                                type: 'bar',
                                y: []
                            }
                        },
                        x: [],
                        isX: false
                    }
                    for (const item in res) {
                        const content = {
                            name: item,
                            x: [],
                            y: []
                        }
                        if (Array.isArray(res[item])) {
                            res[item].forEach(child => {
                                const time = this.$tools.formatTime(child.id, 'YYYY-MM-DD')
                                content.x.push(time + '_')
                                content.y.push(child.count)
                                content.mode = 'line'
                            })
                            returnData.data.push(content)
                        } else {
                            barModel.x.push(item)
                            barModel.isX = true
                            for (const chartItem in barModel.chart) {
                                barModel.chart[chartItem].y.push(res[item][chartItem])
                            }
                        }
                    }
                    if (barModel.isX) {
                        for (const chartItem in barModel.chart) {
                            barModel.chart[chartItem].x = barModel.x
                            returnData.data.push(barModel.chart[chartItem])
                        }
                    }
                }
                return returnData
            },
            drawCharts (charList) {
                charList.forEach(item => {
                    const myDiv = document.getElementById(item.report_type + item.config_id)
                    const data = item.data.data
                    if (!item.hasData) this.$set(item, 'noData', true)
                    const layout = {
                        height: 400,
                        title: item.name,
                        barmode: 'stack'
                    }
                    const options = {
                        displaylogo: false,
                        displayModeBar: false
                    }
                    Plotly.newPlot(myDiv, data, layout, options)
                })
            },
            goDraws () {
                this.drawCharts(this.hostData.disList)
                this.drawCharts(this.instData.disList)
            },
            editCancel () {
                this.editModel = false
                setTimeout(() => {
                    this.goDraws()
                }, 100)
            },
            moveChart (type, dire, key, list) {
                if (dire === 'up' && key !== 0) {
                    list[key] = list.splice(key - 1, 1, list[key])[0]
                } else if (dire === 'down' && key !== list.length - 1) {
                    list[key + 1] = list.splice(key, 1, list[key + 1])[0]
                }
                if (type === 'host') this.hostData.disList = list
                else this.instData.disList = list
                this.updatePosition()
            },
            deleteChart (type, key, list, id) {
                this.$bkInfo({
                    title: this.$tc('Operation["是否确认删除"]'),
                    confirmFn: () => {
                        list.splice(key, 1)
                        if (type === 'host') this.hostData.disList = list
                        else this.instData.disList = list
                        this.deleteOperationChart({
                            id: id
                        })
                        this.updatePosition()
                    }
                })
            },
            async openNew (type, host, data, key) {
                this.editType.hostType = host
                this.editType.key = key
                if (type === 'edit') this.newChart = data
                else {
                    this.newChart = {
                        report_type: 'custom',
                        name: '',
                        config_id: null,
                        bk_obj_id: host,
                        chart_type: 'pie',
                        field: '',
                        width: '50'
                    }
                }
                this.editType.openType = type
                this.isShow = true
            },
            async saveData (data) {
                let editList = []
                if (this.editType.hostType === 'host') editList = this.hostData.disList
                else editList = this.instData.disList
                if (this.editType.openType === 'add') {
                    editList.push(data)
                } else {
                    editList[this.editType.key] = data
                }
                this.isShow = false
                await this.getNavData(data, this.editType.hostType,
                                      this.editType.openType === 'add' ? (editList.length - 1) : this.editType.key,
                                      editList.length)
                this.updatePosition()
                this.newChart = {}
            },
            cancelData () {
                this.isShow = false
                this.newChart = {}
            },
            updatePosition () {
                const data = {
                    'host': [],
                    'inst': []
                }
                this.hostData.disList.forEach(item => {
                    data.host.push(item.config_id)
                })
                this.instData.disList.forEach(item => {
                    data.inst.push(item.config_id)
                })
                this.updateChartPosition({
                    params: {
                        position: data
                    }
                })
                setTimeout(() => {
                    this.goDraws()
                }, 100)
            },
            goRouter (route) {
                this.$router.push(route)
            }
        }
    }
</script>

<style scoped lang="scss">
    *{
        margin:0;
        padding:0
    }
    .operation-top{
        padding: 5px;
        margin: 5px;
        border-bottom: 2px solid #eee;

        .operation-title{
            display: block;
            font-size: 14px;
        }
        .operation-edit{
            display: block;
            text-align: right;
            font-size: 12px;
            color: $cmdbMainBtnColor;
            text-decoration: underline;
            cursor: pointer;
            span{
                margin-left: 10px;
            }
        }
    }
    .operate-menus{
        height: 110px;
        padding: 10px 5px;
        .menu-items{
            height: 100%;
            width: 25%;
            border: 1px solid #eee;
            float: left;
            cursor: pointer;
            padding: 15px 20px;
            box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
            &:first-child{
                 border-right: none;
                .item-left{
                    background: $cmdbMainBtnColor;
                }
            }
            &:last-child{
                 border-left: none;
            }
             .item-left{
                 width: 60px;
                 height: 60px;
                 font-size: 32px;
                 text-align: center;
                 line-height: 60px;
                 float: left;
                 border-radius: 50%;
                 margin-right: 30px;
                 background-color: #f5e41f;
                 i{
                     color: white;
                     vertical-align: inherit;
                 }
             }
             .item-right{
                 height: 100%;
                 text-align: center;
                 width: calc(100% - 120px);
                 float: left;
                 span{
                     display: block;
                     text-align: left;
                     &:first-child{
                        font-weight: bold;
                        font-size: 20px;
                        margin-top: 5px;
                     }
                     &:last-child{
                        margin-top: 5px;
                        font-size: 12px;
                     }
                 }
             }
        }
    }
    .operate-menus-bottom{
        .menu-items{
            &:last-child{
                 border-left: 1px solid #eee;
                .item-left{
                    background: #f33c13;
                }
            }
        }
    }
    .operation-layout{
        padding: 10px 5px;
        display: inline-block;
        position: relative;
    }
    .operation-charts{
        width: 100%;
        border: 1px solid #eee;
        display: table;
        box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
    }
    .charts-options{
        position: absolute;
        right: 15px;
        top: 10px;
        color: $cmdbMainBtnColor;
        i{
            cursor: pointer;
            font-weight: bold;
            margin-left: 3px;
        }
    }
    .null-data {
        position: absolute;
        top: 30%;
        width: 100%;
        span {
            margin: 0 auto;
            display: block;
            height: 120px;
            width: 200px;
            line-height: 120px;
            text-align: center;
            color: #fff;
            background: $cmdbWarningColor;
        }
    }
    .operation-add{
        width: 200px;
        height: 120px;
        border: 1px solid #eee;
        font-size: 100px;
        margin-left: 5px;
        color: $cmdbMainBtnColor;
        box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
        text-align: center;
        cursor: pointer;
        i{
            font-weight: bold;
            vertical-align: baseline;
        }
    }
</style>
