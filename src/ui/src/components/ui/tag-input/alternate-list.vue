<template>
  <div class="cmdb-tag-input-alternate-list-wrapper"
    :class="{
      'is-loading': loading
    }"
    :style="wrapperStyle">
    <ul class="alternate-list" ref="alternateList"
      :style="listStyle"
      @scroll="handleScroll">
      <template v-for="(tag, index) in matchedData">
        <template v-if="tag.hasOwnProperty('children')">
          <!-- eslint-disable vue/require-v-for-key -->
          <li class="alternate-group"
            @click.stop
            @mousedown.left.stop="tagInput.handleGroupMousedown"
            @mouseup.left.stop="tagInput.handleGroupMouseup">
            {{`${tag.value || tag.text}(${tag.children.length})`}}
          </li>
          <!-- eslint-disable vue/valid-v-for -->
          <alternate-item v-for="(child, childIndex) in tag.children"
            ref="alternateItem"
            :index="getIndex(index, childIndex)"
            :tag-input="tagInput"
            :tag="child"
            :keyword="keyword">
          </alternate-item>
        </template>
        <!-- eslint-disable vue/valid-v-for -->
        <alternate-item v-else
          ref="alternateItem"
          :tag-input="tagInput"
          :tag="tag"
          :index="getIndex(index)"
          :keyword="keyword">
        </alternate-item>
      </template>
    </ul>
    <p class="alternate-empty" v-if="!loading && !matchedData.length">{{tagInput.emptyText}}</p>
  </div>
</template>

<script>
  import AlternateItem from './alternate-item.vue'
  import has from 'has'
  export default {
    components: {
      AlternateItem
    },
    data() {
      return {
        tagInput: null,
        keyword: '',
        next: true,
        loading: true,
        matchedData: []
      }
    },
    computed: {
      wrapperStyle() {
        const style = {}
        if (this.tagInput && this.tagInput.panelWidth) {
          style.width = `${parseInt(this.tagInput.panelWidth, 10)}px`
        }
        return style
      },
      listStyle() {
        const style = {
          'max-height': '192px'
        }
        if (this.tagInput) {
          const maxHeight = parseInt(this.tagInput.listScrollHeight, 10)
          if (!isNaN(maxHeight)) {
            style['max-height'] = `${maxHeight}px`
          }
        }
        return style
      }
    },
    watch: {
      keyword() {
        this.$nextTick(() => {
          this.$refs.alternateList.scrollTop = 0
        })
      }
    },
    methods: {
      getIndex(index, childIndex = 0) {
        let flattenedIndex = 0
        this.matchedData.slice(0, index).forEach((tag) => {
          if (has(tag, 'children')) {
            flattenedIndex += tag.children.length
          } else {
            flattenedIndex += 1
          }
        })
        return flattenedIndex + childIndex
      },
      handleScroll() {
        if (this.loading || !this.next) return
        const list = this.$refs.alternateList
        const threshold = 32 // 距离底部2条数据
        if ((list.scrollTop + list.clientHeight) > (list.scrollHeight - threshold)) {
          this.tagInput.search(this.keyword, this.next)
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
.cmdb-tag-input-alternate-list-wrapper {
    width: 190px;
    color: #63656E;
    position: relative;
    background-color: #fff;
    &.is-loading {
        min-height: 32px;
        &:before {
            content: '';
            position: absolute;
            width: 100%;
            height: 100%;
            background-color: rgba(255, 255, 255, .7);
            z-index: 1;
        }
        &:after {
            content: '';
            position: absolute;
            top: 50%;
            left: 50%;
            width: 6px;
            height: 6px;
            margin-left: -30px;
            border-radius: 50%;
            // transform: translate3d(-50%, -50%, 0);
            background-color: transparent;
            box-shadow: 12px 0px 0px 0px #fd6154,
                        24px 0px 0px 0px #ffb726,
                        36px 0px 0px 0px #4cd084,
                        48px 0px 0px 0px #57a3f1;
            animation: tag-input-loading 1s linear infinite;
        }
    }
    .alternate-list{
        margin: 0;
        padding: 0;
        max-height: 162px;
        font-size: 12px;
        line-height: 32px;
        background: #fff;
        overflow-y: auto;
        &::-webkit-scrollbar {
            width: 4px;
            height: 4px;
            &-thumb {
                border-radius: 2px;
                background: #C4C6CC;
                box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
            }
        }
    }
    .alternate-group {
        padding: 0 11px;
        color: #979BA5;
        @include ellipsis;
    }
    .alternate-empty {
        padding: 0;
        margin: 0;
        text-align: center;
        line-height: 44px;
        font-size: 12px;
    }
}
@keyframes tag-input-loading {
   0%{
        box-shadow: 12px 0px 0px 0px #fd6154,
                    24px 0px 0px 0px #ffb726,
                    36px 0px 0px 0px #4cd084,
                    48px 0px 0px 0px #57a3f1;
      }
    14%{
        box-shadow: 12px 0px 0px 1px #fd6154,
                    24px 0px 0px 0px #ffb726,
                    36px 0px 0px 0px #4cd084,
                    48px 0px 0px 0px #57a3f1;
    }
    28% {
        box-shadow: 12px 0px 0px 2px #fd6154,
                    24px 0px 0px 1px #ffb726,
                    36px 0px 0px 0px #4cd084,
                    48px 0px 0px 0px #57a3f1;
    }
    42% {
        box-shadow: 12px 0px 0px 1px #fd6154,
                    24px 0px 0px 2px #ffb726,
                    36px 0px 0px 1px #4cd084,
                    48px 0px 0px 0px #57a3f1;
    }
    56%{
        box-shadow: 12px 0px 0px 0px #fd6154,
                    24px 0px 0px 1px #ffb726,
                    36px 0px 0px 2px #4cd084,
                    48px 0px 0px 1px #57a3f1;
    }
    70% {
        box-shadow: 12px 0px 0px 0px #fd6154,
                    24px 0px 0px 0px #ffb726,
                    36px 0px 0px 1px #4cd084,
                    48px 0px 0px 2px #57a3f1;
    }
    84% {
        box-shadow: 12px 0px 0px 0px #fd6154,
                    24px 0px 0px 0px #ffb726,
                    36px 0px 0px 0px #4cd084,
                    48px 0px 0px 1px #57a3f1;
    }
}
</style>
