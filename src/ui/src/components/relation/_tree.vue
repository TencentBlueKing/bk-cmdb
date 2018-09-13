<template>
    <div class="relation-tree-layout" v-bkloading="{isLoading: $loading()}">
        <template v-if="next.length">
            <div
                v-for="(relation, index) in next" :key="index">
                <div class="tree-row" @click="handleRowClick(relation)">
                    <i class="tree-row-expand-icon icon-cc-triangle-sider"
                        v-if="relation.children.length"
                        :class="{expanded: relation.show}">
                    </i>
                    <i :class="['tree-row-icon', relation['bk_obj_icon']]"></i>
                    <span>{{relation['bk_obj_name']}}</span>
                    <span class="tree-row-count fr" v-if="relation.children.length">{{relation.children.length}}</span>
                </div>
                <div v-show="relation.show">
                    <div class="tree-row tree-row-child" v-for="(inst, index) in relation.children" :key="index">
                        <i :class="['tree-row-icon', relation['bk_obj_icon']]"></i>
                        <span>{{inst['bk_inst_name']}}</span>
                        <span class="tree-row-details fr" @click="handleShowDetails(relation, inst)">{{$t('Common["详情信息"]')}}</span>
                    </div>
                </div>
            </div>
        </template>
        <div class="relation-empty" v-else>
            <img class="empty-image" src="../../assets/images/relevance-empty.png">
            <span class="empty-text">{{$t("Common['当前还未有关联项']")}}</span>
        </div>
        <cmdb-topo-details v-if="details.show"
            :objId="details.objId"
            :instId="details.instId"
            :title="details.title"
            :show.sync="details.show">
        </cmdb-topo-details>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import cmdbTopoDetails from './_details.vue'
    export default {
        components: {
            cmdbTopoDetails
        },
        data () {
            return {
                next: [],
                ignore: ['plat'],
                details: {
                    show: false,
                    objId: null,
                    instId: null,
                    title: ''
                }
            }
        },
        computed: {
            objId () {
                return this.$parent.objId
            },
            instId () {
                return this.$parent.instId
            }
        },
        created () {
            this.getRelationData()
        },
        methods: {
            ...mapActions('objectRelation', ['getInstRelation']),
            async getRelationData () {
                return this.getInstRelation({
                    objId: this.objId,
                    instId: this.instId,
                    config: {
                        requestId: `get_getInstRelation_${this.objId}_${this.instId}`,
                        fromCache: true
                    }
                }).then(data => {
                    const next = data[0].next.filter(obj => !this.ignore.includes(obj['bk_obj_id']))
                    next.forEach(obj => {
                        obj.show = false
                    })
                    this.next = next
                    return data
                })
            },
            handleRowClick (relation) {
                relation.show = !relation.show
            },
            handleShowDetails (obj, inst) {
                this.details.objId = obj['bk_obj_id']
                this.details.instId = inst['bk_inst_id']
                this.details.title = `${obj['bk_obj_name']}-${inst['bk_inst_name']}`
                this.details.show = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-tree-layout {
        min-height: 250px;
    }
    .tree-row {
        position: relative;
        margin: -1px 0 0 0;
        padding: 0 0 0 45px;
        border: 1px solid #bec6de;
        line-height: 36px;
        color: #6b7baa;
        cursor: pointer;
        &.tree-row-child {
            padding: 0 0 0 65px;
        }
        &:hover {
            background-color: #f1f7ff;
            .tree-row-details {
                display: block;
            }
        }
        .tree-row-expand-icon {
            position: absolute;
            left: 30px;
            top: 11px;
            transform: rotate(180deg);
            transition: all .2s;
            &.expanded {
                transform: rotate(225deg);
            }
        }
        .tree-row-icon {
            color: #ffb400;
            margin: 0 6px 0 0;
        }
        .tree-row-count {
            height: 14px;
            padding: 0 8px;
            margin: 11px 13px 11px 10px;
            color: #ffb80f;
            background: #fff7e5;
            border-radius: 7px;
            font-size: 12px;
            line-height: 14px;
        }
        .tree-row-details {
            display: none;
            height: 22px;
            margin: 7px 10px;
            padding: 0 5px;
            line-height: 20px;
            border: 1px solid #c3cdd7;
            border-radius: 2px;
            color: #737987;
            background: #fff;
        }
    }
    .relation-empty {
        text-align: center;
        .empty-image {
            display: block;
            margin: 50px auto 20px;
            width: 130px;
        }
        .empty-text {
            font-size: 14px;
        }
    }
</style>