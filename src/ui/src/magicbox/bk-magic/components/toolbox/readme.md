<script>
    export default {
    	mounted () {
    		this.toolboxShow()
    	},
        methods: {
        	toolboxShow () {
        		for(var j in this.$refs){
        			var parent = this.$refs[j].$el.children[0];
        			var elm = this.$refs[j].$slots['toolbox-content'][0].elm; 
        			parent.style.width = elm.clientWidth + 2 + 'px'
        		}
        	},
            triggerHandler () {
                // do
                for(var j in this.$refs){
        			var parent = this.$refs[j].$el.children[0];
        			parent.className = 'bk-toolbox';
        		}
            }
        }
    }
</script>

## Toolbox

### 基础用法

:::demo `slot[name=toolbox-trigger]`配置触发对象，`slot[name=toolbox-content]`配置弹层

```html
<template>
    <bk-toolbox ref="toolbox">
    	<i class="bk-icon icon-edit bk-toolbox-icon" slot="toolbox-trigger"></i>
    	<div class="bk-toolbox-wrapper" slot="toolbox-content">
			<div class="bk-toolbox-inner">
                <a href="javascript:;">添加</a>
                <a href="javascript:;">删除</a>
            </div>
        </div>
    </bk-toolbox>
</template>
<script>
    export default {
    	mounted () {
    		this.toolboxShow()
    	},
        methods: {
        	toolboxShow () {

        		for(var j in this.$refs){
        			var parent = this.$refs[j].$el.children[0];
        			var elm = this.$refs[j].$slots['toolbox-content'][0].elm; 
        			parent.style.width = elm.clientWidth + 2 + 'px'
        		}
        		
        	}
        }
    }
</script>
```

:::

### 回调用法

:::demo `slot[name=toolbox-trigger]`配置触发对象，`slot[name=toolbox-content]`配置弹层

```html
<template>
    <bk-toolbox ref="toolbox1">
    	<i class="bk-icon icon-edit bk-toolbox-icon" slot="toolbox-trigger"></i>
    	<div class="bk-toolbox-wrapper" slot="toolbox-content">
	    	<div class="bk-toolbox-inner">
	            <a href="javascript:;" @click="triggerHandler">添加</a>
	            <a href="javascript:;" @click="triggerHandler">删除</a>
	            <a href="javascript:;" @click="triggerHandler">收藏</a>
	        </div>
        </div>
    </bk-toolbox>
</template>
<script>
    export default {
    	mounted () {
    		this.toolboxShow()
    	},
        methods: {
        	toolboxShow () {

        		for(var j in this.$refs){
        			var parent = this.$refs[j].$el.children[0];
        			var elm = this.$refs[j].$slots['toolbox-content'][0].elm; 
        			parent.style.width = elm.clientWidth + 2 + 'px'
        		}
        		
        	},
            triggerHandler () {
                // do
                for(var j in this.$refs){
        			var parent = this.$refs[j].$el.children[0];
        			parent.className = 'bk-toolbox';
        		}
            }
        }
    }
</script>
```
:::

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| triggerHandler | 显示时触发此回调函数 | —— |