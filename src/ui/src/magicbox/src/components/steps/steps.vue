<template>
    <div class="bk-steps"
        :class="`bk-steps-${theme}`">
        <div class="bk-step"
            v-for="(step, index) in defaultSteps"
            :class="{done: curStep > (index + 1), current: curStep === (index + 1)}">
            <span class="bk-step-indicator"
                :class="'bk-step-' + (iconType(step) ? 'icon' : 'number')"
                :style="{cursor: controllable ? 'pointer' : ''}"
                @click="jumpTo(index + 1)">
                <i class="bk-icon"
                    v-if="iconType(step)"
                    :class="'icon-' + isIcon(step)">
                </i>
                <span v-else>{{ isIcon(step) }}</span>
            </span>
            <span class="bk-step-title"
                v-if="step.title">
                {{ step.title }}
            </span>
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
     *  @param theme {String} - 组件的主题色
     *  @example
     *  <bk-steps
            :cur-step.sync="curStep"
            :steps="step"
            :theme="theme"
            controllable>
        </bk-steps>
     */

    import locale from '../../mixins/locale'

    export default {
        name: 'bk-steps',
        mixins: [locale],
        props: {
            steps: {
                type: Array,
                default: () => []
                // validator (value) {
                //     return value.length > 1
                // },
                // default () {
                //     return [
                //         {
                //             title: '步骤1',
                //             icon: 1
                //         },
                //         {
                //             title: '步骤2',
                //             icon: 2
                //         },
                //         {
                //             title: '步骤3',
                //             icon: 3
                //         }
                //     ]
                // }
            },
            curStep: {
                type: Number,
                default: 1
            },
            controllable: {
                type: Boolean,
                default: false
            },
            theme: {
                type: String,
                default: 'primary'
            }
        },
        data () {
            return {
                defaultSteps: [
                    {
                        title: this.t('steps.step1'),
                        icon: 1
                    },
                    {
                        title: this.t('steps.step2'),
                        icon: 2
                    },
                    {
                        title: this.t('steps.step3'),
                        icon: 3
                    }
                ]
            }
        },
        created () {
            if (this.steps && this.steps.length) {
                const defaultSteps = []
                this.steps.forEach(step => {
                    if (typeof step === 'string') {
                        defaultSteps.push(step)
                    }
                    else {
                        defaultSteps.push({
                            title: step.title,
                            icon: step.icon
                        })
                    }
                })
                this.defaultSteps.splice(0, this.defaultSteps.length, ...defaultSteps)
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
                if (this.controllable && index !== this.curStep) {
                    this.$emit('update:curStep', index)
                    this.$emit('step-changed', index)
                }
            }
        }
    }
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/steps.scss'
</style>
