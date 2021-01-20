import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import ExportContent from './content.vue'
import RouterQuery from '@/router/query'
import $http from '@/api'
import { $info } from '@/magicbox'
const Component = Vue.extend({
    components: {
        ExportContent
    },
    data () {
        return {
            visible: false
        }
    },
    created () {
        this.unwatch = RouterQuery.watch('*', () => {
            this.close()
        })
    },
    beforeDestroy () {
        this.unwatch()
    },
    methods: {
        toggle (visible) {
            this.visible = visible
        },
        close () {
            document.body.removeChild(this.$el)
            this.$destroy()
        }
    },
    render (h) {
        return (
            <bk-dialog
                width={ 768 }
                value={ this.visible }
                draggable={ false }
                close-icon={ false }
                show-footer={ false }
                transfer={ false }
                on-change={ this.toggle }
                on-after-leave={ this.close }>
                <export-content { ...{ props: this.$options.attrs }} on-close={ () => this.toggle(false) }></export-content>
            </bk-dialog>
        )
    }
})

export default function (data = {}) {
    const props = {
        limit: 10000,
        ...data
    }
    if (props.count <= props.limit) {
        const options = props.options({
            start: 0,
            limit: props.limit
        });
        (async () => {
            const message = $info(i18n.t('正在导出'), 0)
            try {
                await $http.download(options)
            } catch (error) {
                console.error(error)
            } finally {
                message.close()
            }
        })()
    } else {
        const vm = new Component({
            store,
            i18n,
            attrs: props
        })
        vm.$mount()
        document.body.appendChild(vm.$el)
        vm.toggle(true)
    }
}
