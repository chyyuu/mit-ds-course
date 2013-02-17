/**
 * 推广
 * 
 */
(function($){
	$(function(){
		$('<div id="tgimg" style="z-index: 1; height: 200px; position:absolute; width: 200px; left: 1200px; top:30px;"><a href="http://home.51cto.com/activity/index.php?s=/Indextt/index" target="_blank"><img style="" src="http://home.51cto.com/activity/Tpl/default/Public/images/tg.png"></a></div>').each(function(){ $(this).css('top',$('#home_top').position().top+35);}).appendTo('#home_top .top_nav, #home_top .bbs_top_nav');
	});
})(jQuery);