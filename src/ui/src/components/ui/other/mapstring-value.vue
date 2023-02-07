<!-- eslint-disable no-unused-vars -->
<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script>
  import { computed, defineComponent, ref, nextTick, onMounted, onBeforeUnmount, watch } from 'vue'
  import throttle from 'lodash.throttle'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  import { $bkPopover } from '@/magicbox/index.js'

  export default defineComponent({
    name: 'cmdb-mapstring-value',
    props: {
      value: {
        type: [String, Object, Array],
        default: () => ({})
      }
    },
    setup(props) {
      const rootEl = ref(null)
      const listEl = ref(null)
      const ellipsisEl = ref(null)

      let tips = null

      const tags = computed(() => {
        if (!props.value) {
          return []
        }

        let list = props.value
        if (!Array.isArray(props.value)) {
          list = [props.value]
        }

        const labels = []
        list.filter(item => item).forEach((item) => {
          labels.push(...Object.keys(item).map(key => `${key}: ${item[key]}`))
        })
        return labels
      })

      const handleResize = () => {
        removeEllipsisTag()
        if (!tags.value.length) {
          return
        }
        nextTick(() => {
          const items = Array.from(listEl.value.querySelectorAll('.tag-item'))
          const referenceItemIndex = items.findIndex((item, index) => {
            if (index === 0) {
              return false
            }
            const previousItem = items[index - 1]
            return previousItem.offsetTop !== item.offsetTop
          })
          if (referenceItemIndex > -1) {
            insertEllipsisTag(items[referenceItemIndex], referenceItemIndex)
            doubleCheckEllipsisPosition()
          } else {
            removeEllipsisTag()
          }
        })
      }

      const insertEllipsisTag = (reference) => {
        listEl.value.insertBefore(ellipsisEl.value, reference)
      }

      const doubleCheckEllipsisPosition = () => {
        const ellipsis = ellipsisEl.value
        const previous = ellipsis.previousElementSibling
        if (previous && ellipsis.offsetTop !== previous.offsetTop) {
          listEl.value.insertBefore(ellipsis, previous)
        }
        setEllipsisTips()
      }

      const setEllipsisTips = () => {
        const ellipsis = ellipsisEl.value
        const tips = getTipsInstance()
        const tipsNode = listEl.value.cloneNode(false)
        let loopItem = ellipsis
        while (loopItem) {
          const nextItem = loopItem.nextElementSibling
          if (nextItem && nextItem.classList.contains('tag-item')) {
            tipsNode.appendChild(nextItem.cloneNode(true))
            loopItem = nextItem
          } else {
            loopItem = null
          }
        }
        tips.setContent(tipsNode)
      }

      const removeEllipsisTag = () => {
        try {
          listEl.value.removeChild(ellipsisEl.value)
        } catch (e) {}
      }

      const getTipsInstance = () => {
        if (!tips) {
          tips = $bkPopover(ellipsisEl.value, {
            allowHTML: true,
            placement: 'top',
            arrow: true,
            theme: 'light',
            interactive: true
          })
        }
        return tips
      }

      watch(tags, () => {
        handleResize()
      }, { immediate: true })

      const scheduleResize = throttle(handleResize, 300)

      onMounted(() => {
        addResizeListener(rootEl.value, scheduleResize)
      })

      onBeforeUnmount(() => {
        removeResizeListener(rootEl.value, scheduleResize)
      })

      const getCopyValue = () => tags.value.join('\n')

      return {
        tags,
        rootEl,
        listEl,
        ellipsisEl,
        getCopyValue
      }
    }
  })
</script>

<template>
  <div ref="rootEl">
    <ul class="tag-list" ref="listEl" v-if="tags.length">
      <li class="tag-item"
        v-for="(tag, index) in tags"
        :key="index"
        :title="tag">
        {{tag}}
      </li>
      <li class="tag-item ellipsis" ref="ellipsisEl" v-show="tags.length" @click.stop>...</li>
    </ul>
    <span class="tag-empty" v-else>--</span>
  </div>
</template>

<style lang="scss" scoped>
.tag-list {
  flex: 1;
  height: 22px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  overflow: hidden;
  font-size: 12px;
  float: none !important;
  .tag-item {
    display: inline-block;
    max-width: 80px;
    padding: 0 6px;
    border-radius: 2px;
    line-height: 22px;
    color: $textColor;
    background-color: #f0f1f5;
    cursor: default;
    @include ellipsis;
    & ~ .tag-item {
      margin-left: 6px;
    }
    &.ellipsis {
      width: 22px;
      height: 22px;
      text-align: center;
      & ~ .tag-item {
        display: none;
      }
    }
  }
}
</style>
