<template>
    <div class="custom-layout">
        <div class="input-wrapper">
            <textarea class="ip-input"
                :placeholder="$t('请输入IP，换行分隔')"
                :class="{ 'has-error': invalidList.length }"
                v-model="value"
                @focus="handleFocus">
            </textarea>
            <p class="ip-error" v-show="invalidList.length">{{$t('IP不正确或主机不存在')}}</p>
            <bk-button class="ip-confirm" outline theme="primary" @click="handleConfirm">{{$t('添加至列表')}}</bk-button>
        </div>
        <div class="table-wrapper" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
            <host-table :list="hostList" :selected="selected" @select-change="handleHostSelectChange" />
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import HostTable from './host-table.vue'
    export default {
        components: {
            HostTable
        },
        props: {
            selected: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                value: '',
                validList: [],
                invalidList: [],
                hostList: [],
                request: {
                    host: Symbol('host')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['getDefaultSearchCondition'])
        },
        activated () {
            this.value = ''
            this.validList = []
            this.invalidList = []
        },
        methods: {
            handleFocus () {
                this.invalidList = []
            },
            async handleConfirm () {
                try {
                    await this.validateList()
                    if (this.validList.length) {
                        const { info } = await this.$store.dispatch('hostSearch/searchHost', {
                            params: this.getParams(),
                            config: {
                                requestId: this.request.host
                            }
                        })
                        const unexistList = this.validList.filter(ip => {
                            const exist = info.some(target => target.host.bk_host_innerip === ip)
                            return !exist
                        })
                        const newHostList = info.filter(({ host }) => !this.hostList.some(target => target.host.bk_host_id === host.bk_host_id))
                        this.hostList.push(...newHostList)
                        this.invalidList.push(...unexistList)
                    }
                    this.value = this.invalidList.join('\n')
                } catch (e) {
                    console.error(e)
                }
            },
            async validateList () {
                const list = [...new Set(this.value.split('\n').map(ip => ip.trim()).filter(ip => ip.length))]
                const validateQueue = []
                list.forEach(ip => {
                    validateQueue.push(this.$validator.verify(ip, 'ip'))
                })
                const results = await Promise.all(validateQueue)
                const validList = []
                const invalidList = []
                results.forEach(({ valid }, index) => {
                    if (valid) {
                        validList.push(list[index])
                    } else {
                        invalidList.push(list[index])
                    }
                })
                this.validList = validList
                this.invalidList = invalidList
            },
            getParams () {
                return {
                    bk_biz_id: this.bizId,
                    condition: this.getDefaultSearchCondition(),
                    ip: { data: this.validList, exact: 1, flag: 'bk_host_innerip' }
                }
            },
            handleHostSelectChange (data) {
                this.$emit('select-change', data)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .custom-layout {
        position: relative;
        display: flex;
        height: 100%;
        padding-top: 24px;

        .input-wrapper {
            position: relative;
            width: 280px;
        }
        .table-wrapper {
            flex: auto;
            margin-left: 20px;
        }
    }
    .ip-input {
        display: block;
        width: 100%;
        height: calc(100% - 60px);
        padding: 5px 10px;
        font-size: 12px;
        line-height: 20px;
        background-color: #FFF;
        border-radius:2px;
        border:1px solid #C4C6CC;
        cursor: text;
        outline: 0;
        resize: none;
        @include scrollbar;
        &.has-error {
            color: $dangerColor;
            text-decoration: underline;
        }
    }
    .ip-error {
        position: absolute;
        bottom: 45px;
        left: 0;
        line-height: 16px;
        font-size: 12px;
        color: $dangerColor;
    }
    .ip-confirm {
        display: block;
        width: 100%;
        margin: 15px 0;
    }
</style>
