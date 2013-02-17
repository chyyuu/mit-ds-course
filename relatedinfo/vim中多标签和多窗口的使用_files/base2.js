String.prototype.trim = function(){return this.replace(/(^[ |　]*)|([ |　]*$)/g, "");}
function $(s){return document.getElementById(s);}
function $$(s){return document.frames?document.frames[s]:$(s).contentWindow;}
function $c(s){return document.createElement(s);}
function initSendTime(){
	SENDTIME = new Date();
}
var err;
var cnt=0;
function commentSubmit(theform,refurlk){
	
	if(document.readyState=="complete"){

	var sDialog = new dialog();
	sDialog.init();
	if(!theform){
		sDialog.event("请您先 <a href='http://home.51cto.com/index.php?reback="+ encodeURIComponent(encodeURIComponent(refurlk)) +"' style='color:BLUE;'>登录</a> 或 <a href='http://passport.51cto.com/reg.php?reback=" + refurlk + "' style='color:BLUE;'>注册</a> 后再进行此项操作。",'');
		sDialog.button('dialogOk','void 0');
		$('dialogOk').focus();
		return false;
	}

	var smsg =theform.content.value;
	var susername = theform.username.value;
	var sauthnum = theform.authnum.value;
	var sartID = theform.tid.value;

	if(smsg == ''){
		sDialog.event('请输入评论内容!','');
		sDialog.button('dialogOk','void 0');
		$('dialogOk').focus();
		return false;
	}
	if(sauthnum == ''){
		sDialog.event('请输入验证码!','');
		sDialog.button('dialogOk','void 0');
		$('dialogOk').focus();
		return false;
	}
	if(sartID == ''){
		sDialog.event('不是有效的文章','');
		sDialog.button('dialogOk','void 0');
		$('dialogOk').focus();
		return false;
	}

	
	var url = "/commentcheckform.php?authnum="+sauthnum;
	var ajax = new ActiveXObject("MSXML2.XMLHTTP.3.0");
	ajax.open("GET", url, false);
	ajax.send();
	err=ajax.responseText;
	if(!err){
		var ajax = new new XMLHttpRequest();
		ajax.open("GET", url, false);
		ajax.send();
		err=ajax.responseText;
	}
	if(err == "-1"){
		sDialog.event('验证码输入错误!','');
		sDialog.button('dialogOk','void 0');
		$('dialogOk').focus();
		return false;
	}
	//initSendTime();
	//$("src_title").value = $("commentText").innerHTML;
	//$("src_uname").value = $('feedback').submit();
	}
	cnt++;
	if (cnt!=1){
		return false;
	}
}