<template>
    <div class="api-wrapper">
        <div class="filter-wrapper clearfix">
            <cmdb-business-selector
                class="business-selector"
                v-model="filter.bizId"
            ></cmdb-business-selector>
            <bk-button type="primary" class="api-btn" @click="showUserAPISlider('create')">
                {{$t("CustomQuery['新增查询']")}}
            </bk-button>
            <div class="api-input fr">
                <input type="text" class="cmdb-form-input" :placeholder="$t('Inst[\'快速查询\']')" v-model="filter.name" @keyup.enter="getUserAPIList">
            </div>
        </div>
        <cmdb-table
        class="api-table"
        :loading="$loading('searchCustomQuery')"
        :header="table.header"
        :list="table.list"
        :pagination.sync="table.pagination"
        :wrapperMinusHeight="220"
        @handlePageChange="handlePageChange"
        @handleSizeChange="handleSizeChange"
        @handleSortChange="handleSortChange"
        @handleRowClick="showUserAPIDetails">
            <template slot="create_time" slot-scope="{item}">
                {{$tools.formatTime(item['create_time'])}}
            </template>
            <template slot="last_time" slot-scope="{item}">
                {{$tools.formatTime(item['last_time'])}}
            </template>
            <div class="empty-info" slot="data-empty">
                <p>{{$t("Common['暂时没有数据']")}}</p>
                <p>{{$t("CustomQuery['当前业务并无自定义查询，可点击下方按钮新增']")}}</p>
                <bk-button class="process-btn" type="primary" @click="showUserAPISlider('create')">{{$t("CustomQuery['新增查询']")}}</bk-button>
            </div>
        </cmdb-table>
        <cmdb-slider
            :isShow.sync="slider.isShow"
            :hasQuickClose="true"
            :width="430"
            :title="slider.title"
            :beforeClose="handleSliderBeforeClose">
            <v-define slot="content"
                ref="define"
                :id="slider.id"
                :bizId="filter.bizId"
                :type="slider.type"
                @delete="getUserAPIList"
                @create="handleCreate"
                @update="getUserAPIList"
                @cancel="hideUserAPISlider">
            </v-define>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import vDefine from './define'
    export default {
        components: {
            vDefine
        },
        data () {
            return {
                filter: {
                    bizId: '',
                    name: ''
                },
                table: {
                    header: [{
                        id: 'id',
                        name: 'ID'
                    }, {
                        id: 'name',
                        name: this.$t("CustomQuery['查询名称']")
                    }, {
                        id: 'create_user',
                        name: this.$t("CustomQuery['创建用户']")
                    }, {
                        id: 'create_time',
                        name: this.$t("CustomQuery['创建时间']")
                    }, {
                        id: 'modify_user',
                        name: this.$t("CustomQuery['修改人']")
                    }, {
                        id: 'last_time',
                        name: this.$t("CustomQuery['修改时间']")
                    }],
                    list: [],
                    sort: '-last_time',
                    defaultSort: '-last_time',
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    }
                },
                slider: {
                    isShow: false,
                    isCloseConfirmShow: false,
                    type: 'create',
                    id: null,
                    title: this.$t("CustomQuery['新增查询']")
                }
            }
        },
        computed: {
            searchParams () {
                let params = {
                    start: (this.table.pagination.current - 1) * this.table.pagination.size,
                    limit: this.table.pagination.size,
                    sort: this.table.sort
                }
                this.filter.bkBizId ? params['bk_biz_id'] = this.filter.bizId : void (0)
                this.filter.name ? params['condition'] = {'name': this.filter.name} : void (0)
                return params
            }
        },
        watch: {
            'filter.bizId' () {
                this.getUserAPIList()
            }
        },
        methods: {
            ...mapActions('hostCustomApi', [
                'searchCustomQuery'
            ]),
            hideUserAPISlider () {
                this.slider.isShow = false
                this.slider.id = null
            },
            handleSliderBeforeClose () {
                if (this.$refs.define.isCloseConfirmShow()) {
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
            handleCreate (data) {
                this.slider.id = data['id']
                this.slider.type = 'update'
                this.slider.title = this.$t('CustomQuery["编辑查询"]')
                this.handlePageChange(1)
            },
            async getUserAPIList () {
                const res = await this.searchCustomQuery({
                    bizId: this.filter.bizId,
                    params: this.searchParams,
                    config: {
                        requestId: 'searchCustomQuery'
                    }
                })
                if (res.count) {
                    this.table.list = res.info
                } else {
                    this.table.list = []
                }
                this.table.pagination.count = res.count
            },
            showUserAPISlider (type) {
                this.slider.isShow = true
                this.slider.type = type
                this.slider.title = this.$t('CustomQuery["新增查询"]')
            },
            /* 显示自定义API详情 */
            showUserAPIDetails (userAPI) {
                this.slider.isShow = true
                this.slider.type = 'update'
                this.slider.id = userAPI['id']
                this.slider.title = this.$t('CustomQuery["编辑查询"]')
            },
            handlePageChange (current) {
                this.table.pagination.current = current
                this.getUserAPIList()
            },
            handleSizeChange (size) {
                this.table.pagination.size = size
                this.handlePageChange(1)
            },
            handleSortChange (sort) {
                this.table.sort = sort
                this.getUserAPIList()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .api-wrapper {
        .filter-wrapper {
            .business-selector {
                float: left;
                width: 170px;
                margin-right: 10px;
            }
            .api-btn {
                float: left;
            }
            .api-input {
                float: right;
                width: 320px;
            }
        }
        .api-table {
            margin-top: 20px;
        }
    }
</style>

