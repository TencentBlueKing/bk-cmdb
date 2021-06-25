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
      <template v-for="item in list">
        <cmdb-relation-list-table
          ref="associationListTable"
          v-for="association in item.associations"
          :key="association.id"
          :type="item.type"
          :id="item.id"
          :all-instances="instances"
          :association-type="item.associationType"
          :obj-id="objId"
          :inst-id="instId"
          @delete-association="handleDeleteAssociation">
        </cmdb-relation-list-table>
      </template>
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
      uniqueAssociationObject() {
        const ids = this.associationObject.map(association => association.id)
        return [...new Set(ids)].map(id => this.associationObject.find(association => association.id === id))
      },
      list() {
        try {
          const list = []
          this.uniqueAssociationObject.forEach((association) => {
            const isSource = association.bk_obj_id === this.objId
            const modelId = isSource ? association.bk_asst_obj_id : association.bk_obj_id
            list.push({
              type: isSource ? 'source' : 'target',
              id: modelId,
              associationType: this.associationTypes.find(target => target.bk_asst_id === association.bk_asst_id) || {},
              associations: [association]
            })
          })
          return list
        } catch (e) {
          console.log(e)
        }
        return []
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
      this.getAssociation()

      bus.$on('association-change', async () => {
        await this.getAssociation()
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
      async getAssociation() {
        try {
          const sourceCondition = { bk_obj_id: this.objId, bk_inst_id: this.instId }
          const targetCondition = { bk_asst_obj_id: this.objId, bk_asst_inst_id: this.instId }
          const [source, target] = await Promise.all([
            this.$store.dispatch('objectAssociation/searchInstAssociation', {
              params: { condition: sourceCondition },
              config: { requestId: 'getSourceAssociation' }
            }),
            this.$store.dispatch('objectAssociation/searchInstAssociation', {
              params: { condition: targetCondition },
              config: { requestId: 'getTargetAssociation' }
            })
          ])
          this.instances = [...source, ...target]
        } catch (error) {
          console.error(error)
        }
      },
      expandFirstListTable() {
        this.$nextTick(() => {
          if (this.$refs.associationListTable) {
            const firstAssociationListTable = this.$refs.associationListTable.find(listTable => listTable.hasInstance)
            firstAssociationListTable && (firstAssociationListTable.expanded = true)
          }
        })
      },
      handleDeleteAssociation(id) {
        const index = this.instances.findIndex(instance => instance.id === id)
        index > -1 && this.instances.splice(index, 1)
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
