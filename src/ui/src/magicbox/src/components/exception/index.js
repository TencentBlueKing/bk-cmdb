import bkException from './exception'
bkException.install = Vue => {
    Vue.component(bkException.name, bkException)
}
export default bkException
