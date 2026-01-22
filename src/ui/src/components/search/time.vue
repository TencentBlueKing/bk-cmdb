<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <DatePicker
    v-bind="$attrs"
    :enable-format-click="false"
    :model-value="localValue"
    :timezone="timezone"
    @update:modelValue="handleChange"
    @update:timezone="handleChangeTimezone" />
</template>

<script>
  import DatePicker from '@blueking/date-picker/vue2'
  import '@blueking/date-picker/vue2/vue2.css'
  import activeMixin from './mixins/active'
  import { timestampFormatter } from '@/filters/formatter'
  export default {
    name: 'cmdb-search-time',
    components: {
      DatePicker
    },
    mixins: [activeMixin],
    props: {
      // 需为时间戳格式，如果传日期格式，组件内部会在根据时区进行一次转换
      value: {
        type: Array,
        default: () => ([])
      },
      timezone: {
        type: String,
        required: true,
      }
    },
    computed: {
      localValue: {
        get() {
          // 需要转换成时间戳格式
          return [...this.value.map(value => timestampFormatter(value))]
        },
        set(values) {
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    },
    methods: {
      handleChange(timestamp, date) {
        const [{ formatText: startDate }, { formatText: endDate }] = date
        // 将日期转换为时间戳格式
        const newDate = [startDate, endDate].filter(value => !!value)
          .map(value => timestampFormatter(value, this.timezone))
        if (newDate.toString() === this.value.toString()) return
        this.localValue = newDate
      },
      handleChangeTimezone(value) {
        this.$emit('change-timezone', value)
      },
    }
  }
</script>
<style lang="scss">
  .__bk_date_picker__ {
    width: 100%;
    .date-content {
      flex: 1;
    }
}
</style>
