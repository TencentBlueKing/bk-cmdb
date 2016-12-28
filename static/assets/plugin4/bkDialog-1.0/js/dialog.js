(function(){
    function dialog(options){
        var defaultOptions = {
            width: 'auto',
            title: '',
            fixed: false,
            zIndex: 1024,
            quickClose: false,
            content: '',
            okValue: '确定',
            ok: null,
            cancelValue: '取消',
            cancel: null,
            onshow: null,
            onclose: null
        };

        var dialogOptions = $.extend({}, defaultOptions, options);
        var dialogNode = null;

        function _init(){
            var _html = [
                        '<div class="bk-dialog" style="display:none; ">',
                        '<div class="bk-dialog-mask"></div>',
                        '<div class="bk-dialog-box" style="width:'+ dialogOptions.width +'px; z-index:'+ dialogOptions.zIndex +'; ">',
                        '    <div class="bk-dialog-header">',
                        '        <strong class="bk-dialog-title">'+ dialogOptions.title +'</strong>',
                        '        <button class="bk-dialog-close">×</button>',
                        '    </div>',
                        '    <div class="bk-dialog-content">'+ dialogOptions.content +'</div>',
                        '    <div class="bk-dialog-footer">',
                        '        <button type="button" class="bk-dialog-btn bk-dialog-ok">'+ dialogOptions.okValue +'</button>',
                        '        <button type="button" class="bk-dialog-btn bk-dialog-cancel">'+ dialogOptions.cancelValue +'</button>',
                        '    </div>',
                        '</div>',
                        '</div>'
                        ].join('');

            dialogNode = $(_html);

            dialogNode.find('.bk-dialog-close').on('click', function(){
                _remove();
            });

            dialogNode.find('.bk-dialog-ok').on('click', function(){
                if (dialogOptions.ok){
                    if (dialogOptions.ok() === false){
                    }else{
                        _remove();
                    }
                }else{
                    _remove();
                }

            });

            dialogNode.find('.bk-dialog-cancel').on('click', function(){
                dialogOptions.cancel && dialogOptions.cancel();
                _remove();
            });

            _render();
            $('body').append(dialogNode);
        }

        function _render(){
            if (!dialogOptions.ok){
                dialogNode.find('.bk-dialog-ok').remove();  
            }
            if (!dialogOptions.cancel){
                dialogNode.find('.bk-dialog-cancel').remove();  
            }
            if (dialogOptions.cancel === false){
                dialogNode.find('.bk-dialog-close').remove();  
            }
            if (!dialogOptions.title){
                dialogNode.find('.bk-dialog-header').remove(); 
            }
            if (!dialogOptions.ok && !dialogOptions.cancel){
                dialogNode.find('.bk-dialog-footer').remove(); 
            }
            if (dialogOptions.fixed){
                dialogNode.find('.bk-dialog-box').css({position: 'fixed'}); 
            }
            if (dialogOptions.quickClose){
                dialogNode.find('.bk-dialog-mask').show();
                dialogNode.find('.bk-dialog-box').on('click', function(){
                    return false;
                });
                dialogNode.find('.bk-dialog-mask').on('click', function(){
                    _close();
                })
            }
        }

        function _show(){
            dialogNode.show();
            dialogOptions.onshow && dialogOptions.onshow();
        }

        function _showModal(){
            dialogNode.show();
            dialogNode.find('.bk-dialog-mask').css('opacity', '0.7').show();
            dialogOptions.onshow && dialogOptions.onshow();
        }

        function _close(){
            dialogNode.hide();
            dialogOptions.onclose && dialogOptions.onclose();
        }

        function _remove(){
            dialogNode.remove();
            dialogOptions.onclose && dialogOptions.onclose();
        }

        function _content(html){
            dialogNode.find('.bk-dialog-content').html(html);
        }

        function _title(text){
            dialogNode.find('.bk-dialog-title').text(text);
        }

        function _width(width){
            dialogNode.width(width);
        }

        function _height(height){
            dialogNode.height(height);
        }

        function Dialog(){
            _init();
        }
        Dialog.prototype = {
            show: function(){
                _show();
                return this;
            },
            showModal: function(){
                _showModal();
                return this;
            },
            close: function(){
                _close();
                return this;
            },
            remove: function(){
                _remove();
                return this;
            },
            content: function(html){
                _content(html);
                return this;
            },
            title: function(text){
                _text(text);
                return this;
            },
            width: function(width){
                _width(width);
                return this;
            },
            height: function(height){
                _height(height);
                return this;
            }
        };

        return new Dialog();

    }
    window.dialog = dialog;
})();