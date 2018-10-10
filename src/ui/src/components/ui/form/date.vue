<template>
    <div class="form-date">
        <bk-date-picker class="form-date-picker"
            :init-date="initDate"
            :start-date="startDate"
            :end-date="endDate"
            :disabled="disabled"
            @date-selected="handleDateSelected"
            @change="handleChange">
        </bk-date-picker>
        <i class="bk-icon icon-close-circle-shape" @click="handleClear" hidden></i>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-date',
        props: {
            value: {
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
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                initDate: this.value
            }
        },
        watch: {
            value () {
                this.setInitDate()
            },
            initDate (initDate) {
                this.$emit('input', initDate)
                this.$emit('on-select', initDate)
            }
        },
        created () {
            this.setInitDate()
        },
        methods: {
            setInitDate () {
                this.initDate = this.value
            },
            handleDateSelected (date) {
                this.initDate = date
            },
            handleChange (oldVal, newVal) {
                if (oldVal !== newVal) {
                    this.$emit('on-change', newVal, oldVal)
                }
            },
            handleClear () {
                this.initDate = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-date {
        position: relative;
    }
    .icon-close-circle-shape{
        position: absolute;
        right: 40px;
        top: 10px;
        width: 18px;
        height: 18px;
        line-height: 18px;
        border-radius: 50%;
        color: rgb(204, 204, 204);
        text-align: center;
        font-size: 18px;
        opacity: .7;
        transition: opacity linear .2s;
        cursor: pointer;
        &:hover{
            opacity: 1;
        }
    }
    .form-date-picker {
        width: 100%;
    }
</style>