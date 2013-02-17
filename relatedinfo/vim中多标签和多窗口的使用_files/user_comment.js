String.prototype.trim = function(){return this.replace(/(^[ |　]*)|([ |　]*$)/g, "");}
function $(s){return document.getElementById(s);}
function $$(s){return document.frames?document.frames[s]:$(s).contentWindow;}
function $c(s){return document.createElement(s);}
function initSendTime(){
	SENDTIME = new Date();
}
var err;
var cnt=0
function InitAjax(){
	var ajax=false;
	if(window.XMLHttpRequest){
			var ajax = new XMLHttpRequest();
		}
		else if(window.ActiveXObject){
			var ajax = new ActiveXObject("Microsoft.XMLHTTP");
		}
	return ajax;
}
function get_date()
{
	var date = new Date();
	var datearr = [date.getFullYear(), date.getMonth()+1, date.getDate(), date.getHours(), date.getMinutes(), date.getSeconds()];
	for(var i=1;i<datearr.length;i++)
		datearr[i] = (datearr[i]<10) ? "0" + datearr[i] : datearr[i];
	var str = datearr[0] +"-"+ datearr[1] +"-"+ datearr[2] +" "+ datearr[3] +":"+  datearr[4] +":"+ datearr[5];
	return str;
}
function commentSubmitend(){


	var smsg = document.getElementById("commentcontent").value;
	var susername = document.getElementById("commentusername").value;
	var scount = document.getElementById("com_count_ajax").value;
	var sauthnum="";
	var regx = /\n/gi

	if(smsg.search(regx) != -1)
	   smsg = smsg.replace(regx, "|@|");
	var regz = /—/gi
	if(smsg.search(regz)!=-1)
		smsg = smsg.replace(regz,"|)|");
	var regt=/·/gi
	if(smsg.search(regt)!=-1)
		smsg=smsg.replace(regt,'|(|');
	var sartID = document.getElementById("commenttid").value;
	var sniming=0;
	if (document.getElementById("commentniming"))
	{
		if(document.getElementById("commentniming").checked==true)	{
			sniming=1;
		}
	}
	if(smsg == ''){
		alert('请输入评论内容!');
		return false;
	}

	if(sartID == ''){
		alert('不是有效的文章');
		return false;
	}
	if(document.getElementById("commentauthnum")){
	sauthnum = document.getElementById("commentauthnum").value;
	if(sauthnum == ''){
		alert('请输入验证码!');
		return false;
	}
	var url = "/commentcheckform.php?authnum="+sauthnum;
	var ajax =  new InitAjax();
	ajax.open("GET", url, false);
	ajax.send(null);
	err=ajax.responseText;
	if(!err){
		var ajax =  new InitAjax();
		ajax.open("GET", url, false);
		ajax.send();
		err=ajax.responseText;
	}
	if(err == "-1"){
		alert('验证码输入错误!');
		return false;
	}
	}
	//initSendTime();
	//$("src_title").value = $("commentText").innerHTML;
	//$("src_uname").value = $('feedback').submit();
	var url2="/comment2.php?authnum="+sauthnum+"&tid="+sartID+"&username="+encodeURIComponent(susername)+"&content="+encodeURIComponent(smsg)+"&niming="+sniming;
	var ajax2= new InitAjax();
	ajax2.open("GET",url2,false);
	ajax2.send(null);
	result=ajax2.responseText;
	if(result!='添加评论成功！')
	{
		alert(result);
		return false;
	}
	getcommentend(sartID,scount);
	if(document.getElementById("commentauthnum")){
	refimg();

	document.getElementById("commentauthnum").value='';
	}
	document.getElementById("commentcontent").value='';
}

