<template>
    <div class="permission-wrapper">
        <bk-tab class="permission-tab" :active-name.sync="activeTabName">
            <bk-tabpanel name="role" :title="$t('Permission[\'角色\']')">
                <v-role
                    ref="role"
                    @skipToUser="skipToUser"
                    v-if="activeTabName === 'role'"
                ></v-role>
            </bk-tabpanel>
            <bk-tabpanel name="authority" :title="$t('Permission[\'权限\']')">
                <v-authority
                    v-if="activeTabName === 'authority'"
                    :groupId="groupId"
                    @createRole="createRole"
                ></v-authority>
            </bk-tabpanel>
            <bk-tabpanel name="business" :title="$t('Permission[\'业务权限\']')">
                <v-business v-if="activeTabName === 'business'"></v-business>
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import vRole from './role'
    import vAuthority from './authority'
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
            vAuthority,
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
