var favor_app = 'http://home.51cto.com/apps/favorite/';
var favoriteboxurl = favor_app + 'index.php?s=/Index/box';
var favorcloseimg = favor_app + 'Tpl/default/Public/images/close.gif';
var userAgent = navigator.userAgent.toLowerCase();
var is_opera = userAgent.indexOf('opera') != -1 && opera.version();
var is_moz = (navigator.product == 'Gecko') && userAgent.substr(userAgent.indexOf('firefox') + 8, 3);
var is_ie = (userAgent.indexOf('msie') != -1 && !is_opera) && userAgent.substr(userAgent.indexOf('msie') + 5, 3);

function favorId(id) {
	return document.getElementById(id);
}
function favorBox(action) {
	var objs = document.getElementsByTagName("OBJECT");
	if(action == 'open') {
		for(i = 0;i < objs.length; i ++) {
			if(objs[i].style.visibility != 'hidden') {
				objs[i].setAttribute("oldvisibility", objs[i].style.visibility);
				objs[i].style.visibility = 'hidden';
			}
		}
		var clientWidth = document.body.clientWidth;
		var clientHeight = document.documentElement.clientHeight ? document.documentElement.clientHeight : document.body.clientHeight;
		var scrollTop = document.body.scrollTop ? document.body.scrollTop : document.documentElement.scrollTop;
		var favorwidth = 384;
		var favorheight = 315;
		if(!favorId('favorlayer')) {
			div = document.createElement('div');
            div.id = 'favorlayer';
			div.style.width = favorwidth + 'px';
			div.style.height = favorheight + 'px';
			div.style.left = ((clientWidth - favorwidth) / 2) + 'px';
			div.style.position = 'absolute';
			div.style.zIndex = '999';
			document.body.appendChild(div);
			favorId('favorlayer').innerHTML = '<div style=" margin: 5px auto; text-align: left">' +
				'<div style="width: ' + favorwidth + 'px; height: ' + favorheight + 'px; padding: 1px; background: #FFFFFF; border: 8px solid #666; position: relative; left: -3px; top: -3px">' +
				'<div onmouseover="favorDrag(this)" style="cursor: move; position: relative; left: 0px; top: 0px; width: 392px; height: 30px; margin-bottom: -30px;"></div>' +
				'<a href="javascript:favorBox(\'close\');" target="_self"><img style="position: absolute; right: 28px; top: 19px" src="' + favorcloseimg +'" title="Close" border="0" /></a>' +
				'<div id="favormask" style="margin-top: 30px; position: absolute; width: 100%; height: 100%; display: none"></div><iframe id="favorframe" name="favorframe" style="width:' + favorwidth + 'px;height:100%" allowTransparency="true" frameborder="0"></iframe></div></div>';
		}
		favorId('favorlayer').style.display = '';
		favorId('favorlayer').style.top = ((clientHeight - favorheight) / 2 + scrollTop) + 'px';
        favor_url = favor_url.replace(/\//g,'~');
        favor_title = favor_title.replace(/\//g,'~');
        var encode_url = encodeURIComponent(favor_url);
        var encode_title = encodeURIComponent(favor_title);
        var isgbk = typeof(favor_is_gbk) == 'undefined' ? '' : '/isgbk/1';
        var fuid = '/fuid/' + (typeof(favor_fuid) == 'undefined' ? 0 : favor_fuid);
		favorframe.location = favoriteboxurl + '/url/' + encode_url + '/title/' + encode_title + isgbk + fuid;
	} else if(action == 'close') {
		for(i = 0;i < objs.length; i ++) {
			if(objs[i].attributes['oldvisibility']) {
				objs[i].style.visibility = objs[i].attributes['oldvisibility'].nodeValue;
				objs[i].removeAttribute('oldvisibility');
			}
		}
//		hiddenobj = new Array();
		favorId('favorlayer').style.display = 'none';
	}
}

var favordragstart = new Array();
function favorDrag(obj) {
	obj.onmousedown = function(e) {
		if(is_ie) {
			document.body.onselectstart = function() {
				return false;
			}
		}
		favordragstart = is_ie ? [event.clientX, event.clientY] : [e.clientX, e.clientY];
		favordragstart[2] = parseInt(favorId('favorlayer').style.left);
		favordragstart[3] = parseInt(favorId('favorlayer').style.top);
		favorId('favormask').style.display = '';
		favordoane(e);
	}
	favormove(document.body, 'mousemove', function(e) {
		if(favordragstart[0]) {
			var favorwindragnow = is_ie ? [event.clientX, event.clientY] : [e.clientX, e.clientY];
			with(favorId('favorlayer')) {
				style.left = (favordragstart[2] + favorwindragnow[0] - favordragstart[0]) + 'px';
				style.top = (favordragstart[3] + favorwindragnow[1] - favordragstart[1]) + 'px';
			}
			favordoane(e);
		}
	});
	obj.onmouseup = function(e) {
		if(is_ie) {
			document.body.onselectstart = function() {
				return true;
			}
		}
		favordragstart = [];
		favorId('favormask').style.display = 'none';
		favordoane(e);
	}
	obj.onmouseover = null;
}

function favormove(obj, evt, func) {
	if(obj.addEventListener) {
		obj.addEventListener(evt, func, false);
	} else if(obj.attachEvent) {
		obj.attachEvent("on" + evt, func);
	}
}

function favordoane(event) {
	e = event ? event : window.event;
	if(is_ie) {
		e.returnValue = false;
		e.cancelBubble = true;
	} else if(e) {
		e.stopPropagation();
		e.preventDefault();
	}
}