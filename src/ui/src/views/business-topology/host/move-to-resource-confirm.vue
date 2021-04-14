<template>
  <div :class="['confirm-wrapper', { 'has-invalid': hasInvalid }]">
    <h2 class="title">{{$t('确认归还主机池')}}</h2>
    <i18n tag="p" path="确认归还主机池忽略主机数量" class="content" v-if="hasInvalid">
      <span class="count" place="count">{{count}}</span>
      <span class="invalid" place="invalid">{{invalidList.length}}</span>
    </i18n>
    <i18n tag="p" path="确认归还主机池主机数量" class="content" v-else>
      <span class="count" place="count">{{count}}</span>
    </i18n>
    <p class="content">{{$t('归还主机池提示')}}</p>
    <div class="directory">
      {{$t('归还至目录')}}
      <bk-select class="directory-selector ml10"
        v-model="target"
        searchable
        :clearable="false"
        :loading="$loading(request.list)"
        :popover-options="{
          boundary: 'window'
        }">
        <cmdb-auth-option v-for="directory in directories"
          :key="directory.bk_module_id"
          :id="directory.bk_module_id"
          :name="directory.bk_module_name"
          :auth="{ type: $OPERATION.HOST_TO_RESOURCE, relation: [[[bizId], [directory.bk_module_id]]] }">
        </cmdb-auth-option>
      </bk-select>
    </div>
    <invalid-list :title="$t('以下主机不能移除')" :list="invalidList"></invalid-list>
    <div class="options">
      <bk-button class="mr10" theme="primary"
        :disabled="!target"
        @click="handleConfirm">{{$t('确定')}}</bk-button>
      <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import InvalidList from './invalid-list'
  export default {
    name: 'cmdb-move-to-resource-confirm',
    components: {
      InvalidList
    },
    props: {
      count: {
        type: Number,
        default: 0
      },
      bizId: Number,
      invalidList: {
        type: Array,
        default: () => ([])
      }
    },
    data() {
      return {
        target: '',
        directories: [],
        request: {
          list: Symbol('list')
        }
      }
    },
    computed: {
      hasInvalid() {
        return !!this.invalidList.length
      }
    },
    created() {
      this.getDirectories()
    },
    methods: {
      async getDirectories() {
        try {
          const { info } = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
            params: {
              page: {
                sort: 'bk_module_name'
              }
            },
            config: {
              requestId: this.request.list
            }
          })
          this.directories = info
        } catch (error) {
          console.error(error)
        }
      },
      handleConfirm() {
        this.$emit('confirm', this.target)
      },
      handleCancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .confirm-wrapper {
        text-align: center;
        &.has-invalid {
            .content {
                padding: 0 26px;
                text-align: left;
            }
            .directory {
                padding: 0 0 0 26px;
                justify-content: flex-start;
                .directory-selector {
                    width: 514px;
                }
            }
        }
    }
    .title {
        margin: 45px 0 17px;
        line-height: 32px;
        font-size:24px;
        font-weight: normal;
        color: #313238;
    }
    .content {
        line-height:20px;
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
    .directory {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        margin-top: 10px;
        .directory-selector {
            width: 305px;
            margin-left: 10px;
            text-align: left;
        }
    }
    .options {
        margin: 20px 0;
        font-size: 0;
    }
</style>
