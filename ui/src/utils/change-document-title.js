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

import store from '@/store'
import router from '@/router'
import { t } from '@/i18n'

/**
 * 更改文档标题
 * @param {Array} [appendTitles] 追加的标题，会展示在默认名称之后。不传入时会根据当前路由重新生成路径。
 */
export const changeDocumentTitle = (appendTitles = []) => {
  const { name, separator } = store.state.globalConfig.config.site
  const { matched } = router.currentRoute
  let matchedNames = [name]
  matched.forEach((match) => {
    if (match?.meta?.menu?.i18n) {
      matchedNames.push(t(match.meta.menu.i18n))
    }
  })

  if (appendTitles?.length) {
    matchedNames = matchedNames.concat(appendTitles)
  }

  document.title = matchedNames?.join(` ${separator} `)
}
