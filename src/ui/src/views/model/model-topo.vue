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
            :hasCloseConfirm="true"
            :isCloseConfirmShow="slider.isCloseConfirmShow"
            :isShow.sync="slider.isShow" :title="slider.title"
            @closeSlider="closeSlider">
            <v-details slot="content"
                ref="details"
                :isEdit.sync="slider.isEdit"
                @createModel="updateTopo(true)"
                @updateModel="updateTopo(false)"
                @cancel="slider.isShow = false"
            ></v-details>
        </cmdb-slider>
    </div>
</template>

<script>
    import vGlobalModels from './topo/global-models'
    import vTopo from './topo/topo'
    import vTopoList from './topo/topo-list'
    export default {
        components: {
            vGlobalModels,
            vTopo,
            vTopoList
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
                this.pop.isEdit = true
                this.pop.isShow = true
            },
            createClassify () {
                this.pop.isEdit = false
                this.pop.isShow = true
            },
            closeSlider () {
                this.slider.isCloseConfirmShow = this.$refs.details.isCloseConfirmShow()
                if (!this.slider.isCloseConfirmShow) {
                    this.slider.isShow = false
                }
            }
        }
    }
</script>
