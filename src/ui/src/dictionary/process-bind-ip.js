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

import { t } from '@/i18n/index.js'

export const PROCESS_BIND_IPV4_LOOPBACK = {
  1: '127.0.0.1'
}
export const PROCESS_BIND_IPV4_NON = {
  2: '0.0.0.0'
}

export const PROCESS_BIND_IPV6_LOOPBACK = {
  5: '::1'
}
export const PROCESS_BIND_IPV6_NON = {
  6: '::'
}

export const PROCESS_BIND_IPV4_LOOPBACK_AND_NON_MAP = {
  ...PROCESS_BIND_IPV4_LOOPBACK,
  ...PROCESS_BIND_IPV4_NON,
}
export const PROCESS_BIND_IPV6_LOOPBACK_AND_NON_MAP = {
  ...PROCESS_BIND_IPV6_LOOPBACK,
  ...PROCESS_BIND_IPV6_NON,
}

export const PROCESS_BIND_IPV4_MAP = {
  ...PROCESS_BIND_IPV4_LOOPBACK_AND_NON_MAP,
  3: t('第一内网IP'),
  4: t('第一外网IP')
}

export const PROCESS_BIND_IPV6_MAP = {
  ...PROCESS_BIND_IPV6_LOOPBACK_AND_NON_MAP,
  7: t('第一内网IPv6'),
  8: t('第一外网IPv6')
}

export const PROCESS_BIND_IP_ALL_LOOPBACK_AND_NON_MAP = {
  ...PROCESS_BIND_IPV4_LOOPBACK_AND_NON_MAP,
  ...PROCESS_BIND_IPV6_LOOPBACK_AND_NON_MAP,
}

export const PROCESS_BIND_IP_ALL_MAP = {
  ...PROCESS_BIND_IPV4_MAP,
  ...PROCESS_BIND_IPV6_MAP
}
