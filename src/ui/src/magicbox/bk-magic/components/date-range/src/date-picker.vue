<template>
    <div class="date-picker">
        <!-- 日期操作栏Start -->
        <div class="date-top-bar">
            <span class="year-switch-icon pre-year fl" @click="switchToYear('last')"></span>
            <span class="month-switch-icon pre-month fl" @click="switchToMonth('last')"></span>
            <span class="current-date">{{ topBarFormatView.text }}</span>
            <span class="year-switch-icon next-year fr" @click="switchToYear('next')"></span>
            <span class="month-switch-icon next-month fr" @click="switchToMonth('next')"></span>
        </div>
        <!-- 日期操作栏End -->

        <!-- 日期选择面板Start -->
        <div class="date-select-panel">
            <dl>
                <dt>
                    <span class="date-item-view" v-for="day in weekdays" v-html="day"></span>
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
                           'date-range-view': isInRange(currentMonthItem),
                           'today': shouldShowToday(currentMonthItem) === t('datepicker.now'),
                           'selected': shouldBeSelected(currentMonthItem)
                        }"
                        @click.stop.prevent="selectDay(currentMonthItem)">{{ shouldShowToday(currentMonthItem) }}</span>
                </dd>

                <!-- 下个月部分日期显示 -->
                <dd v-for="nextMonthItem in nextMonthList">
                    <span class="date-item-view date-disable-item">{{ nextMonthItem }}</span>
                </dd>

            </dl>
        </div>
        <!-- 日期选择面板End -->

        <div class="time-set-panel" v-if="timer">
            <template v-if="BkDate.setTimer">
                <div class="time-item" v-for="(timeItem, index) in currentTime">
                    <input readonly type="number" name="" :value="timeItem">
                    <span class="time-option fr">
                        <i class="up" @click.prevent.stop="setTime('up', index)"></i>
                        <i class="down" @click.prevent.stop="setTime('down', index)"></i>
                    </span>
                </div>
            </template>
            <template v-else>
                <div class="time-item" v-for="(timeItem, index) in currentTime">
                    <input disabled type="number" name="" :value="timeItem">
                    <span class="time-option fr">
                        <i class="up no-hover"></i>
                        <i class="down no-hover"></i>
                    </span>
                </div>
            </template>
        </div>

    </div>
