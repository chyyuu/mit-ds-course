<!--
var lurl = location.href.split(".");
if(lurl[0]=="http://blog"){
	var domain = lurl[0];
}else{
	var domain = lurl[0]+".blog";
}
function InitAjax(){
	var ajax=false; 
	try { 
		ajax = new ActiveXObject("Msxml2.XMLHTTP.3.0"); 
	} catch (e) { 
		try { 
			ajax = new ActiveXObject("Microsoft.XMLHTTP.3.0"); 
		} catch (E) { 
			ajax = false; 
		} 
	}
	if (!ajax && typeof XMLHttpRequest!='undefined') { 
		ajax = new XMLHttpRequest(); 
	} 
	return ajax;
}
function getcomment(tid,n){
	var url = domain+".51cto.com/comments.php?tid="+tid+"&page="+n;
	var ajax = InitAjax();
	ajax.open("GET", url, true);
	ajax.onreadystatechange = function() {
		if (ajax.readyState == 4 && ajax.status == 200) {
			document.getElementById('artcomment').innerHTML = ajax.responseText;
			document.getElementById('com_top').focus();
		}
	}
	ajax.send(null); 
}
function getcommentend(tid,n){
	var url = domain+".51cto.com/comments.php?tid="+tid+"&page="+n;
	var ajax = InitAjax();
	ajax.open("GET", url, true);
	ajax.onreadystatechange = function() {
		if (ajax.readyState == 4 && ajax.status == 200) {
			document.getElementById('artcomment').innerHTML = ajax.responseText;
			document.getElementById('com_top_top').focus();
		}
	}
	ajax.send(null); 
}
document.onkeydown = function(e){
	var evt = e || window.event ;
	var k = evt.keyCode;
	//var k = event.keyCode;
	if (k == 116){
		window.event.keyCode = 0;
		window.event.returnValue= false;
	}
}
document.onkeyup = function(e){
	var evt = e || window.event ;
	var k = evt.keyCode;
	//var k = event.keyCode;
	if (k == 116){
		window.location.reload();
	}
}
function InitAjax1(){
	var ajax=false; 
	try { 
		ajax = new ActiveXObject("Msxml2.XMLHTTP"); 
	}catch(e){ 
		try { 
			ajax = new ActiveXObject("Microsoft.XMLHTTP"); 
		}catch(E){ 
			ajax = false; 
		}
	}
	if (!ajax && typeof XMLHttpRequest!='undefined') {
		ajax = new XMLHttpRequest(); 
	}
	return ajax;
}

function checkonce(s){
	if(s==1){
		return confirm("你确定删除此篇文章？")
	}else{
		return confirm("你确定删除此条评论？")	
	}
}
function big(o){
	var zoom=parseInt(o.style.zoom, 10)||100;zoom+=window.event.wheelDelta/12;
	if (zoom>0) o.style.zoom=zoom+'%';
	return false; 
}
//-->