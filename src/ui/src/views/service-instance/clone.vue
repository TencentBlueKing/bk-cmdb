<template>
    <div class="clone-layout">
        <div class="host-type clearfix">
            <label class="type-label fl">{{$t('克隆到')}}</label>
            <div class="type-item fl">
                <input class="type-radio"
                    type="radio"
                    id="sourceHost"
                    name="hostTarget"
                    v-model="hostTarget"
                    :value="targetName.source">
                <label for="sourceHost">{{$t('当前主机')}}</label>
            </div>
            <div class="type-item fl">
                <input class="type-radio"
                    type="radio"
                    id="otherHost"
                    name="hostTarget"
                    v-model="hostTarget"
                    :value="targetName.other">
                <label for="otherHost">{{$t('其他主机')}}</label>
            </div>
        </div>
        <component :is="hostTarget"
            :source-processes="processes"
            :module="module">
        </component>
    </div>
</template>

<script>
    import cloneToSource from './children/clone-to-source.vue'
    import cloneToOther from './children/clone-to-other.vue'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        components: {
            [cloneToSource.name]: cloneToSource,
            [cloneToOther.name]: cloneToOther
        },
        data () {
            return {
                hostTarget: cloneToSource.name,
                targetName: {
                    source: cloneToSource.name,
                    other: cloneToOther.name
                },
                module: {},
                processes: []
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            setId () {
                return parseInt(this.$route.params.setId)
            },
            moduleId () {
                return parseInt(this.$route.params.moduleId)
            },
            instanceId () {
                return parseInt(this.$route.params.instanceId)
            }
        },
        async created () {
            try {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('服务拓扑'),
                    route: {
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: {
                            node: 'module-' + this.$route.params.moduleId
                        }
                    }
                }, {
                    label: this.$route.query.title
                }])
                const [module, processes] = await Promise.all([
                    this.getModuleInstance(),
                    this.getServiceInstanceProcesses()
                ])
                this.module = module.info[0]
                this.processes = processes.map(instance => instance.property)
            } catch (e) {
                console.error(e)
            }
        },
        methods: {
            getModuleInstance () {
                return this.$store.dispatch('objectModule/searchModule', {
                    bizId: this.business,
                    setId: this.setId,
                    params: {
                        page: { start: 0, limit: 1 },
                        fields: [],
                        condition: {
                            bk_module_id: this.moduleId,
                            bk_supplier_account: this.$store.getters.supplierAccount
                        }
                    },
                    config: {
                        requestId: 'getModuleInstance'
                    }
                })
            },
            getServiceInstanceProcesses () {
                return this.$store.dispatch('processInstance/getServiceInstanceProcesses', {
                    params: this.$injectMetadata({
                        service_instance_id: this.instanceId
                    }, { injectBizId: true }),
                    config: {
                        requestId: 'getServiceInstanceProcesses'
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .clone-layout {
        padding: 10px 20px 28px;
        font-size: 14px;
    }
    .host-type {
        line-height: 19px;
        .type-label{
            width: 100px;
            &:after {
                content: '*';
                color: $cmdbDangerColor;
            }
        }
        .type-item {
            margin: 0 40px 0 15px;
            font-size: 0;
            .type-radio {
                -webkit-appearance: none;
                width: 16px;
                height: 16px;
                padding: 3px;
                border: 1px solid #979BA5;
                border-radius: 50%;
                background-clip: content-box;
                outline: none;
                cursor: pointer;
                @include inlineBlock;
                &:checked {
                    border-color: #3A84FF;
                    background-color: #3A84FF;
                }
            }
            label {
                padding: 0 0 0 6px;
                font-size: 14px;
                cursor: pointer;
                @include inlineBlock;
            }
        }
    }
</style>
