<template>
    <header class="header-layout clearfix"
        :class="{ 'nav-sticked': navStick }">
        <div class="breadcrumbs fl">
            <i class="breadcrumbs-back icon-cc-arrow" href="javascript:void(0)"
                v-if="showBack"
                @click="back"></i>
            <h2 class="breadcrumbs-current">{{headerTitle}}</h2>
        </div>
        <div class="header-options">
            <cmdb-business-selector
                class="business-selector"
                v-if="!isAdminView">
            </cmdb-business-selector>
            <div class="user" v-click-outside="handleCloseUser">
                <p class="user-name" @click="isShowUserDropdown = !isShowUserDropdown">
                    {{userName}}
                    <i class="user-name-angle bk-icon icon-angle-down"
                        :class="{ dropped: isShowUserDropdown }">
                    </i>
                </p>
                <transition name="toggle-slide">
                    <ul class="user-dropdown" v-show="isShowUserDropdown">
                        <li class="user-dropdown-item" @click="logOut">
                            <i class="icon-cc-logout"></i>
                            {{$t('注销')}}
                        </li>
                    </ul>
                </transition>
            </div>
            <div class="helper" v-click-outside="handleCloseHelper">
                <i class="helper-icon bk-icon icon-question-circle" @click="isShowHelper = !isShowHelper"></i>
                <div class="helper-list" v-show="isShowHelper">
                    <a href="http://docs.bk.tencent.com/product_white_paper/cmdb/" target="_blank" class="helper-link"
                        @click="isShowHelper = false">
                        {{$t('帮助文档')}}
                    </a>
                    <a href="https://github.com/Tencent/bk-cmdb" target="_blank" class="helper-link"
                        @click="isShowHelper = false">
                        {{$t('开源社区')}}
                    </a>
                </div>
            </div>
            <bk-popover
                class="admin-tooltips"
                v-if="hasAdminEntrance && !isAdminView && showTips"
                :always="true"
                :width="275"
                theme="custom-color"
                placement="bottom-end">
                <div slot="content" class="tooltips-main clearfix">
                    <h3>{{$t('管理员后台提示')}}</h3>
                    <p>{{$t('管理员后台描述')}}</p>
                    <span class="fr" @click="handleCloseTips">{{$t('我知道了')}}</span>
                </div>
                <div class="admin" @click="toggleAdminView">
                    {{isAdminView ? $t('返回业务管理') : $t('管理员后台')}}
                </div>
            </bk-popover>
            <div class="admin" v-else-if="hasAdminEntrance" @click="toggleAdminView">
                {{isAdminView ? $t('返回业务管理') : $t('管理员后台')}}
            </div>
        </div>
    </header>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                isShowUserDropdown: false,
                isShowHelper: false
            }
        },
        computed: {
            ...mapGetters([
                'site',
                'userName',
                'admin',
                'navStick',
                'headerTitle',
                'isAdminView',
                'featureTipsParams'
            ]),
            ...mapGetters('objectBiz', ['authorizedBusiness']),
            hasAdminEntrance () {
                return this.$store.state.auth.adminEntranceAuth.is_pass
            },
            userRole () {
                return this.admin ? this.$t('管理员') : this.$t('普通用户')
            },
            showTips () {
                return this.featureTipsParams['adminTips']
            },
            showBack () {
                return this.$route.query.from
            }
        },
        methods: {
            toggleAdminView () {
                this.$store.commit('setAdminView', !this.isAdminView)
            },
            // 回退路由
            back () {
                this.$router.push(this.$route.query.from)
            },
            // 退出登陆
            logOut () {
                this.$http.post(`${window.API_HOST}logout`, {
                    'http_scheme': window.location.protocol.replace(':', '')
                }).then(data => {
                    window.location.href = data.url
                })
            },
            handleCloseUser () {
                this.isShowUserDropdown = false
            },
            handleCloseHelper () {
                this.isShowHelper = false
            },
            handleCloseTips () {
                this.$store.commit('setFeatureTipsParams', 'adminTips')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .header-layout{
        position: relative;
        height: 61px;
        padding: 0 0 0 60px;
        border-bottom: 1px solid $cmdbBorderColor;
        background-color: #fff;
        transition: padding .1s ease-in;
        z-index: 1000;
        &.nav-sticked{
            padding-left: 260px;
        }
    }
    .breadcrumbs{
        line-height: 60px;
        position: relative;
        margin: 0 0 0 12px;
        font-size: 0;
        .breadcrumbs-back{
            display: inline-block;
            vertical-align: middle;
            width: 32px;
            height: 32px;
            line-height: 32px;
            text-align: center;
            font-size: 16px;
            cursor: pointer;
            color: #3c96ff;
            transition: background-color .1s ease-in;
            &:hover{
                background-color: #f0f1f5;
            }
        }
        .breadcrumbs-current{
            margin: 0 0 0 8px;
            padding: 0;
            display: inline-block;
            vertical-align: middle;
            font-size: 16px;
            font-weight: normal;
            color: #313238;
        }
        .icon-info-circle {
            margin-left: 5px;
            font-size: 16px;
        }
    }
    .header-options {
        white-space: nowrap;
        text-align: right;
        font-size: 0;
    }
    .business-selector {
        display: inline-block;
        width: 200px;
        margin: 14px 0 0 20px;
        vertical-align: top;
    }
    .user{
        display: inline-block;
        vertical-align: top;
        font-size: 0;
        line-height: 60px;
        position: relative;
        .user-name{
            margin: 0 5px 0 20px;
            font-size: 14px;
            font-weight: bold;
            color: rgba(115,121,135,1);
            cursor: pointer;
            .user-name-angle{
                display: inline-block;
                font-size: 12px;
                margin: 0 2px;
                color: $cmdbTextColor;
                transition: transform .2s linear;
                &.dropped{
                    transform: rotate(-180deg);
                }
            }
        }
        .user-dropdown{
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
            .user-dropdown-item{
                padding: 0 0 0 12px;
                text-align: left;
                cursor: pointer;
                &:hover{
                    background-color: #f1f7ff;
                    color: #498fe0;
                }
            }
        }
    }
    .helper {
        position: relative;
        display: inline-block;
        width: 50px;
        text-align: center;
        vertical-align: top;
        line-height: 60px;
        .helper-icon {
            font-size: 20px;
            cursor: pointer;
            &:hover {
                color: #0082ff;
            }
        }
        .helper-list {
            position: absolute;
            top: 55px;
            right: 1px;
            text-align: left;
            line-height: 40px;
            background-color: #fff;
            border-radius: 2px;
            box-shadow: 0 1px 5px 0 rgba(12,34,59, .1);
            .helper-link {
                display: block;
                padding: 0 20px;
                font-size: 14px;
                white-space: nowrap;
                &:hover {
                    background-color: #f1f7ff;
                    color: #498fe0;
                }
            }
        }
    }
    .admin {
        display: inline-block;
        padding: 0 30px;
        line-height: 60px;
        font-size: 14px;
        color: #3a84ff;
        border-left: 1px solid #ebf0f5;
        cursor: pointer;
        text-align: center;
        vertical-align: top;
        &:hover {
            background-color: #f7f7f7;
        }
    }
</style>

<style lang="scss">
    .custom-color-theme {
        font-size: 14px !important;
        background-color: #699df4 !important;
        padding: 10px 14px !important;
        .tippy-arrow {
            border-bottom-color: #699df4 !important;
        }
        h3 {
            font-size: 16px;
        }
        p {
            white-space: pre-wrap;
            padding: 4px 0 6px;
        }
        span {
            font-size: 12px;
            padding: 4px 10px;
            background-color: #5d90e4;
            border-radius: 20px;
            cursor: pointer;
            &:hover {
                background-color: #477ad0;
            }
        }
    }
</style>
