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

<script setup>
  import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
  import { $bkPopover } from '@/magicbox/index.js'

  const props = defineProps({
    list: {
      type: Array,
      default: () => ([])
    },
    isLinkStyle: {
      type: Boolean,
      default: false
    },
    isShowLinkIcon: {
      type: Boolean,
      default: false
    },
    isTextStyle: {
      type: Boolean,
      default: false
    },
    maxWidth: {
      type: String,
      default: '400px'
    },
    fontSize: {
      type: String,
      default: '12px'
    },
    gap: {
      type: String,
      default: '4px'
    },
    height: {
      type: String,
      default: '22px'
    }
  })

  const emit = defineEmits(['click'])

  const containerEl = ref(null)
  const plusEl = ref(null)
  let tips = null

  const tags = computed(() => props.list.filter(item => item))
  const gapWidth = computed(() => parseInt(props.gap, 10))
  const ellipsisCount = ref(0)

  const handleClick = (index) => {
    emit('click', index)
  }

  const handleClickIcon = (text, index) => {
    emit('unbind', text, index)
  }

  const getTipsInstance = () => {
    if (!tips) {
      tips = $bkPopover(plusEl.value, {
        allowHTML: true,
        placement: 'top',
        boundary: 'window',
        arrow: true,
        theme: `${props.isLinkStyle ? 'light' : 'dark'} flex-tag-tooltip`,
        interactive: true
      })
    }
    return tips
  }

  const resizeHander = () => {
    if (!tags.value.length) {
      return
    }

    const containerClientWidth = containerEl.value.clientWidth
    const containerScrollWidth = containerEl.value.scrollWidth
    const plusWidth = Math.ceil(plusEl.value.getBoundingClientRect().width)
    const tagItems = Array.from(containerEl.value.querySelectorAll('.tag-item'))

    const tagWidthList = tagItems.map((item, index) => ({
      index,
      el: item,
      width: Math.ceil(item.getBoundingClientRect().width)
    }))

    if (containerScrollWidth > containerClientWidth) {
      // 默认将plus元素宽度计入，这将始终保持plus元素是可见的，但对是否overflow的精度有所降低
      let accWidth = plusWidth

      // 出现断点的那个元素
      let posItem = null

      for (const item of tagWidthList) {
        accWidth = accWidth + item.width + gapWidth.value + 12
        if (accWidth > containerClientWidth) {
          posItem = item
          ellipsisCount.value = tags.value.length - item.index
          break
        }
      }

      // 将plus元素放到断点元素前面并且显示
      if (posItem) {
        containerEl.value.insertBefore(plusEl.value, posItem.el)
        plusEl.value.classList.add('show')
        if (plusEl.value.previousSibling) {
          tagItems.forEach(item => item.classList.remove('is-pos'))
          plusEl.value.previousSibling.classList.add('is-pos')
        }
      }

      // 设置tooltips
      const tooltips = getTipsInstance()
      const contentEl = document.createElement('div')
      const fragment = document.createDocumentFragment()
      contentEl.classList.add('flex-tag-tips-content')
      contentEl.style.setProperty('--fontSize', props.fontSize)

      let tipTags = tags.value.slice(0)
      if (props.isLinkStyle || props.isTextStyle) {
        tipTags = tags.value.slice(tags.value.length - ellipsisCount.value)
      }
      tipTags.forEach((text, index) => {
        const itemEl = document.createElement('div')
        itemEl.classList.add('flex-tag-tips-item')
        if (props.isLinkStyle) {
          itemEl.classList.add('is-link')
          itemEl.addEventListener('click', () => handleClick(index), false)
        }
        itemEl.textContent = text
        fragment.appendChild(itemEl)
      })

      contentEl.appendChild(fragment)
      tooltips.setContent(contentEl)
    } else {
      // 将plus元素放到最后并且隐藏
      containerEl.value.insertBefore(plusEl.value, tagItems[tagItems.length - 1].nextSibling)
      tagItems.forEach(item => item.classList.remove('is-pos'))
      plusEl.value.classList.remove('show')
    }
  }

  const changing = ref(false)

  const execResizeHander = async () => {
    // changing为true时使得所有tagitem可见才能保证宽度的计算准确
    changing.value = true
    await nextTick(resizeHander)
    changing.value = false
  }

  const resizeObserver = new ResizeObserver((entries) => {
    for (const entry of entries) {
      if (entry.target === containerEl.value) {
        execResizeHander()
      }
    }
  })

  watch(tags, () => {
    execResizeHander()
  })

  onMounted(() => {
    resizeObserver.observe(containerEl.value)
    execResizeHander()
  })

  onBeforeUnmount(() => {
    resizeObserver.unobserve(containerEl.value)
  })
