<template>
  <bk-select style="text-align: left;"
    v-model="localSelected"
    ext-popover-cls="business-mix-selector-popover"
    :searchable="true"
    :search-with-pinyin="true"
    :clearable="false"
    :placeholder="$t('请选择业务')"
    :disabled="disabled"
    :popover-options="popoverOptions"
    @toggle="handleSelectToggle">
    <bk-option v-for="(option, index) in sortedList"
      :key="index"
      :id="option.id"
      :name="option.name">
      <div class="option-item-content" :title="option.name">
        <div class="text">
          <span class="item-name">{{option.rawName}}</span>
          <span class="item-id">({{option.rawId}})</span>
        </div>
        <i class="icon icon-cc-business-set" v-if="option.isBizSet"></i>
        <i :class="['icon', 'bk-icon', 'collection', isCollected(option) ? 'icon-star-shape' : 'icon-star']"
          @click.prevent.stop="handleCollect(option)">
        </i>
      </div>
    </bk-option>
    <div class="business-extension" slot="extension" v-if="showApplyPermission || showApplyCreate">
      <template v-if="showApplyPermission">
        <a href="javascript:void(0)" class="extension-link"
          @click="handleApplyBizPermission">
          <i class="bk-icon icon-plus-circle"></i>
          {{$t('申请业务权限')}}
        </a>
        <a href="javascript:void(0)" class="extension-link"
          @click="handleApplyBizSetPermission">
          <i class="bk-icon icon-plus-circle"></i>
          {{$t('申请业务集权限')}}
        </a>
      </template>
      <template v-if="showApplyCreate">
        <a href="javascript:void(0)" class="extension-link"
          @click="handleApplyCreate">
          <i class="bk-icon icon-plus-circle"></i>
          {{$t('申请创建业务')}}
        </a>
        <a href="javascript:void(0)" class="extension-link"
          @click="handleApplyCreate">
          <i class="bk-icon icon-plus-circle"></i>
          {{$t('申请创建业务集')}}
        </a>
      </template>
    </div>
  </bk-select>
</template>

<script>
  import { mapGetters } from 'vuex'
  import businessSetService from '@/service/business-set/index.js'
  import applyPermission from '@/utils/apply-permission.js'
  import { BUSINESS_SELECTOR_COLLECTION } from '@/dictionary/menu-symbol'

  const MAX_COLLECT_COUNT = 8

  export default {
    name: 'cmdb-business-mix-selector',
    props: {
      value: {
        type: String
      },
      disabled: {
        type: Boolean,
        default: false
      },
      popoverOptions: {
        type: Object,
        default() {
          return {}
        }
      },
      showApplyPermission: Boolean,
      showApplyCreate: Boolean
    },
    data() {
      return {
        normalizationList: [],
        sortedList: [],
        requestIds: {
          collection: Symbol()
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId', 'authorizedBusiness']),
      ...mapGetters('userCustom', ['usercustom']),
      collection() {
        return this.usercustom[BUSINESS_SELECTOR_COLLECTION] || []
      },
      localSelected: {
        get() {
          return this.value
        },
        set(value) {
          const [id, type] = value.split('-')
          this.$emit('input', value)
          this.$emit('select', value, Number(id), type === 'bizset')
        }
      }
    },
    async created() {
      const list = await businessSetService.getAuthorizedWithCache()
      const allList = [...this.authorizedBusiness, ...list]
      const normalizationList = []
      allList.forEach((item) => {
        const isBizSet = Boolean(item.bk_biz_set_id)
        const rawId = isBizSet ? item.bk_biz_set_id : item.bk_biz_id
        const rawName = isBizSet ? item.bk_biz_set_name : item.bk_biz_name
        normalizationList.push({
          isBizSet,
          rawId,
          rawName,
          // id值加后缀标明类型
          id: isBizSet ? `${rawId}-bizset` : `${rawId}-biz`,
          name: `${rawName} (${rawId})`
        })
      })
      this.normalizationList = normalizationList
      this.sortedList = this.normalizationList
    },
    methods: {
      isCollected(option) {
        return this.collection.includes(option.id)
      },
      sortList(list) {
        return list.slice().sort((a, b) => {
          if (this.isCollected(a) > this.isCollected(b)) {
            return -1
          }
          return 0
        })
      },
      async handleCollect(option) {
        if (this.$loading(this.requestIds.collection)) {
          return
        }

        let newCollection = []
        const isAdd = !this.collection.some(item => item === option.id)

        if (isAdd && this.collection.length >= MAX_COLLECT_COUNT) {
          this.$warn(this.$t('限制收藏个数提示', { max: MAX_COLLECT_COUNT }))
          return false
        }

        if (isAdd) {
          newCollection = this.collection.concat(option.id)
        } else {
          newCollection = this.collection.filter(item => item !== option.id)
        }

        try {
          await this.$store.dispatch('userCustom/saveUsercustom', {
            [BUSINESS_SELECTOR_COLLECTION]: newCollection
          }, { requestId: this.requestIds.collection })
          this.$success(this.$t(isAdd ? '收藏成功' : '取消收藏成功'))
        } catch (err) {
          this.$error(this.$t(isAdd ? '收藏失败' : '取消收藏失败'))
        }
      },
      handleSelectToggle(isOpen) {
        // 每次下拉展开时重新排序数据，操作收藏时不排序防止顺序跳动
        if (isOpen) {
          this.sortedList = this.sortList(this.normalizationList)
        }
      },
      async handleApplyBizPermission() {
        try {
          await applyPermission({
            type: this.$OPERATION.R_BIZ_RESOURCE,
            relation: []
          })
        } catch (e) {
          console.error(e)
        }
      },
      async handleApplyBizSetPermission() {
        try {
          await applyPermission({
            type: this.$OPERATION.R_BIZ_SET_RESOURCE,
            relation: []
          })
        } catch (e) {
          console.error(e)
        }
      },
      handleApplyCreate() {}
    }
  }
</script>

<style lang="scss" scoped>
  .option-item-content {
    color: #63656E;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .text {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .icon {
      margin-left: 8px;
      &.collection:not(.icon-star-shape) {
        display: none;
      }

      &.icon-star-shape {
        color: #FFB400;
      }
    }

    .item-id {
      color: #C4C6CC;
    }

    &:hover {
      .icon {
        &.collection {
          display: block;
        }
      }
    }
  }
  .business-extension {
    width: calc(100% + 32px);
    margin-left: -16px;
  }
  .extension-link {
    display: block;
    line-height: 38px;
    background-color: #FAFBFD;
    padding: 0 9px;
    font-size: 13px;
    color: #63656E;
    &:hover {
      opacity: .85;
    }
    .bk-icon {
      font-size: 18px;
      color: #979BA5;
      vertical-align: text-top;
    }
  }
</style>

<style lang="scss">
    .bk-select-dropdown-content.business-selector-popover {
        .bk-option-content-default {
            padding: 0;
            .bk-option-name {
                width: 100%;
                overflow: hidden;
                text-overflow: ellipsis;
            }
        }
    }
</style>
