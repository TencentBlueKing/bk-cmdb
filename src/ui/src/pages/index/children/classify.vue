<template>
    <div class="classify-layout clearfix">
        <div class="classify-waterfall fl" 
            v-for="col in classifyColumns.length" 
            :key="col">
            <v-classify-item 
                v-for="classify in classifyColumns[col - 1]"
                :key="classify['bk_classification_id']"
                :classify="classify">
            </v-classify-item>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { bk_host_manage as bkHostManage } from '@/common/json/static_navigation.json'
    import throttle from 'lodash.throttle'
    import vClassifyItem from './classify-item'
    export default {
        components: {
            vClassifyItem
        },
        computed: {
            ...mapGetters('navigation', ['authorizedClassifications', 'authorizedNavigation']),
            hostManageClassification () {
                const hostManageClassification = {
                    'bk_classification_icon': bkHostManage.icon,
                    'bk_classification_id': bkHostManage.id,
                    'bk_classification_name': this.$t(bkHostManage.i18n),
                    'bk_classification_type': 'inner'
                }
                const hostNavigation = this.authorizedNavigation.find(({id}) => id === bkHostManage.id)
                const authorizedHostModels = bkHostManage.children.filter(model => hostNavigation.children.some(nav => nav.id === model.id))
                hostManageClassification['bk_objects'] = authorizedHostModels.map(model => {
                    return {
                        'bk_obj_name': this.$t(model.i18n),
                        'bk_obj_id': model.id,
                        'bk_obj_icon': model.icon,
                        'path': model.path,
                        'bk_classification_id': bkHostManage.id
                    }
                })
                return hostManageClassification
            },
            classifyColumns () {
                const classifies = [this.hostManageClassification, ...this.authorizedClassifications]
                let colHeight = [0, 0, 0, 0]
                let classifyColumns = [[], [], [], []]
                classifies.forEach(classify => {
                    const minColHeight = Math.min(...colHeight)
                    const rowIndex = colHeight.indexOf(minColHeight)
                    classifyColumns[rowIndex].push(classify)
                    colHeight[rowIndex] = colHeight[rowIndex] + this.calcWaterfallHeight(classify)
                })
                return classifyColumns
            }
        },
        methods: {
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
        width: 90%;
        margin: 20px auto 0;
        padding: 0 3px 45px 0;
    }
    .classify-waterfall{
        width: calc((100% - 60px) / 4);
        margin: 0 0 0 20px;
        &:first-child{
            margin: 0;
        }
    }
</style>