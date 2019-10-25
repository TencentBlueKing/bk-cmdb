<template>
    <div :class="['table-stuff', type]">
        <div class="content" v-if="type === 'search'">
            <i class="bk-cc-icon icon-cc-tips"></i>
            <span class="text">搜索内容为空</span>
        </div>
        <div class="content" v-else-if="type === 'permission'">
            <i class="bk-cc-icon icon-cc-no-authority"></i>
            <div>
                <i18n path="抱歉您没有查看权限">
                    <bk-button
                        place="link"
                        text
                        theme="primary"
                        @click="handleApplyPermission"
                    >
                        {{$t('去申请')}}
                    </bk-button>
                </i18n>
            </div>
        </div>
        <div class="content" v-else>
            <img class="img-empty" src="../../../assets/images/empty-content.png" alt="">
            <div>
                <i18n path="您还未XXX" tag="div" v-if="!emptyText">
                    <span place="action">{{action}}</span>
                    <span place="resource">{{resource}}</span>
                    <bk-button
                        text
                        place="link"
                        theme="primary"
                        @click="$emit('create')"
                    >
                        <cmdb-auth :auth="authParams">
                            {{`立即${action}`}}
                        </cmdb-auth>
                    </bk-button>
                </i18n>
                <span v-else>
                    {{emptyText}}
                </span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapState, mapGetters } from 'vuex'
    import permissionMixins from '@/mixins/permission'
    export default {
        name: 'table-stuff',
        mixins: [permissionMixins],
        props: {
            stuff: {
                type: Object,
                default: () => ({
                    type: 'default',
                    payload: {}
                })
            },
            auth: {
                type: String
            }
        },
        data () {
            return {
                permission: this.stuff.payload.permission
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('auth', ['parentMeta']),
            type () {
                return this.stuff.type
            },
            action () {
                return this.stuff.payload.action || this.$t('创建')
            },
            resource () {
                return this.stuff.payload.resource
            },
            emptyText () {
                return this.stuff.payload.emptyText
            },
            payload () {
                return this.stuff.payload
            },
            authParams () {
                return {
                    resource_id: null,
                    bk_biz_id: this.bizId,
                    parent_layers: this.parentMeta.parent_layers,
                    type: this.auth
                }
            }
        },
        watch: {
            stuff: {
                handler (value) {
                    this.permission = value.payload.permission
                },
                deep: true
            }
        },
        methods: {
        }
    }
</script>

<style lang="scss" scoped>
    .table-stuff {
        color: #63656e;
        font-size: 14px;

        .icon-cc-no-authority {
            font-size: 90px;
        }

        .img-empty {
            width: 90px;
        }
    }
</style>
