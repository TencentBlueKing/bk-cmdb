<template>
    <div class="classify-layout clearfix">
        <div class="classify-waterfall fl"
            v-for="col in classifyColumns.length"
            :key="col">
            <cmdb-classify-panel
                v-for="classify in classifyColumns[col - 1]"
                :key="classify['bk_classification_id']"
                :classify="classify"
                :collection="collection"
                :instance-count="instanceCount">
            </cmdb-classify-panel>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import {
        MENU_RESOURCE_COLLECTION,
        MENU_RESOURCE_HOST_COLLECTION,
        MENU_RESOURCE_BUSINESS_COLLECTION
    } from '@/dictionary/menu-symbol'
    import cmdbClassifyPanel from './children/classify-panel'
    export default {
        components: {
            cmdbClassifyPanel
        },
        data () {
            return {
                instanceCount: []
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', ['classifications', 'models']),
            ...mapGetters('userCustom', ['usercustom']),
            collection () {
                const isHostCollected = this.usercustom[MENU_RESOURCE_HOST_COLLECTION] === undefined
                    ? true
                    : this.usercustom[MENU_RESOURCE_HOST_COLLECTION]
                const isBusinessCollected = this.usercustom[MENU_RESOURCE_BUSINESS_COLLECTION] === undefined
                    ? true
                    : this.usercustom[MENU_RESOURCE_BUSINESS_COLLECTION]
                const collection = [...(this.usercustom[MENU_RESOURCE_COLLECTION] || [])]
                if (isHostCollected) {
                    collection.push('host')
                }
                if (isBusinessCollected) {
                    collection.push('biz')
                }
                return collection.filter(modelId => {
                    return this.models.some(model => model.bk_obj_id === modelId)
                })
            },
            filteredClassifications () {
                const result = []
                const filterModels = ['plat', 'process']
                const filterClassify = ['bk_biz_topo']
                this.classifications.forEach(classification => {
                    if (!filterClassify.includes(classification.bk_classification_id)) {
                        const models = classification.bk_objects.filter(model => !filterModels.includes(model.bk_obj_id))
                        if (models.length) {
                            result.push({
                                ...classification,
                                bk_objects: models
                            })
                        }
                    }
                })
                return result
            },
            classifyColumns () {
                const colHeight = [0, 0, 0, 0]
                const classifyColumns = [[], [], [], []]
                this.filteredClassifications.forEach(classify => {
                    const minColHeight = Math.min(...colHeight)
                    const rowIndex = colHeight.indexOf(minColHeight)
                    classifyColumns[rowIndex].push(classify)
                    colHeight[rowIndex] = colHeight[rowIndex] + this.calcWaterfallHeight(classify)
                })
                return classifyColumns
            }
        },
        created () {
            this.getInstanceCount()
        },
        methods: {
            async getInstanceCount () {
                try {
                    this.instanceCount = await this.$store.dispatch('objectCommonInst/getInstanceCount')
                } catch (e) {
                    console.error(e)
                    this.instanceCount = []
                }
            },
            calcWaterfallHeight (classify) {
                // 46px 分类高度
                // 16px 模型列表padding
                // 36 模型高度
                return 46 + 16 + classify['bk_objects'].length * 36
            }
        }
    }
</script>

<style lang="scss" scoped>
    .classify-layout{
        padding: 0 23px 45px 20px;
    }
    .classify-waterfall{
        width: calc((100% - 60px) / 4);
        margin: 0 0 0 20px;
        &:first-child{
            margin: 0;
        }
    }
</style>
