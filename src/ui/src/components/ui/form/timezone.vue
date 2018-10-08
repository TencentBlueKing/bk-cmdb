<template>
    <div class="cmdb-form form-timezone">
        <bk-selector class="form-timezone-selector"
            :searchable="true"
            :list="timezoneList"
            :disabled="disabled"
            :selected.sync="selected">
        </bk-selector>
    </div>
</template>

<script>
    import TIMEZONE from './timezone.json'
    export default {
        name: 'cmdb-form-timezone',
        props: {
            value: {
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            const timezoneList = TIMEZONE.map(timezone => {
                return {
                    id: timezone,
                    name: timezone
                }
            })
            return {
                timezoneList,
                selected: ''
            }
        },
        watch: {
            value (value) {
                if (value !== this.selected) {
                    this.selected = value
                }
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            }
        },
        created () {
            this.selected = this.value ? this.value : 'Asia/Shanghai'
        }
    }
</script>

<style lang="scss" scoped>
    .form-timezone{
        .form-timezone-selector{
            width: 100%;
        }
    }
</style>