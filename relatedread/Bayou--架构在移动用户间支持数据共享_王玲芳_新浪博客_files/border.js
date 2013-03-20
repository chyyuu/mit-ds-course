function startmarquee(lh, speed, delay, index) {
	var t;
	var p = false;
	var o = $(index);
	o.innerHTML += o.innerHTML;
	o.onmouseover = function() {
		p = true
	}
	o.onmouseout = function() {
		p = false
	}
	if (o.scrollHeight != 0) {
		lh = Math.floor(o.scrollHeight / 4);
	}
	//alert(lh)
	o.scrollTop = 0;
	function start() {
		if (!p) o.scrollTop +=1;
		t = setInterval(scrolling, speed);		
	}
	function scrolling()
	 {
		if (o.scrollHeight != 0) {
			lh = Math.floor(o.scrollHeight / 4);
		};		
		if (o.scrollTop % lh != 0)
		 {
			o.scrollTop += 1;
			if (o.scrollTop >= lh * 2)
			 o.scrollTop = 0;
		}
		 else
		 {
			clearInterval(t);
			setTimeout(start, delay);
		}
	}	
	setTimeout(start, delay);
};
function $(id) {
	return document.getElementById(id);
};
function loadScript(b, d, c)
{
		var a = document.createElement("script");
		a.type = "text/javascript";
		if (c) {
			a.charset = c;
		}
		if (d) {
			a.onload = a.onreadystatechange =a.onerror=  function () {
				if (a.readyState && a.readyState != "loaded" && a.readyState != "complete") {
					return;
				}
				a.onreadystatechange = a.onload =a.onerror=  null;
				d();
			};
		};		
		a.src = b;
		document.getElementsByTagName("head")[0].appendChild(a);	
};
function byteLength(str) {
	if (typeof str == "undefined") {
		return 0;
	}
	var aMatch = str.match(/[^\x00-\x80]/g);
	return Math.ceil((str.length + (!aMatch ? 0: aMatch.length)) / 2);
}; 

function shareBlog(u,t,i,c,p,s){	 
	  var f ='http://v.t.sina.com.cn/share/share.php?appkey=2034545597';
	  e=encodeURIComponent;
	  p=['&url=',e(u),'&title=',e(t||document.title),'&ralateUid==',e(i),'&content=',c||'gb2312','&pic=',e(p||''),'&searchPic=',s||false].join('');	 
	  var openWnd = function(){           
		 if(!window.open([f,p].join(''),'mb',['toolbar=0,status=0,resizable=1,width=440,height=430,left=',(screen.width-440)/2,',top=',(screen.height-430)/2].join('')))document.location.href=[f,p].join('');
    }
    if (/Firefox/.test(navigator.userAgent)){
        setTimeout(openWnd, 30);
    } else{
        openWnd();
    };	
};
function sharemicBlog(){
	    clickPFP("http://sina.allyes.com/main/adfclick?db=sina&bid=247819,404562,409875&cid=0,0,0&sid=407237&advid=358&camid=46514&show=ignore")
		shareBlog(data.addata.url,data.addata.title,data.addata.ralateUid,"utf-8",data.addata.pic,data.addata.searchPic);	
	};
function clickPFP(url) {
			var img = document.createElement("img");
			img.src = url;
			img.onload = img.onreadystatechange = img.onerror = function () {
				img = null;
			};
		};
