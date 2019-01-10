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
    import throttle from 'lodash.throttle'
    import vClassifyItem from './classify-item'
    export default {
        components: {
            vClassifyItem
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'authorizedClassifications',
                'activeClassifications'
            ]),
            ...mapGetters('userPrivilege', ['privilege']),
            ...mapGetters(['admin']),
            hostManageClassification () {
                const hostManageClassification = {
                    'bk_classification_icon': 'icon-cc-host',
                    'bk_classification_id': 'bk_collection',
                    'bk_classification_name': this.$t('Hosts["主机管理"]'),
                    'bk_classification_type': 'inner',
                    'bk_objects': []
                }
                // 放开展示权限
                // if (this.admin || (this.privilege['model_config'] || {}).hasOwnProperty('bk_organization')) {
                hostManageClassification['bk_objects'].push({
                    'bk_obj_name': this.$t('Common["业务"]'),
                    'bk_obj_id': 'biz',
                    'bk_obj_icon': 'icon-cc-business',
                    'path': '/business',
                    'bk_classification_id': 'bk_collection'
                })
                // }
                if (this.admin || (this.privilege['sys_config']['global_busi'] || []).includes('resource')) {
                    hostManageClassification['bk_objects'].push({
                        'bk_obj_name': this.$t('Nav["主机"]'),
                        'bk_obj_id': 'resource',
                        'bk_obj_icon': 'icon-cc-host-free-pool',
                        'path': '/resource',
                        'bk_classification_id': 'bk_collection'
                    })
                }
                return hostManageClassification
            },
            classifyColumns () {
                const classifies = [
                    this.hostManageClassification,
                    // ...this.authorizedClassifications // 放开展示权限
                    ...this.activeClassifications
                ].filter(classification => {
                    return classification['bk_classification_id'] !== 'bk_organization' && classification['bk_objects'].length
                })
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