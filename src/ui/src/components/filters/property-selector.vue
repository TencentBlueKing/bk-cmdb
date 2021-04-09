<template>
  <bk-dialog
    v-model="isShow"
    :draggable="false"
    :width="730"
    :transfer="false"
    @after-leave="handleClosed">
    <div class="title" slot="tools">
      <span>{{$t('筛选条件')}}</span>
      <bk-input class="filter-input" v-model.trim="filter" clearable :placeholder="$t('请输入关键字搜索')"></bk-input>
    </div>
    <section class="property-selector">
      <div class="group"
        v-for="group in renderGroups"
        :key="group.id">
        <h2 class="group-title">
          {{group.name}}
          <span class="group-count">（{{group.children.length}}）</span>
        </h2>
        <ul class="property-list clearfix">
          <li class="property-item fl"
            v-for="property in group.children"
            :key="property.bk_property_id">
            <bk-checkbox class="property-checkbox"
              :checked="isChecked(property)"
              @change="handleToggleProperty(property, ...arguments)">
              {{property.bk_property_name}}
            </bk-checkbox>
          </li>
        </ul>
      </div>
    </section>
    <footer class="footer" slot="footer">
      <i18n class="selected-count"
        v-if="selected.length"
        path="已选择条数"
        tag="div">
        <span class="count" place="count">{{selected.length}}</span>
      </i18n>
      <div class="selected-options">
        <bk-button theme="primary" @click="confirm">{{$t('确定')}}</bk-button>
        <bk-button theme="default" @click="close">{{$t('取消')}}</bk-button>
      </div>
    </footer>
  </bk-dialog>
</template>

<script>
  import { mapGetters } from 'vuex'
  import FilterStore from './store'
  import throttle from 'lodash.throttle'
  export default {
    data() {
      return {
        filter: '',
        isShow: false,
        selected: [...FilterStore.selected],
        throttleFilter: throttle(this.handleFilter, 500, { leading: false }),
        renderGroups: []
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      propertyMap() {
        let modelPropertyMap = { ...FilterStore.modelPropertyMap }
        const ignoreHostProperties = ['bk_host_innerip', 'bk_host_outerip', '__bk_host_topology__']
        // eslint-disable-next-line max-len
        modelPropertyMap.host = modelPropertyMap.host.filter(property => !ignoreHostProperties.includes(property.bk_property_id))
        if (!FilterStore.bizId) {
          return modelPropertyMap
        }
        modelPropertyMap = {
          host: modelPropertyMap.host || [],
          module: modelPropertyMap.module || [],
          set: modelPropertyMap.set || []
        }
        return modelPropertyMap
      },
      groups() {
        const sequence = ['host', 'module', 'set', 'biz']
        return Object.keys(this.propertyMap).map((modelId) => {
          const model = this.getModelById(modelId) || {}
          return {
            id: modelId,
            name: model.bk_obj_name,
            children: this.propertyMap[modelId]
          }
        })
          .sort((groupA, groupB) => sequence.indexOf(groupA.id) - sequence.indexOf(groupB.id))
      }
    },
    watch: {
      filter: {
        immediate: true,
        handler() {
          this.throttleFilter()
        }
      }
    },
    methods: {
      handleFilter() {
        if (!this.filter.length) {
          this.renderGroups = this.groups
        } else {
          const filteredGroups = []
          const filter = this.filter.toLowerCase()
          this.groups.forEach((group) => {
            const properties = group.children.filter((property) => {
              const name = property.bk_property_name.toLowerCase()
              return name.indexOf(filter) > -1
            })
            if (properties.length) {
              filteredGroups.push({
                ...group,
                children: properties
              })
            }
          })
          this.renderGroups = filteredGroups
        }
      },
      isChecked(property) {
        return this.selected.some(target => target.id === property.id)
      },
      handleToggleProperty(property, checked) {
        if (checked) {
          this.selected.push(property)
        } else {
          const index = this.selected.findIndex(target => target.id === property.id)
          index > -1 && this.selected.splice(index, 1)
        }
      },
      async confirm() {
        FilterStore.updateSelected(this.selected)
        FilterStore.updateUserBehavior(this.selected)
        this.close()
      },
      handleClosed() {
        this.$emit('closed')
      },
      open() {
        this.isShow = true
      },
      close() {
        this.isShow = false
      }
    }
  }
</script>

<style lang="scss" scoped>
    .title {
        display: flex;
        justify-content: space-between;
        align-items: center;
        vertical-align: middle;
        line-height: 31px;
        font-size: 24px;
        color: #444;
        padding: 15px 0 0 24px;
        .filter-input {
            width: 240px;
            margin-right: 45px;
        }
    }
    .property-selector {
        margin: 0 -24px -24px 0;
        height: 350px;
        @include scrollbar-y;
    }
    .group {
        margin-top: 15px;
        .group-title {
            position: relative;
            padding: 0 0 0 15px;
            line-height: 20px;
            font-size: 15px;
            font-weight: bold;
            color: #63656E;
            &:before {
                content: "";
                position: absolute;
                left: 0;
                top: 3px;
                width: 4px;
                height: 14px;
                background-color: #C4C6CC;
            }
            .group-count {
                color: #C4C6CC;
                font-weight: normal;
            }
        }
    }
    .property-list {
        padding: 10px 0 6px 0;
        .property-item {
            width: 33%;
        }
    }
    .property-checkbox {
        display: block;
        margin: 8px 20px 8px 0;
        @include ellipsis;
        /deep/ {
            .bk-checkbox-text {
                max-width: calc(100% - 25px);
                @include ellipsis;
            }
        }
    }
    .footer {
        display: flex;
        .selected-count {
            font-size: 14px;
            line-height: 32px;
            .count {
                color: #2DCB56;
                padding: 0 4px;
            }
        }
        .selected-options {
            margin-left: auto;
        }
    }
</style>
