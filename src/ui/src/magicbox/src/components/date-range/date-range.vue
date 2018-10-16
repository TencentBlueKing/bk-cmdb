<template>
    <div class="bk-date-range" :class="disabled ? 'disabled' : ''" @click="openDater" v-clickoutside="close" :style="bkDateWidthObj">
        <input type="text" name="date-select" readonly="true" :disabled="disabled" :placeholder="defaultPlaceholder" v-model="selectedDateView">
        <transition :name="transitionName">
            <div :style="panelStyle" :class="['date-dropdown-panel', 'daterange-dropdown-panel', {'has-sidebar': quickSelect}]" v-if="showDatePanel">
                <!-- 开始日期选择面板 -->
                <date-picker class="start-date date-select-container fl"
                    ref="startDater"
                    @date-quick-switch="dateQuickSwitch"
                    @date-selected="triggerSelect"
                    :selected-range="selectedDateRange"
                    :selected-range-tmp="selectedDateRangeTmp"
                    :initDate="initStartDate"
                    :bkDate="bkDateStart"
                    :timer="timer"></date-picker>

                <!-- 结束日期选择面板 -->
                <date-picker class="end-date date-select-container fl"
                    ref="endDater"
                    @date-quick-switch="dateQuickSwitch"
                    @date-selected="triggerSelect"
                    :selected-range="selectedDateRange"
                    :selected-range-tmp="selectedDateRangeTmp"
                    :initDate="initEndDate"
                    :bkDate="bkDateEnd"
                    :timer="timer"></date-picker>

                <!-- 日期快速选择配置 -->
                <div class="range-config fl" v-if="quickSelect">
                    <a href="javascript:;"
                        :class="{'active': shouldBeMatched(range)}"
                        v-for="range in defaultRanges"
                        @click.stop="changeRanges(range)">{{range.text}}</a>
                </div>
                <div class="range-action fl" v-if="quickSelect">
                    <a href="javascript:;" @click.stop="showDatePanel = false">{{t('dateRange.ok')}}</a>
                    <a href="javascript:;" @click.stop="clear">{{t('dateRange.clear')}}</a>
                </div>
            </div>
        </transition>
    </div>
