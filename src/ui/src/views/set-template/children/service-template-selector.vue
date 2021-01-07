<template>
    <section>
        <div class="top" v-if="allTemplates.length">
            <bk-input
                class="search"
                type="text"
                :placeholder="$t('请输入模板名称搜索')"
                clearable
                right-icon="bk-icon icon-search"
                v-model.trim="searchName"
                @enter="handleFilterTemplates"
                @clear="handleFilterTemplates">
            </bk-input>
            <span class="to-template" @click="handleLinkClick">
                <i class="icon-cc-share"></i>
                {{$t('跳转服务模板')}}
            </span>
            <span class="select-all fr" v-if="$parent.$parent.mode !== 'edit'">
                <bk-checkbox :value="isSelectAll" @change="handleSelectAll">全选</bk-checkbox>
            </span>
        </div>
        <ul class="template-list clearfix"
            v-bkloading="{ isLoading: $loading('getServiceTemplate') }"
            :style="{ height: !!templates.length ? '264px' : '306px' }"
            :class="{ 'is-loading': $loading('getServiceTemplate') }">
            <template v-if="templates.length">
                <template v-for="(template, index) in templates">
                    <li
                        class="template-item fl clearfix"
                        :class="{
                            'is-selected': localSelected.includes(template.id),
                            'is-middle': index % 3 === 1,
                            'disabled': $parent.$parent.serviceExistHost(template.id)
                        }"
                        :key="template.id"
                        @click="handleClick(template, $parent.$parent.serviceExistHost(template.id))"
                        @mouseenter="handleShowDetails(template, $event, $parent.$parent.serviceExistHost(template.id))"
                        @mouseleave="handlehideTips">
                        <i class="select-icon bk-icon icon-check-circle-shape fr"></i>
                        <span class="template-name">{{template.name}}</span>
                    </li>
                </template>
            </template>
            <li class="template-empty" v-else>
                <div class="empty-content">
                    <img class="empty-image" src="../../../assets/images/empty-content.png">
                    <i18n class="empty-tips" path="无服务模板提示">
                        <a class="empty-link" href="javascript:void(0)" place="link" @click="handleLinkClick">{{$t('去添加服务模板')}}</a>
                    </i18n>
                </div>
            </li>
        </ul>
        <div ref="templateDetails"
            class="template-details"
            v-bkloading="{ isLoading: $loading(processRequestId) }"
            v-show="tips.show">
            <div class="disabled-tips" v-show="processInfo.disabled">{{$t('该模块下有主机不可取消')}}</div>
            <div class="info-item">
                <span class="label">{{$t('模板名称')}} ：</span>
                <div class="details">{{curTemplate.name}}</div>
            </div>
            <div class="info-item">
                <span class="label">{{$t('服务分类')}} ：</span>
                <div class="details">{{processInfo.cagetory}}</div>
            </div>
            <div class="info-item">
                <span class="label">{{$t('服务进程')}} ：</span>
                <div class="details">
                    <p v-for="(item, index) in processInfo.processes" :key="index">{{item}}</p>
                    <template v-if="!processInfo.processes.length">
                        <p>{{$t('模板没配置进程')}}</p>
                    </template>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
    import { MENU_BUSINESS_SERVICE_TEMPLATE } from '@/dictionary/menu-symbol'
    import { mapGetters } from 'vuex'
    export default {
        name: 'serviceTemplateSelector',
        props: {
            selected: {
                type: Array,
                default: () => []
            },
            servicesHost: {
                type: Array,
                default: () => []
            }
        },
        data () {
            return {
                allTemplates: [],
                templates: [],
                localSelected: [...this.selected],
                searchName: '',
                templateDetailsData: {},
                processRequestId: Symbol('processDetails'),
                tips: {
                    show: false,
                    instance: null
                },
                curTemplate: {},
                cagetory: [],
                processInfo: {
                    disabled: false,
                    cagetory: '',
                    processes: []
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            isSelectAll () {
                return this.localSelected.length === this.allTemplates.length
            }
        },
        async created () {
            this.getTemplates()
            await this.getServiceCategory()
        },
        methods: {
            async getTemplates () {
                try {
                    const data = await this.$store.dispatch('serviceTemplate/searchServiceTemplate', {
                        params: { bk_biz_id: this.bizId },
                        config: {
                            requestId: 'getServiceTemplate'
                        }
                    })
                    this.templates = data.info.map(datum => datum.service_template).sort((A, B) => {
                        return A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' })
                    }).sort((A, B) => {
                        const weightA = this.selected.includes(A.id) ? 1 : 0
                        const weightB = this.selected.includes(B.id) ? 1 : 0
                        return weightB - weightA
                    })
                    this.allTemplates = this.templates
                } catch (e) {
                    console.error(e)
                    this.templates = []
                }
            },
            handleClick (template, disabled) {
                if (disabled) return
                const index = this.localSelected.indexOf(template.id)
                if (index > -1) {
                    this.localSelected.splice(index, 1)
                } else {
                    this.localSelected.push(template.id)
                }
            },
            getSelectedServices () {
                return this.localSelected.map(id => this.allTemplates.find(template => template.id === id))
            },
            handleLinkClick () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_SERVICE_TEMPLATE
                })
            },
            handleFilterTemplates () {
                this.templates = this.allTemplates.filter(template => template.name.indexOf(this.searchName) > -1)
            },
            handleSelectAll (checked) {
                if (checked) {
                    this.localSelected = this.allTemplates.map(template => template.id)
                } else {
                    this.localSelected = []
                }
            },
            async handleShowDetails (template = {}, event, disabled) {
                this.curTemplate = template
                this.processInfo.disabled = disabled
                const curInfo = this.templateDetailsData[template.id]
                if (curInfo) {
                    this.setProcessInfo(curInfo, event)
                    return
                }
                try {
                    const data = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: {
                            bk_biz_id: this.bizId,
                            service_template_id: template.id
                        },
                        config: {
                            requestId: this.processRequestId,
                            cancelPrevious: true
                        }
                    })
                    this.setProcessInfo(data.info, event)
                    this.templateDetailsData[template.id] = data.info
                } catch (e) {
                    console.error(e)
                }
            },
            setProcessInfo (data = [], event) {
                this.processInfo.processes = data.map(process => {
                    const port = process.property
                        ? process.property.port ? process.property.port.value : ''
                        : ''
                    return `${process.bk_process_name}${port ? `:${port}` : ''}`
                })
                const subCagetory = this.cagetory.find(item => item.id === this.curTemplate.service_category_id) || {}
                const cagetory = this.cagetory.find(item => item.id === subCagetory.bk_parent_id) || {}
                this.processInfo.cagetory = subCagetory && cagetory
                    ? `${cagetory.name} / ${subCagetory.name}`
                    : '-- / --'
                this.tips.instance && this.tips.instance.destroy()
                this.tips.instance = this.$bkPopover(event.target, {
                    content: this.$refs.templateDetails,
                    delay: 300,
                    zIndex: 9999,
                    width: 'auto',
                    trigger: 'manual',
                    boundary: 'window',
                    arrow: true
                })
                this.tips.show = true
                this.$nextTick(() => {
                    this.tips.instance.show()
                })
            },
            handlehideTips () {
                this.tips.instance && this.tips.instance.destroy()
                this.tips.instance = null
            },
            async getServiceCategory () {
                try {
                    const data = await this.$store.dispatch('serviceClassification/searchServiceCategoryWithoutAmout', {
                        params: { bk_biz_id: this.bizId },
                        config: {
                            requestId: 'getServiceCategoryWithoutAmount'
                        }
                    })
                    this.cagetory = data.info
                } catch (e) {
                    console.error(e)
                    this.cagetory = []
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .top {
        margin-bottom: 10px;
        .search {
            @include inlineBlock;
            width: 240px;
        }
        .to-template {
            @include inlineBlock;
            color: #3A84FF;
            margin-left: 10px;
            cursor: pointer;
            .icon-cc-share {
                margin-top: -2px;
            }
        }
        .select-all {
            @include inlineBlock;
            line-height: 32px;
        }
    }
    .template-list {
        height: 264px;
        @include scrollbar-y;
        &.is-loading {
            min-height: 144px;
        }
        .template-item {
            width: calc((100% - 20px) / 3);
            height: 32px;
            margin: 0 0 16px 0;
            padding: 0 6px 0 10px;
            line-height: 30px;
            border-radius: 2px;
            border: 1px solid #DCDEE5;
            color: #63656E;
            cursor: pointer;
            &.is-middle {
                margin: 0 10px 16px;
            }
            &.is-selected {
                background-color: #E1ECFF;
                .select-icon {
                    font-size: 18px;
                    border: none;
                    border-radius: initial;
                    background-color: initial;
                    color: #3A84FF;
                }
            }
            &.disabled {
                cursor: not-allowed;
                .select-icon {
                    color: #C4C6CC;
                    cursor: not-allowed;
                }
            }
            .select-icon {
                width: 18px;
                height: 18px;
                font-size: 0px;
                margin: 6px 0;
                color: #fff;
                background-color: #fff;
                border-radius: 50%;
                border: 1px solid #979BA5;
            }
            .template-name {
                display: block;
                max-width: calc(100% - 18px);
                @include ellipsis;
            }
            .bk-tooltip {
                width: 100%;
                /deep/ .bk-tooltip-ref {
                    width: 100%;
                }
            }
        }
    }
    .template-empty {
        height: 280px;
        &:before {
            content: "";
            height: 100%;
            width: 0;
            @include inlineBlock;
        }
        .empty-content {
            width: 100%;
            @include inlineBlock;
            .empty-image {
                display: block;
                margin: 0 auto;
            }
            .empty-tips {
                display: block;
                text-align: center;
            }
            .empty-link {
                color: #3A84FF;
            }
        }
    }
    .template-details {
        line-height: 20px;
        .disabled-tips {
            border-bottom: 1px solid #FFFFFF;
            padding: 0 0 6px;
            margin-bottom: 6px;
            font-size: 12px;
        }
        .info-item {
            display: flex;

            .label {
                font-size: 12px;
                font-weight: 700;
            }
        }
        .details {
            font-size: 12px;
        }
    }
</style>
