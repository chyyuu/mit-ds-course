var lang = new Array();
var userAgent = navigator.userAgent.toLowerCase();
var is_opera = userAgent.indexOf('opera') != -1 && opera.version();
var is_moz = (navigator.product == 'Gecko') && userAgent.substr(userAgent.indexOf('firefox') + 8, 3);
var is_ie = (userAgent.indexOf('msie') != -1 && !is_opera) && userAgent.substr(userAgent.indexOf('msie') + 5, 3);

function $(id) {
	return document.getElementById(id);
}
function doane(event) {
	e = event ? event : window.event;
	if(is_ie) {
		e.returnValue = false;
		e.cancelBubble = true;
	} else if(e) {
		e.stopPropagation();
		e.preventDefault();
	}
}
var hiddenobj = new Array();
var pmwinposition = new Array();
function pmwin(action, param) {
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
		var pmwidth = 800;
		var pmheight = clientHeight * 0.9;
		if(!$('pmlayer')) {
			div = document.createElement('div');div.id = 'pmlayer';
			div.style.width = pmwidth + 'px';
			div.style.height = pmheight + 'px';
			div.style.left = ((clientWidth - pmwidth) / 2) + 'px';
			div.style.position = 'absolute';
			div.style.zIndex = '999';
			$('append_parent').appendChild(div);
			$('pmlayer').innerHTML = '<div style="width: 800px; background: #666666; margin: 5px auto; text-align: left">' +
				'<div style="width: 800px; height: ' + pmheight + 'px; padding: 1px; background: #FFFFFF; border: 1px solid #7597B8; position: relative; left: -6px; top: -3px">' +
				'<div onmouseover="pmwindrag(this)" style="cursor: move; position: relative; left: 0px; top: 0px; width: 800px; height: 30px; margin-bottom: -30px;"></div>' +
				'<a href="javascript:void(0)" onclick="pmwin(\'close\')" title="¹Ø±Õ" target="_self"><div style="cursor:pointer;position: absolute; right: 20px; top: 15px;border: 0px ">X</div></a>' +
				'<div id="pmwinmask" style="margin-top: 30px; position: absolute; width: 100%; height: 100%; display: none"></div><iframe id="pmframe" name="pmframe" style="width:' + pmwidth + 'px;height:100%" allowTransparency="true" frameborder="0"></iframe></div></div>';
		}
		$('pmlayer').style.display = '';
		$('pmlayer').style.top = ((clientHeight - pmheight) / 2 + scrollTop) + 'px';
		if(!param) {
			pmframe.location = 'http://blog.51cto.com/pm.php';
		} else {
			pmframe.location = 'http://blog.51cto.com/pm.php?msgto='+param;
		}
		dis_select=document.getElementsByTagName('select');
		for (i=0;i<dis_select.length ; i++)
		{
			dis_select[i].style.display="none";
		}
	} else if(action == 'close') {
		for(i = 0;i < objs.length; i ++) {
			if(objs[i].attributes['oldvisibility']) {
				objs[i].style.visibility = objs[i].attributes['oldvisibility'].nodeValue;
				objs[i].removeAttribute('oldvisibility');
			}
		}
		hiddenobj = new Array();
		$('pmlayer').style.display = 'none';
		dis_select=document.getElementsByTagName('select');
		for (i=0;i<dis_select.length ; i++)
		{
			dis_select[i].style.display="";
		}
	}
}

var pmwindragstart = new Array();
function pmwindrag(obj) {
	obj.onmousedown = function(e) {
		if(is_ie) {
			document.body.onselectstart = function() {
				return false;
			}
		}
		pmwindragstart = is_ie ? [event.clientX, event.clientY] : [e.clientX, e.clientY];
		pmwindragstart[2] = parseInt($('pmlayer').style.left);
		pmwindragstart[3] = parseInt($('pmlayer').style.top);
		$('pmwinmask').style.display = '';
		doane(e);
	}
	_attachEvent(document.body, 'mousemove', function(e) {
		if(pmwindragstart[0]) {
			var pmwindragnow = is_ie ? [event.clientX, event.clientY] : [e.clientX, e.clientY];
			with($('pmlayer')) {
				style.left = (pmwindragstart[2] + pmwindragnow[0] - pmwindragstart[0]) + 'px';
				style.top = (pmwindragstart[3] + pmwindragnow[1] - pmwindragstart[1]) + 'px';
			}
			doane(e);
		}
	});
	obj.onmouseup = function(e) {
		if(is_ie) {
			document.body.onselectstart = function() {
				return true;
			}
		}
		pmwindragstart = [];
		$('pmwinmask').style.display = 'none';
		doane(e);
	}
	obj.onmouseover = null;
}
function _attachEvent(obj, evt, func) {
	if(obj.addEventListener) {
		obj.addEventListener(evt, func, false);
	} else if(obj.attachEvent) {
		obj.attachEvent("on" + evt, func);
	}
}