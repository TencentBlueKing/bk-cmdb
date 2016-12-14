<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<div id="panelContain_9899" class="sidebar-panel " style="transform: translate3d(0px, 0px, 0px); height: 100%; width: 860px;">
        <div class="sidebar-panel-container  sidebar-panel-container-new" data-child="innerContainer">
            <h4 class="lh24 col-md-12 pl0 host-details-top">
                <span>主机详情页</span>
                <a class="btn-close btn pull-right" href="javascript:void(0)" id="closePanleBtn">关闭</a>
            </h4>            


            <!-- <div class="col-md-12">                 -->
                <div class="col-md-6 detail-box basis">   
                    <p class="host-detail-title">基础属性</p>
                    <div class="icon-detail-gear public-cloud-attribute" data-content=".basis,.network,.handwarea,.os"></div>       
                    <div class="detail-row-add">
                        <span class="detail-icon icon-detail-plus"></span>
                        <input type="hidden" class="select2-eln" data-type="basis">
                    </div>
                </div>
            <!-- </div> -->
            <!-- <div class="clearfix"></div> -->                
                <div class="col-md-6 detail-box ccattribute">     
                    <p class="host-detail-title">自定义属性</p>            
                    <div class="icon-detail-gear public-cloud-attribute" data-content=".ccattribute"></div>                        
                    <div class="detail-row-add">
                        <span class="detail-icon icon-detail-plus"></span>
                        <input type="hidden" class="select2-eln" data-type="ccattribute">
                    </div>
                </div>                
        </div>
    </div>
</div>

<script>
$('document').ready(function(){
    function HostDetail(attritubeData){
        this.attritubeData= attritubeData;
        this.select2Data = [];
    }
    HostDetail.prototype={
        constructor : HostDetail,
        init:function(){
            this.initData();
            this.initAddNewSelect2();
            this.initEvent();
        },
        
        initData:function(){
            var data = this.attritubeData;
            var keys = Object.keys(data);
            var topThis = this;
            var temp ={};
            for(var i =0;i<keys.length;i++){
                temp[keys[i]]=[];              
                data[keys[i]].map(function(v,j){
                    if(v.selected && v.selected == true){
                        v.disabled = true;
                        topThis.addNewAttritube(v,keys[i],'init');
                    }else{
                        v.disabled = false;
                        temp[keys[i]].push(v);
                    }
                    return v;
                });                
            } 
            this.select2Data = temp;           
        },
        formatAttritubeData:function(type,id,event){            
            var topThis = this;
            topThis.attritubeData[type].map(function(v,i){
                if(v.id == id){
                    if(event === "change"){
                        topThis.select2Data[type].forEach(function(v1,i1){
                            if(id == v1.id){                            
                                topThis.select2Data[type].splice(i1,1);
                                return v1;
                            }
                        });
                    }else{
                        topThis.select2Data[type].push(v);
                    }
                    v.disabled = !v.disabled;
                    v.selected = !v.selected;
                }
            });
        },     
        initAddNewSelect2:function(){
            var topThis = this;
            $('.select2-eln').each(function(i,v){
                var type = $(v).data('type');
                var data = topThis.select2Data[type];
                if(data != undefined){
                    $(v).select2({
                        data: data,
                        width:"130px",
                        placeholder:'请选择属性',
                        formatNoMatches:"无"
                    }).on('change',function(e){
                        var selectData = $(this).select2('data');
                        topThis.addNewAttritube(selectData,type);
                        topThis.formatAttritubeData(type,selectData.id,'change');
                        $(v).select2('val',"");
                        $.ajax({
                            url: "/host/setDefaultField",
                            type: "POST",
                            data: "key="+selectData.key + "&type=a",
                            dataType: "json"
                        });
                    });
                }         
            })         
        },
        addNewAttritube:function(data,type,mode){            
            var temp = '<div class="detail-row" id="#type#_#id#_#key#">'+
                            '<span class="detail-icon icon-detail-minus"></span> '+
                            '<span class="detail-title">#title#</span> '+
                            '<span class="detail-content" #popoverTemp# >#content#</span>'+
                        '</div>';
            
            temp = temp.replace('#id#',data.id).replace('#title#',data.text).replace('#content#',data.content).replace('#type#',type).replace('#key#',data.key);
            if(data.tips == true){
                var popoverTemp ='data-toggle="popover" data-placement="bottom" data-container="body" data-content="'+data.content+'"';
                temp = temp.replace('#popoverTemp#',popoverTemp);
            }
            temp = $(temp); 

            if(mode !== 'init'){
                temp.css('left','28px');
                temp.find('.detail-icon').css('opacity','1');
            }
            $('.'+type).find('.detail-row-add').before(temp);

            if(data.tips == true){
                $('[data-toggle="popover"]').popover();
            }
        },
        getSelectedAttr:function(){
            var data = this.attritubeData;
            var keys = Object.keys(data);            
            var newData = {};
            for(var i =0;i<keys.length;i++){
                var items = data[keys[i]];
                newData[keys[i]] =[];
                for(j = 0; j<items.length;j++){
                    var v = data[keys[i]][j];
                    if(v.selected && v.selected == true){
                        newData[keys[i]].push({
                            id:v.id,text:v.text
                        });
                    }                
                }
            }            
            return newData;
        },
        initEvent:function(){
            var topThis = this;         
            $('.sidebar-panel-container-new').on('click','.public-cloud-attribute',function(){
                var $this = $(this);
                var content = $($this.data('content'));
                if($this.hasClass('collapse')){
                    content.find(".detail-row").animate({left:'-20px'});
                    content.find(".detail-icon").animate({opacity: "0.2"});
                    content.find(".detail-row-add").animate({left:'-165px'});
                    $this.removeClass('collapse');                    
                }else{
                    content.find(".detail-row").animate({left:'28px'});
                    content.find(".detail-icon").animate({opacity: "1"});
                    content.find(".detail-row-add").animate({left:'28px'});
                    $this.addClass('collapse');
                }
            }).on('click','.icon-detail-minus',function(){
                var data = $(this).parent().prop('id').split('_');
                $(this).parent().remove();
                topThis.formatAttritubeData(data[0],data[1],'delete');
                $.ajax({
                            url: "/host/setDefaultField",
                            type: "POST",
                            data: "key="+data[2] + "&type=d",
                            dataType: "json",
                        });
            }).on('click',function(){
                $('[data-toggle="popover"]').popover('hide');
            }).on('click','.detail-content',function(e){
                e.preventDefault();
                e.stopPropagation();
            });
            $(document).on('scroll',function(e){
               $('[data-toggle="popover"]').popover('hide');
            });
        }
    }


    var data = <?php echo $host;?>;
    var hh = new HostDetail(data);
    hh.init();  //初始化表格
    hh.getSelectedAttr() //获取编辑后的属性列表
});
</script>