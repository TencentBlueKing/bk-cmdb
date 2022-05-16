/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { mapState } from 'vuex'
export default {
  computed: {
    ...mapState('hostDetails', ['info']),
    HOST_AUTH() {
      if (!this.info) {
        return { U_HOST: null, D_SERVICE_INSTANCE: null }
      }
      const { biz, module, host } = this.info
      // 已分配主机
      if (biz[0].default === 0) {
        const bizId = biz[0].bk_biz_id
        return {
          U_HOST: {
            type: this.$OPERATION.U_HOST,
            relation: [bizId, host.bk_host_id]
          },
          D_SERVICE_INSTANCE: {
            type: this.$OPERATION.D_SERVICE_INSTANCE,
            relation: [bizId]
          },
          U_SERVICE_INSTANCE: {
            type: this.$OPERATION.U_SERVICE_INSTANCE,
            relation: [bizId]
          }
        }
      }
      return {
        U_HOST: {
          type: this.$OPERATION.U_RESOURCE_HOST,
          relation: [module[0].bk_module_id, host.bk_host_id]
        },
        D_SERVICE_INSTANCE: null,
        U_SERVICE_INSTANCE: null
      }
    }
  }
}
