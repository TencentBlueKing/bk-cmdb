<template>
    <div class="bk-date" @click="openDater" v-click-outside="close" :style="bkDateWidthObj">
        <input type="text" name="date-select" readonly="true" :disabled="disabled" :placeholder="placeholder?placeholder:t('datepicker.selectDate')" v-model="selectedDateView">
        <transition :name="transitionName">
            <div :style="panelStyle" :class="['date-dropdown-panel', 'daterange-dropdown-panel', {'has-sidebar': quickSelect || isShowConfirm}]" v-if="showDatePanel">
                <!-- 开始日期选择面板 -->
                <date-picker class="start-date date-select-container fl"
                    ref="startDater"
                    @date-quick-switch="dateQuickSwitch(...arguments, 'startDater')"
                    @date-selected="triggerSelect"
                    :selected-range="selectedDateRange"
                    :selected-range-tmp="selectedDateRangeTmp"
                    :initDate="initStartDate"
                    :startDate="minDate"
                    :endDate="maxDate"
                    :bkDate="bkDateStart"
                    :timer="timer"></date-picker>

                <!-- 结束日期选择面板 -->
                <date-picker class="end-date date-select-container fl"
                    ref="endDater"
                    @date-quick-switch="dateQuickSwitch(...arguments, 'endDater')"
                    @date-selected="triggerSelect"
                    :selected-range="selectedDateRange"
                    :selected-range-tmp="selectedDateRangeTmp"
                    :initDate="initEndDate"
                    :startDate="minDate"
                    :endDate="maxDate"
                    :bkDate="bkDateEnd"
                    :timer="timer"></date-picker>

                <!-- 日期快速选择配置 -->
                <div class="range-config fl" v-if="quickSelect">
                    <a href="javascript:;"
                        :class="{'active': shouldBeMatched(range)}"
                        v-for="range in defaultRanges"
                        @click.stop="changeRanges(range)">{{range.text}}</a>
                </div>
                <div class="range-action fl" v-if="isShowConfirm">
                    <a href="javascript:;" @click.stop="confirm">确定</a>
                    <a href="javascript:;" @click.stop="clear">清空</a>
                </div>
            </div>
        </transition>
    </div>
