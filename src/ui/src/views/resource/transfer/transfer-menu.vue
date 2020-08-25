<template>
    <div class="transfer-menu">
        <bk-dropdown-menu
            trigger="click"
            font-size="medium"
            :disabled="disabled"
            @show="handleMenuToggle(true)"
            @hide="handleMenuToggle(false)">
            <div class="dropdown-trigger-btn" style="padding-left: 19px;" slot="dropdown-trigger">
                <span>{{$t('转移到')}}</span>
                <i :class="['bk-icon icon-angle-down', { 'icon-flip': isShow }]"></i>
            </div>
            <ul class="bk-dropdown-list" slot="dropdown-content">
                <li><a href="javascript:;" @click="transferToIdleModule">{{$t('空闲模块')}}</a></li>
                <li><a href="javascript:;" @click="transferToBizModule">{{$t('业务模块')}}</a></li>
                <cmdb-auth tag="li">
                    <a href="javascript:;" @click="transferToResourcePool">{{$t('资源池')}}</a>
                </cmdb-auth>
            </ul>
        </bk-dropdown-menu>
        <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height">
            <component
                :is="dialog.component"
                v-bind="dialog.props"
                @cancel="handleDialogCancel"
                @confirm="handleDialogConfirm">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import HostStore from './host-store'
    import ModuleSelector from '../../business-topology/host/module-selector'
    export default {
        components: {
            [ModuleSelector.name]: ModuleSelector
        },
        data () {
            return {
                isShow: false,
                dialog: {
                    width: 720,
                    height: 460,
                    show: false,
                    component: null,
                    props: {}
                }
            }
        },
        computed: {
            disabled () {
                return !HostStore.isSelected
            }
        },
        methods: {
            handleMenuToggle (isShow) {
                this.isShow = isShow
            },
            validateSameBiz () {
                if (!HostStore.isSameBiz) {
                    this.$error(this.$t('该功能仅支持对相同业务下的主机进行操作'))
                    return false
                }
                return true
            },
            transferToIdleModule () {
                const valid = this.validateSameBiz()
                if (!valid) {
                    return false
                }
                const props = {
                    moduleType: 'idle',
                    title: this.$t('转移主机到空闲模块')
                }
                this.dialog.props = props
                this.dialog.width = 720
                this.dialog.height = 460
                this.dialog.component = ModuleSelector.name
                this.dialog.show = true
            },
            transferToBizModule () {
                const valid = this.validateSameBiz()
                if (!valid) {
                    return false
                }
            },
            transferToResourcePool () {
                const valid = this.validateSameBiz()
                if (!valid) {
                    return false
                }
                const isAllIdleModule = HostStore.isAllIdleModule
                if (!isAllIdleModule) {
                    this.$error(this.$t('改功能仅支持对空闲机模块下的主机进行操作'))
                    return false
                }
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDialogConfirm () {}
        }
    }
</script>

<style lang="scss" scoped>
    .transfer-menu {
        display: inline-block;
    }
    .dropdown-trigger-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #c4c6cc;
        height: 32px;
        min-width: 68px;
        border-radius: 2px;
        padding: 0 15px;
        color: #63656E;
        font-size: 14px;
    }
    .dropdown-trigger-btn.bk-icon {
        font-size: 18px;
    }
    .dropdown-trigger-btn .bk-icon {
        font-size: 22px;
    }
    .dropdown-trigger-btn:hover {
        cursor: pointer;
        border-color: #979ba5;
    }
</style>
