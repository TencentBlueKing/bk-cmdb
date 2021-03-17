<template>
  <section class="across-confirm">
    <h1 class="title">{{$t('转移主机到其他业务')}}</h1>
    <i18n tag="p" path="确认跨业务转移忽略主机数量" class="content">
      <span class="count" place="count">{{count}}</span>
      <span class="invalid" place="invalid">{{invalidList.length}}</span>
    </i18n>
    <invalid-list :title="$t('以下主机不能移除')" :list="invalidList"></invalid-list>
    <div class="footer">
      <bk-button theme="primary" @click="next">{{$t('下一步')}}</bk-button>
      <bk-button class="ml10" theme="default" @click="cancel">{{$t('取消')}}</bk-button>
    </div>
  </section>
</template>

<script>
  import InvalidList from './invalid-list'
  export default {
    name: 'across-business-confirm',
    components: {
      InvalidList
    },
    props: {
      count: {
        type: Number,
        default: 0
      },
      invalidList: {
        type: Array,
        default: () => ([])
      }
    },
    methods: {
      next() {
        this.$emit('confirm')
      },
      cancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .across-confirm {
        .title {
            text-align: center;
            margin: 45px 0 17px;
            line-height: 32px;
            font-size:24px;
            font-weight: normal;
            color: #313238;
        }
        .content {
            padding: 0 20px;
            line-height: 20px;
            font-size:14px;
            color: $textColor;
            .count {
                font-weight: bold;
                color: $successColor;
                padding: 0 4px;
            }
            .invalid {
                font-weight: bold;
                color: $dangerColor;
                padding: 0 4px;
            }
        }
        .footer {
            display: flex;
            margin: 20px 0 0 0;
            align-items: center;
            justify-content: flex-end;
            height: 50px;
            padding: 8px 20px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
        }
    }
</style>
