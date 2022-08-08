<script>
  import { computed, defineComponent } from '@vue/composition-api'
  import store from '@/store'
  import { t } from '@/i18n'
  import routerActions from '@/router/actions'
  import CmdbLoading from '@/components/loading/loading'
  import { MENU_BUSINESS_SET_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      CmdbLoading
    },
    props: {
      node: {
        type: Object,
        default: () => ({})
      },
      data: {
        type: Object,
        default: () => ({})
      },
      isBlueKing: Boolean,
      editable: Boolean,
      nodeCountType: String
    },
    setup(props, { emit }) {
      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const nodeIconMap = {
        1: 'icon-cc-host-free-pool',
        2: 'icon-cc-host-breakdown',
        default: 'icon-cc-host-free-pool'
      }

      const getInternalNodeClass = (node, data) => {
        const classNames = []
        classNames.push(nodeIconMap[data.default] || nodeIconMap.default)
        if (node.selected) {
          classNames.push('is-selected')
        }
        return classNames
      }

      const isTemplate = node => node.data.service_template_id || node.data.set_template_id

      const isShowCreate = (node, data) => {
        const isModule = data.bk_obj_id === 'module'
        const isIdleSet = data.is_idle_set
        const isContainer = data.is_container
        return !isModule && !isIdleSet && !isContainer
      }

      const getSetNodeTips = (node) => {
        const tips = document.createElement('div')
        const span = document.createElement('span')
        span.innerText = t('需在集群模板中新建')
        const link = document.createElement('a')
        link.innerText = t('立即跳转')
        link.href = 'javascript:void(0)'
        link.style.color = '#3a84ff'
        link.addEventListener('click', () => {
          routerActions.redirect({
            name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
            params: {
              templateId: node.data.set_template_id
            },
            history: true
          })
        })
        tips.appendChild(span)
        tips.appendChild(link)
        return tips
      }

      const getNodeCount = (data) => {
        const count = data[props.nodeCountType]
        if (typeof count === 'number') {
          return count
        }
        return 0
      }

      const handleSetNodeTipsToggle = (tips) => {
        const element = tips.reference.parentElement
        if (tips.state.isVisible) {
          element.classList.remove('hovering')
        } else {
          element.classList.add('hovering')
        }
        return true
      }

      const handleCreate = (node) => {
        emit('create', node)
      }

      return {
        bizId,
        getInternalNodeClass,
        isTemplate,
        isShowCreate,
        getSetNodeTips,
        getNodeCount,
        handleSetNodeTipsToggle,
        handleCreate
      }
    }
  })
</script>

<template>
  <div :class="['topology-tree-node', { 'is-container': node.data.is_container, 'is-selected': node.selected }]">
    <i class="internal-node-icon"
      v-if="data.default !== 0"
      :class="getInternalNodeClass(node, data)">
    </i>
    <i v-else
      :class="['node-icon', { 'is-selected': node.selected, 'is-template': isTemplate(node) }]">
      {{data.icon_text || data.bk_obj_name[0]}}
    </i>

    <span class="node-name" :title="node.name">{{node.name}}</span>

    <div class="node-extra">
      <cmdb-auth v-if="isShowCreate(node, data)"
        class="node-create-trigger"
        :auth="{ type: $OPERATION.C_TOPO, relation: [bizId] }">
        <template #default="{ disabled }">
          <i v-if="isBlueKing && !editable"
            class="node-button disabled-node-button"
            v-bk-tooltips.top="{ content: $t('蓝鲸业务拓扑节点提示'), interactive: false }">
            {{$t('新建')}}
          </i>
          <i v-else-if="data.set_template_id"
            class="node-button disabled-node-button"
            v-bk-tooltips.top="{
              content: getSetNodeTips(node),
              interactive: true,
              onShow: handleSetNodeTipsToggle,
              onHide: handleSetNodeTipsToggle
            }">
            {{$t('新建')}}
          </i>
          <bk-button v-else class="node-button" v-test-id="'createNode'"
            theme="primary"
            :disabled="disabled"
            @click.stop="handleCreate(node)">
            {{$t('新建')}}
          </bk-button>
        </template>
      </cmdb-auth>

      <cmdb-loading :class="['node-count', { 'is-selected': node.selected }]"
        :loading="['pending', undefined].includes(data.status)">
        {{getNodeCount(data)}}
      </cmdb-loading>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.topology-tree-node {
  display: flex;

  &:hover {
    .node-extra {
      .node-create-trigger {
        display: inline-block;
        & ~ .node-count {
          display: none;
        }
      }
    }
  }

  &.is-container {
    .node-icon {
      border-radius: 4px;
    }
  }

  .internal-node-icon {
    width: 20px;
    height: 20px;
    line-height: 20px;
    text-align: center;
    margin: 8px 4px 8px 0;
    &.is-selected {
      color: #FFB400;
    }
  }

  .node-icon {
    flex: none;
    width: 20px;
    height: 20px;
    margin: 8px 4px 8px 0;
    border-radius: 50%;
    background-color: #C4C6CC;
    line-height: 1.666667;
    text-align: center;
    font-size: 12px;
    font-style: normal;
    color: #FFF;
    &.is-template {
      background-color: #97aed6;
    }
    &.is-selected {
      background-color: #3A84FF;
    }
  }

  .node-name {
    display: block;
    height: 36px;
    line-height: 36px;
    overflow: hidden;
    @include ellipsis;
  }

  .node-extra {
    margin-left: auto;
    display: flex;

    .node-create-trigger {
      display: none;
      font-size: 0;
      &.hovering {
        display: inline-block;
        & ~ .node-count {
          display: none;
        }
      }
      .node-button {
        height: 24px;
        padding: 0 6px;
        margin: 0 20px 0 4px;
        line-height: 22px;
        border-radius: 4px;
        font-size: 12px;
        min-width: auto;
        &.disabled-node-button {
          @include inlineBlock;
          line-height: 24px;
          font-style: normal;
          background-color: #dcdee5;
          color: #ffffff;
          outline: none;
          cursor: not-allowed;
        }
      }
    }

    .node-count {
      padding: 0 5px;
      margin: 9px 20px 9px 4px;
      height: 18px;
      line-height: 17px;
      border-radius: 2px;
      background-color: #f0f1f5;
      color: #979ba5;
      font-size: 12px;
      text-align: center;
      &.is-selected {
        background-color: #a2c5fd;
        color: #fff;
      }
      &.loading {
        background-color: transparent;
      }
    }
  }
}
</style>