</template>
<script>
    /**
    * bk-date
    * 参数配置：
    * @param initDate: 'YYYY-MM-DD' -初始化值，默认为空
    * @param startDate: 'YYYY-MM-DD' -开始日期
    * @param endDate: 'YYYY-MM-DD' -结束日期
    * @param selectedRange: [] -已选日期集合
    */

    export default{
        name: 'bk-date',
        props: {
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
            selectedRange: {
                type: Array,
                default: () => {
                    return []
                }
            },
            // selectedDateRange 的副本，用于 date-picker 中 shouldBeSelected 高亮判断仅需要判断日期而不需要判断时间
            selectedRangeTmp: {
                type: Array,
                default: () => {
                    return []
                }
            },
            timer: {
                type: Boolean,
                default: false
            },
            bkDate: {
                type: Object,
                default: () => {
                    return {}
                }
            }
        },
        data () {
            return {
                weekdays: [
                    this.t('datepicker.weeks.sun'),
                    this.t('datepicker.weeks.mon'),
                    this.t('datepicker.weeks.tue'),
                    this.t('datepicker.weeks.wed'),
                    this.t('datepicker.weeks.thu'),
                    this.t('datepicker.weeks.fri'),
                    this.t('datepicker.weeks.sat')
                ],
                BkDate: this.bkDate,
                selectedValue: this.initDate
            }
        },
        watch: {
            initDate (val) {
            }
        },
        computed: {
            // 切换列表年月日期显示
            topBarFormatView () {
                return {
                    text: `${this.BkDate.year}${this.t('datepicker.year')}${this.t('datepicker.month' + this.BkDate.month)}`,
                    value: `${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-01`
                }
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
        methods: {
            // 选择日期
            selectDay (value) {
                if (!this.isAvailableDate(value)) return
                // this.BkDate.setDate(`${this.BkDate.year}-${this.BkDate.month}-${value}`)
                // this.selectedValue = this.BkDate.getFormatDate()
                // this.$emit('date-selected', this.selectedValue)

                let newSelectedDate = `${this.BkDate.year}-${this.BkDate.month}-${value}`
                this.BkDate.setDate(newSelectedDate)

                let selectedDate = `${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(this.BkDate.day)}`
                let selectedTime = ''
                if (this.timer) {
                    selectedTime = ` ${this.formatValue(this.BkDate.currentTime.hour)}:${this.formatValue(this.BkDate.currentTime.minute)}:${this.formatValue(this.BkDate.currentTime.second)}`
                } else {
                    selectedTime = ''
                }
                this.selectedValue = `${selectedDate}${selectedTime}`
                this.$emit('date-selected', this.selectedValue)
            },

            // 月日小于10补零
            formatValue (value) {
                return parseInt(value) < 10 ? `0${value}` : value
            },

            // 高亮显示已选日期
            shouldBeSelected (value) {
                let selectedDate = `${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(value)}`

                // 高亮日期判断不需要加入时间
                // let selectedTime = ''
                // if (this.timer) {
                //     selectedTime = ` ${this.formatValue(this.BkDate.currentTime.hour)}:${this.formatValue(this.BkDate.currentTime.minute)}:${this.formatValue(this.BkDate.currentTime.second)}`
                // } else {
                //     selectedTime = ''
                // }
                // let triggerDate = `${selectedDate}${selectedTime}`

                let triggerDate = `${selectedDate}`
                return this.selectedRangeTmp.indexOf(triggerDate) >= 0
            },

            // 标记 '今天'
            shouldShowToday (value) {
                const currentSelectedDate = {
                    year: this.BkDate.year,
                    month: this.BkDate.month,
                    day: value
                }
                const currentDate = new Date()
                const current = {
                    year: currentDate.getFullYear(),
                    month: currentDate.getMonth() + 1,
                    day: currentDate.getDate()
                }
                let isToday = JSON.stringify(currentSelectedDate) === JSON.stringify(current)
                if (isToday) {
                    return this.t('datepicker.now')
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
                this.BkDate.setDate(`${toMonthDate.year}-${toMonthDate.month}-${this.BkDate.day}`)
                this.$emit('date-quick-switch', {
                    type: type,
                    value: `${toMonthDate.year}-${this.formatValue(toMonthDate.month)}-01`
                })
            },

            // 切换年份
            switchToYear (type) {
                const toYearDate = {
                    month: this.BkDate.month
                }
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

                this.BkDate.setDate(`${toYearDate.year}-${this.BkDate.month}-${this.BkDate.day}`)
                this.$emit('date-quick-switch', {
                    type: type,
                    value: `${toYearDate.year}-${this.formatValue(toYearDate.month)}-01`
                })
            },

            // 判断日期是否可选
            isAvailableDate (day) {
                let cmpTime = new Date(`${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(day)}`).getTime()
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

            // 判断日期是否在已选范围 -- 已选日期范围添加背景
            isInRange (day) {
                if (!this.selectedRange[0]) return false
                let dayTime = new Date(`${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(day)}`).getTime()
                let startDateTime = new Date(this.selectedRange[0]).getTime()
                let endDateTime = new Date(this.selectedRange[1]).getTime()
                return dayTime - startDateTime > 0 && dayTime - endDateTime < 0
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

                let selectedDate = `${this.BkDate.year}-${this.formatValue(this.BkDate.month)}-${this.formatValue(this.BkDate.day)}`
                let selectedTime = ''
                if (this.timer) {
                    selectedTime = ` ${this.formatValue(this.BkDate.currentTime.hour)}:${this.formatValue(this.BkDate.currentTime.minute)}:${this.formatValue(this.BkDate.currentTime.second)}`
                } else {
                    selectedTime = ''
                }
                this.selectedValue = `${selectedDate}${selectedTime}`

                this.$emit('date-selected', this.selectedValue, this.BkDate.index)
            }
        }
    }
</script>