</template>
<script>
    import clickoutside from '../../directives/clickoutside'
    import datepicker from './date-picker'
    // import moment from 'moment'

    // import {
    //     format,
    //     subDays, subMonths, addMonths,
    //     getHours, getMinutes, getSeconds,
    //     setSeconds, setMinutes, setHours,
    //     isAfter, isBefore, isSameYear, isSameMonth, isSameDay, isSameHour, isSameMinute, isSameSecond
    // } from 'date-fns'

    import format from 'date-fns/format'
    import subDays from 'date-fns/sub_days'
    import subMonths from 'date-fns/sub_months'
    import addMonths from 'date-fns/add_months'
    import getHours from 'date-fns/get_hours'
    import getMinutes from 'date-fns/get_minutes'
    import getSeconds from 'date-fns/get_seconds'
    import setSeconds from 'date-fns/set_seconds'
    import setMinutes from 'date-fns/set_minutes'
    import setHours from 'date-fns/set_hours'
    import isAfter from 'date-fns/is_after'
    import isBefore from 'date-fns/is_before'
    import isSameYear from 'date-fns/is_same_year'
    import isSameMonth from 'date-fns/is_same_month'
    import isSameDay from 'date-fns/is_same_day'
    import isSameHour from 'date-fns/is_same_hour'
    import isSameMinute from 'date-fns/is_same_minute'
    import isSameSecond from 'date-fns/is_same_second'
    import differenceInMonths from 'date-fns/difference_in_months'
    import differenceInYears from 'date-fns/difference_in_years'

    import locale from '../../mixins/locale'

    const oneOf = (value, validList) => {
        for (let i = 0; i < validList.length; i++) {
            if (value === validList[i]) {
                return true
            }
        }
        return false
    }

    class BkDate {
        constructor (flag, weekdays, time) {
            this.weekdays = weekdays

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
                second: dater.getSeconds()
            }
            // 日期选择器默认显示
            this.year = this.currentDay.year
            this.month = this.currentDay.month
            this.day = this.currentDay.day

            this.setTimer = false

            this.index = flag === 'start' ? 0 : 1

            // 年份向前按钮是否禁用
            this.preYearDisabled = false
            // 月份向前按钮是否禁用
            this.preMonthDisabled = false
            // 年份向后按钮是否禁用
            this.nextYearDisabled = false
            // 月份向后按钮是否禁用
            this.nextMonthDisabled = false
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
            return parseInt(value) < 10 ? (0 + '' + value) : value
        }

        // 获取当天格式化日期
        getFormatToday () {
            // return `${this.currentDay.year}-${this.formatDateString(this.currentDay.month)}-${this.formatDateString(this.currentDay.day)}`
            return this.currentDay.year
                + '-'
                + this.formatDateString(this.currentDay.month)
                + '-'
                + this.formatDateString(this.currentDay.day)
        }

        // 获取当前格式化日期
        getFormatDate () {
            // return `${this.year}-${this.formatDateString(this.month)}-${this.formatDateString(this.day)}`
            return this.year
                + '-'
                + this.formatDateString(this.month)
                + '-'
                + this.formatDateString(this.day)
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
        name: 'bk-date-range',
        mixins: [locale],
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
            ranges: {
                type: Object,
                default: () => {}
                // default: () => {
                //     return {
                //         // '昨天': [moment().subtract(1, 'days'), moment()],
                //         // '最近一周': [moment().subtract(7, 'days'), moment()],
                //         // '最近一个月': [moment().subtract(1, 'month'), moment()],
                //         // '最近三个月': [moment().subtract(3, 'month'), moment()]
                //         '昨天': [subDays(new Date(), 1), new Date()],
                //         '最近一周': [subDays(new Date(), 7), new Date()],
                //         '最近一个月': [subMonths(new Date(), 1), new Date()],
                //         '最近三个月': [subMonths(new Date(), 3), new Date()]
                //     }
                // }
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
            // const sdr = [moment().format('YYYY-MM-DD HH:mm:ss')]
            // const sdrt = [moment(sdr[0]).format('YYYY-MM-DD')]
            const sdr = [format(new Date(), 'YYYY-MM-DD HH:mm:ss')]
            const sdrt = [format(sdr[0], 'YYYY-MM-DD')]

            // const initStartDate = this.startDate || moment().subtract(1, 'month').format('YYYY-MM-DD')
            // const initEndDate = this.endDate || moment().format('YYYY-MM-DD')

            const weekdays = [
                this.t('dateRange.datePicker.weekdays.sun'),
                this.t('dateRange.datePicker.weekdays.mon'),
                this.t('dateRange.datePicker.weekdays.tue'),
                this.t('dateRange.datePicker.weekdays.wed'),
                this.t('dateRange.datePicker.weekdays.thu'),
                this.t('dateRange.datePicker.weekdays.fri'),
                this.t('dateRange.datePicker.weekdays.sat')
            ]
            const bkDateStart = this.startDate ? new BkDate('start', weekdays, this.startDate) : new BkDate('start', weekdays)
            const bkDateEnd = this.endDate ? new BkDate('end', weekdays, this.endDate) : new BkDate('end', weekdays)

            let initStartDate = this.startDate ? format(this.startDate, 'YYYY-MM-DD') : format(subMonths(new Date(), 1), 'YYYY-MM-DD')
            let initStartDateCopy = this.startDate ? format(this.startDate, 'YYYY-MM') : format(subMonths(new Date(), 1), 'YYYY-MM')
            let initEndDate = this.endDate ? format(this.endDate, 'YYYY-MM-DD') : format(new Date(), 'YYYY-MM-DD')
            let initEndDateCopy = this.endDate ? format(this.endDate, 'YYYY-MM') : format(new Date(), 'YYYY-MM')
            if (initStartDateCopy === initEndDateCopy) {
                initEndDate = format(addMonths(initEndDate, 1), 'YYYY-MM-DD')
                bkDateStart.nextMonthDisabled = true
                bkDateStart.nextYearDisabled = true
                bkDateEnd.preMonthDisabled = true
                bkDateEnd.preYearDisabled = true
            }
            else {
                // 右边和左边相差 大于 12 个月即一年
                if (differenceInMonths(initEndDate, initStartDate) > 12) {
                    bkDateEnd.preMonthDisabled = false
                    bkDateEnd.preYearDisabled = false
                    bkDateStart.nextMonthDisabled = false
                    bkDateStart.nextYearDisabled = false
                }
                else {
                    bkDateEnd.preYearDisabled = true
                    bkDateStart.nextYearDisabled = true
                    if (differenceInMonths(initEndDate, initStartDate) > 1) {
                        bkDateEnd.preMonthDisabled = false
                        bkDateStart.nextMonthDisabled = false
                    }
                    else {
                        bkDateEnd.preMonthDisabled = true
                        bkDateStart.nextMonthDisabled = true
                    }
                }
            }

            bkDateStart.setDate(initStartDate)
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
                defaultPlaceholder: this.t('dateRange.selectDate'),
                defaultRanges: [ // 默认快捷选择菜单栏
                    {
                        text: this.t('dateRange.yestoday'),
                        // value: [moment().subtract(1, 'days'), moment()]
                        value: [subDays(new Date(), 1), new Date()]
                    },
                    {
                        text: this.t('dateRange.lastweek'),
                        // value: [moment().subtract(7, 'days'), moment()]
                        value: [subDays(new Date(), 7), new Date()]
                    },
                    {
                        text: this.t('dateRange.lastmonth'),
                        // value: [moment().subtract(1, 'month'), moment()]
                        value: [subMonths(new Date(), 1), new Date()]
                    },
                    {
                        text: this.t('dateRange.last3months'),
                        // value: [moment().subtract(3, 'month'), moment()]
                        value: [subMonths(new Date(), 3), new Date()]
                    }
                ]
            }
        },
        directives: {
            clickoutside
        },
        created () {
            if (this.ranges && Object.keys(this.ranges).length) {
                const defaultRanges = []
                Object.keys(this.ranges).forEach(range => {
                    defaultRanges.push({
                        text: range,
                        value: this.ranges[range]
                    })
                })
                this.defaultRanges.splice(0, this.defaultRanges.length, ...defaultRanges)
            }

            if (this.placeholder) {
                this.defaultPlaceholder = this.placeholder
            }

            // if (this.startDate) {
            //     this.selectedDateRange.unshift(this.startDate)
            // }
            // if (this.endDate) {
            //     this.selectedDateRange.pop()
            //     this.selectedDateRange.push(this.endDate)
            // }

            let hour = ''
            let minute = ''
            let second = ''

            if (this.startDate) {
                hour = getHours(new Date(this.startDate))
                minute = getMinutes(new Date(this.startDate))
                second = getSeconds(new Date(this.startDate))
                this.selectedDateRange.unshift(
                    this.timer
                        // ? moment(this.startDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                        ? format(setHours(setMinutes(setSeconds(this.startDate, second), minute), hour), 'YYYY-MM-DD HH:mm:ss')
                        : this.startDate
                )
                this.selectedDateRangeTmp.unshift(
                    this.timer
                        // ? moment(this.startDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD')
                        ? format(setHours(setMinutes(setSeconds(this.startDate, second), minute), hour), 'YYYY-MM-DD')
                        : this.startDate
                )
            }
            if (this.endDate) {
                hour = getHours(new Date(this.endDate))
                minute = getMinutes(new Date(this.endDate))
                second = getSeconds(new Date(this.endDate))
                this.selectedDateRange.pop()
                this.selectedDateRange.push(
                    this.timer
                        // ? moment(this.endDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                        ? format(setHours(setMinutes(setSeconds(this.endDate, second), minute), hour), 'YYYY-MM-DD HH:mm:ss')
                        : this.endDate
                )
                this.selectedDateRangeTmp.pop()
                this.selectedDateRangeTmp.push(
                    this.timer
                        // ? moment(this.endDate).hour(hour).minute(minute).second(second).format('YYYY-MM-DD')
                        ? format(setHours(setMinutes(setSeconds(this.endDate, second), minute), hour), 'YYYY-MM-DD')
                        : this.endDate
                )
            }
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
            showDatePanel (val) {
                if (!val) {
                    this.$emit('close', this.selectedDateView)
                } else {
                    this.$emit('show', this.selectedDateView)
                }
            },
            selectedDateRange () {
                const seed = this.timer ? this.selectedDateRange : this.selectedDateRangeTmp
                if (seed.length === 2) {
                    // let newSelectedDate = seed.join(` ${this.rangeSeparator} `)
                    let newSelectedDate = seed.join(' ' + this.rangeSeparator + ' ')
                    // 判断此次选择的日期范围是否发生改变
                    if (this.selectedDateView !== newSelectedDate) {
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
                // let formatDateStart = moment(this.selectedDateRange[0]).format('YYYY-MM')
                // let formatDateEnd = moment(this.selectedDateRange[1]).format('YYYY-MM')
                let formatDateStart = format(this.selectedDateRange[0], 'YYYY-MM')
                let formatDateEnd = format(this.selectedDateRange[1], 'YYYY-MM')

                // this.initStartDate = this.selectedDateView.split(` ${this.rangeSeparator} `)[0]
                this.initStartDate = this.selectedDateView.split(' ' + this.rangeSeparator + ' ')[0]
                // 已选日期范围是同一个月时 重置两个月份选择面板 --  显示相邻月份
                // if (moment(formatDateStart).isSame(formatDateEnd)) {
                if (isSameYear(formatDateStart, formatDateEnd) && isSameMonth(formatDateStart, formatDateEnd)) {
                    // this.initEndDate = moment(this.selectedDateRange[0]).add(1, 'month').format(
                    //     this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'
                    // )
                    this.initEndDate = format(
                        addMonths(this.selectedDateRange[0], 1),
                        this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'
                    )
                } else {
                    this.initEndDate = this.selectedDateView.split(' ' + this.rangeSeparator + ' ')[1]
                }
            }
        },
        methods: {
            // 日期快速切换回调
            dateQuickSwitch (date) {
                // 左边日期面板信息
                let startDateInfo = this.$refs.startDater
                let startTopDate = startDateInfo.topBarFormatView.value

                // 右边日期面板信息
                let endDateInfo = this.$refs.endDater
                let endTopDate = endDateInfo.topBarFormatView.value

                if (startTopDate === endTopDate) {
                    switch (date.type) {
                        case 'next':
                            // endDateInfo.BkDate.setDate(moment(endTopDate).add(1, 'month').format('YYYY-MM'))
                            endDateInfo.BkDate.setDate(format(addMonths(endTopDate, 1), 'YYYY-MM'))
                            break
                        case 'last':
                            // startDateInfo.BkDate.setDate(moment(startTopDate).add(-1, 'month').format('YYYY-MM'))
                            startDateInfo.BkDate.setDate(format(addMonths(startTopDate, -1), 'YYYY-MM'))
                            break
                        default:
                            break
                    }
                }

                // 右边和左边相差 大于 12 个月即一年
                if (differenceInMonths(endTopDate, startTopDate) > 12) {
                    this.bkDateEnd.preMonthDisabled = false
                    this.bkDateEnd.preYearDisabled = false
                    this.bkDateStart.nextMonthDisabled = false
                    this.bkDateStart.nextYearDisabled = false
                }
                else {
                    this.bkDateEnd.preYearDisabled = true
                    this.bkDateStart.nextYearDisabled = true
                    if (differenceInMonths(endTopDate, startTopDate) > 1) {
                        this.bkDateEnd.preMonthDisabled = false
                        this.bkDateStart.nextMonthDisabled = false
                    }
                    else {
                        this.bkDateEnd.preMonthDisabled = true
                        this.bkDateStart.nextMonthDisabled = true
                    }
                }
            },

            changeRanges (range) {
                let rangeStartDate = format(range.value[0], 'YYYY-MM-DD HH:mm:ss')
                let rangeEndDate = format(range.value[1], 'YYYY-MM-DD HH:mm:ss')
                let rangeStartDateTmp = format(range.value[0], 'YYYY-MM-DD')
                let rangeEndDateTmp = format(range.value[1], 'YYYY-MM-DD')

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
                // let formatDateStart = moment(rangeStartDate).format('YYYY-MM')
                // let formatDateEnd = moment(rangeEndDate).format('YYYY-MM')
                let formatDateStart = format(rangeStartDate, 'YYYY-MM')
                let formatDateEnd = format(rangeEndDate, 'YYYY-MM')
                // if (moment(formatDateStart).isSame(formatDateEnd)) {
                if (isSameYear(formatDateStart, formatDateEnd) && isSameMonth(formatDateStart, formatDateEnd)) {
                    // this.initEndDate = moment(this.selectedDateRange[0]).add(1, 'month').format(
                    //     this.timer ? 'YYYY-MM-DD HH:mm:ss' : 'YYYY-MM-DD'
                    // )
                    this.initEndDate = format(
                        addMonths(this.selectedDateRange[0], 1),
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
                    // const hour = moment(date).hour()
                    // const minute = moment(date).minute()
                    // const second = moment(date).second()
                    const hour = getHours(date)
                    const minute = getMinutes(date)
                    const second = getSeconds(date)
                    this.selectedDateRange.splice(
                        bkDateIndex,
                        1,
                        // moment(this.selectedDateRange[bkDateIndex])
                        //     .hour(hour).minute(minute).second(second).format('YYYY-MM-DD HH:mm:ss')
                        format(
                            setHours(
                                setMinutes(
                                    setSeconds(
                                        this.selectedDateRange[bkDateIndex],
                                        second
                                    ), minute
                                ), hour
                            ), 'YYYY-MM-DD HH:mm:ss'
                        )
                    )
                } else {
                    const selectedLen = this.selectedDateRange.length
                    // date = moment(date).format('YYYY-MM-DD HH:mm:ss')
                    // const dateTmp = moment(date).format('YYYY-MM-DD')
                    date = format(date, 'YYYY-MM-DD HH:mm:ss')
                    const dateTmp = format(date, 'YYYY-MM-DD')

                    switch (selectedLen) {
                        case 0:
                            this.selectedDateRange.push(date)
                            this.selectedDateRangeTmp.push(dateTmp)
                            break
                        case 1:
                            // 首先验证第二次选择日期是否合格
                            // if (moment(date).isSame(this.selectedDateRange[0]) || moment(date).isAfter(this.selectedDateRange[0])) {
                            if (
                                (
                                    isSameYear(date, this.selectedDateRange[0])
                                        && isSameMonth(date, this.selectedDateRange[0])
                                        && isSameDay(date, this.selectedDateRange[0])
                                        && isSameHour(date, this.selectedDateRange[0])
                                        && isSameMinute(date, this.selectedDateRange[0])
                                        && isSameSecond(date, this.selectedDateRange[0])
                                )
                                || isAfter(date, this.selectedDateRange[0])
                            ) {
                                this.selectedDateRange.push(date)
                                this.selectedDateRangeTmp.push(dateTmp)
                            }
                            // if (moment(date).isBefore(this.selectedDateRange[0])) {
                            if (isBefore(date, this.selectedDateRange[0])) {
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
                let isMatched = this.selectedDateRangeTmp[0] === format(range.value[0], 'YYYY-MM-DD')
                    && this.selectedDateRangeTmp[1] === format(range.value[1], 'YYYY-MM-DD')
                return isMatched
            },

            // 控制选择器显示隐藏
            openDater () {
                if (this.disabled) {
                    return
                }
                this.showDatePanel = true
            },

            close () {
                // todo: 已选情况取消了选中状态然后点击空白处关闭日期选择器时，还原默认值选中状态
                if (this.selectedDateView && this.selectedDateRange.length === 0) {
                    // this.selectedDateRange = this.selectedDateView.split(` ${this.rangeSeparator} `)
                    this.selectedDateRange = this.selectedDateView.split(' ' + this.rangeSeparator + ' ')
                }
                this.showDatePanel = false
            },

            clear () {
                this.$emit('change', this.selectedDateView, '')

                this.selectedDateView = ''
                this.selectedDateRange = []
                this.selectedDateRangeTmp = []

                // const date = moment().format('YYYY-MM-DD HH:mm:ss')
                // const dateTmp = moment(date).format('YYYY-MM-DD')
                const date = format(new Date(), 'YYYY-MM-DD HH:mm:ss')
                const dateTmp = format(date, 'YYYY-MM-DD')
                this.selectedDateRange.push(date)
                this.selectedDateRangeTmp.push(dateTmp)

                // this.initStartDate = this.startDate || moment().subtract(1, 'month').format('YYYY-MM-DD')
                // this.initEndDate = this.endDate || moment().format('YYYY-MM-DD')
                this.initStartDate = this.startDate || format(subMonths(new Date(), 1), 'YYYY-MM-DD')
                this.initEndDate = this.endDate || format(new Date(), 'YYYY-MM-DD')
                this.bkDateStart.setDate(this.initStartDate)
                this.bkDateEnd.setDate(this.initEndDate)

                this.bkDateStart.currentTime = this.bkDateEnd.currentTime = {
                    // hour: moment().hour(),
                    // minute: moment().minute(),
                    // second: moment().second()
                    hour: getHours(new Date()),
                    minute: getMinutes(new Date()),
                    second: getSeconds(new Date())
                }

                this.showDatePanel = false
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/date-range.scss'
</style>
