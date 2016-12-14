/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

function _init() {
    $.AdminLTE.layout = {
        activate: function() {
            var a = this;
            a.fix(), a.fixSidebar(), $(window, ".wrapper").resize(function() {
                a.fix(), a.fixSidebar()
            })
        },
        fix: function() {
            var a = $(".main-header").outerHeight() + $(".main-footer").outerHeight(),
                b = $(window).height(),
                c = $(".sidebar").height();
            if ($("body").hasClass("fixed")) $(".content-wrapper, .right-side").css("min-height", b - $(".main-footer").outerHeight());
            else {
                var d;
                b >= c ? ($(".content-wrapper, .right-side").css("min-height", b - a), d = b - a) : ($(".content-wrapper, .right-side").css("min-height", c), d = c);
                var e = $($.AdminLTE.options.controlSidebarOptions.selector);
                "undefined" != typeof e && e.height() > d && $(".content-wrapper, .right-side").css("min-height", e.height())
            }
        },
        fixSidebar: function() {
            return $("body").hasClass("fixed") ? ("undefined" == typeof $.fn.slimScroll && console && console.error("Error: the fixed layout requires the slimscroll plugin!"), void($.AdminLTE.options.sidebarSlimScroll && "undefined" != typeof $.fn.slimScroll && ($(".sidebar").slimScroll({
                destroy: !0
            }).height("auto"), $(".sidebar").slimscroll({
                height: $(window).height() - $(".main-header").height() + "px",
                color: "rgba(0,0,0,0.2)",
                size: "3px"
            })))) : void("undefined" != typeof $.fn.slimScroll && $(".sidebar").slimScroll({
                destroy: !0
            }).height("auto"))
        }
    }, $.AdminLTE.pushMenu = {
        activate: function(a) {
            var b = $.AdminLTE.options.screenSizes;
            $(a).on("click", function(a) {
                a.preventDefault(), $(window).width() > b.sm - 1 ? $("body").hasClass("sidebar-collapse") ? $("body").removeClass("sidebar-collapse").trigger("expanded.pushMenu") : $("body").addClass("sidebar-collapse").trigger("collapsed.pushMenu") : $("body").hasClass("sidebar-open") ? $("body").removeClass("sidebar-open").removeClass("sidebar-collapse").trigger("collapsed.pushMenu") : $("body").addClass("sidebar-open").trigger("expanded.pushMenu")
            }), $(".content-wrapper").click(function() {
                $(window).width() <= b.sm - 1 && $("body").hasClass("sidebar-open") && $("body").removeClass("sidebar-open")
            }), ($.AdminLTE.options.sidebarExpandOnHover || $("body").hasClass("fixed") && $("body").hasClass("sidebar-mini")) && this.expandOnHover()
        },
        expandOnHover: function() {
            var a = this,
                b = $.AdminLTE.options.screenSizes.sm - 1;
            $(".main-sidebar").hover(function() {
                $("body").hasClass("sidebar-mini") && $("body").hasClass("sidebar-collapse") && $(window).width() > b && a.expand()
            }, function() {
                $("body").hasClass("sidebar-mini") && $("body").hasClass("sidebar-expanded-on-hover") && $(window).width() > b && a.collapse()
            })
        },
        expand: function() {
            $("body").removeClass("sidebar-collapse").addClass("sidebar-expanded-on-hover")
        },
        collapse: function() {
            $("body").hasClass("sidebar-expanded-on-hover") && $("body").removeClass("sidebar-expanded-on-hover").addClass("sidebar-collapse")
        }
    }, $.AdminLTE.tree = function(a) {
        var b = this,
            c = $.AdminLTE.options.animationSpeed;
        $("li a", $(a)).on("click", function(a) {
            var d = $(this),
                e = d.next();
            if (e.is(".treeview-menu") && e.is(":visible")) e.slideUp(c, function() {
                e.removeClass("menu-open")
            }), e.parent("li").removeClass("active");
            else if (e.is(".treeview-menu") && !e.is(":visible")) {
                var f = d.parents("ul").first(),
                    g = f.find("ul:visible").slideUp(c);
                g.removeClass("menu-open");
                var h = d.parent("li");
                e.slideDown(c, function() {
                    e.addClass("menu-open"), f.find("li.active").removeClass("active"), h.addClass("active"), b.layout.fix()
                })
            }
            e.is(".treeview-menu") && a.preventDefault()
        })
    }, $.AdminLTE.controlSidebar = {
        activate: function() {
            var a = this,
                b = $.AdminLTE.options.controlSidebarOptions,
                c = $(b.selector),
                d = $(b.toggleBtnSelector);
            d.on("click", function(d) {
                d.preventDefault(), c.hasClass("control-sidebar-open") || $("body").hasClass("control-sidebar-open") ? a.close(c, b.slide) : a.open(c, b.slide)
            });
            var e = $(".control-sidebar-bg");
            a._fix(e), $("body").hasClass("fixed") ? a._fixForFixed(c) : $(".content-wrapper, .right-side").height() < c.height() && a._fixForContent(c)
        },
        open: function(a, b) {
            b ? a.addClass("control-sidebar-open") : $("body").addClass("control-sidebar-open")
        },
        close: function(a, b) {
            b ? a.removeClass("control-sidebar-open") : $("body").removeClass("control-sidebar-open")
        },
        _fix: function(a) {
            var b = this;
            $("body").hasClass("layout-boxed") ? (a.css("position", "absolute"), a.height($(".wrapper").height()), $(window).resize(function() {
                b._fix(a)
            })) : a.css({
                position: "fixed",
                height: "auto"
            })
        },
        _fixForFixed: function(a) {
            a.css({
                position: "fixed",
                "max-height": "100%",
                overflow: "auto",
                "padding-bottom": "50px"
            })
        },
        _fixForContent: function(a) {
            $(".content-wrapper, .right-side").css("min-height", a.height())
        }
    }, $.AdminLTE.boxWidget = {
        selectors: $.AdminLTE.options.boxWidgetOptions.boxWidgetSelectors,
        icons: $.AdminLTE.options.boxWidgetOptions.boxWidgetIcons,
        animationSpeed: $.AdminLTE.options.animationSpeed,
        activate: function(a) {
            var b = this;
            a || (a = document), $(a).find(b.selectors.collapse).on("click", function(a) {
                a.preventDefault(), b.collapse($(this))
            }), $(a).find(b.selectors.remove).on("click", function(a) {
                a.preventDefault(), b.remove($(this))
            })
        },
        collapse: function(a) {
            var b = this,
                c = a.parents(".box").first(),
                d = c.find("> .box-body, > .box-footer, > form  >.box-body, > form > .box-footer");
            c.hasClass("collapsed-box") ? (a.children(":first").removeClass(b.icons.open).addClass(b.icons.collapse), d.slideDown(b.animationSpeed, function() {
                c.removeClass("collapsed-box")
            })) : (a.children(":first").removeClass(b.icons.collapse).addClass(b.icons.open), d.slideUp(b.animationSpeed, function() {
                c.addClass("collapsed-box")
            }))
        },
        remove: function(a) {
            var b = a.parents(".box").first();
            b.slideUp(this.animationSpeed)
        }
    }
}
if ("undefined" == typeof jQuery) throw new Error("AdminLTE requires jQuery");
$.AdminLTE = {}, $.AdminLTE.options = {
    navbarMenuSlimscroll: !0,
    navbarMenuSlimscrollWidth: "3px",
    navbarMenuHeight: "200px",
    animationSpeed: 500,
    sidebarToggleSelector: "[data-toggle='offcanvas']",
    sidebarPushMenu: !0,
    sidebarSlimScroll: !0,
    sidebarExpandOnHover: !1,
    enableBoxRefresh: !0,
    enableBSToppltip: !0,
    BSTooltipSelector: "[data-toggle='tooltip']",
    enableFastclick: !0,
    enableControlSidebar: !0,
    controlSidebarOptions: {
        toggleBtnSelector: "[data-toggle='control-sidebar']",
        selector: ".control-sidebar",
        slide: !0
    },
    enableBoxWidget: !0,
    boxWidgetOptions: {
        boxWidgetIcons: {
            collapse: "fa-minus",
            open: "fa-plus",
            remove: "fa-times"
        },
        boxWidgetSelectors: {
            remove: '[data-widget="remove"]',
            collapse: '[data-widget="collapse"]'
        }
    },
    directChat: {
        enable: !0,
        contactToggleSelector: '[data-widget="chat-pane-toggle"]'
    },
    colors: {
        lightBlue: "#3c8dbc",
        red: "#f56954",
        green: "#00a65a",
        aqua: "#00c0ef",
        yellow: "#f39c12",
        blue: "#0073b7",
        navy: "#001F3F",
        teal: "#39CCCC",
        olive: "#3D9970",
        lime: "#01FF70",
        orange: "#FF851B",
        fuchsia: "#F012BE",
        purple: "#8E24AA",
        maroon: "#D81B60",
        black: "#222222",
        gray: "#d2d6de"
    },
    screenSizes: {
        xs: 480,
        sm: 768,
        md: 992,
        lg: 1200
    }
}, $(function() {
    "undefined" != typeof AdminLTEOptions && $.extend(!0, $.AdminLTE.options, AdminLTEOptions);
    var a = $.AdminLTE.options;
    _init(), $.AdminLTE.layout.activate(), $.AdminLTE.tree(".sidebar"), a.enableControlSidebar && $.AdminLTE.controlSidebar.activate(), a.navbarMenuSlimscroll && "undefined" != typeof $.fn.slimscroll && $(".navbar .menu").slimscroll({
        height: a.navbarMenuHeight,
        alwaysVisible: !1,
        size: a.navbarMenuSlimscrollWidth
    }).css("width", "100%"), a.sidebarPushMenu && $.AdminLTE.pushMenu.activate(a.sidebarToggleSelector), a.enableBSToppltip && $("body").tooltip({
        selector: a.BSTooltipSelector
    }), a.enableBoxWidget && $.AdminLTE.boxWidget.activate(), a.enableFastclick && "undefined" != typeof FastClick && FastClick.attach(document.body), a.directChat.enable && $(a.directChat.contactToggleSelector).on("click", function() {
        var a = $(this).parents(".direct-chat").first();
        a.toggleClass("direct-chat-contacts-open")
    }), $('.btn-group[data-toggle="btn-toggle"]').each(function() {
        var a = $(this);
        $(this).find(".btn").on("click", function(b) {
            a.find(".btn.active").removeClass("active"), $(this).addClass("active"), b.preventDefault()
        })
    })
}),
    function(a) {
        a.fn.boxRefresh = function(b) {
            function c(a) {
                a.append(f), e.onLoadStart.call(a)
            }

            function d(a) {
                a.find(f).remove(), e.onLoadDone.call(a)
            }
            var e = a.extend({
                    trigger: ".refresh-btn",
                    source: "",
                    onLoadStart: function(a) {},
                    onLoadDone: function(a) {}
                }, b),
                f = a('<div class="overlay"><div class="fa fa-refresh fa-spin"></div></div>');
            return this.each(function() {
                if ("" === e.source) return void(console && console.log("Please specify a source first - boxRefresh()"));
                var b = a(this),
                    f = b.find(e.trigger).first();
                f.on("click", function(a) {
                    a.preventDefault(), c(b), b.find(".box-body").load(e.source, function() {
                        d(b)
                    })
                })
            })
        }
    }(jQuery),
    function(a) {
        a.fn.activateBox = function() {
            a.AdminLTE.boxWidget.activate(this)
        }
    }(jQuery),
    function(a) {
        a.fn.todolist = function(b) {
            var c = a.extend({
                onCheck: function(a) {},
                onUncheck: function(a) {}
            }, b);
            return this.each(function() {
                "undefined" != typeof a.fn.iCheck ? (a("input", this).on("ifChecked", function(b) {
                    var d = a(this).parents("li").first();
                    d.toggleClass("done"), c.onCheck.call(d)
                }), a("input", this).on("ifUnchecked", function(b) {
                    var d = a(this).parents("li").first();
                    d.toggleClass("done"), c.onUncheck.call(d)
                })) : a("input", this).on("change", function(b) {
                    var d = a(this).parents("li").first();
                    d.toggleClass("done"), c.onCheck.call(d)
                })
            })
        }
    }(jQuery);

