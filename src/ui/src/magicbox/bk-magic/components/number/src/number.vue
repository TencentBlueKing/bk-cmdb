<template>
    <div :class="['bk-number', {'focus': isFocus, 'disabled':disabled}]" :style="exStyle">
        <div class="bk-number-content " :class="[{'bk-number-larger':size === 'large','bk-number-small':size === 'small'}]">
            <input type="number" 
                :disabled="disabled" 
                :placeholder="placeholder" 
                class="bk-number-input" 
                id="bk-number-input" 
                @focus="focus"
                @blur="blur" 
                v-model="counter">
            <div class="bk-number-icon-content" v-if="!hideOperation">
                <div :class="['bk-number-icon-top', {'btn-disabled': isMax}]" @click="add">
                    <i class="bk-icon icon-angle-up"></i>
                </div>
                <div :class="['bk-number-icon-lower', {'btn-disabled': isMin}]" @click="minus">
                    <i class="bk-icon icon-angle-down"></i>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
    export default {
        name: 'bk-number-input',
        props: {
            value: {
                type: [Number, String],
                default: 0
            },
            hideOperation: {
                type: Boolean,
                default: false
            },
            exStyle: {
                type: Object,
                default () {
                    return {}
                }
            },
            placeholder: {
                type: String,
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            min: {
                type: Number,
                detault: Number.NEGATIVE_INFINITY
            },
            max: {
                type: Number,
                detault: Number.POSITIVE_INFINITY
            },
            steps: {
                type: Number,
                detault: 1
            },
            size: {
                type: String,
                default: 'large',
                validator (value) {
                    return [
                        'large',
                        'small'
                    ].indexOf(value) > -1
                }
            }
        },
        data () {
            return {
                isMax: false,
                isMin: false,
                counter: this.value,
                isFocus: false,
                maxNumber: this.max,
                minNumber: this.min
            }
        },
        watch: {
            'counter': function (val) {
                val = this.checkInput(val)
                this.$emit('update:value', val)
                this.$emit('change', val)
                this.counter = val
            }
        },
        computed: {
            prevNumber () {
                if (this.counter === '') {
                    if (this.min) {
                        return this.min
                    } else {
                        return 0
                    }
                } else {
                    return this.counter
                }
            }
        },
        methods: {
            focus () {
                this.isFocus = true
            },
            blur () {
                this.isFocus = false
            },
            checkInput (val) {
                val = val + ''
                val = val.replace(/[^\d|^-]/g, '')
                if (val === '') {
                    return val
                } else {
                    val = parseInt(val)
                    val = this.checkMinMax(val)
                    return val
                }
            },
            keyup (event) {
                if (event.code === 'ArrowUp') {
                    this.add()
                } else if (event.code === 'ArrowDown') {
                    this.minus()
                }
            },
            checkMinMax (val) {
                if (val <= this.minNumber) {
                    val = this.minNumber
                    this.isMin = true
                } else {
                    this.isMin = false
                }
                if (val >= this.maxNumber) {
                    val = this.maxNumber
                    this.isMax = true
                } else {
                    this.isMax = false
                }
                return val
            },
            add: function () {
                if (this.disabled) return
                if (String(this.steps).indexOf('.') > 0) {
                    this.counter = (parseFloat(this.counter) + this.steps).toFixed(1)
                } else {
                    this.counter = parseFloat(this.counter) + this.steps
                }
                // this.checkMinMax()
            },
            minus: function () {
                if (this.disabled) return
                if (String(this.steps).indexOf('.') > 0) {
                    this.counter = (parseFloat(this.counter) - this.steps).toFixed(1)
                } else {
                    this.counter = parseFloat(this.counter) - this.steps
                }
                // this.checkMinMax()
            }
        }
    }
</script>
