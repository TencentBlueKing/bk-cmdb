<template>
  <span class="filter-tag" @click="handleClick">
    <label class="tag-name">{{property.bk_property_name}}</label>
    <span class="tag-colon" v-if="showColon">:</span>
    <component
      v-if="['foreignkey', 'service-template'].includes(property.bk_property_type)"
      :is="`cmdb-search-${property.bk_property_type}`"
      v-bk-overflow-tips="tipsConfig"
      display-type="info"
      :value="value">
      <template slot="info-prepend">{{operatorSymbol}}</template>
    </component>
    <span class="tag-value" v-else
      v-bk-overflow-tips="tipsConfig">
      {{displayText}}
    </span>
    <i class="tag-delete bk-icon icon-close" @mouseenter.prevent.stop @click.stop="handleRemove"></i>
  </span>
</template>

<script>
  import FilterStore from './store'
  import formatter from '@/filters/formatter'
  import FilterTagForm from './filter-tag-form'
  import Vue from 'vue'
  import i18n from '@/i18n'
  import store from '@/store'
  import Tippy from 'bk-magic-vue/lib/utils/tippy'
  import Utils from './utils'
  export default {
    props: {
      property: {
        type: Object,
        default: () => ({})
      },
      operator: {
        type: String,
        default: '$eq'
      },
      value: {
        type: [String, Array, Number, Boolean],
        default: ''
      }
    },
    data() {
      return {
        tipsConfig: {
          triggerTarget: null,
          interactive: false,
          hideOnClick: false,
          allowHTML: true
        }
      }
    },
    computed: {
      transformedValue() {
        let { value } = this
        if (!Array.isArray(value)) {
          value = [value]
        }
        return value.map(value => formatter(value, this.property))
      },
      showColon() {
        return this.operator === '$range'
      },
      operatorSymbol() {
        return Utils.getOperatorSymbol(this.operator)
      },
      displayText() {
        if (this.operator === '$range') {
          const [start, end] = this.transformedValue
          return `${start} ~ ${end}`
        }
        return `${this.operatorSymbol} ${this.transformedValue.join(' | ')}`
      }
    },
    mounted() {
      this.tipsConfig.triggerTarget = this.$el
    },
    beforeDestroy() {
      this.tagFormInstance && this.tagFormInstance.destroy()
      this.tagFormViewModel && this.tagFormViewModel.$destroy()
    },
    methods: {
      handleClick() {
        if (this.tagFormInstance) {
          this.tagFormInstance.show()
        } else {
          const self = this
          this.tagFormViewModel = new Vue({
            i18n,
            store,
            render(h) {
              return h(FilterTagForm, {
                ref: 'filterTagForm',
                props: {
                  property: self.property
                },
                on: {
                  confirm: self.handleHideTagForm,
                  cancel: self.handleHideTagForm
                }
              })
            }
          })
          this.tagFormViewModel.$mount()
          this.tagFormInstance = this.$bkPopover(this.$el, {
            content: this.tagFormViewModel.$el,
            theme: 'light',
            allowHTML: true,
            placement: 'bottom',
            trigger: 'manual',
            interactive: true,
            arrow: true,
            zIndex: window.__bk_zIndex_manager.nextZIndex(), // eslint-disable-line no-underscore-dangle
            onHide: () => !this.tagFormViewModel.$refs.filterTagForm.active,
            onHidden: () => {
              this.tagFormViewModel.$refs.filterTagForm.resetCondition()
            }
          })
          this.tagFormInstance.show()
        }
        Tippy.hideAll({ exclude: this.tagFormInstance })
      },
      handleHideTagForm() {
        this.tagFormInstance && this.tagFormInstance.hide()
      },
      handleRemove() {
        FilterStore.resetValue(this.property)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-tag {
        display: inline-flex;
        align-items: center;
        margin: 0 3px 10px;
        padding: 0 0 0 5px;
        border-radius: 2px;
        font-size: 12px;
        background: #f0f1f5;
        line-height: 22px;
        cursor: pointer;
        &:hover {
            background-color: #DCDEE5;
        }
        .tag-name {
            max-width: 150px;
            padding-right: 5px;
            color: #63656E;
            cursor: pointer;
            @include ellipsis;
        }
        .tag-colon {
            padding-right: 5px;
        }
        .tag-value {
            max-width: 220px;
            color: #313238;
            @include ellipsis;
        }
        .tag-delete {
            font-size: 20px;
            color: #9b9ea8;
            &:hover {
                color: #313238;
            }
        }
    }
</style>
