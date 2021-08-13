<template>
  <div class="association-list" v-bkloading="{ isLoading: loading }">
    <div class="association-empty" v-if="!hasRelation">
      <div class="empty-content">
        <i class="bk-icon icon-empty">
        </i>
        <span>{{$t('暂无关联关系')}}</span>
      </div>
    </div>
    <template v-else>
      <cmdb-relation-list-table
        ref="associationListTable"
        v-for="item in list"
        :key="item.id"
        :target-obj-id="item.modelId"
        :type="item.association.type"
        :association-instances="item.instances"
        :association-type="item.associationType"
        :obj-id="objId"
        @delete-association="handleDeleteAssociation">
      </cmdb-relation-list-table>
    </template>
    <div class="association-empty" v-if="hasRelation && !hasRelationInstance">
      <div class="empty-content">
        <i class="bk-icon icon-empty">
        </i>
        <span>{{$t('暂无关联数据')}}</span>
      </div>
    </div>
  </div>
</template>

<script>
  import bus from '@/utils/bus.js'
  import { mapGetters, mapActions } from 'vuex'
  import cmdbRelationListTable from './list-table.vue'
  export default {
    name: 'cmdb-relation-list',
    components: {
      cmdbRelationListTable
    },
    props: {
      associationObject: {
        type: Array,
        required: true
      },
      associationTypes: {
        type: Array,
        required: true
      }
    },
    data() {
      return {
        // 基于关联关系创建的实例关联列表，实例数据自身不包含所属关联关系字段数据
        // 需要在实例关联列表中通过 bk_obj_asst_id 识别
        instances: []
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectModelClassify', ['models']),
      objId() {
        return this.$parent.objId
      },
      instId() {
        return this.$parent.formatedInst.bk_inst_id
      },
      hasRelation() {
        return !!this.$parent.hasRelation
      },
      hasRelationInstance() {
        return !!this.instances.length
      },
      list() {
        const list = []
        // 基于当前模型对象的所有关联关系数据，组装出关联列表
        // 外层是关联关系，里层是关联下创建的实例关联列表
        this.associationObject.forEach((association) => {
          const isSource = association.type === 'source'

          // 用于展示关联关系中的模型名称等
          const modelId = isSource ? association.bk_asst_obj_id : association.bk_obj_id
          const associationType = this.associationTypes.find(item => item.bk_asst_id === association.bk_asst_id) || {}

          // 关联关系的唯一标识，用于匹配关联实例
          const objAsstId = isSource
            ? `${this.objId}_${associationType.bk_asst_id}_${modelId}`
            : `${modelId}_${associationType.bk_asst_id}_${this.objId}`

          list.push({
            // 关联关系id和源或目标的关系（指向）组成唯一性
            id: `${association.id}-${association.type}`,
            modelId,
            association,
            associationType,
            // 此关联关系下同一指向的关联实例并且关联id是匹配的
            instances: this.instances.filter((item) => {
              const sameType = item.bk_asst_id === association.bk_asst_id && item.type === association.type
              const matchAsst = item.bk_obj_asst_id === objAsstId
              return sameType && matchAsst
            })
          })
        })

        // 过滤掉无关联实例的关联
        return list.filter(item => item.instances.length)
      },
      loading() {
        return this.$loading([
          'getSourceAssociation',
          'getTargetAssociation'
        ])
      },
      resourceType() {
        return this.$parent.resourceType
      }
    },
    watch: {
      list() {
        this.expandFirstListTable()
      }
    },
    created() {
      this.getInstAssociation()

      bus.$on('association-change', async () => {
        await this.getInstAssociation()
      })

      this.expandFirstListTable()
    },
    beforeDestroy() {
      bus.$off('association-change')
    },
    methods: {
      ...mapActions('objectAssociation', [
        'searchObjectAssociation'
      ]),
      async getInstAssociation() {
        try {
          const sourceCondition = { bk_obj_id: this.objId, bk_inst_id: this.instId }
          const targetCondition = { bk_asst_obj_id: this.objId, bk_asst_inst_id: this.instId }
          let [source, target] = await Promise.all([
            this.$store.dispatch('objectAssociation/searchInstAssociation', {
              params: { condition: sourceCondition, bk_obj_id: this.objId },
              config: { requestId: 'getSourceAssociation' }
            }),
            this.$store.dispatch('objectAssociation/searchInstAssociation', {
              params: { condition: targetCondition, bk_obj_id: this.objId },
              config: { requestId: 'getTargetAssociation' }
            })
          ])
          source = source.map(item => ({ ...item, type: 'source' }))
          target = target.map(item => ({ ...item, type: 'target' }))
          this.instances = [...source, ...target]
        } catch (error) {
          console.error(error)
        }
      },
      expandFirstListTable() {
        this.$nextTick(() => {
          if (this.$refs.associationListTable) {
            const [firstAssociationListTable] = this.$refs.associationListTable
            firstAssociationListTable && (firstAssociationListTable.expanded = true)
          }
        })
      },
      handleDeleteAssociation() {
        // 重新获取以刷新数据
        this.getInstAssociation()
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
