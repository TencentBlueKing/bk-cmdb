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
import { useHttp } from '@/api'

export default function useSearchUser() {
  const http = useHttp()

  const search = async (val) => {
    let result = []
    if (window.Site.userManageUrl) {
      const api = `${window.Site.userManageUrl}/api/v3/open-web/tenant/users/-/search/?keyword=${val}`
      const data = await http.get(api, { globalHeaders: false, globalError: false, headers: { 'X-Bk-Tenant-Id': window.Site.tenantId } })
      if (!data) {
        console.error('fetch user failed')
        return []
      }
      result = data.map(item => ({
        id: item.login_name,
        username: item.bk_username,
        name: item.display_name
      }))
    } else {
      const data = await http.get(`${window.API_HOST}user/list`, {
        params: {
          fuzzy_lookups: val
        },
        config: {
          cancelPrevious: true
        }
      })
      result = (data || []).map(user => ({
        id: user.english_name,
        username: user.english_name,
        name: user.chinese_name
      }))
    }

    return result
  }

  const lookup = async (values) => {
    const api = `${window.Site.userManageUrl}/api/v3/open-web/tenant/users/-/lookup/?lookups=${values}&lookup_fields=bk_username`
    const data = await http.get(api, { globalHeaders: false, globalError: false, headers: { 'X-Bk-Tenant-Id': window.Site.tenantId } })
    return data
  }

  return {
    search,
    lookup
  }
}
