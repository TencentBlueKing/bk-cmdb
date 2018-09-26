<template>
    <div class="topo-wrapper">
        <v-topo-list 
            v-if="bkClassificationId === 'bk_biz_topo'"
            @createModel="createModel"
            @editModel="editModel"
            ref="topo"
        ></v-topo-list>
        <v-topo
            v-else-if="bkClassificationId !== void 0"
            @createModel="createModel"
            @editModel="editModel"
            @editClassify="editClassify"
            ref="topo"
        ></v-topo>
        <v-global-models v-else></v-global-models>
        <cmdb-slider
            :isShow.sync="slider.isShow" :title="slider.title"
            :beforeClose="handleSliderBeforeClose">
            <v-details slot="content"
                ref="details"
                :isEdit.sync="slider.isEdit"
                @createModel="updateTopo(true)"
                @updateModel="updateTopo"
                @cancel="slider.isShow = false"
            ></v-details>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapMutations, mapActions } from 'vuex'
    import vGlobalModels from './topo/global-models'
    import vTopo from './topo/topo'
    import vTopoList from './topo/topo-list'
    import vDetails from './details'
    export default {
        components: {
            vGlobalModels,
            vTopo,
            vTopoList,
            vDetails
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    title: '',
                    isEdit: false,
                    isCloseConfirmShow: false
                }
            }
        },
        computed: {
            bkClassificationId () {
                return this.$route.params.classifyId
            }
        },
        methods: {
            ...mapActions('objectModelClassify', [
                'searchClassificationsObjects'
            ]),
            ...mapMutations('objectModel', [
                'setActiveModel'
            ]),
            createModel (prevModelId) {
                this.slider.title = this.$t('ModelManagement["新增模型"]')
                this.slider.isShow = true
                this.slider.isEdit = false
                this.setActiveModel({
                    bk_classification_id: this.bkClassificationId,
                    bk_asst_obj_id: prevModelId
                })
            },
            editModel (model) {
                this.slider.title = model['bk_obj_name']
                this.setActiveModel(model)
                this.slider.isEdit = true
                this.slider.isShow = true
            },
            editClassify () {
                this.$emit('editClassify')
            },
            handleSliderBeforeClose () {
                if (this.$refs.details.isCloseConfirmShow()) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
                            confirmFn: () => {
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                return true
            },
            async updateTopo (isShow) {
                await this.searchClassificationsObjects({})
                this.slider.isShow = isShow
                this.$refs.topo.initTopo()
            }
        }
    }
</script>
