import Vue from 'vue'
import ViewModel from './scroll-bar.vue'

let Model = Vue.extend(ViewModel)

let install = Vue => {
    /**
     *  初始化插件配置
     */
    function initConfig (el) {
        el._scrollAttrs = {
            defaultScrollBarWidth: calcDefaultScrollBarWidth()
        }

        return el._scrollAttrs
    }

    /**
     *  计算默认滚动条的宽度
     */
    function calcDefaultScrollBarWidth () {
        return `${window.innerWidth - document.body.clientWidth}px`
    }

    /**
     *  complie template
     *  @param html {String} string of DOM
     *  @return DOM {DOMNode} DOM
     */
    function compiler (html) {
        let temp = document.createElement('div')
        let children = null
        let fragment = document.createDocumentFragment()

        temp.innerHTML = html
        children = temp.childNodes

        for (let child of children) {
            fragment.appendChild(child.cloneNode(true))
        }

        return fragment
    }

    /**
     *  初始化滚动条样式
     */
    function initStyle (el, binding) {
        let children = el.childNodes
        let childrenTpl = ''
        let wrapper = document.createElement('div')

        wrapper.classList.add('bk-scrollbar-wrapper')

        // el.insertBefore(wrapper, el)
        el.appendChild(wrapper)
        wrapper.style.height = '100%'
    }

    /**
     *  初始化指令
     */
    function initScrollBar (el, binding) {
        initStyle(el, binding)
    }

    Vue.directive('bkscrollbar', {
        inserted (el, binding) {
            console.log(el.firshChild)
            initScrollBar(el, binding)
            // let value = binding.value || {}
            // let computedStyle = getComputedStyle(el)
            // let position = computedStyle.position
            // let { vertical = true, horizontal = false } = value
            // let style = el.style
            //
            // if (!position || (position !== 'relative' && position !== 'absolute')) {
            //     style.position = 'relative'
            // }
            //
            // console.log(vertical, horizontal)
            //
            // let scrollbar = new Model({
            //     data: value
            // })
            //
            // el[`${vertical ? 'vertical' : 'horizontal'}bar`] = scrollbar
            // el.$vm = scrollbar.$mount()
            // el.appendChild(el.$vm.$el)
            // console.log(computedStyle.height)
        }
    })
}

export default install
