<template>
    <bk-select
        v-model="localValue"
        v-bind="$attrs"
        :clearable="false">
        <bk-option v-for="(option, index) in options"
            :key="index"
            v-bind="option">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            },
            type: {
                type: String,
                default: 'singlechar'
            }
        },
        computed: {
            options () {
                const EQ = '$eq'
                const NE = '$ne'
                const IN = '$in'
                const NIN = '$nin'
                const LTE = '$lte'
                const GTE = '$gte'
                const REGEX = '$regex'
                const RANGE = '$range'
                const optionMap = {
                    [EQ]: this.$t('等于'),
                    [NE]: this.$t('不等于'),
                    [IN]: this.$t('属于'),
                    [NIN]: this.$t('不属于'),
                    [LTE]: this.$t('小于等于'),
                    [GTE]: this.$t('大于等于'),
                    [RANGE]: this.$t('范围'), // 前端构造的操作符，真实数据中会拆分数据为gte, lte向后台传递
                    [REGEX]: 'Like'
                }
                const typeMap = {
                    bool: [EQ, NE],
                    date: [GTE, LTE],
                    enum: [IN, NIN],
                    float: [EQ, NE, RANGE],
                    int: [EQ, NE, RANGE],
                    list: [IN, NIN],
                    longchar: [IN, NIN],
                    objuser: [IN, NIN],
                    organization: [IN, NIN],
                    singlechar: [IN, NIN],
                    time: [GTE, LTE],
                    timezone: [IN, NIN],
                    'service-template': [IN]
                }
                return typeMap[this.type].map(operator => ({ id: operator, name: optionMap[operator] }))
            },
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        }
    }
</script>
