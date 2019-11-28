<template>
    <section>
        <bk-input v-if="templates.length"
            class="search"
            type="text"
            :placeholder="$t('请输入模板名称搜索')"
            clearable
            right-icon="bk-icon icon-search"
            v-model.trim="searchName"
            @enter="hanldeFilterTemplates"
            @clear="hanldeFilterTemplates">
        </bk-input>
        <ul class="template-list clearfix"
            v-bkloading="{ isLoading: $loading('getServiceTemplate') }"
            :style="{ height: !!templates.length ? '264px' : '306px' }"
            :class="{ 'is-loading': $loading('getServiceTemplate') }">
            <template v-if="templates.length">
                <template v-for="(template, index) in templates">
                    <li v-if="$parent.$parent.serviceExistHost(template.id)"
                        class="template-item disabled fl clearfix"
                        :class="{
                            'is-selected': localSelected.includes(template.id),
                            'is-middle': index % 3 === 1
                        }"
                        :key="template.id"
                        v-bk-tooltips="$t('该模块下有主机不可取消')">
                        <i class="select-icon bk-icon icon-check-circle-shape fr"></i>
                        <span class="template-name" :title="template.name">{{template.name}}</span>
                    </li>
                    <li v-else
                        class="template-item fl clearfix"
                        :class="{
                            'is-selected': localSelected.includes(template.id),
                            'is-middle': index % 3 === 1
                        }"
                        :key="template.id"
                        @click="handleClick(template)">
                        <i class="select-icon bk-icon icon-check-circle-shape fr"></i>
                        <span class="template-name" :title="template.name">{{template.name}}</span>
                    </li>
                </template>
            </template>
            <li class="template-empty" v-else>
                <div class="empty-content">
                    <img class="empty-image" src="../../../assets/images/empty-content.png">
                    <i18n class="empty-tips" path="无服务模板提示">
                        <a class="empty-link" href="javascript:void(0)" place="link" @click="handleLinkClick">{{$t('跳转添加')}}</a>
                    </i18n>
                </div>
            </li>
        </ul>
    </section>
</template>

<script>
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
                searchName: ''
            }
        },
        created () {
            this.getTemplates()
        },
        methods: {
            async getTemplates () {
                try {
                    const data = await this.$store.dispatch('serviceTemplate/searchServiceTemplate', {
                        params: this.$injectMetadata({}),
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
            handleClick (template) {
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
                this.$router.push({
                    name: 'operationalTemplate'
                })
            },
            hanldeFilterTemplates () {
                this.templates = this.allTemplates.filter(template => template.name.indexOf(this.searchName) > -1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search {
        width: 240px;
        margin-bottom: 10px;
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
</style>
