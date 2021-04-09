<template>
  <bk-select
    v-model="localValue"
    v-bind="$attrs"
    v-on="listeners"
    :clearable="false">
    <bk-option v-for="(option, index) in options"
      class="operator-option"
      :key="index"
      :id="option.id"
      :name="option.name">
      {{option.name}}
      <span class="operator-description">({{option.title}})</span>
    </bk-option>
  </bk-select>
</template>

<script>
  import Utils from './utils'
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
      listeners() {
        const internalEvent = ['input', 'change']
        const listeners = {}
        Object.keys(this.$listeners).forEach((key) => {
          if (!internalEvent.includes(key)) {
            listeners[key] = this.$listeners[key]
          }
        })
        return listeners
      },
      options() {
        const EQ = '$eq'
        const NE = '$ne'
        const IN = '$in'
        const NIN = '$nin'
        const LT = '$lt'
        const GT = '$gt'
        const LTE = '$lte'
        const GTE = '$gte'
        const RANGE = '$range' // 前端构造的操作符，真实数据中会拆分数据为gte, lte向后台传递
        const typeMap = {
          bool: [EQ, NE],
          date: [GTE, LTE],
          enum: [IN, NIN],
          float: [EQ, NE, GT, LT, RANGE],
          int: [EQ, NE, GT, LT, RANGE],
          list: [IN, NIN],
          longchar: [IN, NIN],
          objuser: [IN, NIN],
          organization: [IN, NIN],
          singlechar: [IN, NIN],
          time: [GTE, LTE],
          timezone: [IN, NIN],
          foreignkey: [IN, NIN],
          table: [IN, NIN],
          'service-template': [IN]
        }
        const nameDescription = {
          [EQ]: this.$t('等于'),
          [NE]: this.$t('不等于'),
          [LT]: this.$t('小于'),
          [GT]: this.$t('大于'),
          [IN]: this.$t('包含'),
          [NIN]: this.$t('不包含'),
          [RANGE]: this.$t('数值范围'),
          [LTE]: this.$t('小于等于'),
          [GTE]: this.$t('大于等于')
        }
        return typeMap[this.type].map(operator => ({
          id: operator,
          name: Utils.getOperatorSymbol(operator),
          title: nameDescription[operator]
        }))
      },
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .operator-option {
        &:hover {
            .operator-description {
                display: initial;
            }
        }
        .operator-description {
            display: none;
        }
    }
</style>
