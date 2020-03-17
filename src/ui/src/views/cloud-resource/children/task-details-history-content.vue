<template>
    <div class="content-layout">
        <div class="content-group"
            v-for="(group, groupIndex) in groupedList"
            :key="groupIndex">
            <h3 class="group-title">
                <span class="title-type">{{group.type}}</span>
                <span class="title-count">({{group.list.length}})</span>
                <i class="title-copy icon-cc-details-copy"
                    @click="handleCopy($event, group)">
                    <span class="copy-tips">{{$t('复制成功')}}</span>
                </i>
            </h3>
            <ul class="host-list">
                <li class="host-item"
                    v-for="(host, hostIndex) in group.list"
                    :key="hostIndex">
                    {{host.bk_host_innerip}}
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'task-details-history-content',
        props: {
            list: {
                type: Array,
                default: () => ([])
            }
        },
        computed: {
            groupedList () {
                const created = []
                const updated = []
                this.list.forEach(item => {
                    item.isCreate ? created.push(item) : updated.push(item)
                })
                const groupedList = [{ type: this.$t('新增'), list: created }, { type: this.$t('更新'), list: updated }]
                return groupedList.filter(group => group.list.length)
            }
        },
        methods: {
            async handleCopy (event, group) {
                try {
                    await this.$copyText(group.list.map(item => item.bk_host_innerip).join('\n'))
                    const target = event.currentTarget
                    target.classList.add('show-tips')
                    setTimeout(() => {
                        target.classList.remove('show-tips')
                    }, 500)
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .content-layout {
        max-height: 200px;
        padding: 0 18px;
        @include scrollbar-y;
    }
    .content-group {
        margin: 10px 0 15px;
    }
    .group-title {
        font-size: 12px;
        .title-count {
            font-size: 12px;
            font-weight: normal;
            color: #C4C6CC;
        }
        .title-copy {
            position: relative;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                opacity: .8;
            }
            &.show-tips {
                .copy-tips {
                    animation: copy-tips-animation .5s linear;
                }
            }
            .copy-tips {
                position: absolute;
                bottom: 100%;
                left: 50%;
                padding: 0 8px;
                transform: translate3d(-50%, 0, 0);
                line-height: 26px;
                font-size: 12px;
                color: #ffffff;
                text-align: center;
                background-color: #000;
                border-radius: 2px;
                white-space: nowrap;
                transition: all .3s linear;
                pointer-events: none;
                opacity: 0;
            }
        }
    }
    .host-list {
        display: flex;
        margin: 5px 0 0 0;
        flex-wrap: wrap;
        line-height: 16px;
        .host-item {
            flex: 20% 0 0;
            padding: 0 25px 0 0;
            @include ellipsis;
        }
    }
    @keyframes copy-tips-animation {
        0% {
            transform: translate3d(-50%, 0, 0);
            opacity: 0;
        }
        25% {
            transform: translate3d(-50%, -4px, 0);
            opacity: .4;
        }
        50% {
            transform: translate3d(-50%, -8px, 0);
            opacity: .8;
        }
        75% {
            transform: translate3d(-50%, -12px, 0);
            opacity: .4;
        }
        100% {
            transform: translate3d(-50%, -16px, 0);
            opacity: 0;
        }
    }
</style>
