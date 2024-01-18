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

export const MULTIPLE_IP_REGEXP = /^((1?\d{1,2}|2[0-4]\d|25[0-5])[.]){3}(1?\d{1,2}|2[0-4]\d|25[0-5])(,((1?\d{1,2}|2[0-4]\d|25[0-5])[.]){3}(1?\d{1,2}|2[0-4]\d|25[0-5]))*$/

// 所有可能是IP的
export const ALL_PROBABLY_IP = /([a-zA-Z0-9.:[]]*)+/g

export const LT_REGEXP = /</g

// 常规IPV4
export const IPV4_IP = /(?<!\d)(?<ip>((?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.]){3}(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5]))(?!\d)/g

// 常规IPV6
export const IPV6_IP = /(?<![0-9a-fA-F:])(?<ip>([0-9a-fA-F]{4}:){7}[0-9a-fA-F]{4})(?![0-9a-fA-F:])/g

// 管控区域加IPV4
export const AREA_IPV4_IP = /(?<![a-zA-Z0-9.])(?<ip>((\d+):((?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.]){3}(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5]))|((\d+):\[((?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.]){3}(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\]))(?!\d)/g

// 管控区域加IPV6
export const AREA_IPV6_IP = /(?<![a-zA-Z0-9:])(?<ip>(\d+):\[(([0-9a-fA-F]{4}:){7}[0-9a-fA-F]{4})\])(?![0-9a-fA-F:])/g

