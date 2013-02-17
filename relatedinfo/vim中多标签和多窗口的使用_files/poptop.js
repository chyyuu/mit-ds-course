function ChgTab(strDivIdBlock,strListIdBlock,strDivStyleBlock,strDivStyleNone,displayNumber,allNumber)
{
    for (var i = 1; i <= allNumber; i++) {
    	document.getElementById(strDivIdBlock + i).className = strDivStyleNone;
		document.getElementById(strListIdBlock + i).style.display = "none";
    }
	document.getElementById(strDivIdBlock + displayNumber).className = strDivStyleBlock;
	document.getElementById(strListIdBlock + displayNumber).style.display = "block";
}



//导航BEGIN

var TP_menu_popup;
var LSN_menu_popup;

function show(str,obj){
	popupNoClose();

	var w3c=(document.getElementById)? true:false;
	var ns6=(w3c && (navigator.appName=="Netscape"))? true: false;
	var left,top,ele;
	
	if (!ns6){
		var nLt = 0;
		var nTp = 0;
		var offsetParent = obj;
		while (offsetParent!=null && offsetParent!=document.body) {
			nLt += offsetParent.offsetLeft;
			nTp += offsetParent.offsetTop;
			offsetParent=offsetParent.offsetParent;
		}
		left = nLt;
		top = nTp + obj.offsetHeight + 3;
	} else {
		left = obj.offsetLeft;
		top = obj.offsetTop + obj.offsetHeight + 3;
	}
	

		document.getElementById("s3").style.display = "none";
	
	ele = document.getElementById("s"+str);
	TP_menu_popup = ele;
	
	ele.style.display = "block";
	ele.style.top = top + "px";
	ele.style.left = left + "px";
}
function popupClose(){
	LSN_menu_popup = window.setTimeout(function(){
		TP_menu_popup.style.display = "none";
	},100);
}
function popupNoClose(){
	if(LSN_menu_popup) window.clearTimeout(LSN_menu_popup);
}
//导航OVER