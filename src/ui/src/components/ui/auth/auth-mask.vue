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
    name: 'cmdb-auth-mask',
    props: {
      auth: [Object, Array],
      authorized: Boolean,
      tag: {
        type: String,
        default: 'span'
      },
      ignore: Boolean,
      callbackUrl: String
    },
    data() {
      return {
        useIAM: this.$Site.authscheme === 'iam'
      }
    },
    render(h) {
      if (!this.useIAM || this.ignore) {
        return this.$scopedSlots.default({
          disabled: false
        })
      }

      return h(this.tag, {
        directives: [
          {
            name: 'cursor',
            value: {
              auth: this.auth,
              active: !this.authorized,
              callbackUrl: this.callbackUrl
            }
          }
        ],
        staticClass: 'auth-mask',
      }, [
        this.$scopedSlots.default({
          disabled: !this.authorized
        })
      ])
    }
  })
</script>
