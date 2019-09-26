<template>
    <section>
        <cmdb-tips class="tips">{{$t('服务模板选择器提示语')}}</cmdb-tips>
        <ul class="template-list clearfix"
            v-bkloading="{ isLoading: $loading('getServiceTemplate') }"
            :class="{ 'is-loading': $loading('getServiceTemplate') }">
            <li class="template-item fl clearfix"
                v-for="(template, index) in templates"
                :class="{
                    'is-selected': selected.includes(template.id),
                    'is-middle': index % 3 === 1
                }"
                :key="template.id"
                @click="handleClick(template)">
                <i class="select-icon bk-icon icon-check-circle-shape fr"></i>
                <span class="template-name" :title="template.name">{{template.name}}</span>
            </li>
        </ul>
    </section>
</template>

<script>
    export default {
        name: 'serviceTemplateSelector',
        data () {
            return {
                visible: false,
                templates: [],
                selected: []
            }
        },
        created () {
            this.getTemplates()
        },
        methods: {
            show (selected) {
                this.selected = selected
                this.visible = true
            },
            async getTemplates () {
                try {
                    const data = await this.$store.dispatch('serviceTemplate/searchServiceTemplate', {
                        params: this.$injectMetadata({}),
                        config: {
                            requestId: 'getServiceTemplate'
                        }
                    })
                    this.templates = data.info.map(datum => datum.service_template)
                } catch (e) {
                    console.error(e)
                    this.templates = []
                }
            },
            handleClick (template) {
                const index = this.selected.indexOf(template.id)
                if (index > -1) {
                    this.selected.splice(index, 1)
                } else {
                    this.selected.push(template.id)
                }
            },
            getSelectedServices () {
                return this.templates.filter(template => this.selected.includes(template.id))
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tips {
        margin-top: -4px;
    }
    .template-list {
        max-height: 340px;
        @include scrollbar-y;
        &.is-loading {
            min-height: 144px;
        }
        .template-item {
            width: calc((100% - 20px) / 3);
            height: 32px;
            margin: 16px 0 0 0;
            padding: 0 6px 0 10px;
            line-height: 30px;
            border-radius: 2px;
            border: 1px solid #DCDEE5;
            color: #63656E;
            cursor: pointer;
            &.is-middle {
                margin: 16px 10px 0;
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
</style>
