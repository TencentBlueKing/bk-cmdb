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

import Form from '@/components/service/form/form.js'
import Bus from '../common/bus'
export default {
  methods: {
    handleAddProcess() {
      Form.show({
        type: 'create',
        title: `${this.$t('添加进程')}(${this.row.name})`,
        hostId: this.row.bk_host_id,
        bizId: this.bizId,
        submitHandler: this.createSubmitHandler
      })
    },
    async createSubmitHandler(values) {
      try {
        await this.$store.dispatch('processInstance/createServiceInstanceProcess', {
          params: {
            bk_biz_id: this.bizId,
            service_instance_id: this.row.id,
            processes: [{
              process_info: values
            }]
          }
        })
        this.$emit('refresh-count', this.row, this.row.process_count + 1)
        Bus.$emit('refresh-process-list', this.row)
      } catch (error) {
        console.error(error)
      }
    }
  }
}
