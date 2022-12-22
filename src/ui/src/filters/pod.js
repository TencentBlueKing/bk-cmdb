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

export default function (item, propertyId) {
  if (!propertyId) {
    return null
  }

  // 读取pod的ref字段中的name，暂无其它展示场景，如之后有更丰富的场景可考虑将ref的展示拆分为独立的组件
  if (propertyId === 'ref') {
    return item?.[propertyId]?.name
  }
  return item?.[propertyId]
}
