/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <bk-select 
        :filterable="true"
        :selected.sync="timezoneSelected"
        :disabled="disabled">
        <bk-select-option v-for="(item, index) in timezoneData"
            :key="index"
            :label="item"
            :value="item">
        </bk-select-option>
    </bk-select>
</template>

<script>
    import {mapGetters} from 'vuex'
    const timezoneData = require('../../common/json/timezone.json')
    export default {
        props: {
            disabled: {
                type: Boolean,
                default: false
            },
            selected: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                timezoneData: [],
                timezoneSelected: ''
            }
        },
        computed: {
            ...mapGetters([
                'timezoneList'
            ])
        },
        watch: {
            selected (selected) {
                this.timezoneSelected = selected
            },
            timezoneSelected (val) {
                this.$emit('update:selected', val)
            }
        },
        mounted () {
            if (!this.timezoneList.length) {
                this.$store.commit('setTimezoneList', timezoneData)
            }
            this.timezoneSelected = this.selected || 'Asia/Shanghai'
            this.timezoneData = timezoneData
        }
    }
</script>

<style lang="scss" >

</style>
