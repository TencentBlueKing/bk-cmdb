<template>
    <div class="recently-layout clearfix">
        <template v-if="recently.length">
            <div class="recently-browse fl"
                v-for="index in recentlyCount"
                :key="index"
                @click="gotoRecently(recentlyModels[index - 1])">
                <template v-if="recentlyModels[index - 1]">
                    <i :class="['recently-icon', recentlyModels[index - 1].icon]"></i>
                    <div class="recently-info">
                        <strong class="recently-name">{{getRecentlyName(recentlyModels[index - 1])}}</strong>
                        <span class="recently-inst">数量：21</span>
                    </div>
                    <i class="recently-delete" @click.stop="deleteRecently(recentlyModels[index - 1])"></i>
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
            return {}
        },
        computed: {
            ...mapGetters('usercustom', ['usercustom', 'recentlyKey']),
            ...mapGetters('navigation', ['authorizedNavigation']),
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
            recentlyModelsPath () {
                return this.recentlyModels.map(model => model.path)
            },
            recentlyCount () {
                return this.recentlyModels.length > 4 ? 8 : 4
            }
        },
        methods: {
            getRouteModel (path) {
                let model
                for (let i = 0; i < this.authorizedNavigation.length; i++) {
                    const models = this.authorizedNavigation[i]['children'] || []
                    model = models.find(model => model.path === path)
                    if (model) break
                }
                return model
            },
            getRecentlyName (model) {
                return model.i18n ? this.$t(model.i18n) : model.name
            },
            gotoRecently (model) {
                this.$router.push(model.path)
            },
            deleteRecently (model) {
                const deletedRecently = this.recentlyModelsPath.filter(path => path !== model.path)
                this.$store.dispatch('usercustom/updateUserCustom', {
                    [this.recentlyKey]: deletedRecently
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .recently-layout{
        width: calc(90% + 6px);
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
        &:hover{
            border-color: #aaccff;
            z-index: 1;
            .recently-delete{
                display: block;
            }
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
            .recently-name{
                display: block;
                margin: 0 0 5px 0;
                font-size: 16px;
                color: #333c48;
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
            width: 10px;
            height: 10px;
            cursor: pointer;
            &:hover:before,
            &:hover:after{
                border-color: #aaccff;
            }
            &:before,
            &:after{
                content: "";
                position: absolute;
                left: -2px;
                top: 5px;
                width: 14px;
                border-top: 1px solid #c3cdd7;
            }
            &:before{
                transform: rotate(-45deg);
            }
            &:after{
                transform: rotate(45deg);
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