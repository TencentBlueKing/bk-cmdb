<template>
    <div class="list-wrapper">
        <ul class="list-box">
            <li :class="{'active': activeClassify['bk_classification_id'] === classify['bk_classification_id']}" 
            @click="changeClassify(classify)" 
            v-for="(classify, index) in localClassify"
            :key="index">
                <i :class="classify['bk_classification_icon']"></i>
                <span class="text">{{classify['bk_classification_name']}}</span>
            </li>
        </ul>
        <div class="add-btn-wrapper" @click="showPop(false)">
            <a class="add-btn" href="javascript:;">
                <span>新增</span>
            </a>
        </div>
        <v-pop
            ref="pop"
            :isShow.sync="isPopShow"
            :classification="activeClassify"
            :isEdit="isEdit"
            @confirm="saveClassify"
        ></v-pop>
    </div>
</template>

<script type="text/javascript">
    import bus from '@/eventbus/bus'
    import vPop from './pop'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        data () {
            return {
                isPopShow: false,                               // 弹窗显示状态
                isEdit: false,                                  // 弹窗是否处于编辑状态
                activeClassify: {                               // 当前分类信息
                    bk_classification_icon: '',
                    bk_classification_id: 'bk_host_manage',
                    bk_classification_name: '',
                    bk_classification_type: ''
                }
            }
        },
        computed: {
            ...mapGetters([
                'allClassify'
            ]),
            /*
                分类列表
            */
            localClassify: {
                get () {
                    return this.$deepClone(this.allClassify)
                },
                set () {
                    this.localClassify = this.$deepClone(this.allClassify)
                }
            }
        },
        watch: {
            allClassify (newVal, oldVal) {
                if (!oldVal.length) {
                    this.initActiveClassify()
                }
            }
        },
        methods: {
            /*
                显示分类弹窗
                isEdit: 是否为编辑状态
            */
            showPop (isEdit) {
                this.isPopShow = true
                this.isEdit = isEdit
            },
            /*
                切换分类
                classify: 当前分类
            */
            changeClassify (classify) {
                this.activeClassify = classify
                this.$emit('changeClassify', classify)
            },
            /*
                弹窗确认按钮回调
                classification: 分类信息
            */
            saveClassify (classification) {
                if (this.isEdit) {
                    this.editClassify(classification)
                } else {
                    this.createClassify(classification)
                }
            },
            /*
                编辑分组
                classification: 分类信息
            */
            async editClassify (classification) {
                let params = {
                    bk_classification_icon: classification['bk_classification_icon'],
                    bk_classification_name: classification['bk_classification_name']
                }
                try {
                    await this.$axios.put(`object/classification/${this.activeClassify['id']}`, params)
                    this.isPopShow = false
                    this.activeClassify['bk_classification_icon'] = classification['bk_classification_icon']
                    this.activeClassify['bk_classification_name'] = classification['bk_classification_name']
                    this.$store.commit('updateClassify', this.activeClassify)
                } catch (e) {
                    console.error(e)
                    this.$alertMsg(e.data['bk_error_msg'])
                }
            },
            /*
                新增分组
                classification: 分类信息
            */
            async createClassify (classification) {
                let params = {
                    bk_classification_icon: classification['bk_classification_icon'],
                    bk_classification_id: classification['bk_classification_id'],
                    bk_classification_name: classification['bk_classification_name']
                }
                try {
                    const res = await this.$axios.post(`object/classification`, params)
                    this.isPopShow = false
                    this.activeClassify = {
                        ...params,
                        bk_classification_type: '',
                        bk_objects: [],
                        id: res.data.id
                    }
                    this.$store.commit('createClassify', this.activeClassify)
                    this.changeClassify(this.activeClassify)
                } catch (e) {
                    console.error(e)
                    this.$alertMsg(e.data['bk_error_msg'])
                }
            },
            /*
                删除分类
            */
            deleteModelClassify () {
                let self = this
                this.$bkInfo({
                    title: '确认要删除此分类？',
                    confirmFn () {
                        self.deletes()
                    }
                })
            },
            async deletes () {
                try {
                    await this.$axios.delete(`object/classification/${this.activeClassify['id']}`)
                    this.$store.commit('deleteClassify', this.activeClassify)
                    this.changeClassify(this.$deepClone(this.allClassify[0]))
                    this.$alertMsg('删除成功', 'success')
                } catch (e) {
                    console.error(e)
                    this.$alertMsg(e.data.message)
                }
            },
            initActiveClassify () {
                this.changeClassify(this.$deepClone(this.allClassify[0]))
            }
        },
        created () {
            bus.$on('editModelClass', () => {
                this.showPop(true)
            })
            bus.$on('deleteModelClass', () => {
                this.deleteModelClassify()
            })
            if (this.allClassify.length) {
                this.initActiveClassify()
            }
        },
        components: {
            vPop
        }
    }
</script>

<style lang="scss" scoped>
    $primaryColor: #737987;
    $primaryHoverColor: #3c96ff; 
    $primaryHoverBgColor: #f1f7ff;
    $primaryActiveBgColor: #e2efff;
    $white: #fff;
    $borderColor: #dde4eb;
    $btnColor: #c3cdd7;
    .list-wrapper{
        width:188px;
        float:left;
        border-left: none;
        border-top: none;
        height: 100%;
        overflow-y: auto;
        border-right: 1px solid $borderColor;
        @include scrollbar;
        .list-box{
            >li{
                height: 48px;
                line-height: 48px;
                padding: 0 30px 0 44px;
                width: 100%;
                cursor: pointer;
                font-size: 14px;
                color: $primaryColor;
                font-size: 14px;
                position: relative;
                white-space: nowrap;
                text-overflow: ellipsis;
                overflow: hidden;
                i{
                    font-size: 16px;
                }
                .icon-left{
                    margin-left: -12px;
                }
                &:hover{
                    color: $primaryHoverColor;
                    background: $primaryHoverBgColor;
                }
                .text{
                    padding: 0 3px 0 5px;
                    min-width: 64px;
                    vertical-align: top;
                }
                &.active{
                    color: $primaryHoverColor;
                    background: $primaryActiveBgColor;
                }
            }
        }
        .add-btn-wrapper{
            width: 148px;
            height: 32px;
            background: $white;
            cursor: pointer;
            font-size: 0;
            margin: 10px auto;
            .add-btn{
                display: block;
                height: 32px;
                line-height: 30px;
                border-radius: 2px;
                color: $btnColor;
                border: dashed 1px $btnColor;
                text-align: center;
                font-size: 14px;
                &:hover{
                    border-color: $primaryHoverColor;
                    color: $primaryHoverColor;
                }
            }
        }
    }
</style>