window.CC = window.CC || {};

/**
 * CC.rightPanel
 */
void (function(w,name) {
    var a = $,
    n = {
         supportTransform: function() {
            var t = document.body.style;
            return "WebkitTransform" in t || "MozTransform" in t || "OTransform" in t || "Transform" in t || "transform" in t ? !0: !1
        },
        setTransformTransitionForElem: function(e, t, n) {
            e && (void 0 !== t && (e.style.WebkitTransform = t, e.style.MozTransform = t, e.style.OTransform = t, e.style.msTransform = t, e.style.Transform = t, e.style.transform = t), void 0 !== n && (e.style.WebkitTransition = n, e.style.MozTransition = n, e.style.OTransition = n, e.style.msTransition = n, e.style.Transition = n, e.style.transition = n))
        }
    },
    o = {
        main: '<div class="sidebar-panel-container" data-child="innerContainer"></div>'
    },
    r = "panelContain_" + Math.floor(1e4 * Math.random()),
    l = {
        rendTo: $('body'),
        width: "860px",
        height: "100%",
        onBeforeShow: null,
        onAfterShow: null,
        onBeforeHide: null,
        onAfterHide: null,
        module: {},
        data: {},
        speed: 300,
        transition: null
    },
    c = {},
    d = (document.body.style, n.supportTransform()),
    p = 0,
    u = !1,
    time,
    h = n.setTransformTransitionForElem,
    m = function(e) {
        if((typeof e=='undefined') || a(e.target).attr('id')==='closePanleBtn'){
            w[name].hide();
        }
    };
    $(document).bind("onmousedown" in document?"mousedown": "click", function(e){
        var target  = $(e.target);
            if(target.closest(".sidebar-panel").length == 0){
                time=setTimeout(m,300);
            }
    });
    w[name] = {
        show: function(t, e) {
            clearTimeout(time),
            t || (t = {}),
            c = {};
            for (var i in l) c[i] = l[i];
            c = a.extend(c, t);
            var n = document.getElementById(r);
            if (n || (n = document.createElement("div"), a(n).attr("id", r), n.className = "sidebar-panel", d && h(n, "translate3d(100%,0px,0px)", ""), a(n).html(o.main), c.rendTo.append(n), u = !0), n.style.height = c.height, "always" == c.transition ? this.hide() : c.module && c.module.destroy && c.module.destroy(), a("#" + r + " [data-child=innerContainer]").empty(), "panel" != c.rendTo.attr("data-parent") && (a(n).parent().length && (a(n).parent().removeAttr("data-parent"), a(n).parent()[0].removeChild(n)), c.rendTo.attr("data-parent", "panel"), c.rendTo.append(n), u = !0), c.module && c.module.render && c.module.render(a("#" + r + " [data-child=innerContainer]"), c.data), !e && c.onBeforeShow && c.onBeforeShow(), d) {
                var s = (100 / parseInt(c.speed)).toFixed(1);
                p && clearTimeout(p),
                u ? setTimeout(function() {
                    n.style.width = c.width,
                    h(n, "translate3d(0px,0px,0px)", "all " + s + "s ease-out")
                },
                10) : (n.style.width = c.width, h(n, "translate3d(0px,0px,0px)", "all " + s + "s ease-out")),
                p = setTimeout(function() {
                    h(n, void 0, ""),
                    !e && c.onAfterShow && c.onAfterShow()
                },
                1e3 * s)
            } else a(n).show(),
            n.style.width = "0px",
            a(n).stop(!0, !0),
            a(n).animate({
                width: c.width
            },
            c.speed, "swing", 
            function() { ! e && c.onAfterShow && c.onAfterShow()
            })

            a(n).bind("click", m)
        },
        hide: function(t) {
            c.module && c.module.destroy && c.module.destroy();
            var e = document.getElementById(r);
            if (e) if (!t && c.onBeforeHide && c.onBeforeHide(), d) {
                var i = (100 / parseInt(c.speed)).toFixed(1);
                p && clearTimeout(p),
                h(e, "translate3d(100%,0px,0px)", "all " + i + "s ease-out"),
                p = setTimeout(function() { ! t && c.onAfterHide && c.onAfterHide(),
                    h(e, void 0, "")
                },
                1e3 * i)
            } else a(e).stop(!0, !0),
            a(e).animate({
                width: "0px"
            },
            c.speed, "swing", 
            function() {
                a(e).hide(),
                !t && c.onAfterHide && c.onAfterHide()
            })
        },
        render: function(t) {
            var e = document.getElementById(r);
            e && (t ? a("#" + r + " [data-child=innerContainer]").html(t) : e.empty())
        }
    }
})(window.CC,'rightPanel');

