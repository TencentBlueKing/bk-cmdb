<template>
    <div class="classify">
        <h4 class="classify-name">{{`${classify['bk_classification_name']}(${classify['bk_objects'].length})`}}</h4>
        <ul class="models-list">
            <li class="models-item" v-for="(model, index) in classify['bk_objects']" :key="index">
                <router-link exact class="model-link"
                    :to="getModelLink(model)"
                    :title="model['bk_obj_name']">
                    <i :class="['model-icon','icon', model['bk_obj_icon']]"></i>
                    <span class="model-name">{{model['bk_obj_name']}}</span>
                </router-link>
                <i class="model-star bk-icon icon-star"></i>
            </li>
        </ul>
    </div>
</template>

<script>
    export default {
        props: {
            classify: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                notModelClassify: ['bk_host_manage', 'bk_back_config']
            }
        },
        methods: {
            getModelLink (model) {
                if (this.notModelClassify.includes(model['bk_classification_id'])) {
                    return model.path
                }
                return `/organization/${model['bk_obj_id']}`
            }
        }
    }
</script>

<style lang="scss" scoped>
    .classify{
        margin: 0 0 20px 0;
        padding: 0 20px;
        background-color: #fff;
        border: 1px solid #ebf0f5;
    }
    .classify-name{
        padding: 13px 0;
        margin: 0;
        line-height: 20px;
        font-size: 14px;
        color: $textColor;
        border-bottom: 1px solid #ebf0f5;
    }
    .models-list{
        padding: 8px 0;
        .models-item{
            height: 36px;
            padding: 6px 0;
            position: relative;
            .model-star{
                position: absolute;
                right: 10px;
                top: 10px;
                font-size: 16px;
                cursor: pointer;
            }
        }
    }
    .model-link{
        display: block;
        width: calc(100% - 40px);
        font-size: 0;
        &:before{
            content: "";
            display: inline-block;
            height: 100%;
            vertical-align: middle;
        }
        &:hover .model-icon,
        &:hover .model-name{
            color: #0082ff;
        }
        .model-icon,
        .model-name{
            display: inline-block;
            vertical-align: middle;
        }
        .model-icon{
            font-size: 16px;
            color: $textColor;
        }
        .model-name{
            margin: 0 0 0 12px;
            font-size: 14px;
            line-height: 24px;
            color: $textColor;
        }
    }
</style>