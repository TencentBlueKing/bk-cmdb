<template>
    <bk-sideslider
        v-transfer-dom
        :width="514"
        :title="title"
        :is-show.sync="isShow"
        @close="handleClose">
        <div class="association-layout" slot="content" v-if="isShow">
            <div class="form-group">
                <label class="form-label"
                    :class="{
                        required: !isViewMode
                    }">
                    {{$t('源模型')}}
                </label>
                <bk-input type="text"class="cmdb-form-input"
                    disabled
                    :value="getModelName(info.source)">
                </bk-input>
            </div>
            <div class="form-group">
                <label class="form-label"
                    :class="{
                        required: !isViewMode
                    }">
                    {{$t('目标模型')}}
                </label>
                <bk-input type="text"class="cmdb-form-input"
                    disabled
                    :value="getModelName(info.target)">
                </bk-input>
            </div>
            <div class="form-group" v-if="!isViewMode">
                <label class="form-label required">
                    {{$t('关联类型')}}
                </label>
                <ul class="association-list clearfix">
                    <li class="association-item fl"
                        v-for="(association, index) in localAssociationList"
                        :key="index"
                        :class="{
                            selected: association.id === info.association.id
                        }"
                        @click="handleSelectAssociation(association)">
                        {{getAssociationItemName(association)}}
                    </li>
                </ul>
            </div>
            <div class="form-group">
                <label class="form-label"
                    :class="{
                        required: !isViewMode
                    }">
                    {{$t('关联描述')}}
                </label>
                <bk-input type="text"class="cmdb-form-input"
                    :disabled="isViewMode"
                    name="description"
                    v-validate="'required|singlechar|length:256'"
                    v-model="info.description">
                </bk-input>
                <p class="form-error" v-if="errors.has('description')">{{errors.first('description')}}</p>
            </div>
            <div class="form-group">
                <label class="form-label"
                    :class="{
                        required: !isViewMode
                    }">
                    {{$t('源-目标约束')}}
                </label>
                <cmdb-selector
                    name="constraint"
                    v-validate="'required'"
                    v-model="info.constraint"
                    :auto-select="false"
                    :list="constraintList"
                    :disabled="isViewMode">
                </cmdb-selector>
                <p class="form-error" v-if="errors.has('constraint')">{{errors.first('constraint')}}</p>
            </div>
            <div class="button-group" v-if="isEditMode && !info.ispre">
                <bk-button class="form-button"
                    theme="primary"
                    :loading="$loading()"
                    @click="handleSave">
                    {{$t('确定')}}
                </bk-button>
                <bk-button class="form-button"
                    v-if="isViewMode"
                    theme="danger"
                    :loading="$loading()"
                    @click="handleDelete">
                    {{$t('删除关联')}}
                </bk-button>
                <bk-button class="form-button"
                    v-else
                    theme="default"
                    @click="handleCancel">
                    {{$t('取消')}}
                </bk-button>
            </div>
        </div>
    </bk-sideslider>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    const DEFAULT_INFO = {
        source: null,
        target: null,
        description: '',
        constraint: '',
        ispre: false,
        association: {}
    }
    export default {
        name: 'cmdb-graphics-assoication',
        inject: ['parentRefs'],
        data () {
            return {
                title: '',
                info: { ...DEFAULT_INFO },
                constraintList: [{
                    id: 'n:n',
                    name: 'N-N'
                }, {
                    id: '1:n',
                    name: '1-N'
                }, {
                    id: '1:1',
                    name: '1-1'
                }]
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', ['models']),
            ...mapGetters('objectAssociation', ['associationList']),
            ...mapGetters('globalModels', [
                'association',
                'addEdgePromise',
                'isEditMode'
            ]),
            localAssociationList () {
                const filterFlag = ['bk_mainline']
                return this.associationList.filter(association => {
                    return !filterFlag.includes(association['bk_asst_id'])
                })
            },
            isShow: {
                get () {
                    return this.association.show
                },
                set (val) {
                    this.$store.commit('globalModels/setAssociation', {
                        show: val,
                        edge: null
                    })
                }
            },
            isViewMode () {
                const edge = this.association.edge || {}
                const data = edge.data || {}
                return typeof data['bk_inst_id'] === 'number'
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    if (this.isViewMode) {
                        this.getAssociationInfo()
                    } else {
                        this.setCreateInfo()
                    }
                } else {
                    this.reset()
                }
            }
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchObjectAssociation',
                'createObjectAssociation',
                'updateObjectAssociation',
                'deleteObjectAssociation'
            ]),
            async getAssociationInfo () {
                const id = this.association.edge.data['bk_inst_id']
                const [data] = await this.searchObjectAssociation({
                    params: this.$injectMetadata({
                        condition: { id }
                    })
                })
                this.info = {
                    source: data['bk_obj_id'],
                    target: data['bk_asst_obj_id'],
                    description: data['bk_obj_asst_name'],
                    constraint: data['mapping'],
                    ispre: data['ispre']
                }
                this.title = this.getTitle(data['bk_asst_id'])
            },
            setCreateInfo () {
                const edge = this.association.edge
                this.info.source = edge.from
                this.info.target = edge.to
                this.info.association = this.localAssociationList[0] || {}
                this.title = this.$t('新建关联')
            },
            async createAssociation () {
                try {
                    const isValid = await this.$validator.validateAll()
                    if (isValid) {
                        const { source, target, description, constraint, association } = this.info
                        const data = await this.createObjectAssociation({
                            params: this.$injectMetadata({
                                'bk_obj_id': source,
                                'bk_asst_obj_id': target,
                                'bk_asst_id': association['bk_asst_id'],
                                'bk_obj_asst_id': `${source}_${association['bk_asst_id']}_${target}`,
                                'bk_obj_asst_name': description,
                                'mapping': constraint
                            })
                        })
                        const createdAssociation = {
                            bk_asst_inst_id: association.id,
                            bk_asst_name: association['bk_asst_name'],
                            bk_asst_type: '',
                            bk_inst_id: data.id,
                            bk_obj_id: target,
                            label: (data.metadata || {}).label || {},
                            node_type: 'obj'
                        }
                        this.$store.commit('globalModels/addAssociation', {
                            id: source,
                            association: createdAssociation
                        })
                        this.addEdgePromise.resolve(createdAssociation)
                        this.isShow = false
                    }
                } catch (e) {
                    console.error(e)
                    this.addEdgePromise.reject(e)
                }
            },
            handleSave () {
                try {
                    if (!this.isViewMode) {
                        this.createAssociation()
                    } else {
                        this.isShow = false
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            handleDelete () {
                this.$bkInfo({
                    title: this.$t('确定删除关联关系?'),
                    confirmFn: async () => {
                        try {
                            const edge = this.association.edge
                            const associationId = edge.data['bk_inst_id']
                            await this.deleteObjectAssociation({
                                id: associationId
                            })
                            this.$store.commit('globalModels/deleteAssociation', associationId)
                            this.parentRefs.graphics.instance.deleteEdge(edge.id)
                            this.isShow = false
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            handleSelectAssociation (association) {
                this.info.association = association
            },
            handleCancel () {
                this.handleClose()
                this.isShow = false
            },
            handleClose () {
                const { reject } = this.addEdgePromise
                reject && reject(new Error(false))
            },
            reset () {
                this.info = { ...DEFAULT_INFO }
            },
            getTitle (asstId) {
                const data = this.localAssociationList.find(data => data['bk_asst_id'] === asstId) || {}
                return data['bk_asst_name']
            },
            getModelName (modelId) {
                const model = this.models.find(model => model['bk_obj_id'] === modelId) || {}
                return model['bk_obj_name']
            },
            getAssociationItemName (association) {
                return `${association['bk_asst_name']}(${association['bk_asst_id']})`
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association-layout {
        height: 100%;
        padding: 0 20px;
        @include scrollbar-y;
    }
    .form-group {
        position: relative;
        .form-label {
            display: block;
            margin: 15px 0 0 0;
            line-height: 36px;
            &.required:after {
                display: inline-block;
                padding: 4px 0 0 0;
                vertical-align: top;
                content: "*";
                line-height: 32px;
                color: $cmdbDangerColor;
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            font-size: 12px;
            color: $cmdbDangerColor;
        }
        .association-list {
            .association-item {
                height: 26px;
                padding: 0 8px;
                margin: 5px 7px 5px 0;
                border: 1px solid $cmdbBorderColor;
                line-height: 24px;
                font-size: 12px;
                cursor: pointer;
                background-color: #f5f7f9;
                &:hover {
                    background: #fafafa;
                }
                &.selected {
                    border-color: $cmdbBorderFocusColor;
                    background: $cmdbBorderFocusColor;
                    color: #fff;
                }
            }
        }
    }
    .button-group {
        margin: 30px 0 0 0;
        font-size: 0;
        .form-button {
            margin: 0 10px 0 0;
        }
    }
</style>
