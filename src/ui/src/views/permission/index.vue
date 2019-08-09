<template>
    <div class="permission-wrapper">
        <bk-tab class="permission-tab" type="unborder-card" :active.sync="activeTabName">
            <bk-tab-panel name="role" :label="$t('角色')">
                <v-role
                    ref="role"
                    @skipToUser="skipToUser"
                    v-if="activeTabName === 'role'"
                ></v-role>
            </bk-tab-panel>
            <bk-tab-panel name="business" :label="$t('业务权限')">
                <v-business v-if="activeTabName === 'business'"></v-business>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import vRole from './role'
    import vBusiness from './business'
    export default {
        data () {
            return {
                activeTabName: 'role',
                groupId: ''
            }
        },
        watch: {
            activeTabName (name) {
                if (name !== 'authority') {
                    this.groupId = ''
                }
            }
        },
        methods: {
            skipToUser (groupId) {
                this.activeTabName = 'authority'
                this.groupId = groupId
            },
            setRoles () {
                
            },
            createRole () {
                this.activeTabName = 'role'
                this.$nextTick(() => {
                    this.$refs.role.createRole()
                })
            }
        },
        components: {
            vRole,
            vBusiness
        }
    }
</script>

<style lang="scss" scoped>
    .permission-wrapper{
        padding: 0 0 20px;
        height: 100%;
    }
</style>

<style lang="scss">
    .permission-tab.bk-tab2{
        height: 100%;
        border: none;
        .bk-tab2-head{
            .bk-tab2-nav {
                .tab2-nav-item{
                    text-align: center;
                }
            }
        }
        .bk-tab2-content{
            @include scrollbar;
            >section{
                height: 100%;
            }
        }
    }
</style>