function commentSubmitend2(){

	Timer();
	var smsg = document.getElementById("commentcontent").value;
	var susername = document.getElementById("commentusername").value;
	//var scount = document.getElementById("com_count_ajax").value;
	var scount = 'e';
	var sauthnum="";
	var parentid=document.getElementById("parentid").value;
	var regx = /\n/gi

	if(smsg.search(regx) != -1)
	   smsg = smsg.replace(regx, "|@|");
	var regt= /\+/g;
	if(smsg.search(regt) != -1)
		smsg=smsg.replace(regt,"＋");
	var sartID = document.getElementById("commenttid").value;
	var sniming=0;
	if (document.getElementById("commentniming"))
	{
		if(document.getElementById("commentniming").checked==true)	{
			sniming=1;
		}
	}
	if(smsg == ''){
		alert('请输入评论内容!');
		return false;
	}

	if(sartID == ''){
		alert('不是有效的文章');
		return false;
	}
	if(document.getElementById("commentyanzheng").style.display!="none"){
	sauthnum = document.getElementById("commentauthnum").value;
	if(sauthnum == ''){
		alert('请输入验证码!');
		return false;
	}
	sauthnum=sauthnum.toLowerCase();
	var url = "/commentcheckform.php?authnum="+sauthnum;
	var ajax =  new InitAjax();
	ajax.open("GET", url, false);
	ajax.send(null);
	err=ajax.responseText;
	if(!err){
		var ajax =  new InitAjax();
		ajax.open("GET", url, false);
		ajax.send();
		err=ajax.responseText;
	}
	if(err == "-1"){
		alert('验证码输入错误!');
		return false;
	}
	}
	//initSendTime();
	//$("src_title").value = $("commentText").innerHTML;
	//$("src_uname").value = $('feedback').submit();
	if(document.getElementById("favour").checked == true) {
		var favour = '1';
	}
	else
	{
		var favour = '0';
	}
	var url2="/comment3.php";
	var re="authnum="+sauthnum+"&tid="+sartID+"&username="+encodeURIComponent(susername)+"&content="+encodeURIComponent(smsg)+"&niming="+sniming + "&favour=" + favour+'&parentid='+parentid;
	var ajax2= new InitAjax();
	ajax2.open("POST",url2,false);
	ajax2.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
	ajax2.send(re);
	result=ajax2.responseText;
	if(result!='添加评论成功！')
	{
		alert(result);
		return false;
	}
	getcommentend(sartID,scount);
	if(document.getElementById("commentauthnum")){
	refimg();

	document.getElementById("commentauthnum").value='';
	}
	document.getElementById("commentcontent").value='';
}
function Timer()
{
	if(document.getElementById("commentsubmit").disabled){
		document.getElementById("commentsubmit").disabled = false;
		return;
	}
	document.getElementById("commentsubmit").disabled = true;
	window.setTimeout("Timer();",4000);
}
function delcomment_ajax(id,selid,reid,recid)
{
	if(checkonce(2)==false)
	{
		return false;
	}
	var url = "/user_index.php?action=delcomment&jop=del&selid="+selid;
	var ajax = new InitAjax();
	ajax.open("GET",url,false);
	ajax.send(null);
	document.getElementById(id).style.display="none";
	document.getElementById(reid).style.display="none";
	document.getElementById(recid).style.display="none";

}
function display_recomment(id,count)
{
	if(count==1)
	{
		document.getElementById(id).style.display="";
	}
	else
	{
		document.getElementById(id).style.display="none";
	}
}
function recommentsubmit(id,reid,contentid,resultid,message,lou)
{
	var content=document.getElementById(contentid).value;
	if(content=="")
	{
		alert("请填写回复内容");
	}
	content2=content;
	var regz = /—/gi
	if(content2.search(regz)!=-1)
		content2 = content2.replace(regz,"|)|");
	var regt=/·/gi
	if(content2.search(regt)!=-1)
		content2=content2.replace(regt,'|(|');
	var bozhumessage=0;
	if(document.getElementById(message).checked==true)
	{
		bozhumessage=1;
	}
	var url2="/recomment.php";
	var re="cid="+reid+"&content="+encodeURIComponent(content2)+"&bozhumessage="+bozhumessage+"&lou="+lou;
	var ajax2= new InitAjax();
	ajax2.open("POST",url2,false);
	ajax2.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
	ajax2.send(re);
	result=ajax2.responseText;
	if(isNaN(result))
	{
		alert(result);
		return false;
	}
	var regx = /\n/gi

	if(content.search(regx) != -1)
	   content = content.replace(regx, "<br/>");
	document.getElementById(id).style.display="none";
	document.getElementById(resultid).innerHTML="<br/><b>博主回复:</b><br/>"+content+"<br/>"+get_date()+"<a href='javascript:void(0);'  class='operlink' onclick='delrecomment_ajax(this,"+result+")' >删除</a>";
}
function delrecomment_ajax(id,tid,hrefid)
{
	if(checkonce(2)==false)
	{
		return false;
	}
	var url = "/user_index.php?action=delrecomment&jop=del&selid="+tid;
	var ajax = new InitAjax();
	ajax.open("GET",url,false);
	ajax.send(null);
	result=ajax.responseText;
	if(result!='')
	{
		alert(result);
	}
	id.parentNode.innerHTML = "";
	if(document.getElementById(hrefid))
	{document.getElementById(hrefid).style.display="";}
	//p = id.parentNode.parentNode;
	//p.removeChild(p.lastChild);
}