$(function(){
    $('#speed_search').click(function(){
        var content=$('#search-textarea').val();
        /*没有输入任何值*/
        if(!content){
            return true;
        }

        /*多值*/
        if (content.indexOf("\n") > 0 || content.indexOf(",") > 0) {
            var content = content.replace(/\s+/g, ',');/*换行转逗号*/
            var data = content.split(',');
            var parame = '';
            if (data[0].indexOf(".") > 0) {/*输入的值为IP*/
                var ip = data[0].split('.');
                var arr = ["10", "172", "192"];
                var result = $.inArray(ip[0], arr);
                if(result >= 0){
                    $('#search-textarea').attr('name', 'InnerIP').val(content);
                }else{
                    $('#search-textarea').attr('name', 'OuterIP').val(content);
                }
            }else{/*非IP*/
                $('#search-textarea').attr('name', 'AssetID').val(content);
            }
            $('#quick_search').submit();
            return true;
        }else{/*单值*/
            var parame = '';
            var flag = false;/*非法标志*/
            if (content.indexOf(".") > 0) {/*IP*/
                var ip = content.split('.');
                if(ip.length != 4){
                    flag = true;
                }

                $.each(ip, function(i, item){
                    if(item > 255 || item < 0){
                        flag = true;
                    }

                    if(isNaN(item)){
                        flag = true;
                    }
                });
                
                if(flag){
                    var arr = ["10", "172", "192"];
                    var result = $.inArray(ip[0], arr);
                    if(result >= 0){
                        $('#search-textarea').attr('name', 'InnerIP').val(content);
                    }else{
                        $('#search-textarea').attr('name', 'OuterIP').val(content);
                    }

                    $('#quick_search').submit();
                    return true;
                }

                var arr = ["10", "172", "192"];
                var result = $.inArray(ip[0], arr);
                if(result >= 0){
                    parame = 'InnerIP=' + content;
                }else{
                    parame = 'OuterIP=' + content;
                }
            }else{/*固资*/
                parame = 'AssetID=' + content;
            }

            parame = parame + '&ApplicationID='+ cookie.get('defaultAppId');
            $.ajax({
                url : "/host/details",
                type: "POST",
                data: parame,
                dataType:"html",
                success : function(data){
                    
                    if(Object.prototype.toString.call(data) === '[object Array]'){
                        var obj = $.parseJSON(data);
                        var setDefaultApp = $('#speed_search_input').get(0);
                        var diaCopyMsg = dialog({
                            quickClose: true,// 点击空白处快速关闭
                            align:'left',
                            //fixed:true,
                            padding:'5px 5px 5px 30px',
                            skin:'c-Popuplayer-remind-left',
                            width:160,
                            content: '<span style="color:#fff">' + obj.errInfo + '</span>'
                        });
                        diaCopyMsg.show(setDefaultApp);
                        setTimeout(function () {
                            diaCopyMsg.close().remove();
                        }, 3000);
                        return true;
                    }

                    CC.rightPanel.show();
                    CC.rightPanel.render(data);
                }
            });
        }
    });

    /*  顶栏搜索文本域触发下拉事件 */
    $('.search-textarea').focus(function(){
        $(this).addClass('search-textarea-down');
    }).blur(function(){
        $(this).removeClass('search-textarea-down');
    });

    $('.dropdown-toggle').dropdown();

    $.fn.divLoad = function(options){
        var $this = $(this);
        $this.addClass('divLoad-parent');
        $this.find('.loading-warp').remove();
        var width = $this.css('width').replace('px','');
        var height = $this.css('height').replace('px','');
        if(options == "show"){
            var html = '<div class="loading-warp divLoad-shade"></div>';
            $this.append(html);
            $this.find('.loading-warp').css('height',height).css('width',width);
        }else if(options == "hide"){
            $this.find('.loading-warp').remove();
            $this.removeClass('divLoad-parent');
        }
    };

    $('.search-textarea').click(function(){
        $(this).attr('placeholder','内网IP/外网IP/固资编号');
    });

    $('.search-textarea').blur(function(){
        var value = $(this).val();
        if(!value){
            $(this).attr('placeholder','快速查询');
        }
    });
});