<template>
    <div class="header-wrapper clearfix">
        <div class="breadcrumbs fl">
            <a class="breadcrumbs-back bk-icon icon-arrows-left" href="javascript:void(0)"
                v-if="historyCount > 1"
                @click="back"></a>
            <h2 class="breadcrumbs-current">{{currentName}}</h2>
        </div>
        <div class="user fr">
            <p class="user-name">
                {{userName}}
                <i class="user-name-angle bk-icon icon-angle-down"></i>
            </p>
            <ul class="user-dropdown">
                <li class="user-dropdown-item" @click="logOut">
                    <i class="icon-cc-logout"></i>
                    {{$t("Common['注销']")}}
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapMutations } from 'vuex'
    export default {
        data () {
            return {
                userName: window.userName,
                singleClassifies: ['/index']
            }
        },
        computed: {
            ...mapGetters('navigation', ['fold', 'historyCount', 'authorizedNavigation']),
            ...mapGetters('usercustom', ['usercustom', 'recentlyKey']),
            currentPath () {
                return this.$route.path
            },
            currentName () {
                if (this.singleClassifies.includes(this.currentPath)) {
                    const classify = this.authorizedNavigation.find(navigation => navigation.path === this.currentPath)
                    return classify ? classify.i18n ? this.$t(classify.i18n) : classify.name : null
                } else {
                    const model = this.getRouteModel(this.currentPath)
                    return model ? model.i18n ? this.$t(model.i18n) : model.name : null
                }
            },
            recently () {
                return this.usercustom[this.recentlyKey] || []
            }
        },
        watch: {
            currentPath (path) {
                this.updateRecently(path)
            }
        },
        methods: {
            ...mapMutations('navigation', ['updateHistoryCount']),
            back () {
                this.$store.commit('navigation/updateHistoryCount', -2)
                this.$router.back()
            },
            logOut () {
                window.location.href = window.siteUrl + 'logout'
            },
            getRouteModel (path) {
                let model
                for (let i = 0; i < this.authorizedNavigation.length; i++) {
                    const models = this.authorizedNavigation[i]['children'] || []
                    model = models.find(model => model.path === this.currentPath)
                    if (model) break
                }
                return model
            },
            updateRecently (path) {
                if (!this.recently.includes(path)) {
                    const model = this.getRouteModel(path)
                    if (model) {
                        this.$store.dispatch('usercustom/updateUserCustom', {
                            [this.recentlyKey]: [...this.recently, path]
                        })
                    }
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .header-wrapper{
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 61px;
        padding: 0 0 0 60px;
        border-bottom: 1px solid $borderColor;
        z-index: 1200;
        background-color: #fff;
    }
    .breadcrumbs{
        line-height: 60px;
        position: relative;
        margin: 0 0 0 25px;
        &-back{
            font-size: 12px;
            &:hover{
                color: #3c96ff;
            }
        }
        &-current{
            margin: 0;
            padding: 0;
            display: inline-block;
            font-size: 16px;
            font-weight: normal;
        }
    }
    .user{
        font-size: 0;
        line-height: 60px;
        position: relative;
        &-name{
            padding: 0;
            margin: 0 20px 0 0;
            font-size: 14px;
            font-weight: bold;
            color: rgba(115,121,135,1);
            cursor: pointer;
            &:hover ~ .user-dropdown{
                display: block;
            }
            &-angle{
                font-size: 12px;
                margin: 0 2px;
                color: $textColor;
            }
        }
        &-dropdown{
            display: none;
            position: absolute;
            width: 100px;
            top: 100%;
            right: 0;
            padding: 10px 0;
            line-height: 45px;
            font-size: 14px;
            background-color: #fff;
            box-shadow: 0 1px 5px 0 rgba(12,34,59, .1);
            z-index: 1;
            &:hover{
                display: block;
            }
            &-item{
                padding: 0 0 0 12px;
                cursor: pointer;
                &:hover{
                    background-color: #f1f7ff;
                    color: #498fe0;
                }
            }
        }
    }
</style>