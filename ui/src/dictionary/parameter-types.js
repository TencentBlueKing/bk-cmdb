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

/**
 * 后台定义的字段类型的 key 和对应的中文翻译。
 * 因为后端返回的 key 是英文，而前端 UI 中展示的是需要国际化的文案，
 * 所以需要先将 key 转换成中文，再拿中文作为 key 来进行国际化的转换，
 * 这样即使后台返回的 key 发生改变，国际化的 key 也不会失效，只要在这里将 key 更改即可。
 */
export const PARAMETER_TYPES = {
  number: '数字',
  float: '浮点',
  singlechar: '短字符',
  longchar: '长字符',
  associationId: '关联类型唯一标识',
  classifyId: '模型分组唯一标识',
  modelId: '模型唯一标识',
  fieldId: '模型字段唯一标识',
  enumId: '枚举 ID',
  enumName: '枚举名称',
  namedCharacter: '服务分类名称',
  instanceTagKey: '服务实例标签 Key',
  instanceTagValue: '服务实例标签 Value',
  businessTopoInstNames: '集群/模块/实例名称',
}
