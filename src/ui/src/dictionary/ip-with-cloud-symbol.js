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

import { t } from '@/i18n'

export const IPWithCloudSymbol = '#IPWithCloud'
export const IPv6WithCloudSymbol = '#IPv6WithCloud'
export const IPv46WithCloudSymbol = '#IPv46WithCloud'
export const IPv64WithCloudSymbol = '#IPv64WithCloud'

export const IPWithCloudFields = {
  [IPWithCloudSymbol]: `${t('管控区域')}ID:IPv4`,
  [IPv6WithCloudSymbol]: `${t('管控区域')}ID:IPv6`,
  [IPv46WithCloudSymbol]: `${t('管控区域')}ID:IP(${t('IPv4优先')})`,
  [IPv64WithCloudSymbol]: `${t('管控区域')}ID:IP(${t('IPv6优先')})`
}
