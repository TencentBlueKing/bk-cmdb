
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

export const LANG_COOKIE_NAME = 'blueking_language'

export const LANG_KEYS = Object.freeze({
  ZH_CN: 'zh_CN',
  EN: 'en'
})

export const LANG_SET = Object.freeze([
  {
    id: LANG_KEYS.EN,
    name: 'English',
    icon: 'english'
  },
  {
    id: LANG_KEYS.ZH_CN,
    alias: ['zh-cn'],
    name: '中文',
    icon: 'chinese'
  }
])
