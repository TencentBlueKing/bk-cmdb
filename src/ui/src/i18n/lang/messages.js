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

// 通用
import en from './en.json'
import cn from './cn.json'

// 全局配置
import globalConfigZhCN from '@/views/global-config/i18n/zh-CN.json'
import globalConfigEn from '@/views/global-config/i18n/en.json'

// 模型管理
import modelManageZhCN from '@/views/model-manage/i18n/zh-CN.json'
import modelManageEn from '@/views/model-manage/i18n/en.json'

export default {
  en: {
    ...en,
    ...globalConfigEn,
    ...modelManageEn
  },
  zh_CN: {
    ...cn,
    ...globalConfigZhCN,
    ...modelManageZhCN
  }
}
