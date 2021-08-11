<template>
  <div class="association-list" v-bkloading="{ isLoading: loading }">
    <div class="association-empty" v-if="!hasAssociation">
      <div class="empty-content">
        <i class="bk-icon icon-empty">
        </i>
        <span>{{$t('暂无关联关系')}}</span>
      </div>
    </div>
    <template v-else>
      <cmdb-host-association-list-table
        ref="associationListTable"
        v-for="item in list"
        :key="item.id"
        :type="item.type"
        :id="item.modelId"
        :association-instances="item.instances"
        :association-type="item.associationType"
        @delete-association="handleDeleteAssociation">
      </cmdb-host-association-list-table>
    </template>
  </div>
</template>

<script>
  import bus from '@/utils/bus.js'
  import { mapGetters } from 'vuex'
  import cmdbHostAssociationListTable from './association-list-table.vue'
  export default {
    name: 'cmdb-host-association-list',
    components: {
      cmdbHostAssociationListTable
    },
    computed: {
      ...mapGetters('hostDetails', [
        'mainLine',
        'source',
        'target',
        'allInstances',
        'associationTypes'
      ]),
      id() {
        return parseInt(this.$route.params.id, 10)
      },
      hasAssociation() {
        return this.allInstances.length
      },
      list() {
        try {
          const list = []
          const associations = [...this.source, ...this.target]
          associations.forEach((association, index) => {
            const isSource = index < this.source.length
            const type = isSource ? 'source' : 'target'

            const modelId = isSource ? association.bk_asst_obj_id : association.bk_obj_id
            const associationType = this.associationTypes.find(item => item.bk_asst_id === association.bk_asst_id) || {}

            // 关联关系的唯一标识，用于匹配关联实例
            const objAsstId = isSource
              ? `host_${associationType.bk_asst_id}_${modelId}`
              : `${modelId}_${associationType.bk_asst_id}_host`

            list.push({
              // 关联关系id和源或目标的关系（指向）组成唯一性
              id: `${association.id}-${type}`,
              type,
              modelId,
              associationType,
              // 此关联关系下同一指向的关联实例并且关联id是匹配的
              instances: this.allInstances.filter((item) => {
                const sameType = item.bk_asst_id === association.bk_asst_id && item.type === type
                const matchAsst = item.bk_obj_asst_id === objAsstId
                return sameType && matchAsst
              })
            })
          })

          // 过滤掉无关联实例的关联
          return list.filter(item => item.instances.length)
        } catch (e) {
          console.error(e)
          return []
        }
      },
      loading() {
        return this.$loading([
          'getAssociation',
          'getMainLine',
          'getAssociationType',
          'getSourceAssociation',
          'getTargetAssociation'
        ])
      }
    },
    watch: {
      list() {
        this.$nextTick(() => {
          if (this.$refs.associationListTable) {
            const [firstAssociationListTable] = this.$refs.associationListTable
            firstAssociationListTable && (firstAssociationListTable.expanded = true)
          }
        })
      }
    },
    created() {
      this.getData()
      bus.$on('association-change', () => {
        this.getData()
      })
    },
    beforeDestroy() {
      bus.$off('association-change')
    },
    methods: {
      async getData() {
        try {
          const [source, target, mainLine, associationTypes, instances] = await Promise.all([
            this.getAssociation({ bk_obj_id: 'host' }),
            this.getAssociation({ bk_asst_obj_id: 'host' }),
            this.getMainLine(),
            this.getAssociationType(),
            this.getInstAssociation()
          ])
          const mainLineModels = mainLine.filter(model => !['biz', 'host'].includes(model.bk_obj_id))
          const availabelSource = this.getAvailableAssociation(source, [], mainLineModels)
          const availabelTarget = this.getAvailableAssociation(target, [], mainLineModels)
          this.setState({
            source: availabelSource,
            target: availabelTarget,
            mainLine: mainLineModels,
            associationTypes,
            instances
          })
        } catch (e) {
          this.setState({
            source: [],
            target: [],
            mainLine: [],
            associationTypes: [],
            instances: []
          })
          console.error(e)
        }
      },
      setState({ source, target, mainLine, associationTypes, instances }) {
        this.$store.commit('hostDetails/setAssociation', { type: 'source', association: source })
        this.$store.commit('hostDetails/setAssociation', { type: 'target', association: target })
        this.$store.commit('hostDetails/setMainLine', mainLine)
        this.$store.commit('hostDetails/setInstances', instances)
        this.$store.commit('hostDetails/setAssociationTypes', associationTypes)
      },
      getAssociation(condition) {
        return this.$store.dispatch('objectAssociation/searchObjectAssociation', {
          params: { condition },
          config: {
            requestId: 'getAssociation'
          }
        })
      },
      getMainLine() {
        return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {}, {
          config: {
            requestId: 'getMainLine'
          }
        })
      },
      async getAssociationType() {
        const { info } = await this.$store.dispatch('objectAssociation/searchAssociationType', {}, {
          config: {
            requestId: 'getAssociationType'
          }
        })
        return Promise.resolve(info)
      },
      async getInstAssociation() {
        try {
          const sourceCondition = { bk_obj_id: 'host', bk_inst_id: this.id }
          const targetCondition = { bk_asst_obj_id: 'host', bk_asst_inst_id: this.id }
          let [source, target] = await Promise.all([
            this.$store.dispatch('objectAssociation/searchInstAssociation', {
              params: { condition: sourceCondition, bk_obj_id: 'host' },
              config: { requestId: 'getSourceAssociation' }
            }),
            this.$store.dispatch('objectAssociation/searchInstAssociation', {
              params: { condition: targetCondition, bk_obj_id: 'host' },
              config: { requestId: 'getTargetAssociation' }
            })
          ])
          source = source.map(item => ({ ...item, type: 'source' }))
          target = target.map(item => ({ ...item, type: 'target' }))
          return [...source, ...target]
        } catch (error) {
          console.error(error)
          return []
        }
      },
      getAvailableAssociation(data, reference = [], mainLine = []) {
        return data.filter((association) => {
          const sourceId = association.bk_obj_id
          const targetId = association.bk_asst_obj_id
          const isMainLine = mainLine.some(model => [sourceId, targetId].includes(model.bk_obj_id))
          const isExist = reference.some(target => target.id === association.id)
          return !isMainLine && !isExist
        })
      },
      handleDeleteAssociation() {
        // 重新获取以刷新数据
        this.getData()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .association-list {
        height: 100%;
    }
    .association-empty {
        height: 100%;
        text-align: center;
        font-size: 14px;
        &:before {
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
            content: "";
        }
        .empty-content {
            display: inline-block;
            vertical-align: middle;
            .bk-icon {
                display: inline-block;
                margin: 0 0 10px 0;
                font-size: 65px;
                color: #c3cdd7;
            }
            span {
                display: inline-block;
                width: 100%;
            }
        }
    }
</style>
