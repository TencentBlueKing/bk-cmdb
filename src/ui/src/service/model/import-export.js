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

import http from '@/api'

const importService = async (objId, event) => {
  try {
    const [file] = event.target.files
    const form = new FormData()
    form.append('file', file)
    return http.post(`${window.API_HOST}object/object/${objId}/import`, form)
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const exportService = objId => http.download({ url: `${window.API_HOST}object/object/${objId}/export`, data: {} })

/**
 * 批量导入模型文件解析
 * @param {BinaryData} file 导入的模型文件
 * @param {String} [password] 加密文件密码
 * @returns {Promise}
 */
export const batchImportFileAnalysis = ({
  file,
  password,
}) => {
  const importData = new FormData()

  importData.append('file', file)

  if (password) {
    importData.append('params', JSON.stringify({ password }))
  }

  return http.post(`${window.API_HOST}object/importmany/analysis`, importData, {
    globalError: false,
    transformData: false
  })
    .then((res) => {
      if (res.bk_error_code === 0) {
        return Promise.resolve(res.data)
      }

      return Promise.reject(res)
    })
}

/**
 * 批量导入模型
 * @param {Object} import_object 需要导入的模型全量数据
 * @param {Object} import_asst 需要导入的关联关系类型的全量数据
 * @returns {Promise}
 */
export const batchImport = ({
  import_object,
  import_asst
}) => http.post(`${window.API_HOST}object/importmany`, {
  import_object,
  import_asst
})

/**
 * 批量导出模型
 * @param {Array} object_id 需要导出的模型对应的 id 列表，注意不是 bk_obj_id 而是 id
 * @param {string} file_name 导出文件名
 * @param {Array} [excluded_asst_id] 不需要导出的模型关系的关联 bk_asst_id 列表
 * @param {string} [password] 文件密码
 * @param {number} [expiration] 文件有效期，单位为天，无限期为 0
 * @returns {Promise}
 */
export const batchExport = ({
  object_id,
  excluded_asst_id,
  password,
  expiration,
  file_name
}) => http.download({
  url: `${window.API_HOST}object/exportmany`,
  data: {
    object_id,
    excluded_asst_id,
    password,
    expiration,
    file_name
  }
})

export default {
  import: importService,
  export: exportService,
  batchExport,
  batchImportFileAnalysis,
  batchImport
}
