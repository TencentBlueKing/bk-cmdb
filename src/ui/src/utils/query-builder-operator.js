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

const mapping = {
  $eq: 'equal',
  $ne: 'not_equal',
  $in: 'in',
  $nin: 'not_in',
  $lt: 'less',
  $lte: 'less_or_equal',
  $gt: 'greater',
  $gte: 'greater_or_equal',
  $range: 'between',
  $nrange: 'not_between',
  $regex: 'contains'
}

export default operator => mapping[operator]
