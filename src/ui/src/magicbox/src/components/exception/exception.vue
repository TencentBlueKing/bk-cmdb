<template>
    <div class="bk-exception">
        <img :src="images[type]" :alt="type"/>
        <div class="exception-text">
            <slot>{{tipText[type]}}</slot>
        </div>
    </div>
</template>
<script>
    import Building from '../../bk-magic-ui/src/images/building.png'
    import notFound from '../../bk-magic-ui/src/images/404.png'
    import permissions from '../../bk-magic-ui/src/images/403.png'
    import maintain from '../../bk-magic-ui/src/images/500.png'
    export default {
        name: 'bk-exception',
        props: {
            type: {
                type: [String, Number],
                default: 404,
                validator (value) {
                    return [
                        404,
                        403,
                        500,
                        '404',
                        '403',
                        '500',
                        'building'
                    ].indexOf(value) > -1
                }
            }
        },
        data () {
            return {
                images: {
                    403: permissions,
                    404: notFound,
                    500: maintain,
                    building: Building
                },
                tipText: {
                    404: '页面找不到了！',
                    403: 'Sorry，您的权限不足',
                    500: '服务维护中，请稍后...',
                    building: '功能正在建设中···'
                }
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/exception.scss'
</style>
