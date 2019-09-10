<template>
    <ul class="icon-set">
        <li class="icon"
            ref="iconItem"
            :class="{ 'active': icon.value === value }"
            v-bk-tooltips="{ content: language === 'zh_CN' ? icon.nameZh : icon.nameEn }"
            v-for="(icon, index) in curIconList"
            :key="index"
            @click="handleChooseIcon(icon.value)">
            <i :class="icon.value"></i>
            <span class="checked-status"></span>
        </li>
    </ul>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            value: {
                type: String,
                default: 'icon-cc-default'
            },
            iconList: {
                type: Array,
                default: () => []
            },
            filterIcon: {
                type: String,
                default: ''
            }
        },
        computed: {
            ...mapGetters([
                'language'
            ]),
            curIconList () {
                if (this.filterIcon) {
                    return this.iconList.filter(icon => icon.nameZh.toLowerCase().indexOf(this.filterIcon.toLowerCase()) > -1 || icon.nameEn.toLowerCase().indexOf(this.filterIcon.toLowerCase()) > -1)
                }
                return this.iconList
            }
        },
        methods: {
            handleChooseIcon (value) {
                this.$emit('input', value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .icon-set {
        width: 560px;
        display: flex;
        flex-wrap: wrap;
        padding-bottom: 10px;
        .icon {
            position: relative;
            display: flex;
            justify-content: center;
            align-items: center;
            flex: 0 0 10%;
            height: 50px;
            font-size: 24px;
            outline: 0;
            cursor: pointer;
            &:hover {
                color: #3a84ff;
                background-color: #ebf4ff;
            }
            &.active {
                color: #3a84ff;
                background-color: #ebf4ff;
                border: 1px dashed #3a84ff;
                .checked-status {
                    display: block;
                }
            }
            .checked-status {
                display: none;
                position: absolute;
                bottom: -6px;
                right: -6px;
                width: 18px;
                height: 18px;
                background-color: #2dcb56;
                border-radius: 50%;
                z-index: 2;
                &::before {
                    content: '';
                    position: absolute;
                    bottom: 5px;
                    right: 0;
                    width: 14px;
                    height: 7px;
                    border-bottom: 3px solid #ffffff;
                    border-left: 3px solid #ffffff;
                    transform: rotate(-45deg) scale(.5);
                }
            }
        }
    }
</style>
