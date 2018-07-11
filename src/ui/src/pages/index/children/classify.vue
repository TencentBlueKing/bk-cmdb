<template>
    <div class="classify-layout clearfix" :class="`recently-total-row-${recentlyRow}`">
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
        data () {
            const hostManageClassification = {
                'bk_classification_icon': bkHostManage.icon,
                'bk_classification_id': bkHostManage.id,
                'bk_classification_name': this.$t(bkHostManage.i18n),
                'bk_classification_type': 'inner',
                'bk_objects': bkHostManage.children.map(nav => {
                    return {
                        'bk_obj_name': this.$t(nav.i18n),
                        'bk_obj_id': nav.id,
                        'bk_obj_icon': nav.icon,
                        'path': nav.path,
                        'bk_classification_id': bkHostManage.id
                    }
                })
            }
            return {
                hostManageClassification
            }
        },
        computed: {
            ...mapGetters('navigation', ['authorizedClassifications', 'authorizedNavigation']),
            ...mapGetters('usercustom', ['usercustom', 'recentlyKey']),
            recently () {
                return this.usercustom[this.recentlyKey] || []
            },
            recentlyModels () {
                let models = []
                this.recently.forEach(path => {
                    const model = this.getRouteModel(path)
                    if (model) {
                        models.push(model)
                    }
                })
                return models
            },
            recentlyRow () {
                return this.recentlyModels.length > 4 ? 2 : 1
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
            },
            getRouteModel (path) {
                let model
                for (let i = 0; i < this.authorizedNavigation.length; i++) {
                    const models = this.authorizedNavigation[i]['children'] || []
                    model = models.find(model => model.path === path)
                    if (model) break
                }
                return model
            }
        }
    }
</script>

<style lang="scss" scoped>
    .classify-layout{
        height: calc(100% - 300px);
        overflow-y: auto;
        @include scrollbar;
        width: 90%;
        margin: 20px auto 0;
        &.recently-total-row-2{
            height: calc(100% - 400px);
        }
    }
    .classify-waterfall{
        width: calc((100% - 60px) / 4);
        margin: 0 0 0 20px;
        &:first-child{
            margin: 0;
        }
    }
</style>