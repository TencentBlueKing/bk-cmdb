<template>
    <div class="recently-layout clearfix">
        <template v-if="visibleRecentlyModels.length">
            <div class="recently-browse fl"
                v-for="index in recentlyCount"
                :key="index"
                :class="{'recently-browse-model': !!recentlyModels[index - 1]}"
                @click="gotoRecently(recentlyModels[index - 1])">
                <template v-if="recentlyModels[index - 1]">
                    <i :class="['recently-icon', recentlyModels[index - 1].icon]"></i>
                    <div class="recently-info">
                        <strong class="recently-name">{{getRecentlyName(recentlyModels[index - 1])}}</strong>
                        <span class="recently-inst">数量：{{getRecentlyCount(recentlyModels[index - 1])}}</span>
                    </div>
                    <i class="recently-delete bk-icon icon-close" @click.stop="deleteRecently(recentlyModels[index - 1])"></i>
                </template>
            </div>
        </template>
        <div class="recently-empty" v-else>{{$t("Index['暂无使用记录']")}}</div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                modelInstCount: {}
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            ...mapGetters('usercustom', ['usercustom', 'recentlyKey']),
            ...mapGetters('navigation', ['authorizedNavigation']),
            recently () {
                return this.usercustom[this.recentlyKey] || []
            },
            // 最近浏览的所有通用模型
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
            // 最近浏览的所有通用模型路由
            recentlyModelsPath () {
                return this.recentlyModels.map(model => model.path)
            },
            // 展示前8个最近浏览
            visibleRecentlyModels () {
                return this.recentlyModels.slice(0, 8)
            },
            // 最多展示的最近浏览的个数
            recentlyCount () {
                return this.visibleRecentlyModels.length > 4 ? 8 : 4
            }
        },
        watch: {
            // 最近浏览变更时，重新加载最近浏览模型的实例数量
            visibleRecentlyModels (models) {
                models.forEach(model => {
                    if (!this.modelInstCount.hasOwnProperty(model.id)) {
                        this.loadInst(model.id)
                    }
                })
            }
        },
        created () {
            // 首次加载最近浏览模型实例数量
            this.visibleRecentlyModels.forEach(model => {
                if (!this.modelInstCount.hasOwnProperty(model.id)) {
                    this.loadInst(model.id)
                }
            })
        },
        methods: {
            // 获取路由对应的模型
            getRouteModel (path) {
                let model
                for (let i = 0; i < this.authorizedNavigation.length; i++) {
                    const models = this.authorizedNavigation[i]['children'] || []
                    model = models.find(model => model.path === path)
                    if (model) break
                }
                return model
            },
            // 最近浏览模型的名称
            getRecentlyName (model) {
                return model.i18n ? this.$t(model.i18n) : model.name
            },
            // 最近浏览模型的实例数量
            getRecentlyCount (model) {
                return this.modelInstCount.hasOwnProperty(model.id) ? this.modelInstCount[model.id] : '--'
            },
            // 导航至最近浏览的模型
            gotoRecently (model) {
                if (model) {
                    this.$store.commit('navigation/updateHistoryCount', 2)
                    this.$router.push(model.path)
                }
            },
            // 删除最近浏览的模型
            deleteRecently (model) {
                const deletedRecently = this.recentlyModelsPath.filter(path => path !== model.path)
                this.$store.dispatch('usercustom/updateUserCustom', {
                    [this.recentlyKey]: deletedRecently
                })
            },
            // 加载最近浏览模型的实例数量
            loadInst (id) {
                const funcMaps = {
                    'biz': this.loadBizInst,
                    'default': this.loadCommonInst
                }
                let loadFunc = funcMaps.hasOwnProperty(id) ? funcMaps[id] : funcMaps['default']
                loadFunc(id).then(res => {
                    if (res.result) {
                        this.$set(this.modelInstCount, id, res.data.count)
                    }
                })
            },
            // 加载业务实例数量
            loadBizInst () {
                return this.$axios.post(`/biz/search/${this.bkSupplierAccount}`, {
                    condition: {
                        'bk_data_status': {
                            '$ne': 'disabled'
                        }
                    },
                    fields: [],
                    page: {
                        start: 0,
                        limit: 1
                    }
                })
            },
            // 加载通用模型实例数量
            loadCommonInst (id) {
                return this.$axios.post(`inst/association/search/owner/${this.bkSupplierAccount}/object/${id}`, {
                    condition: {},
                    fields: {},
                    page: {
                        start: 0,
                        limit: 1
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .recently-layout{
        width: 90%;
        margin: 0 auto;
    }
    .recently-browse{
        position: relative;
        width: 25%;
        height: 100px;
        margin: -1px 0 0 -1px;
        padding: 0 30px;
        font-size: 0;
        background-color: #fff;
        border: 1px solid #ebf0f5;
        white-space: nowrap;
        &-model{
            cursor: pointer;
            &:hover{
                border-color: #aaccff;
                z-index: 1;
                .recently-delete{
                    display: block;
                }
            }
        }
        &:nth-child(1),
        &:nth-child(5){
            margin-left: 0;
        }
        &:before{
            content: "";
            display: inline-block;
            width: 0;
            height: 100%;
            vertical-align: middle;
        }
        .recently-icon{
            display: inline-block;
            width: 42px;
            height: 42px;
            line-height: 44px;
            text-align: center;
            font-size: 24px;
            background-color: #4c84ff;
            color: #fff;
            border-radius: 50%;
        }
        .recently-info{
            display: inline-block;
            vertical-align: middle;
            margin-left: 16px;
            width: calc(100% - 42px - 16px);
            .recently-name{
                display: block;
                margin: 0 0 5px 0;
                font-size: 16px;
                color: #333c48;
                @include ellipsis;
            }
            .recently-inst{
                font-size: 12px;
                color: #9ba2b2;
            }
        }
        .recently-delete{
            display: none;
            position: absolute;
            right: 5px;
            top: 5px;
            width: 20px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            font-size: 12px;
            color: #c3cdd7;
            cursor: pointer;
            transform: scale(0.8);
            &:hover{
                color: #4c84ff;
            }
        }
    }
    .recently-empty{
        height: 100px;
        line-height: 96px;
        text-align: center;
        font-weight: bold;
        border: 2px dashed #dde4eb;
        background-color: rgba(51,60,72,0.02);
    }
</style>