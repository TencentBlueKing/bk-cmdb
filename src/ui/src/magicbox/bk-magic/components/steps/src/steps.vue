<template>
    <div class="bk-flows">
        <div class="bk-flow" v-for="(step, index) in steps" :class="{done: curStep > (index + 1), current: curStep === (index + 1)}">
            <span :class="'bk-flow-' + (iconType(step) ? 'icon' : 'number')" @click="jumpTo(index + 1)" :style="{cursor: controllable ? 'pointer' : ''}">
        <i class="bk-icon"
          :class="'icon-' + isIcon(step)"
          v-if="iconType(step)"></i>
        <span v-else>{{ isIcon(step) }}</span>
            </span>
            <span class="bk-flow-title" v-if="step.title">{{ step.title }}</span>
        </div>
    </div>
</template>
<script>
    /**
     *  bk-steps
     *  @module components/steps
     *  @desc 步骤条组件
     *  @param steps {Array} - 组件步骤内容，数组中的元素可以是对象，可以是数字，可以是字符串，也可以是三者混合元素是对象时，有两个可选的key：title和icon；当元素是数字或字符串时，组件会将其解析成对象形式下的icon的值；素是字符串时，使用蓝鲸icon
     *  @param curStep {Number} - 当前步骤的索引值，从1开始；支持.sync修饰符
     *  @param controllable {Boolean} - 步骤可否被控制前后跳转，默认为false
     *  @example
     *  <bk-steps
          :cur-step.sync="curStep"
          :steps="step"
          controllable></bk-steps>
     */
    export default {
        name: 'bk-steps',
        props: {
            steps: {
                type: Array,
                validator (value) {
                    return value.length > 1
                },
                default () {
                    return [{
                        title: '步骤1',
                        icon: 1
                    }, {
                        title: '步骤2',
                        icon: 2
                    }, {
                        title: '步骤3',
                        icon: 3
                    }]
                }
            },
            curStep: {
                type: Number,
                default: 1
            },
            controllable: {
                type: Boolean,
                default: false
            }
        },
        methods: {
            iconType (step) {
                let icon = step.icon

                if (icon) {
                    return typeof icon === 'string'
                } else {
                    return typeof step === 'string'
                }
            },
            isIcon (step) {
                return step.icon ? step.icon : step
            },
            jumpTo (index) {
                if (this.controllable) {
                    this.$emit('update:curStep', index)
                    this.$emit('step-changed', index)
                }
            }
        }
    }
</script>
