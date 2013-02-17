 function dialog(){
	var titile = '';
	var width = 300;
	var height = 180;
	var src = "";
	var path = "http://blog.51cto.com/images/dialog/";
	var sFunc = '<input id="dialogOk" type="button" style="width:62px;height:22px;border:0;background:url(\'http://www.51cto.com/php/image/smb_btn_bg.gif\');line-height:20px;" value="确认" onclick="new dialog().reset();" /> <input id="dialogCancel" type="button" style="width:62px;height:22px;border:0;background:url(\'http://www.51cto.com/php/image/smb_btn_bg.gif\');line-height:20px;" value="取消" onclick="new dialog().reset();" />';
	var sClose = '<input type="image" id="dialogBoxClose" onclick="new dialog().reset();" src="' + path + 'dialogClose0.gif" border="0" width="17" height="17" onmouseover="this.src=\'' + path + 'dialogCloseF.gif\';" onmouseout="this.src=\'' + path + 'dialogClose0.gif\';" align="absmiddle" />';
	var sBody = '\
		<table id="dialogBodyBox" border="0" align="center" cellpadding="0" cellspacing="0">\
			<tr height="10"><td colspan="4"></td></tr>\
			<tr>\
				<td width="10"></td>\
				<td width="80" align="center" valign="absmiddle"><img id="dialogBoxFace" src="' + path + '3.gif" /></td>\
				<td id="dialogMsg" style="font-size:12px;color:#000;"></td>\
				<td width="10"></td>\
			</tr>\
			<tr height="10"><td colspan="4" align="center"></td></tr>\
			<tr><td id="dialogFunc" colspan="4" align="center">' + sFunc + '</td></tr>\
			<tr height="10"><td colspan="4" align="center"></td></tr>\
		</table>\
	';
	var sBox = '\
		<table id="dialogBox" width="' + width + '" border="0" cellpadding="0" cellspacing="0" style="border:1px solid #000;display:none;z-index:10;">\
			<tr height="1" bgcolor="#D6E3EB"><td></td></tr>\
			<tr height="25" bgcolor="#6795B4">\
				<td>\
					<table onselectstart="return false;" style="-moz-user-select:none;" width="100%" border="0" cellpadding="0" cellspacing="0">\
						<tr>\
							<td width="6"></td>\
							<td id="dialogBoxTitle" onmousedown="new dialog().moveStart(event, \'dialogBox\')" style="color:#fff;cursor:move;font-size:12px;font-weight:bold;">系统提示信息</td>\
							<td id="dialogClose" width="27" align="right" valign="middle">\
								' + sClose + '\
							</td>\
							<td width="6"></td>\
						</tr>\
					</table>\
				</td>\
			</tr>\
			<tr height="2" bgcolor="#EDEDED"><td></td></tr>\
			<tr id="dialogHeight" style="height:' + height + '">\
				<td id="dialogBody" style="background:#fff;color:#000;">' + sBody + '</td>\
			</tr>\
		</table>\
		<div id="dialogBoxShadow" style="display:none;z-index:9;"></div>\
	';
	function $(_sId){return document.getElementById(_sId)}
	this.show = function(){$('dialogBodyBox') ? function(){} : this.init();this.middle('dialogBox');this.shadow();}
	this.reset = function(){this.hideModule('select', '');$('dialogBox').style.display='none';$('dialogBoxShadow').style.display = "none";$('dialogBody').innerHTML = sBody;}
	this.html = function(_sHtml){$("dialogBody").innerHTML = _sHtml;this.show();}
	this.init = function(){
		$('dialogCase') ? $('dialogCase').parentNode.removeChild($('dialogCase')) : function(){};
		var oDiv = document.createElement('span');
		oDiv.id = "dialogCase";
		oDiv.innerHTML = sBox;
		document.body.appendChild(oDiv);
	}
	this.button = function(_sId, _sFuc){
		if($(_sId)){
			$(_sId).style.display = '';
			if($(_sId).addEventListener){
				if($(_sId).act){$(_sId).removeEventListener('click', function(){eval($(_sId).act)}, false);}
				$(_sId).act = _sFuc;
				$(_sId).addEventListener('click', function(){eval(_sFuc)}, false);
			}else{
				if($(_sId).act){$(_sId).detachEvent('onclick', function(){eval($(_sId).act)});}
				$(_sId).act = _sFuc;
				$(_sId).attachEvent('onclick', function(){eval(_sFuc)});
			}
		}
	}
	this.shadow = function(){
		var oShadow = $('dialogBoxShadow');
		var oDialog = $('dialogBox');
		oShadow['style']['position'] = "absolute";
		oShadow['style']['background']	= "#000";
		oShadow['style']['display']	= "";
		oShadow['style']['opacity']	= "0.2";
		oShadow['style']['filter'] = "alpha(opacity=20)";
		oShadow['style']['top'] = oDialog.offsetTop + 6;
			
		oShadow['style']['left'] = oDialog.offsetLeft + 6;
		oShadow['style']['width'] = oDialog.offsetWidth;
		oShadow['style']['height'] = oDialog.offsetHeight;
		
	}
	this.open = function(_sUrl, _sMode){
		this.show();
		if(!_sMode || _sMode == "no" || _sMode == "yes"){
			$("dialogBody").innerHTML = "<iframe id='dialogFrame' width='100%' height='100%' frameborder='0' scrolling='" + _sMode + "'></iframe>";
			$("dialogFrame").src = _sUrl;
		}
	}
	this.showWindow = function(_sUrl, _iWidth, _iHeight, _sMode){
		var oWindow;
		var sLeft = (screen.width) ? (screen.width - _iWidth)/2 : 0;
		var sTop = (screen.height) ? (screen.height - _iHeight)/2 : 0;
		if(window.showModalDialog && _sMode == "m"){
			oWindow = window.showModalDialog(_sUrl,"","dialogWidth:" + _iWidth + "px;dialogheight:" + _iHeight + "px");
		} else {
			oWindow = window.open(_sUrl, '', 'height=' + _iHeight + ', width=' + _iWidth + ', top=' + sTop + ', left=' + sLeft + ', toolbar=no, menubar=no, scrollbars=' + _sMode + ', resizable=no,location=no, status=no');
		}
	}
	this.event = function(_sMsg, _sOk, _sCancel, _sClose){
		$('dialogFunc').innerHTML = sFunc;
		$('dialogClose').innerHTML = sClose;
		$('dialogBodyBox') == null ? $('dialogBody').innerHTML = sBody : function(){};
		$('dialogMsg') ? $('dialogMsg').innerHTML = _sMsg  : function(){};
		this.show();
		_sOk ? this.button('dialogOk', _sOk) | $('dialogOk').focus() : $('dialogOk').style.display = 'none';
		_sCancel ? this.button('dialogCancel', _sCancel) : $('dialogCancel').style.display = 'none';
		_sClose ? this.button('dialogBoxClose', _sClose) : function(){};
		//_sOk ? this.button('dialogOk', _sOk) : _sOk == "" ? function(){} : $('dialogOk').style.display = 'none';
		//_sCancel ? this.button('dialogCancel', _sCancel) : _sCancel == "" ? function(){} : $('dialogCancel').style.display = 'none';
	}
	this.set = function(_oAttr, _sVal){
		var oShadow = $('dialogBoxShadow');
		var oDialog = $('dialogBox');
		var oHeight = $('dialogHeight');

		if(_sVal != ''){
			switch(_oAttr){
				case 'title':
					$('dialogBoxTitle').innerHTML = _sVal;
					title = _sVal;
					break;
				case 'width':
					oDialog['style']['width'] = _sVal;
					width = _sVal;
					break;
				case 'height':
					oHeight['style']['height'] = _sVal;
					height = _sVal;
					break;
				case 'src':
					if(parseInt(_sVal) > 0){
						$('dialogBoxFace') ? $('dialogBoxFace').src = path + _sVal + '.gif' : function(){};
					}else{
						$('dialogBoxFace') ? $('dialogBoxFace').src = _sVal : function(){};
					}
					src = _sVal;
					break;
			}
		}
		this.middle('dialogBox');
		oShadow['style']['top'] = oDialog.offsetTop + 6;
		oShadow['style']['left'] = oDialog.offsetLeft + 6;
		oShadow['style']['width'] = oDialog.offsetWidth;
		oShadow['style']['height'] = oDialog.offsetHeight;
	}
	this.moveStart = function (event, _sId){
		var oObj = $(_sId);
		oObj.onmousemove = mousemove;
		oObj.onmouseup = mouseup;
		oObj.setCapture ? oObj.setCapture() : function(){};
		oEvent = window.event ? window.event : event;
		var dragData = {x : oEvent.clientX, y : oEvent.clientY};
		var backData = {x : parseInt(oObj.style.top), y : parseInt(oObj.style.left)};
		function mousemove(){
			var oEvent = window.event ? window.event : event;
			var iLeft = oEvent.clientX - dragData["x"] + parseInt(oObj.style.left);
			var iTop = oEvent.clientY - dragData["y"] + parseInt(oObj.style.top);
			oObj.style.left = iLeft;
			oObj.style.top = iTop;
			$('dialogBoxShadow').style.left = iLeft + 6;
			$('dialogBoxShadow').style.top = iTop + 6;
			dragData = {x: oEvent.clientX, y: oEvent.clientY};
			
		}
		function mouseup(){
			var oEvent = window.event ? window.event : event;
			oObj.onmousemove = null;
			oObj.onmouseup = null;
			if(oEvent.clientX < 1 || oEvent.clientY < 1 || oEvent.clientX > document.body.clientWidth || oEvent.clientY > document.body.clientHeight){
				oObj.style.left = backData.y;
				oObj.style.top = backData.x;
				$('dialogBoxShadow').style.left = backData.y + 6;
				$('dialogBoxShadow').style.top = backData.x + 6;
			}
			oObj.releaseCapture ? oObj.releaseCapture() : function(){};
		}
	}
	this.hideModule = function(_sType, _sDisplay){
		var aIframe = parent.document.getElementsByTagName("iframe");aIframe=0;
		var aType = document.getElementsByTagName(_sType);
		var iChildObj, iChildLen;
		for (var i = 0; i < aType.length; i++){
			aType[i].style.display	= _sDisplay;
		}
		for (var j = 0; j < aIframe.length; j++){
			iChildObj = document.frames ? document.frames[j] : aIframe[j].contentWindow;
			iChildLen = iChildObj.document.body.getElementsByTagName(_sType).length;
			for (var k = 0; k < iChildLen; k++){
				iChildObj.document.body.getElementsByTagName(_sType)[k].style.display = _sDisplay;
			}
		}
	}
	this.middle = function(_sId){
		document.getElementById(_sId)['style']['display'] = '';
		document.getElementById(_sId)['style']['position'] = "absolute";
		document.getElementById(_sId)['style']['left'] = (document.body.clientWidth / 2) - (document.getElementById(_sId).offsetWidth / 2);
		//var range=getRange();		
		document.getElementById(_sId)['style']['top'] = (300+ document.documentElement.scrollTop) - (document.getElementById(_sId).offsetHeight / 2);
		
	}
}
