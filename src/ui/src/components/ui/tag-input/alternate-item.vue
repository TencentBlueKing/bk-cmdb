<template>
  <li class="alternate-item"
    :class="{
      highlight: index === tagInput.highlightIndex,
      disabled: disabled && !tagInput.renderList
    }"
    @click.stop
    @mousedown.left.stop="tagInput.handleTagMousedown(tag, disabled)"
    @mouseup.left.stop="tagInput.handleTagMouseup(tag, disabled)">
    <template v-if="tagInput.renderList">
      <render-list
        :tag-input="tagInput"
        :keyword="keyword"
        :tag="tag"
        :disabled="disabled">
      </render-list>
    </template>
    <template v-else>
      <span class="item-name" :title="getTitle()" v-html="getItemContent()"></span>
    </template>
  </li>
</template>
<script>
  import RenderList from './render-list.js'
  export default {
    name: 'alternate-item',
    components: {
      RenderList
    },
    // eslint-disable-next-line
    props: ['tagInput', 'tag', 'keyword', 'index'],
    computed: {
      disabled() {
        return this.tagInput.disabledData.includes(this.tag.value)
      }
    },
    methods: {
      getItemContent() {
        let displayText = this.tag.text || this.tag.value
        if (this.keyword) {
          // eslint-disable-next-line no-underscore-dangle
          displayText = displayText.replace(new RegExp(this.keyword, 'ig'), `<span ${this.$options._scopeId}>$&</span>`)
        }
        return displayText
      },
      getTitle() {
        return this.tagInput.getDisplayText(this.tag)
      }
    }
  }
</script>

<style lang="scss" scoped>
.alternate-item {
    padding: 0 10px;
    justify-content: space-between;
    cursor: pointer;
    &.highlight,
    &:hover{
        background-color: #f1f7ff;
    }
    &.disabled {
        opacity: .5;
        cursor: not-allowed;
    }
    .item-name {
        display: block;
        @include ellipsis;
        span {
            color: #3A84FF;
        }
    }
}
</style>
