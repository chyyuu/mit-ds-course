function dfa(tid,uid,domains,homeLogUrl){
	var url = "http://"+domains+".blog.51cto.com/joinbowen.php?tidd="+tid+"&uuid="+uid+"&domains="+domains ;
	postRequest(url,homeLogUrl);
	
	return false;
}

function dfanologin(homelogurl){
	if(homelogurl)
	{
		window.location.href = homelogurl;	
		return false;
	}
	else
	{	
		var scrollTop = window.pageYOffset 
                || document.documentElement.scrollTop 
                || document.body.scrollTop 
                || 0;
		document.getElementById("add_sys_type_divs").style.top=(scrollTop+140)+"px";
		document.getElementById("add_sys_type_divs").innerHTML="<table width='491' height='100' border='0' align='center' cellpadding='0' cellspacing='0'><tr><td width='491' height=53><div align='center'><strong><font color=#000000>对不起,您还没登录,<a href='http://blog.51cto.com/login.php'><font color=#ff0000><u>请登录</u></font></a>后再推送</font></strong></div></td></tr><tr><td height='47'><div align='center'>&nbsp;<label><input type='button' name='Submit22' value='关闭' onclick=\"javascript:document.getElementById('add_sys_type_divs').style.display='none';\"></label></div></td></tr></table>";		
		document.getElementById("add_sys_type_divs").style.display="";
		return false;
	}	
}

function dass(tid,uid,domains){
	var j=0;
	var ppp="";
	document.formsp.sos.value="";
	document.formsp.sosqu.value="";
	document.formsp.sosming.value="";
	for(var i=0;i<document.formsp.elements.length;i++){
		var e=document.formsp.elements[i];
		if ( e.checked==true && e.name=="gidi") { 
			if(j==0) {
				ppp= ""; 
			}else{
				ppp= ",";
			}
			document.formsp.sos.value += ppp+e.value; 
			if(document.formsp.elements[i+1].value==""){
				alert("没选择"+document.formsp.elements[i+2].value+"文章区");
				document.formsp.elements[i+1].focus();
				return false;
			}
			document.formsp.sosqu.value += ppp+document.formsp.elements[i+1].value; 
			document.formsp.sosming.value += ppp+document.formsp.elements[i+2].value; 
			j++;
		} 

	}
	if(j==0){
		alert("您没有选择任何技术圈");
		return false;
	}
	var url = "http://"+domains+".blog.51cto.com/joinbowen.php?tidd="+tid+"&uuid="+uid+"&jobd=1&sos="+document.formsp.sos.value+"&sosqu="+document.formsp.sosqu.value+"&domains="+domains+"&sosming="+document.formsp.sosming.value;
	postRequesttwo(url);
	return false;
}

function postRequest(strURL,homelogurl){
	var xmlHttp;
	if(!homelogurl)
	{
		homelogurl = '';	
	}
	
	if(window.XMLHttpRequest){ 
		var xmlHttp = new XMLHttpRequest();
	}else if(window.ActiveXObject){ 
		var xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
	}
	xmlHttp.open('GET', strURL, true);
	xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xmlHttp.onreadystatechange = function(){
		if (xmlHttp.readyState == 4 ){
			updatepage(xmlHttp.responseText,homelogurl);
		}
	}
	xmlHttp.send(null);
}
	
function updatepage(str,homeLogUrl){
	if(str == 'nologin')
	{
		alert("你还没有登录");
		window.location.href = homeLogUrl;		
	}
	else
	{	
		var scrollTop = window.pageYOffset 
                || document.documentElement.scrollTop 
                || document.body.scrollTop 
                || 0;


	document.getElementById("add_sys_type_div").style.top=(scrollTop+140)+"px";
		//document.getElementById("add_sys_type_div").style.top=(document.documentElement.scrollTop+140)+"px";
		document.getElementById("add_sys_type_div").innerHTML=str;		
		document.getElementById("add_sys_type_div").style.display="";
	}	
}
  
function postRequesttwo(strURL){
	var xmlHttp;
	if(window.XMLHttpRequest){ 
		var xmlHttp = new XMLHttpRequest();
	}
	else if(window.ActiveXObject){ 
		var xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
	}
	xmlHttp.open('GET', strURL, true);
	xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xmlHttp.onreadystatechange = function(){
		if (xmlHttp.readyState == 4 ){
			updatepagetwo(xmlHttp.responseText);
		}
	}
	xmlHttp.send(null);
}
	
function updatepagetwo(str){
	document.getElementById("add_sys_type_div").innerHTML=str;		
	document.getElementById("add_sys_type_div").style.display="";
}


function joingroups(tid,domains){
	var url = "http://"+domains+".blog.51cto.com/group/joingroups.php?tid="+tid+"&domains="+domains ;
	postRequestthree(url);
	return false;
}
function postRequestthree(strURL){
	var xmlHttp;
	if(window.XMLHttpRequest){ 
		var xmlHttp = new XMLHttpRequest();
	}
	else if(window.ActiveXObject){ 
		var xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
	}
	xmlHttp.open('GET', strURL, true);
	xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xmlHttp.onreadystatechange = function(){
		if (xmlHttp.readyState == 4 ){
			
			updatepagethree(xmlHttp.responseText);
			
		}
	}
	xmlHttp.send(null);
}
	
function updatepagethree(str){
		var scrollTop = window.pageYOffset 
                || document.documentElement.scrollTop 
                || document.body.scrollTop 
                || 0;


	document.getElementById("add_sys_type_divs").style.top=(scrollTop+140)+"px";
	//document.getElementById("add_sys_type_divs").style.top=(document.documentElement.scrollTop+140)+"px";
	document.getElementById("add_sys_type_divs").innerHTML=str;		
	document.getElementById("add_sys_type_divs").style.display="";
}
function thissubmit(){
	var bgObj;
	var title;
	var msgObj;
	bgObj=document.getElementById("bgDiv");
	msgObj=document.getElementById("msgDiv");
	title=document.getElementById("msgTitle");
	document.body.removeChild(bgObj);
	document.getElementById("msgDiv").removeChild(title);
	document.body.removeChild(msgObj);
}