</script>
<script>
  export default {
    name: 'flex-tag'
  }
</script>

<template>
  <ul
    :class="['flex-tag', { changing, 'is-link-style': isLinkStyle, 'is-text-style': isTextStyle }]"
    ref="containerEl"
    :style="{
      '--fontSize': fontSize,
      '--gap': gap,
      '--maxWidth': maxWidth,
      '--height': height
    }">
    <li class="tag-item" v-bk-overflow-tips
      v-for="(tag, index) in tags"
      :key="index"
      @click="handleClick(index)">
      <span>{{tag.name || tag}}</span>
      <bk-icon type="chain" v-if="isShowLinkIcon" @click="handleClickIcon(tag,index)" v-bk-tooltips="'解绑模版'" />
    </li>
    <li class="more-plus" ref="plusEl" v-show="ellipsisCount">+{{ellipsisCount}}</li>
  </ul>
</template>

<style lang="scss" scoped>
  .flex-tag {
    display: flex;
    gap: var(--gap);
    font-size: var(--fontSize);

    .tag-item {
      color: #63656E;
      background: #F0F1F5;
      border-radius: 2px;
      padding: 0 .6em;
      white-space: nowrap;
      flex: none;
      overflow: hidden;
      text-overflow: ellipsis;
      height: var(--height);
      line-height: var(--height);
      max-width: var(--maxWidth);
      display: flex;
      align-items: center;
    }

    &.is-link-style {
      .tag-item {
        color: #3A84FF;
        background: none;
        border-radius: 0;
        padding: 0;
        cursor: pointer;

        &::after {
          content: ',';
          color: #63656E;
        }

        // 倒数第2个元素为tag-item的最后一个元素，倒数第1个元素为plus
        &:nth-last-of-type(2) {
          &::after {
            display: none;
          }
        }
        &.is-pos {
          &::after {
            content: '...';
          }
        }
      }
    }

    &.is-text-style {
      .tag-item {
        background: none;
        border-radius: 0;
        padding: 0;

        &::after {
          content: '|';
          color: #63656E;
        }

        &:nth-last-of-type(2) {
          &::after {
            display: none;
          }
        }
        &.is-pos {
          &::after {
            content: '...';
          }
        }
      }
    }

    .more-plus {
      display: none;
      align-items: center;
      justify-content: center;
      font-size: var(--fontSize);
      height: var(--height);
      white-space: nowrap;
      background: #F0F1F5;
      border-radius: 2px;
      padding: 0 .6em;
      cursor: pointer;

      & ~ .tag-item {
        display: none;
      }

      &.show {
        display: flex;
      }
    }

    &.changing {
      .more-plus {
        display: flex;
        & ~ .tag-item {
          display: initial;
          visibility: hidden;
        }
      }
    }
  }
</style>
<style lang="scss">
  .flex-tag-tooltip-theme {
    .flex-tag-tips-content {
      padding: .2em;
      font-size: var(--fontSize);
      display: flex;
      flex-direction: column;
      gap: 6px;

      .flex-tag-tips-item {
        &.is-link {
          color: #3A84FF;
          cursor: pointer;
        }
      }
    }
  }
</style>
