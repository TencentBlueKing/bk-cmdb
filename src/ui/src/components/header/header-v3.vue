<template>
    <div class="header-wrapper clearfix" 
        :class="{'nav-sticked': navStick}">
        <div class="breadcrumbs fl">
            <i class="breadcrumbs-back bk-icon icon-arrows-left" href="javascript:void(0)"
                v-if="historyCount > 0"
                @click="back"></i>
            <h2 class="breadcrumbs-current">{{currentName}}</h2>
        </div>
        <div class="user fr">
            <p class="user-name" @click="isShowUserDropdown = !isShowUserDropdown">
                {{userName}}({{userRole}})
                <i class="user-name-angle bk-icon icon-angle-down"
                    :class="{dropped: isShowUserDropdown}">
                </i>
            </p>
            <transition name="toggle-slide">
                <ul class="user-dropdown" v-show="isShowUserDropdown">
                    <li class="user-dropdown-item" @click="logOut">
                        <i class="icon-cc-logout"></i>
                        {{$t("Common['注销']")}}
                    </li>
                </ul>
            </transition>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapMutations } from 'vuex'
    export default {
        data () {
            return {
                userRole: window.isAdmin === '1' ? this.$t('Common["管理员"]') : this.$t('Common["普通用户"]'),
                userName: window.userName,
                singleClassifies: ['/index'],
                isShowUserDropdown: false
            }
        },
        computed: {
            ...mapGetters('navigation', ['navStick', 'historyCount', 'authorizedNavigation']),
            ...mapGetters('usercustom', ['usercustom', 'recentlyKey']),
            currentPath () {
                return this.$route.path
            },
            // 根据当前路由路径获取路由名称
            currentName () {
                if (this.singleClassifies.includes(this.currentPath)) {
                    const classify = this.authorizedNavigation.find(navigation => navigation.path === this.currentPath)
                    return classify ? classify.i18n ? this.$t(classify.i18n) : classify.name : null
                } else {
                    const model = this.getRouteModel(this.currentPath)
                    return model ? model.i18n ? this.$t(model.i18n) : model.name : null
                }
            },
            // 最近浏览的通用模型
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
            // 回退路由
            back () {
                this.$store.commit('navigation/updateHistoryCount', -2)
                this.$router.back()
            },
            // 退出登陆
            logOut () {
                window.location.href = window.siteUrl + 'logout'
            },
            // 获取当前路由对应的模型
            getRouteModel (path) {
                let model
                for (let i = 0; i < this.authorizedNavigation.length; i++) {
                    const models = this.authorizedNavigation[i]['children'] || []
                    model = models.find(model => model.path === this.currentPath)
                    if (model) break
                }
                return model
            },
            // 更新最近浏览记录
            updateRecently (path) {
                const recently = this.recently.filter(oldPath => oldPath !== path)
                const model = this.getRouteModel(path)
                if (model && !['bk_host_manage', 'bk_back_config'].includes(model.classificationId)) {
                    this.$store.dispatch('usercustom/updateUserCustom', {
                        [this.recentlyKey]: [path, ...recently]
                    })
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
        transition: padding .1s ease-in;
        &.nav-sticked{
            padding-left: 240px;
        }
    }
    .breadcrumbs{
        line-height: 60px;
        position: relative;
        margin: 0 0 0 25px;
        font-size: 0;
        &-back{
            display: inline-block;
            vertical-align: middle;
            width: 24px;
            height: 24px;
            line-height: 24px;
            text-align: center;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            &:hover{
                color: #3c96ff;
            }
        }
        &-current{
            margin: 0;
            padding: 0;
            display: inline-block;
            vertical-align: middle;
            font-size: 16px;
            font-weight: normal;
        }
    }
    .user{
        font-size: 0;
        line-height: 60px;
        position: relative;
        &-name{
            padding: 0 20px;
            margin: 0;
            font-size: 14px;
            font-weight: bold;
            color: rgba(115,121,135,1);
            cursor: pointer;
            &-angle{
                display: inline-block;
                font-size: 12px;
                margin: 0 2px;
                color: $textColor;
                transition: transform .2s linear;
                &.dropped{
                    transform: rotate(-180deg);
                }
            }
        }
        &-dropdown{
            position: absolute;
            width: 100px;
            top: 55px;
            right: 20px;
            padding: 10px 0;
            line-height: 45px;
            font-size: 14px;
            background-color: #fff;
            box-shadow: 0 1px 5px 0 rgba(12,34,59, .1);
            z-index: 1;
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