function main(){
	if(typeof $("wbblog")=="undefined" || $("wbblog")==null){
		return;
		}
	if(typeof data == "undefined"){
		return;
	};	
	if(data.isUseWb ==false){
		$("wbblog").innerHTML ="<p>·<a href="+data.txturl+" target='_blank'>"+data.txtcontent+"</a></p>";
		$("wbblog").style.display = "block";		
		return;
	};
	var adbordata = {};
	adbordata=data.bloger;
	data.addata.ralateUid=data.bloger.uid;
	adbordata.pdpsid="PDPS000000011727";	
	$("wbblog").style.display = "block";
		var vhtml = "";
		if (adbordata.isV) {
			vhtml = '<img width="15px" height="11px" style="background:url(http://img.t.sinajs.cn/t4/style/images/common/ico_user.png?id=1319610952671) no-repeat 0 -'+(125+12.5*adbordata.vType)+'px scroll;" src="http://img.t.sinajs.cn/t4/style/images/common/transparent.gif"/>'
		};	
		var top="4";
		try{
		if (navigator.appName.indexOf("Microsoft") != -1) {
		top="6";
	    }else {
		top="4";
	    }
		}catch(e){
			top="4";
			}
		if (byteLength(adbordata.name) < 6) {
			$("wbblog").innerHTML = '<span style="line-height:20px;float:left;padding-top:'+top+'px;">·欢迎关注  <a href="http://sina.allyes.com/main/adfclick?db=sina&bid=247819,409445,414758&cid=0,0,0&sid=412238&advid=358&camid=46514&show=ignore&url=http://www.weibo.com/' + adbordata.id + '" target="_blank"><strong>' + adbordata.name + '</strong></a>' + vhtml + '</span><span style="margth-top:0px;float:left;"><iframe src="http://d3.sina.com.cn/litong/kuaijieweibo/yafeng/boke/iframe/followBtn.html?uid=' + adbordata.id + '&activeid=' + adbordata.pdpsid + '&clickurl='+encodeURIComponent("http://sina.allyes.com/main/adfclick?db=sina&bid=247819,404561,409874&cid=0,0,0&sid=407236&advid=358&camid=46514&show=ignore")+'&sclick='+encodeURIComponent("http://sina.allyes.com/main/adfclick?db=sina&bid=247819,409446,414759&cid=0,0,0&sid=412239&advid=358&camid=46514&show=ignore")+'"  frameborder="0"  scrolling="no" width="62px" height="22px"></iframe></span><span  style="line-height:20px;float:left;padding-top:3px;padding-left:4px"><a href="javascript:void(0)" onClick="javascript:sharemicBlog();">分享</a></span>';			
		} else {
			$("wbblog").innerHTML = '<span style="line-height:20px;float:left;padding-top:'+top+'px;">·<a href="http://sina.allyes.com/main/adfclick?db=sina&bid=247819,409445,414758&cid=0,0,0&sid=412238&advid=358&camid=46514&show=ignore&url=http://www.weibo.com/' + adbordata.id + '" target="_blank"><strong>' + adbordata.name + '</strong></a>' + vhtml + '</span><span style="margth-top:0px;float:left;"><iframe src="http://d3.sina.com.cn/litong/kuaijieweibo/yafeng/boke/iframe/followBtn.html?uid=' + adbordata.id + '&activeid=' + adbordata.pdpsid + '&clickurl='+encodeURIComponent("http://sina.allyes.com/main/adfclick?db=sina&bid=247819,404561,409874&cid=0,0,0&sid=407236&advid=358&camid=46514&show=ignore")+'&sclick='+encodeURIComponent("http://sina.allyes.com/main/adfclick?db=sina&bid=247819,409446,414759&cid=0,0,0&sid=412239&advid=358&camid=46514&show=ignore")+'"  frameborder="0"  scrolling="no" width="62px" height="22px"></iframe></span><span style="line-height:20px;float:left;padding-top:3px;padding-left:4px"><a href="javascript:void(0)" onClick="javascript:sharemicBlog();" >分享</a></span>';
		};	
}		
(function() {    
	startmarquee(20, 25, 5000, "scroll");
	var url="http://ba.sass.sina.com.cn/front/deliver?psId=PDPS000000011727";
	if(window.location.href.match(/test/)!= undefined){		
		url="http://d1.sina.com.cn/litong/kuaijieweibo/yafeng/boke/js/addata.js";
	};
	//var url="http://d1.sina.com.cn/litong/kuaijieweibo/yafeng/boke/js/addata.js";
	loadScript(url,main,"gb2312")
})();