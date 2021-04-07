<template>
  <div class="service-category">
    <span class="title">{{$t('服务分类')}}</span>
    <div class="selector-item mt10 clearfix">
      <bk-select class="category-selector fl"
        :clearable="false"
        v-model="parent"
        @change="handleParentChange">
        <bk-option v-for="category in parentList"
          :key="category.id"
          :id="category.id"
          :name="category.name">
        </bk-option>
      </bk-select>
      <bk-select class="category-selector fl"
        :clearable="false"
        v-model="child"
        v-validate="'required'"
        name="secondCategory">
        <bk-option v-for="category in childList"
          :key="category.id"
          :id="category.id"
          :name="category.name">
        </bk-option>
      </bk-select>
      <span class="second-category-errors" v-if="errors.has('secondCategory')">{{errors.first('secondCategory')}}</span>
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  export default {
    name: 'form-service-category',
    props: {
      instance: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        parent: '',
        child: '',
        list: []
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      parentList() {
        return this.list.filter(category => !category.bk_parent_id)
      },
      childList() {
        return this.list.filter(category => category.bk_parent_id === this.parent)
      }
    },
    watch: {
      child(child) {
        this.$emit('change', child)
      }
    },
    async created() {
      await this.getList()
      this.setupValue()
    },
    methods: {
      async getList() {
        try {
          const { info = [] } = await this.$store.dispatch('serviceClassification/searchServiceCategoryWithoutAmout', {
            params: { bk_biz_id: this.bizId }
          })
          this.list = info
        } catch (error) {
          console.error(error)
        }
      },
      setupValue() {
        const {
          bk_parent_id: parent = '',
          id: child = ''
        } = this.list.find(category => category.id === this.instance.service_category_id) || {}
        this.parent = parent
        this.child = child
      },
      handleParentChange() {
        const [firstChild] = this.childList
        this.child = firstChild ? firstChild.id : ''
      }
    }
  }
</script>

<style lang="scss" scoped>
.service-category {
    font-size: 12px;
    padding: 20px 0 24px 36px;
    margin: 0 20px;
    border-bottom: 1px solid #dcdee5;
    .selector-item {
        position: relative;
        width: 50%;
        max-width: 554px;
        padding-right: 54px;
    }
    .category-selector {
        width: calc(50% - 5px);
        & + .category-selector {
            margin-left: 10px;
        }
    }
    .second-category-errors {
        position: absolute;
        top: 100%;
        left: 0;
        margin-left: calc(50% - 18px);
        line-height: 14px;
        font-size: 12px;
        color: #ff5656;
        max-width: 100%;
        @include ellipsis;
    }
}
</style>
