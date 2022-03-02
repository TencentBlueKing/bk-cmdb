<template>
  <bk-select style="text-align: left;"
    v-model="localSelected"
    ext-popover-cls="business-mix-selector-popover"
    :searchable="true"
    :clearable="false"
    :placeholder="$t('请选择业务')"
    :disabled="disabled"
    :popover-options="popoverOptions">
    <bk-option v-for="(option, index) in normalizationList"
      :key="index"
      :id="option.id"
      :name="option.name">
      <div class="option-item-content" :title="option.name">
        <span class="text">{{option.name}}</span>
        <i class="icon icon-cc-business-set" v-if="option.isBizSet"></i>
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
        authorizedBusinessSet: []
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId', 'authorizedBusiness']),
      localSelected: {
        get() {
          return this.value
        },
        set(value) {
          const [id, type] = value.split('-')
          this.$emit('input', value)
          this.$emit('select', value, Number(id), type === 'bizset')
        }
      },
      normalizationList() {
        const list = [...this.authorizedBusiness, ...this.authorizedBusinessSet]
        const normalizationList = []
        list.forEach((item) => {
          const isBizSet = Boolean(item.bk_biz_set_id)
          const rawId = isBizSet ? item.bk_biz_set_id : item.bk_biz_id
          const rawName = isBizSet ? item.bk_biz_set_name : item.bk_biz_name
          normalizationList.push({
            isBizSet,
            rawId,
            // id值加后缀标明类型
            id: isBizSet ? `${rawId}-bizset` : `${rawId}-biz`,
            name: `[${rawId}] ${rawName}`
          })
        })
        return normalizationList
      }
    },
    async created() {
      const list = await businessSetService.getAuthorizedWithCache()
      this.authorizedBusinessSet = Object.freeze(list)
    },
    methods: {
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
