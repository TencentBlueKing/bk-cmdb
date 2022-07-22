<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="apply-edit">
    <component
      :is="currentView"
      :mode="mode"
      :ids="targetIds"
      :action="action">
    </component>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import has from 'has'
  import multiConfig from './children/multi-config'
  import singleConfig from './children/single-config'
  import { MENU_BUSINESS_HOST_APPLY_CONFIRM  } from '@/dictionary/menu-symbol'
  import serviceTemplateService, { CONFIG_MODE } from '@/services/service-template/index.js'

  export default {
    components: {
      multiConfig,
      singleConfig
    },
    data() {
      return {
        currentView: '',
        moduleMap: {},
        templateMap: new Map()
      }
    },
    provide() {
      return {
        getModuleName: this.getModuleName,
        getModulePath: this.getModulePath,
        getTemplateName: this.getTemplateName
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      mode() {
        return this.$route.params.mode
      },
      targetIds() {
        const { id } = this.$route.query
        let targetIds = []
        if (id) {
          targetIds = String(id).split(',')
            .map(id => Number(id))
        }
        return targetIds
      },
      isBatch() {
        return has(this.$route.query, 'batch')
      },
      action() {
        return this.$route.query.action
      },
      title() {
        let title
        if (this.isBatch) {
          title = this.$t(this.action === 'batch-del' ? '批量删除' : '批量编辑')
        } else {
          const getName = {
            [CONFIG_MODE.MODULE]: this.getModuleName,
            [CONFIG_MODE.TEMPLATE]: this.getTemplateName
          }
          title = `${this.$t('编辑')} ${getName[this.mode](this.targetIds[0])}`
        }
        return title
      }
    },
    async created() {
      if (this.mode === CONFIG_MODE.MODULE) {
        await this.initTopoData()
      } else if (this.mode === CONFIG_MODE.TEMPLATE) {
        await this.initTemplateData()
      }

      this.setBreadcrumbs()
      this.currentView = this.isBatch ? multiConfig.name : singleConfig.name
    },
    beforeRouteLeave(to, from, next) {
      if (to.name !== MENU_BUSINESS_HOST_APPLY_CONFIRM) {
        this.$store.commit('hostApply/clearRuleDraft')
      }
      next()
    },
    methods: {
      async initTopoData() {
        try {
          const topopath = await this.getTopopath()
          const moduleMap = {}
          topopath.nodes.forEach((node) => {
            moduleMap[node.topo_node.bk_inst_id] = node.topo_path
          })
          this.moduleMap = Object.freeze(moduleMap)
        } catch (e) {
          console.log(e)
        }
      },
      async initTemplateData() {
        try {
          const templateMap = new Map()
          const list = await serviceTemplateService.findAllByIds(this.targetIds, { bk_biz_id: this.bizId })
          list.forEach(item => templateMap.set(item.id, item))
          this.templateMap = templateMap
        } catch (e) {
          console.log(e)
        }
      },
      setBreadcrumbs() {
        this.$store.commit('setTitle', this.title)
      },
      getTopopath() {
        return this.$store.dispatch('hostApply/getTopopath', {
          bizId: this.bizId,
          params: {
            topo_nodes: this.targetIds.map(id => ({ bk_obj_id: 'module', bk_inst_id: id }))
          }
        })
      },
      getModulePath(id) {
        const info = this.moduleMap[id] || []
        const path = info.map(node => node.bk_inst_name).reverse()
          .join(' / ')
        return path
      },
      getModuleName(id) {
        const topoInfo = this.moduleMap[id] || []
        const target = topoInfo.find(target => target.bk_obj_id === 'module' && target.bk_inst_id === id) || {}
        return target.bk_inst_name
      },
      getTemplateName(id) {
        return this.templateMap.get(id)?.name
      }
    }
  }
</script>

<style lang="scss" scoped>
  .apply-edit {
    padding: 0;
  }
</style>
