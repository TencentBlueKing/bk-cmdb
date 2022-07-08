/**
 * 内置到 Vue.proottype 的全局变量，命名规则为 $ 开头
 */
import Vue from 'vue'

/**
 * @global Site 是后台编译时内置在 builder/config/index.js 中变量，随着前端编译后会渲染到 index.html 并保存在全局变量 window.Site 中。
 */
Vue.prototype.$Site = window.Site
