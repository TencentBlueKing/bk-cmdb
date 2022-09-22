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

import { ref, watch, reactive } from '@vue/composition-api'

export const categories = ref([])

export const sizes = reactive({
  expanded: false,
  showMore: true,
  rows: 1,
  itemMarginRight: 20,
  itemMarginBottom: 12,
  itemHeight: 0,
  containerHeight: 0
})

export default function useTab(aggregations, root) {
  const getModelById = root.$store.getters['objectModelClassify/getModelById']

  const calculateSizes = () => {
    const $categories = document.querySelector('.categories')
    const $toggleAnchor = $categories.querySelector('.toggle-anchor')
    const $categoryItems = $categories.querySelectorAll('.categories .category-item')
    const containerWidth = $categories.offsetWidth
    const itemHeight = $categoryItems[0].offsetHeight + sizes.itemMarginBottom
    const anchorWidth = $toggleAnchor.offsetWidth

    const maxlength = $categoryItems.length
    let oneRowMaxWidth = 0
    let maxWidthIndex = maxlength

    const breakLineIndex = []
    for (let i = 0; i < maxlength; i++) {
      const nextItemWidth = Math.ceil($categoryItems[i].offsetWidth) + sizes.itemMarginRight
      oneRowMaxWidth += nextItemWidth
      if (oneRowMaxWidth > containerWidth) {
        breakLineIndex.push([i, oneRowMaxWidth])
        oneRowMaxWidth = nextItemWidth
      }
    }

    // 取第1行的值
    maxWidthIndex = breakLineIndex.length ? breakLineIndex[0][0] : maxWidthIndex

    // 是否需要展开
    sizes.showMore = maxWidthIndex !== maxlength

    sizes.itemHeight = itemHeight

    // 放置展开切换标签
    if (!sizes.expanded && breakLineIndex.length) {
      // eslint-disable-next-line max-len
      if (breakLineIndex[0][1] - Math.ceil($categoryItems[maxWidthIndex].offsetWidth) + anchorWidth - sizes.itemMarginRight > containerWidth) {
        maxWidthIndex -= 1
      }
      // 将显示更多放置在一行的末尾
      $categoryItems[maxWidthIndex - 1].after($toggleAnchor)
    } else {
      $categoryItems[$categoryItems.length - 1].after($toggleAnchor)
    }

    // 得出总行数
    let lastRowMaxWidth = 0
    if (breakLineIndex.length > 0) {
      for (let i = breakLineIndex[breakLineIndex.length - 1][0]; i < maxlength; i++) {
        const nextItemWidth = Math.ceil($categoryItems[i].offsetWidth) + sizes.itemMarginRight
        lastRowMaxWidth += nextItemWidth
      }
    }
    sizes.rows = breakLineIndex.length + 1 + (lastRowMaxWidth + anchorWidth > containerWidth ? 1 : 0)
  }

  watch(aggregations, (list) => {
    const models = []
    const instances = []

    list.forEach(({ key, kind, count }) => {
      const { bk_obj_id: id, bk_obj_name: name } = getModelById(key) || {}
      const item = { id, name, kind, count }
      if (item.kind === 'model') {
        models.push(item)
      } else if (item.kind === 'instance') {
        instances.push({ ...item, count: item.count > 999 ? '999+' : item.count, total: item.count })
      }
    })

    categories.value = [...instances]

    // 将所有模型聚合为一个“模型”分类
    const modelCount = models.reduce((acc, cur) => acc + cur.count, 0)
    if (models.length && modelCount > 0) {
      categories.value.unshift({
        id: models.map(({ id }) => id).join(','),
        name: root.$t('模型'),
        kind: 'model',
        count: modelCount > 999 ? '999+' : modelCount,
        total: modelCount
      })
    }
  })

  return {
    categories,
    calculateSizes
  }
}
