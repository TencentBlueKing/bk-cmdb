<template>
  <bk-select class="filter-collection"
    ref="selector"
    searchable
    multiple
    :popover-width="220"
    :disabled="!loaded"
    font-size="normal"
    v-model="selected"
    v-bk-tooltips="$t('已收藏的条件')"
    @click.native="loadCollections">
    <cmdb-loading class="filter-loading" :loading="loadingCollections"
      slot="trigger"
      slot-scope="scopeProps"
      :value="scopeProps.value">
      <icon-button class="filter-trigger" icon="icon-cc-star"
        :class="{ 'is-selected': !!storageCollection }"></icon-button>
    </cmdb-loading>
    <bk-option v-for="collection in collections"
      :key="collection.id"
      :id="collection.id"
      :name="collection.name"
      :disabled="isPending(collection)">
      <cmdb-loading :class="['collection-item', { 'is-editing': editState.id === collection.id }]"
        :loading="isPending(collection)"
        @click.native="handleNativeSelect($event, collection)">
        <template v-if="collection.id === editState.id">
          <bk-input class="colletion-form" size="small"
            v-focus
            v-model.trim="editState.value"
            @blur="handleUpdateCollection"
            @enter="handleUpdateCollection">
          </bk-input>
        </template>
        <template v-else>
          <i class="collection-state bk-icon icon-check-1" v-if="selected.includes(collection.id)"></i>
          <span class="collection-name">{{collection.name}}</span>
          <span class="collection-options">
            <i class="option-icon option-edit icon-cc-edit" @click.stop="handleEdit(collection)"></i>
            <bk-popconfirm
              trigger="click"
              placement="top"
              :tippy-options="{
                boundary: 'window'
              }"
              :on-show="deactiveAllowClose"
              :on-hide="activeAllowClose"
              :title="$t('收藏条件删除提示')"
              :confirm-button-is-text="false"
              @click.native.stop
              @confirm="handleRemove(collection)">
              <i class="option-icon option-delete bk-icon icon-close"></i>
            </bk-popconfirm>
          </span>
        </template>
      </cmdb-loading>
    </bk-option>
    <div class="business-extension" slot="extension">
      <a href="javascript:void(0)" class="extension-link"
        @click="handleCreate">
        <i class="bk-icon icon-plus-circle"></i>
        {{$t('新增收藏条件')}}
      </a>
    </div>
  </bk-select>
</template>

<script>
  import CmdbLoading from '@/components/loading/loading'
  import FilterStore from './store'
  import FilterForm from './filter-form.js'
  export default {
    components: {
      CmdbLoading
    },
    directives: {
      focus: {
        inserted: (el) => {
          const input = el.querySelector('input')
          input.focus()
        }
      }
    },
    data() {
      return {
        loaded: false,
        allowClose: true,
        editState: {
          raw: null,
          id: null,
          value: ''
        }
      }
    },
    computed: {
      selected: {
        get() {
          return this.storageCollection ? [this.storageCollection.id] : []
        },
        set(value = []) {
          this.handleApply(value)
        }
      },
      collections() {
        return FilterStore.collections || []
      },
      storageCollection() {
        return FilterStore.activeCollection
      },
      loadingCollections() {
        return this.$loading(FilterStore.request.collections)
      }
    },
    mounted() {
      const { instance } = this.$refs.selector.$refs.selectDropdown
      const bindedHideFunc = instance.props.onHide
      instance.set({
        placement: 'bottom-end',
        onHide: (instance) => {
          if (this.allowClose) {
            bindedHideFunc(instance)
            return true
          }
          return false
        }
      })
    },
    methods: {
      async loadCollections() {
        if (this.loaded || this.loadingCollections) {
          return false
        }
        try {
          await FilterStore.loadCollections()
          this.loaded = true
          this.$nextTick(() => {
            this.$refs.selector.show()
          })
        } catch (error) {
          console.error(error)
        }
      },
      handleApply(value) {
        this.$refs.selector.close()
        let selected
        if (value.length === 2) {
          const [, now] = value
          selected = now
        } else {
          const [now] = value
          selected = now
        }
        const collection = selected ? this.collections.find(collection => collection.id === selected) : null
        FilterStore.setActiveCollection(collection)
      },
      deactiveAllowClose() {
        this.allowClose = false
      },
      activeAllowClose() {
        this.allowClose = true
      },
      handleEdit(collection) {
        this.deactiveAllowClose()
        this.editState.id = collection.id
        this.editState.value = collection.name
        this.editState.raw = collection
      },
      isPending(collection) {
        return this.$loading([
          FilterStore.request.deleteCollection(collection.id),
          FilterStore.request.updateCollection(collection.id)
        ])
      },
      handleNativeSelect(event, collection) {
        if (collection.id === this.editState.id) {
          event.stopPropagation()
        }
      },
      async handleUpdateCollection() {
        if (!this.editState.raw) {
          return false
        }
        const hasChange = this.editState.raw.name !== this.editState.value
        const data = {
          ...this.editState.raw,
          bk_biz_id: FilterStore.bizId,
          name: this.editState.value
        }
        this.editState.raw = null
        this.editState.id = null
        this.editState.value = ''
        this.deactiveAllowClose()
        if (hasChange) {
          try {
            await FilterStore.updateCollection(data)
          } catch (error) {
            console.error(error)
          }
        }
        setTimeout(this.activeAllowClose, 300)
      },
      async handleRemove(collection) {
        try {
          this.deactiveAllowClose()
          await FilterStore.removeCollection(collection)
          this.activeAllowClose()
          if (this.selected.includes(collection.id)) {
            this.selected = []
          }
        } catch (error) {
          console.error(error)
        }
      },
      handleCreate() {
        FilterStore.setActiveCollection(null)
        FilterForm.show()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-collection {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        border: none;
        width: 32px;
        height: 32px;
        overflow: hidden;
        &.is-disabled {
            cursor: pointer;
        }
        /deep/ {
            .bk-tooltip-ref {
                display: flex !important;
                align-items: center;
                justify-content: center;
            }
        }
    }
    .filter-loading.loading {
        width: 32px;
        height: 32px;
        border: 1px solid #c4c6cc;
        border-radius: 2px;
    }
    .filter-trigger.icon-button {
        &:hover,
        &.is-selected {
            color: $primaryColor;
        }
        /deep/ {
            .icon-wrapper:before {
                font-size: 18px;
            }
        }
    }
    .collection-item {
        display: flex;
        align-items: center;
        padding: 0 16px;
        margin: 0 -16px;
        &.is-editing {
            height: 32px;
            padding: 0 6px;
        }
        &.loading {
            vertical-align: middle;
        }
        &:hover {
            .collection-options {
                display: initial;
            }
        }
        .collection-state {
            font-size: 24px;
            margin-left: -14px;
            & ~ .collection-name {
                margin-left: initial;
            }
        }
        .collection-name {
            margin-left: -6px;
            @include ellipsis;
        }
        .collection-options {
            display: none;
            margin-right: -10px;
            margin-left: auto;
            .option-icon {
                width: 24px;
                height: 24px;
                display: inline-flex;
                align-items: center;
                justify-content: center;
                color: $textColor;
                &:hover {
                    color: $primaryColor;
                }
            }
            .option-icon.option-edit {
                font-size: 12px;
            }
            .option-icon.option-delete {
                font-size: 22px;
            }
        }
        .collection-form {
            width: 100%;
        }
    }
    .extension-link {
        display: block;
        line-height: 38px;
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
