<template>
    <div class="bk-date-picker" :class="disabled ? 'disabled' : ''" @click="openDater" v-clickoutside="close">
        <input type="text" name="date-select" readonly :disabled="disabled" :placeholder="t('datePicker.selectDate')" v-model="selectedValue">
        <transition :name="transitionName">
            <div :style="panelStyle" class="date-dropdown-panel" v-if="showDatePanel">
                <!-- 日期操作栏 Start -->
                <div class="date-top-bar">
                    <span class="year-switch-icon pre-year fl" @click="switchToYear('last')"></span>
                    <span class="month-switch-icon pre-month fl" @click="switchToMonth('last')"></span>
                    <span class="current-date">{{topBarFormatView}}</span>
                    <span class="year-switch-icon next-year fr" @click="switchToYear('next')"></span>
                    <span class="month-switch-icon next-month fr" @click="switchToMonth('next')"></span>
                </div>
                <!-- 日期操作栏 End -->

                <!-- 日期选择面板 Start -->
                <div class="date-select-panel">
                    <dl>
                        <dt>
                            <span class="date-item-view" v-for="day in BkDate.weekdays" v-html="day"></span>
                        </dt>
                        <!-- 上月部分日期显示 -->
                        <dd v-for="lastMonthItem in lastMonthList">
                            <span class="date-item-view date-disable-item">{{ lastMonthItem }}</span>
                        </dd>

                        <!-- 本月高亮日期显示 -->
                        <dd v-for="currentMonthItem in BkDate.getCurrentMouthDays()">
                            <span
                              :class="{
                                 'date-table-item': isAvailableDate(currentMonthItem),
                                 'date-item-view date-disable-item': !isAvailableDate(currentMonthItem),
                                 'selected': shouldBeSelected(currentMonthItem),
                                 'today': shouldShowToday(currentMonthItem) === t('datePicker.today')
                              }"
                              @click.stop.prevent="selectDay(currentMonthItem)">{{shouldShowToday(currentMonthItem)}}</span>
                        </dd>

                        <!-- 下个月部分日期显示 -->
                        <dd v-for="nextMonthItem in nextMonthList">
                            <span class="date-item-view date-disable-item">{{ nextMonthItem }}</span>
                        </dd>

                    </dl>
                </div>
                <!-- 日期选择面板End -->

                <div class="time-set-panel" v-if="timer">
                    <div class="time-item" v-for="(timeItem, index) in currentTime"><input type="number" name="" :value="timeItem">
                        <span class="time-option fr">
                            <i class="up" @click="setTime('up', index)"></i>
                            <i class="down" @click="setTime('down', index)"></i>
                        </span>
                    </div>
                </div>
            </div>
        </transition>
  </div>
</template>

