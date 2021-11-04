<template>
  <bk-dropdown-menu trigger="click" :disabled="disabled" font-size="medium">
    <bk-button class="clipboard-trigger" theme="default" slot="dropdown-trigger" v-test-id="'copy'"
      :disabled="disabled">
      {{$t('复制')}}
      <i class="bk-icon icon-angle-down"></i>
    </bk-button>
    <ul class="clipboard-list" slot="dropdown-content" v-test-id="'copy'">
      <li v-for="(item, index) in list"
        class="clipboard-item"
        :key="index"
        @click="handleClick(item)">
        {{item[labelKey]}}
      </li>
    </ul>
  </bk-dropdown-menu>
</template>

<script>
  export default {
    name: 'cmdb-clipboard-selector',
    props: {
      disabled: {
        type: Boolean,
        default: false
      },
      list: {
        type: Array,
        default() {
          return []
        }
      },
      idKey: {
        type: String,
        default: 'id'
      },
      labelKey: {
        type: String,
        default: 'name'
      }
    },
    methods: {
      handleClick(item) {
        this.$emit('on-copy', item)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .clipboard-trigger{
        padding: 0 16px;
        .icon-angle-down {
            font-size: 20px;
            margin: 0 -4px;
        }
    }
    .clipboard-list{
        width: 100%;
        font-size: 14px;
        line-height: 32px;
        // 漏出半个 item，引导用户下拉
        max-height: calc(160px + (32px / 2));
        @include scrollbar-y;
        &::-webkit-scrollbar{
            width: 3px;
            height: 3px;
        }
        .clipboard-item{
            padding: 0 15px;
            cursor: pointer;
            @include ellipsis;
            &:hover{
                background-color: #ebf4ff;
                color: #3c96ff;
            }
        }
    }
</style>
