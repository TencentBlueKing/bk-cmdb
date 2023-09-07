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

import { unref } from 'vue'
import { t } from '@/i18n'
import { $bkInfo } from '@/magicbox/index.js'
import fieldTemplateService from '@/service/field-template'

export default function useTemplate(templates) {
  const handleDelete = (template, successCallback, errorCallback) => {
    const templateList = unref(templates)
    if (!templateList.some(item => item.id === template.id)) {
      console.error('invalid template!')
      return
    }
    $bkInfo({
      title: t('确认要删除', { name: template.name }),
      confirmLoading: true,
      confirmFn: async () => {
        try {
          await fieldTemplateService.deleteTemplate({
            data: {
              id: template.id
            }
          })
          successCallback?.()
        } catch (error) {
          console.error(error)
          errorCallback?.()
          return false
        }
      }
    })
  }

  return {
    handleDelete
  }
}
