/**
 * @file https://github.com/freeze-component/vue-popper
 */

import Popper from 'popper.js'

export default {
    props: {
        placement: {
            type: String,
            default: 'bottom'
        },
        boundariesPadding: {
            type: Number,
            default: 5
        },
        reference: Object,
        popper: Object,
        offset: {
            default: 0
        },
        value: {
            type: Boolean,
            default: false
        },
        transition: String,
        options: {
            type: Object,
            default () {
                return {
                    gpuAcceleration: false,
                    boundariesElement: 'body'
                }
            }
        }
        // visible: {
        //     type: Boolean,
        //     default: false
        // }
    },
    data () {
        return {
            visible: this.value
        }
    },
    watch: {
        value: {
            immediate: true,
            handler (val) {
                this.visible = val
                this.$emit('input', val)
            }
        },
        visible (val) {
            if (val) {
                this.updatePopper()
                this.$emit('on-show', this)
            } else {
                this.destroyPopper()
                this.$emit('on-hide', this)
            }
            this.$emit('input', val)
        }
    },
    methods: {
        createPopper () {
            if (!/^(top|bottom|left|right)(-start|-end)?$/g.test(this.placement)) {
                return
            }

            const options = this.options
            const popper = this.popper || this.$refs.popper
            const reference = this.reference || this.$refs.reference

            if (!popper || !reference) {
                return
            }

            if (this.popperJS && this.popperJS.hasOwnProperty('destroy')) {
                this.popperJS.destroy()
            }

            options.placement = this.placement
            options.offset = this.offset

            // this.popperJS = new Popper(reference, popper, {
            //     ...options,
            //     onCreate: popper => {
            //         this.resetTransformOrigin(popper.instance.popper)
            //         this.$nextTick(this.updatePopper)
            //         this.$emit('created', this)
            //     }
            // })

            this.popperJS = new Popper(reference, popper, Object.assign({}, options, {
                onCreate: popper => {
                    this.resetTransformOrigin(popper.instance.popper)
                    this.$nextTick(this.updatePopper)
                    this.$emit('created', this)
                }
            }))
        },
        updatePopper () {
            this.popperJS ? this.popperJS.update() : this.createPopper()
        },
        doDestroy () {
            if (this.visible) {
                return
            }
            this.popperJS.destroy()
            this.popperJS = null
        },
        destroyPopper () {
            if (this.popperJS) {
                this.resetTransformOrigin(this.popperJS.popper)
            }
        },
        resetTransformOrigin (popperNode) {
            let placementMap = {top: 'bottom', bottom: 'top', left: 'right', right: 'left'}
            let placement = popperNode.getAttribute('x-placement').split('-')[0]
            let origin = placementMap[placement]
            popperNode.style.transformOrigin = ['top', 'bottom'].indexOf(placement) > -1
                ? 'center ' + origin
                : origin + ' center'
        }
    },
    beforeDestroy () {
        if (this.popperJS) {
            this.popperJS.destroy()
        }
    }
}
