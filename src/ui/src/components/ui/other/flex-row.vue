<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script>
  import { defineComponent } from 'vue'

  export default defineComponent({
    props: {
      gap: {
        type: [String, Number],
        default: 0
      },
      colon: {
        type: Boolean,
        default: false
      },
      colonGap: {
        type: [String, Number],
        default: '14px'
      }
    },
    setup(props, { slots }) {
      return () => <div class={['flex-row', { colon: props.colon }]}>{ slots.default() }</div>
    }
  })
</script>

<style lang="scss" scoped>
  .flex-row {
    display: flex;
    align-items: center;
    gap: v-bind(gap);

    &.colon {
      :first-child {
        position: relative;
        padding-right: v-bind(colonGap);
        &::after {
          position: absolute;
          right: 0;
          content: "：";
        }
      }
    }
  }
</style>
