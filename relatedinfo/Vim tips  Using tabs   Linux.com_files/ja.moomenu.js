/*------------------------------------------------------------------------
# JA Vauxite for Joomla 1.5 - Version 1.0 - Licence Owner JA85339
# ------------------------------------------------------------------------
# Copyright (C) 2004-2008 J.O.O.M Solutions Co., Ltd. All Rights Reserved.
# @license - Copyrighted Commercial Software
# Author: J.O.O.M Solutions Co., Ltd
# Websites:  http://www.joomlart.com -  http://www.joomlancers.com
# This file may not be redistributed in whole or significant part.
-------------------------------------------------------------------------*/
if (typeof(MooTools) != 'undefined') {
 
    var subnav = new Array();
    Element.implement(
    {
        hide: function(timeout)
        {
            this.status = 'hide';
            clearTimeout(this.timeout);
            if (timeout)
            {
                this.timeout = setTimeout(this.anim.bind(this), timeout);
            } else {
                this.anim();
            }
        },
 
        show: function(timeout)
        {
            this.status = 'show';
            clearTimeout(this.timeout);
            if (timeout)
            {
                this.timeout = setTimeout(this.anim.bind(this), timeout);
            } else {
                this.anim();
            }
        },
 
        setActive: function() {
            this.className+='sfhover';
        },
 
        setDeactive: function() {
            this.className=this.className.replace(new RegExp('sfhover\\b'), '');
        },
 
        anim: function() {
            if ((this.status == 'hide' && this.style.left != 'auto') || (this.status == 'show' && this.style.left == 'auto' && !this.hidding)) return;
            this.setStyle('overflow', 'hidden');
            if (this.status == 'show') {
                this.hidding = 0;
                this.hideAll();
            }
 
            if (this.status == 'hide')
            {
                this.hidding = 1;
                this.myFx2.cancel();
                if (this.parent._id) this.myFx2.start('width', this.offsetWidth, 0);
                    else this.myFx2.start('height', this.offsetHeight, 0);
            } else {
                this.setStyle('left', 'auto');
                this.myFx2.cancel();
                if (this.parent._id) this.myFx2.start('width', 0, this.mw);
                    else this.myFx2.start('height', 0, this.mh);
            }
        },
 
        init: function() {
            this.mw = this.clientWidth;
            this.mh = this.clientHeight;
            this.myFx2 = new Fx.Tween(this, {
                              duration: 300,
                              link: 'cancel'
                              });
            if (this.parent._id) {
                this.myFx2 = new Fx.Tween(this, {
                              duration: 300,
                              link: 'cancel'
                              });
        this.myFx2.set('width',0);
            }else {
                this.myFx2 = new Fx.Tween(this, {
                             duration: 300,
                             link: 'cancel'
                             });
             this.myFx2.set('height',0);
            }
            this.setStyle('left', '-999em');
            animComp = function() {
                if (this.status == 'hide') {
                    this.setStyle('left', '-999em');
                    this.hidding = 0;
                }
                this.setStyle('overflow', '');
            }
            this.myFx2.addEvent('onComplete', animComp.bind(this));
        },
 
        hideAll: function() {
            for (var i = 0; i < subnav.length; i++) {
                if (!this.isChild(subnav[i])) {
                    subnav[i].hide(0);
                }
            }
        },
 
        isChild: function(_obj) {
            obj = this;
            while (obj.parent) {
                if (obj._id == _obj._id) {
                    return true;
                }
                obj = obj.parent;
            }
            return false;
        }
 
 
    });
 
 
    var DropdownMenu = new Class({
        initialize: function(element) {
            $A($(element).childNodes).each(function(el) {
                if (el.nodeName.toLowerCase() == 'li') {
                    $A($(el).childNodes).each(function(el2) {
                        if (el2.nodeName.toLowerCase() == 'ul') {
                            $(el2)._id = subnav.length + 1;
                            $(el2).parent = $(element);
                            subnav.push($(el2));
                            el2.init();
                            el.addEvent('mouseenter', function() {
                                el.setActive();
                                el2.show(0);
                                return false;
                            });
 
                            el.addEvent('mouseleave', function() {
                                el.setDeactive();
                                el2.hide(20);
                            });
                            new DropdownMenu(el2);
                            el.hasSub = 1;
                        }
                    });
                    if (!el.hasSub) {
                        el.addEvent('mouseenter', function() {
                            el.setActive();
                            return false;
                        });
 
                        el.addEvent('mouseleave', function() {
                            el.setDeactive();
                        });
                    }
                }
            });
            return this;
        }
    });
 
    window.addEvent('domready', function() {
        new DropdownMenu($('ja-cssmenu'))
    });
 
} else {
 
    sfHover = function() {
        var sfEls = document.getElementById("ja-cssmenu").getElementsByTagName("li");
        for (var i = 0; i < sfEls.length; ++i) {
            sfEls[i].onmouseover = function() {
                this.className += "sfhover";
            }
            sfEls[i].onmouseout = function() {
                this.className = this.className.replace(new RegExp("sfhover\\b"), "");
            }
        }
    }
    if (window.attachEvent) window.attachEvent("onload", sfHover);
}