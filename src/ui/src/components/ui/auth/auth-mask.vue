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
      auth: Object,
      authorized: Boolean,
      tag: {
        type: String,
        default: 'span'
      }
    },
    render(h) {
      return h(this.tag, {
        directives: [
          {
            name: 'cursor',
            value: {
              auth: this.auth,
              active: !this.authorized
            }
          }
        ],
        staticClass: 'auth-mask'
      }, this.$slots.default)
    }
  })
</script>