<script>
    /**
     * bk-date
     * 参数配置：
     * @param timer: true/false  -是否配置时间设置项 默认配置
     * @param initDate: 'YYYY-MM-DD' -初始化值，默认为空
     * @param startDate: 'YYYY-MM-DD' -开始日期
     * @param endDate: 'YYYY-MM-DD' -结束日期
     * @param autoClose: true/false -选择日期后是否自动关闭
    */

    import format from 'date-fns/format'

    import clickoutside from '../../directives/clickoutside'
    import locale from '../../mixins/locale'
    import {formatMonth} from '../../util'

    const oneOf = (value, validList) => {
        for (let i = 0; i < validList.length; i++) {
            if (value === validList[i]) {
                return true
            }
        }
        return false
    }

    class BkDate {
        constructor (weekdays) {
            this.weekdays = weekdays

            const dater = new Date()

            // 当前日期
            this.currentDay = {
                year: dater.getFullYear(),
                month: dater.getMonth() + 1,
                day: dater.getDate()
            }
            // 当前时间
            this.currentTime = {
                hour: dater.getHours(),
                minute: dater.getMinutes() + 1,
                second: dater.getSeconds()
            }
            // 日期选择器默认显示
            this.year = this.currentDay.year
            this.month = this.currentDay.month
            this.day = this.currentDay.day
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
            return parseInt(value) < 10 ? '0' + value : value
        }

        // 获取当前格式化日期
        getFormatDate () {
            // return `${this.year}-${this.formatDateString(this.month)}-${this.formatDateString(this.day)}`
            return this.year + '-' + this.formatDateString(this.month) + '-' + this.formatDateString(this.day)
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

    export default {
        name: 'bk-date-picker',
        mixins: [locale],
        directives: {
            clickoutside
        },
        props: {
            autoClose: {
                type: Boolean,
                default: true
            },
            disabled: {
                type: Boolean,
                default: false
            },
            timer: {
                type: Boolean,
                default: false
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
            position: {
                validator (value) {
                    return oneOf(value, ['top', 'bottom'])
                },
                default: 'bottom'
            }
        },
        data () {
            let transitionName = 'toggle-slide'
            const panelStyle = {}
            const positionArr = this.position.split('-')
            if (positionArr.indexOf('top') > -1) {
                panelStyle.bottom = '38px'
                transitionName = 'toggle-slide2'
            } else {
                panelStyle.top = '38px'
            }

            const weekdays = [
                this.t('datePicker.weekdays.sun'),
                this.t('datePicker.weekdays.mon'),
                this.t('datePicker.weekdays.tue'),
                this.t('datePicker.weekdays.wed'),
                this.t('datePicker.weekdays.thu'),
                this.t('datePicker.weekdays.fri'),
                this.t('datePicker.weekdays.sat')
            ]

            const bkDate = new BkDate(weekdays)

            return {
                panelStyle: panelStyle,
                transitionName: transitionName,
                BkDate: bkDate,
                selectedValue: this.initDate || '',
                currentDate: new Date(),
                showDatePanel: false,
                isSetTimer: false
            }
        },
        computed: {
            // 切换列表年月日期显示
            topBarFormatView () {
                // console.log(this.t('datePicker.topBarFormatView'))
                // const date = new Date(this.BkDate.year, this.BkDate.month - 1)
                // console.log(date)
                // console.log(formatDateLabels('en-US', this.t('datePicker.topBarFormatView'), date))
                // console.log(this.t('datePicker.test', {vari: 'dddd', ccc: '是'}))
                // console.log(this.t('datePicker.test1', {mm: 'dddd', yyyy: '是'}))
                // console.log(this.t('datePicker.test2', {mm: this.BkDate.month >= 10 ? this.BkDate.month : '0' + this.BkDate.month, yyyy: this.BkDate.year}))
                // console.error(this.t('datePicker.test1', {mm: formatMonth(12, this.t('lang'), true), yyyy: this.BkDate.year}))
                // console.error(this.t('datePicker.test3', {mmmm: formatMonth(12, this.t('lang')), yyyy: this.BkDate.year}))

                const month = this.BkDate.month >= 10 ? this.BkDate.month : '0' + this.BkDate.month

                // return `${this.BkDate.year}年${this.BkDate.month >= 10 ? this.BkDate.month : '0' + this.BkDate.month}月`
                // return this.BkDate.year + '年' + (this.BkDate.month >= 10 ? this.BkDate.month : '0' + this.BkDate.month + '月')
                return this.t('datePicker.topBarFormatView', {
                    mmmm: formatMonth(month, this.t('lang')),
                    mm: formatMonth(month, this.t('lang'), true),
                    yyyy: this.BkDate.year
                })
            },
            // 上个月部分日期显示列表
            lastMonthList () {
                let lastMonthVisibleNum = this.BkDate.getCurrentMonthBeginWeek()
                let lastMonthDays = this.BkDate.getLastMouthDays()
                let lastMonthVisibleList = []
                for (let i = lastMonthVisibleNum - 1; i >= 0; i--) {
                    lastMonthVisibleList.push(lastMonthDays - i)
                }
                return lastMonthVisibleList
            },
            // 下个月部分日期显示列表
            nextMonthList () {
                let lastMonthVisibleNum = this.BkDate.getCurrentMonthBeginWeek()
                let currentMonthDays = this.BkDate.getCurrentMouthDays()
                let nextMonthVisibleList = 42 - lastMonthVisibleNum - currentMonthDays
                return nextMonthVisibleList
            },
            currentTime () {
                const time = [
                    this.formatValue(this.BkDate.currentTime.hour),
                    this.formatValue(this.BkDate.currentTime.minute),
                    this.formatValue(this.BkDate.currentTime.second)
                ]
                return time
            }
        },
        watch: {
            initDate () {
                this.BkDate.setDate(this.initDate)
                if (this.selectedValue !== this.initDate) {
                    this.$emit('change', this.selectedValue, this.initDate)
                }
                this.showDate()
                this.$emit('date-selected', this.selectedValue)
                // 是否关闭日期选择
                if (this.autoClose && !this.isSetTimer) {
                    this.close()
                }
            }
        },
        created () {
            this.BkDate.setDate(this.initDate)
        },
        methods: {
            // 选择日期
            selectDay (value) {
                if (!this.isAvailableDate(value)) {
                    return
                }
                // let newSelectedDate = `${this.BkDate.year}-${this.BkDate.month}-${value}`
                let newSelectedDate = this.BkDate.year
                    + '-'
                    + this.formatValue(this.BkDate.month)
                    + '-'
                    + this.formatValue(value)

                if (this.timer) {
                    newSelectedDate += ' '
                        + this.formatValue(this.BkDate.currentTime.hour)
                        + ':'
                        + this.formatValue(this.BkDate.currentTime.minute)
                        + ':'
                        + this.formatValue(this.BkDate.currentTime.second)
                }

                // change 回调
                if (this.selectedValue !== newSelectedDate) {
                    this.$emit('change', this.selectedValue, newSelectedDate)
                }

                this.BkDate.setDate(newSelectedDate)
                this.showDate()

                this.$emit('date-selected', this.selectedValue)

                // 是否关闭日期选择
                if (this.autoClose) {
                    this.close()
                }
            },

            // 同步显示日期
            showDate () {
                // let selectedDate = `${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(this.BkDate.day)}`
                let selectedDate = this.BkDate.year
                    + '-'
                    + this.formatValue(this.BkDate.month)
                    + '-'
                    + this.formatValue(this.BkDate.day)

                let selectedTime
                if (this.timer) {
                    // selectedTime = ` ${this.formatValue(this.BkDate.currentTime.hour)}:${this.formatValue(this.BkDate.currentTime.minute)}:${this.formatValue(this.BkDate.currentTime.second)}`
                    selectedTime = ' '
                        + this.formatValue(this.BkDate.currentTime.hour)
                        + ':'
                        + this.formatValue(this.BkDate.currentTime.minute)
                        + ':'
                        + this.formatValue(this.BkDate.currentTime.second)
                } else {
                    selectedTime = ''
                }
                this.selectedValue = selectedDate + '' + selectedTime
            },

            formatValue (value) {
                return parseInt(value) < 10 ? '0' + value : value
            },

            // 高亮显示已选日期
            shouldBeSelected (value) {
                return this.BkDate.day === value
            },

            // 标记今天
            shouldShowToday (value) {
                const currentSelectedDate = {
                    year: this.BkDate.year,
                    month: this.BkDate.month,
                    day: value
                }
                const current = {
                    year: this.currentDate.getFullYear(),
                    month: this.currentDate.getMonth() + 1,
                    day: this.currentDate.getDate()
                }
                let isToday = JSON.stringify(currentSelectedDate) === JSON.stringify(current)
                if (isToday) {
                    return this.t('datePicker.today')
                }
                return value
            },

            // 切换月份
            switchToMonth (type) {
                const toMonthDate = {}
                let year = this.BkDate.year
                let month = this.BkDate.month
                switch (type) {
                    case 'last':
                        toMonthDate.year = month - 1 > 0 ? year : year - 1
                        toMonthDate.month = month - 1 > 0 ? month - 1 : 12
                        break
                    case 'next':
                        toMonthDate.year = month + 1 > 12 ? year + 1 : year
                        toMonthDate.month = month + 1 > 12 ? 1 : month + 1
                        break
                    default:
                        break
                }
                // this.BkDate.setDate(`${toMonthDate.year}-${toMonthDate.month}-${this.BkDate.day}`)
                this.BkDate.setDate(toMonthDate.year + '-' + toMonthDate.month + '-' + this.BkDate.day)
            },

            // 切换年份
            switchToYear (type) {
                const toYearDate = {}
                let year = this.BkDate.year
                switch (type) {
                    case 'last':
                        toYearDate.year = year - 1 > 0 ? year - 1 : 0
                        break
                    case 'next':
                        toYearDate.year = year + 1
                        break
                    default:
                        break
                }

                // this.BkDate.setDate(`${toYearDate.year}-${this.BkDate.month}-${this.BkDate.day}`)
                this.BkDate.setDate(toYearDate.year + '-' + this.BkDate.month + '-' + this.BkDate.day)
            },

            // 时间设置
            setTime (type, index) {
                const option = ['hour', 'minute', 'second'][index]
                let defaultTime = {...this.BkDate.currentTime}
                switch (option) {
                    case 'hour':
                        if (type === 'up') {
                            defaultTime.hour = defaultTime.hour + 1 < 24 ? defaultTime.hour + 1 : 0
                        }
                        if (type === 'down') {
                            defaultTime.hour = defaultTime.hour - 1 >= 0 ? defaultTime.hour - 1 : 23
                        }
                        break
                    case 'minute':
                        if (type === 'up') {
                            defaultTime.minute = defaultTime.minute + 1 < 60 ? defaultTime.minute + 1 : 0
                        }
                        if (type === 'down') {
                            defaultTime.minute = defaultTime.minute - 1 >= 0 ? defaultTime.minute - 1 : 59
                        }
                        break
                    case 'second':
                        if (type === 'up') {
                            defaultTime.second = defaultTime.second + 1 < 60 ? defaultTime.second + 1 : 0
                        }
                        if (type === 'down') {
                            defaultTime.second = defaultTime.second - 1 >= 0 ? defaultTime.second - 1 : 59
                        }
                        break
                    default:
                }

                this.BkDate.currentTime = {...defaultTime}
                this.showDate()
                this.isSetTimer = true
                if (this.selectedValue !== this.initDate) {
                    this.$emit('change', this.selectedValue, this.initDate)
                }
                this.$emit('date-selected', this.selectedValue)
            },

            // 判断日期是否可选
            isAvailableDate (day) {
                // let cmpTime = new Date(`${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(day)}`).getTime()
                let cmpTime = new Date(this.BkDate.year + '-' + this.formatValue(this.BkDate.month) + '-' + this.formatValue(day)).getTime()
                let startTime, endTime
                let checkStartTime = true
                let checkEndTime = true
                if (this.startDate) {
                    startTime = new Date(this.startDate).getTime()
                    checkStartTime = cmpTime >= startTime
                }
                if (this.endDate) {
                    endTime = new Date(this.endDate).getTime()
                    checkEndTime = (cmpTime <= endTime)
                }
                return checkStartTime && checkEndTime
            },

            // 控制选择器显示隐藏
            openDater () {
                if (this.disabled) return
                this.showDatePanel = true
            },

            close () {
                this.showDatePanel = false
                this.isSetTimer = false
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/date-picker.scss'
</style>