</template>
<script>
    import clickoutside from './../../../utils/clickoutside'
    import datepicker from './date-picker'
    import moment from 'moment'
    import {getActualTop, getActualLeft} from '../../../utils/utils'

    const oneOf = (value, validList) => {
        for (let i = 0; i < validList.length; i++) {
            if (value === validList[i]) {
                return true
            }
        }
        return false
    }

    class BkDate {
        constructor (flag, time) {
            const dater = time ? new Date(time) : new Date()
            
            // 当前日期
            this.currentDay = {
                year: dater.getFullYear(),
                month: dater.getMonth() + 1,
                day: dater.getDate()
            }
            // 当前时间
            this.currentTime = {
                hour: dater.getHours(),
                minute: dater.getMinutes(),
                // minute: dater.getMinutes() + 1,
                second: dater.getSeconds()
            }
            // 日期选择器默认显示
            this.year = this.currentDay.year
            this.month = this.currentDay.month
            this.day = this.currentDay.day

            this.setTimer = false

            this.index = flag === 'start' ? 0 : 1
        }

        // 日期选择器选择重置
        setDate (date) {
            let dateItems = date.split('-')
            if (dateItems[0]) {
                this.year = parseInt(dateItems[0])
            }
            if (dateItems[1]) {
                this.month = parseInt(dateItems[1])
            }
            if (dateItems[2]) {
                this.day = parseInt(dateItems[2])
            }
        }

        // 格式化日期字符串
        formatDateString (value) {
            return parseInt(value) < 10 ? `0${value}` : value
        }

        // 获取当天格式化日期
        getFormatToday () {
            return `${this.currentDay.year}-${this.formatDateString(this.currentDay.month)}-${this.formatDateString(this.currentDay.day)}`
        }

        // 获取当前格式化日期
        getFormatDate () {
            return `${this.year}-${this.formatDateString(this.month)}-${this.formatDateString(this.day)}`
        }

        // 获取当前月天数
        getCurrentMouthDays () {
            return new Date(this.year, this.month, 0).getDate()
        }

        // 获取上一个月天数
        getLastMouthDays () {
            return new Date(this.year, this.month - 1, 0).getDate()
        }

        // 获取当前月份1号是星期几
        getCurrentMonthBeginWeek () {
            return new Date(this.year, this.month - 1, 1).getDay()
        }
    }

    /**
    * bk-date
    * 参数配置：
    * @param placeholder: '请选择日期' -默认占位内容
    * @param disabled: 'false' -是否禁用
    * @param align: 'left/center/right' - 日期面板对齐方式
    * @param rangeSeparator: '-' - 日期范围分隔符
    * @param startDate: 'YYYY-MM-DD' -默认开始日期
    * @param endDate: 'YYYY-MM-DD' -默认结束日期
    * @param quickSelect: true/false -快捷选择开关
    * @param ranges: true/false -选择日期后是否自动关闭
    */

    export default {
        name: 'bk-daterangepicker',
        components: {
            'date-picker': datepicker
        },
        props: {
            placeholder: {
                type: String,
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            align: {
                type: String,
                default: 'left'
            },
            quickSelect: {
                type: Boolean,
                default: true
            },
            rangeSeparator: {
                type: String,
                default: '-'
            },
            initDate: {
                type: String,
                default: ''
            },
            startDate: {
                type: String,
                default: ''
            },
            endDate: {
                type: String,
                default: ''
            },
            minDate: {
                type: String,
                default: ''
            },
            maxDate: {
                type: String,
                default: ''
            },
            isShowConfirm: {
                type: Boolean,
                default: false
            },
            ranges: {
                type: Object,
                default: () => {
                    return {
                        '昨天': [moment().subtract(1, 'days'), moment()],
                        '最近一周': [moment().subtract(7, 'days'), moment()],
                        '最近一个月': [moment().subtract(1, 'month'), moment()],
                        '最近三个月': [moment().subtract(3, 'month'), moment()]
                    }
                }
            },
            timer: {
                type: Boolean,
                default: false
            },
            position: {
                validator (value) {
                    return oneOf(
                        value,
                        [
                            // top <=> top-right, bottom <=> bottom-right
                            'top', 'bottom', 'left', 'right',
                            'top-left', 'top-right', 'bottom-left', 'bottom-right'
                        ]
                    )
                },
                default: 'bottom-right'
            }
        },
        data () {
            window.moment = moment
            const sdr = [moment().format('YYYY-MM-DD HH:mm:ss')]
            const sdrt = [moment(sdr[0]).format('YYYY-MM-DD')]

            const initEndDate = this.endDate || moment().format('YYYY-MM-DD')
            const initStartDate = moment(initEndDate).subtract(1, 'month').format('YYYY-MM-DD')

            const bkDateStart = this.startDate ? new BkDate('start', this.startDate) : new BkDate('start')
            bkDateStart.setDate(initStartDate)

            const bkDateEnd = this.endDate ? new BkDate('end', this.endDate) : new BkDate('end')
            bkDateEnd.setDate(initEndDate)

            let transitionName = 'toggle-slide'
            const panelStyle = {}
            const positionArr = this.position.split('-')
            if (positionArr.indexOf('top') > -1) {
                panelStyle.bottom = '38px'
                transitionName = 'toggle-slide2'
            } else {
                panelStyle.top = '38px'
            }

            if (positionArr.indexOf('left') > -1) {
                panelStyle.right = 0
            } else {
                panelStyle.left = 0
            }

            return {
                panelStyle: panelStyle,
                transitionName: transitionName,
                bkDateWidthObj: this.timer ? {width: '350px'} : {},
                showDatePanel: false, // 日期选择面板开关
                selectedDateView: '', // 已选日期范围显示
                selectedRange: '', // 已选择的快速选择
                initStartDate: initStartDate,
                initEndDate: initEndDate,
                bkDateStart: bkDateStart,
                bkDateEnd: bkDateEnd,
                selectedDateRange: sdr, // 已选日期数据
                selectedDateRangeTmp: sdrt, // selectedDateRange 的副本，用于 date-picker 中 shouldBeSelected 高亮判断仅需要判断日期而不需要判断时间
                defaultRanges: [ // 默认快捷选择菜单栏
                    {
                        text: '昨天',
                        value: [moment().subtract(1, 'days'), moment()]
                    },
                    {
                        text: '最近一周',
                        value: [moment().subtract(7, 'days'), moment()]
                    },
                    {
                        text: '最近一个月',
                        value: [moment().subtract(1, 'month'), moment()]
                    },
                    {
                        text: '最近三个月',
                        value: [moment().subtract(3, 'month'), moment()]
                    }
                ]
            }
        },
        directives: {
            clickoutside
        },
        created () {
            this.init()
        },
        updated () {
            // if (!this.selectedDateView) {
            //     return
            // }
            // const [initStartDate, initEndDate] = this.selectedDateView.split(` ${this.rangeSeparator} `)
            // this.bkDateStart.setDate(initStartDate)
            // this.bkDateEnd.setDate(initEndDate)
        },
        watch: {
            startDate (val) {
                let defaultRanges = []
                for (let item in this.ranges) {
                    defaultRanges.push({
                        text: item,
                        value: this.ranges[item]
                    })
                }
                this.defaultRanges = defaultRanges

                let hour = ''
                let minute = ''
                let second = ''
                if (this.timer) {
                    hour = moment(val).hour()
                    minute = moment(val).minute()
                    second = moment(val).second()
                }

                if (this.startDate) {
                    this.selectedDateRange.shift()
                    this.selectedDateRange.unshift(
                        this.timer
                            ? moment(this.startDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                            : this.startDate
                    )
                    this.selectedDateRangeTmp.shift()
                    this.selectedDateRangeTmp.unshift(
                        this.timer
                            ? moment(this.startDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD')
                            : this.startDate
                    )
                }
            },
            endDate (val) {
                let defaultRanges = []
                for (let item in this.ranges) {
                    defaultRanges.push({
                        text: item,
                        value: this.ranges[item]
                    })
                }
                this.defaultRanges = defaultRanges

                let hour = ''
                let minute = ''
                let second = ''
                if (this.timer) {
                    hour = moment(val).hour()
                    minute = moment(val).minute()
                    second = moment(val).second()
                }
                
                if (this.endDate) {
                    this.selectedDateRange.pop()
                    this.selectedDateRange.push(
                        this.timer
                            ? moment(this.endDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                            : this.endDate
                    )
                    this.selectedDateRangeTmp.pop()
                    this.selectedDateRangeTmp.push(
                        this.timer
                            ? moment(this.endDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD')
                            : this.endDate
                    )
                }
            },
            selectedDateRange () {
                const seed = this.timer ? this.selectedDateRange : this.selectedDateRangeTmp
                if (seed.length === 2) {
                    let newSelectedDate = seed.join(` ${this.rangeSeparator} `)
                    // 判断此次选择的日期范围是否发生改变
                    if ((this.selectedDateView !== newSelectedDate) && !this.isShowConfirm) {
                        this.$emit('update:callbackSet', newSelectedDate)
                        this.$emit('change', this.selectedDateView, newSelectedDate)
                    }
                    this.selectedDateView = newSelectedDate
                    this.bkDateStart.setTimer = true
                    this.bkDateEnd.setTimer = true
                    // this.close()
                } else {
                    this.bkDateStart.setTimer = false
                    this.bkDateEnd.setTimer = false
                }
            },
            selectedDateView () {
                let formatDateStart = moment(this.selectedDateRange[0]).format('YYYY-MM')
                let formatDateEnd = moment(this.selectedDateRange[1]).format('YYYY-MM')
                this.initStartDate = this.selectedDateView.split(` ${this.rangeSeparator} `)[0]
                // 已选日期范围是同一个月时 重置两个月份选择面板 --  显示相邻月份
                if (moment(formatDateStart).isSame(formatDateEnd)) {
                    this.initEndDate = moment(this.selectedDateRange[0]).add(1, 'month').format(
                        this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'
                    )
                } else {
                    this.initEndDate = this.selectedDateView.split(` ${this.rangeSeparator} `)[1]
                }
            }
        },
        methods: {
            // 日期快速切换回调
            dateQuickSwitch (date, dater) {
                // 左边日期面板信息
                let startDateInfo = this.$refs.startDater
                // 右边日期面板信息
                let endDateInfo = this.$refs.endDater
                let startTopDate = startDateInfo.topBarFormatView.value
                let endTopDate = endDateInfo.topBarFormatView.value
                switch (dater) {
                    case 'startDater':
                        endDateInfo.BkDate.setDate(moment(startTopDate).add(1, 'month').format('YYYY-MM'))
                        break
                    case 'endDater':
                        startDateInfo.BkDate.setDate(moment(endTopDate).add(-1, 'month').format('YYYY-MM'))
                        break
                    default:
                        break
                }
            },

            changeRanges (range) {
                let rangeStartDate = range.value[0].format('YYYY-MM-DD HH:mm:ss')
                let rangeEndDate = range.value[1].format('YYYY-MM-DD HH:mm:ss')
                let rangeStartDateTmp = range.value[0].format('YYYY-MM-DD')
                let rangeEndDateTmp = range.value[1].format('YYYY-MM-DD')

                // 重复点击时不进行切换，比较时不需要比较时间，只需要比较日期
                if (rangeStartDateTmp === this.selectedDateRangeTmp[0] && rangeEndDateTmp === this.selectedDateRangeTmp[1]) {
                    return
                }
                this.selectedDateRange = [rangeStartDate, rangeEndDate]
                this.selectedDateRangeTmp = [rangeStartDateTmp, rangeEndDateTmp]
                this.initStartDate = this.timer ? rangeStartDate : rangeStartDateTmp
                this.initEndDate = this.timer ? rangeEndDate : rangeEndDateTmp
                this.selectedRange = range.text
                // this.showDatePanel = false
                // 已选日期范围是同一个月时 重置两个月份选择面板 --  显示相邻月份
                let formatDateStart = moment(rangeStartDate).format('YYYY-MM')
                let formatDateEnd = moment(rangeEndDate).format('YYYY-MM')
                if (moment(formatDateStart).isSame(formatDateEnd)) {
                    this.initEndDate = moment(this.selectedDateRange[0]).add(1, 'month').format(
                        this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'
                    )
                }

                this.bkDateStart.setDate(this.initStartDate)
                this.bkDateEnd.setDate(this.initEndDate)
            },

            // 选择回调传值
            triggerSelect (date, bkDateIndex) {
                // bkDateIndex 不为 undefined 说明变化的是 timer 而不是 date
                if (bkDateIndex !== undefined) {
                    const hour = moment(date).hour()
                    const minute = moment(date).minute()
                    const second = moment(date).second()
                    this.selectedDateRange.splice(
                        bkDateIndex,
                        1,
                        moment(this.selectedDateRange[bkDateIndex])
                            .hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                    )
                } else {
                    const selectedLen = this.selectedDateRange.length
                    date = moment(date).format('YYYY-MM-DD HH:mm:ss')
                    const dateTmp = moment(date).format('YYYY-MM-DD')
                    switch (selectedLen) {
                        case 0:
                            this.selectedDateRange.push(date)
                            this.selectedDateRangeTmp.push(dateTmp)
                            break
                        case 1:
                            // 首先验证第二次选择日期是否合格
                            if (moment(date).isSame(this.selectedDateRange[0]) || moment(date).isAfter(this.selectedDateRange[0])) {
                                this.selectedDateRange.push(date)
                                this.selectedDateRangeTmp.push(dateTmp)
                            }
                            if (moment(date).isBefore(this.selectedDateRange[0])) {
                                this.selectedDateRange = [date]
                                this.selectedDateRangeTmp = [dateTmp]
                            }
                            break
                        case 2:
                            this.selectedDateRange = [date]
                            this.selectedDateRangeTmp = [dateTmp]
                            break
                        default:
                    }
                }
            },

            // 快速选择栏标记
            shouldBeMatched (range) {
                let isMatched = this.selectedDateRangeTmp[0] === range.value[0].format('YYYY-MM-DD') && this.selectedDateRangeTmp[1] === range.value[1].format('YYYY-MM-DD')
                return isMatched
            },

            // 控制选择器显示隐藏
            openDater () {
                // if (this.disabled) {
                //     return
                // }
                // this.showDatePanel = true
                if (this.disabled) {
                    return
                }

                // 判断当前展示空间是否够默认展示
                let distanceLeft = getActualLeft(event.currentTarget)
                let distanceTop = getActualTop(event.currentTarget)
                let winWidth = document.body.clientWidth
                let winHeight = document.body.clientHeight
                let xSet = {}
                let ySet = {}

                if (distanceTop + 18 < winHeight / 2) {
                    ySet = {
                        top: '36px',
                        bottom: 'auto'
                    }
                } else {
                    ySet = {
                        top: 'auto',
                        bottom: '36px'
                    }
                    this.transitionName = 'toggle-slide2'
                }

                if (winWidth - distanceLeft < 660) {
                    xSet = {
                        left: 'auto',
                        right: 0
                    }
                } else {
                    xSet = {
                        left: 0,
                        right: 'auto'
                    }
                }

                this.panelStyle = {...xSet, ...ySet}

                this.showDatePanel = true
            },

            close () {
                // todo: 已选情况取消了选中状态然后点击空白处关闭日期选择器时，还原默认值选中状态
                if (this.selectedDateView && this.selectedDateRange.length === 0) {
                    this.selectedDateRange = this.selectedDateView.split(` ${this.rangeSeparator} `)
                }
                this.showDatePanel = false
            },
            confirm () {
                const seed = this.timer ? this.selectedDateRange : this.selectedDateRangeTmp
                if (seed.length === 2) {
                    let newSelectedDate = seed.join(` ${this.rangeSeparator} `)
                    // 判断此次选择的日期范围是否发生改变
                    this.$emit('update:callbackSet', newSelectedDate)
                    this.$emit('change', this.selectedDateView, newSelectedDate)
                }
                this.showDatePanel = false
            },
            clear () {
                this.$emit('change', this.selectedDateView, '')

                this.selectedDateView = ''
                this.selectedDateRange = []
                this.selectedDateRangeTmp = []

                const date = moment().format('YYYY-MM-DD HH:mm:ss')
                const dateTmp = moment(date).format('YYYY-MM-DD')
                this.selectedDateRange.push(date)
                this.selectedDateRangeTmp.push(dateTmp)

                this.initStartDate = this.startDate || moment().subtract(1, 'month').format('YYYY-MM-DD')
                this.initEndDate = this.endDate || moment().format('YYYY-MM-DD')
                this.bkDateStart.setDate(this.initStartDate)
                this.bkDateEnd.setDate(this.initEndDate)

                this.bkDateStart.currentTime = this.bkDateEnd.currentTime = {
                    hour: moment().hour(),
                    minute: moment().minute(),
                    second: moment().second()
                }

                this.showDatePanel = false
            },
            init () {
                let defaultRanges = []
                for (let item in this.ranges) {
                    defaultRanges.push({
                        text: item,
                        value: this.ranges[item]
                    })
                }
                this.defaultRanges = defaultRanges

                let hour = ''
                let minute = ''
                let second = ''
                if (this.timer) {
                    hour = moment().hour()
                    minute = moment().minute()
                    second = moment().second()
                }

                if (this.startDate) {
                    hour = moment(this.startDate).hour()
                    minute = moment(this.startDate).minute()
                    second = moment(this.startDate).second()
                    this.selectedDateRange.unshift(
                        this.timer
                            ? moment(this.startDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                            : this.startDate
                    )
                    this.selectedDateRangeTmp.unshift(
                        this.timer
                            ? moment(this.startDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD')
                            : this.startDate
                    )
                }
                if (this.endDate) {
                    hour = moment(this.endDate).hour()
                    minute = moment(this.endDate).minute()
                    second = moment(this.endDate).second()
                    this.selectedDateRange.pop()
                    this.selectedDateRange.push(
                        this.timer
                            ? moment(this.endDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                            : this.endDate
                    )
                    this.selectedDateRangeTmp.pop()
                    this.selectedDateRangeTmp.push(
                        this.timer
                            ? moment(this.endDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD')
                            : this.endDate
                    )
                }
            }
        }
    }
</script>
