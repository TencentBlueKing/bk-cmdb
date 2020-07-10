<template>
    <div class="footer">
        <p class="concat">
            <template v-for="(link, index) in links">
                <i class="gap" v-if="index > 0" :key="index"></i>
                <bk-link
                    theme="primary"
                    target="_blank"
                    :key="link.value"
                    :href="link.value">
                    {{$i18n.locale === 'en' ? link.i18n.en : link.i18n.cn}}
                </bk-link>
            </template>
        </p>
        <p class="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. {{site.buildVersion}}
        </p>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                year: (new Date()).getFullYear()
            }
        },
        computed: {
            ...mapGetters(['site']),
            links () {
                const { footer = {} } = this.site
                const links = footer.links || []
                return links.filter(link => link.enabled)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .footer {
        position: absolute;
        width: calc(100vw - 50px);
        left: 25px;
        bottom: 0;
        padding: 10px 0;
        font-size: 12px;
        text-align: center;
        color: #C4C6CC;
        border-top: 1px solid #DCDEE5;
        background-color: #F5F6FA;
        z-index: 2;
        .concat {
            display: flex;
            justify-content: center;
            align-items: center;
            .gap {
                display: inline-flex;
                width: 2px;
                height: 16px;
                background-color: #c4c6cc;
                margin: 0 10px;
            }
        }
        .copyright {
            line-height: 24px;
        }
    }
</style>
