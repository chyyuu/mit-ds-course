(function(){  

	var ua=navigator.userAgent.toLowerCase();  

	var is=(ua.match(/\b(chrome|opera|safari|msie|firefox)\b/) || ['','mozilla'])[1];  

	var r='(?:'+is+'|version)[\\/: ]([\\d.]+)';  

	var v=(ua.match(new RegExp(r)) ||[])[1];  

	jQuery.browser.is=is;  

	jQuery.browser.ver=v;  

	jQuery.browser[is]=true;  

})(); 

(function (jQuery){

/*

 * jQuery Plugin - Messager

 * Author: corrie	Mail: corrie@sina.com	Homepage: www.corrie.net.cn

 * Copyright (c) 2008 corrie.net.cn

 * @license http://www.gnu.org/licenses/gpl.html [GNU General Public License]

 *

 * $Date: 2008-12-26 

 * $Vesion: 1.4

 @ how to use and example: Please Open demo.html

 */

	this.version = '@1.3';

	this.layer = {'width' : 389, 'height': 196};//200 100

	this.title = '51cto博客';

	this.time = 150000;//4000

	this.anims = {'type' : 'slide', 'speed' : 600};
	this.timer1 = null;


	this.inits = function(title, text){

		if(jQuery("#message").is("div")){ return; }
		var topHeight = document.documentElement.scrollTop + document.documentElement.clientHeight - this.layer.height -10;
		jQuery(document.body).prepend('<div id="message" class="blogpop" style="background-color:#FFFFFF;z-index:100;width:'+this.layer.width+'px;height:'+this.layer.height+'px;position:absolute; display:none; bottom:0; top:'+ topHeight +'px;right:0; overflow:hidden;"><iframe style="width:800px;height:536px;top:0px;left:0px;position:absolute;visibility:inherit;z-index:-1;" frameborder=0 ></iframe><div style="height:36px;background:url(http://blog.51cto.com/image/blog_top/popbgx.jpg) repeat-x; position:relative;width:387px;border-right-color:#2668B7;border-right-width:2px;border-right-style:solid 2px #2668B7;"><span style="right:10px;padding:12px 0 0;position:absolute;top:0"><span style="margin-right:80px;color:#828282;width:190px;padding:10px 0 0 10px;font-size:14px">每日博报 精彩不止一点</span><img title="关闭" id="message_close" src="http://blog.51cto.com/image/blog_top/close.jpg" alt="关闭" /></span><div style="width:387px;height:36px;left:0;overflow:hidden; text-indent:-999em;background:url(http://blog.51cto.com/image/blog_top/hdtit.jpg) no-repeat;border-right:solid 2px #2668B7;"><a style="display: block; height: 36px;width:138px;" target="_blank" href="http://blog.51cto.com"></a>'+title+'</div><div style="clear:both;"></div></div><div id="message_content" style="border-top:none;font-family:Microsoft YaHei;font-size:12px;width:389px;height:160pxpx;text-align:left;overflow:hidden;">'+text+'</div></div>');
		

		jQuery("#message_close").click(function(){		
//            jQuery("#message").remove();
			setTimeout('this.mclose()', 1);

		});
		jQuery("#message").hover(function(){
			clearTimeout(timer1);
			timer1 = null;
		},function(){
			timer1 = setTimeout('this.mclose()', time);
			//alert(timer1);
		});

	};

	this.showblogtop = function(title, text, time){

		if(jQuery("#message").is("div")){ return; }

		if(title==0 || !title)title = this.title;

		this.inits(title, text);

		if(time>=0)this.time = time;

		switch(this.anims.type){

			case 'slide':jQuery("#message").slideDown(this.anims.speed);break;

			case 'fade':jQuery("#message").fadeIn(this.anims.speed);break;

			case 'show':jQuery("#message").show(this.anims.speed);break;

			default:jQuery("#message").slideDown(this.anims.speed);break;

		}

		if(jQuery.browser.is=='chrome'){

			setTimeout(function(){

				jQuery("#message").remove();

				this.inits(title, text);

				jQuery("#message").css("display","block");

			},this.anims.speed-(this.anims.speed/5));

		}

		//$("#message").slideDown('slow');

		this.rmmessage(this.time);

	};

	this.lays = function(width, height){

		if(jQuery("#message").is("div")){ return; }

		if(width!=0 && width)this.layer.width = width;

		if(height!=0 && height)this.layer.height = height;

	}

	this.anim = function(type,speed){

		if(jQuery("#message").is("div")){ return; }

		if(type!=0 && type)this.anims.type = type;

		if(speed!=0 && speed){

			switch(speed){

				case 'slow' : ;break;

				case 'fast' : this.anims.speed = 200; break;

				case 'normal' : this.anims.speed = 400; break;

				default:					

					this.anims.speed = speed;

			}			

		}

	}

	this.rmmessage = function(time){

		if(time>0){

			timer1 = setTimeout('this.mclose()', time);

			//setTimeout('$("#message").remove()', time+1000);

		}

	};
	this.mclose = function(){
		switch(this.anims.type){
			case 'slide':jQuery("#message").slideUp(this.anims.speed);break;
			case 'fade':jQuery("#message").fadeOut(this.anims.speed);break;
			case 'show':jQuery("#message").hide(this.anims.speed);break;
			default:jQuery("#message").slideUp(this.anims.speed);break;
		};
		setTimeout('jQuery("#message").remove();', this.anims.speed);
		//this.original();	
	}

	this.original = function(){	

		//this.layer = {'width' : 200, 'height': 100};390 198
		this.layer = {'width' : 389, 'height': 196};

		this.title = '51cto博客';

		//this.time = 4000;
		this.time = 150000;

		this.anims = {'type' : 'slide', 'speed' : 600};

	};

    jQuery.messager = this;

    return jQuery;

})(jQuery);

jQuery(window).scroll( function() {
	if(navigator.userAgent.indexOf("Chrome")>0){
	  var height = document.body.scrollTop;
	}else{
	  var height = document.documentElement.scrollTop;
	}
	var topHeight = height + document.documentElement.clientHeight - this.layer.height -10;
	var max = jQuery(document.body).height() - this.layer.height -10;

	if (topHeight > max) {
		topHeight = max;
	}

	//jQuery("#message").animate({"top":topHeight+"px"},90);						   
	jQuery("#message").css("top",topHeight+"px");						
});
