/*jshint browser:true */
/**
	@file 前端用户行为跟踪
	@example

<script type="text/javascript" charset="utf-8" src="http://csdnimg.cn/pubfooter/js/tracking.js" defer=""></script>


	@author 曹宇 <caoyu@csdn.net>
	@version 20130107
 */
;(function(window, undefined) {
	'use strict';
	var doc = document,
		nav = navigator,
		loc = location,
		/**
			转换数组
			@param {Array} array - 数组，或数组类似的对象
			@param {Function(ele, i)} callback - 转换用的回调函数，返回值将作为转换的结果
			@returns {Array} 转换结果
		 */
		map = function(array, callback) {
			var i, len = array.length,
				ret = new Array(len);
			for(i = 0; i < array.length; i++) {
				ret[i] = callback(array[i], i);
			}
			return ret;
		},
		/**
			在map中使用，把HTML元素转换为其所指示的tag，会清理头尾的空白字符
			@param {HTMLElement} ele - HTML元素
			@param {Number} i - 当前元素在数组中的索引
			@returns {String} 当前元素的innerHTML
		 */
		elementTag = function(ele) {
			return ele.innerHTML.replace(/^\s+|\s+$/g, '');
		},
		/**
			编码执行的参数对象，用于附加到Url结尾的查询字符串
			@param {Object} params - 参数对象
			@returns {String} 可直接附加到Url结尾的查询字符串，不包含?符号
		 */
		encode = function(params) {
			var i, ret = [];
			for(i in params) {
				if(params.hasOwnProperty(i)) {
					ret.push(i + '=' + encodeURIComponent(params[i]));
				}
			}
			return ret.join('&');
		},
		/**
			发送跨域HTTP GET请求
			@param {String} url - 请求的Url
			@param {Object} params - 请求要附加的参数对象
		 */
		crossdomainGet = function(url, params) {
			var i = new Image();
			i.onload = i.onerror = function() {
				i.onload = i.onerror = null;
				i.removeAttribute('src');
				i = null;
			};
			i.src = url + (url.indexOf('?') >= 0 ? '&' : '?') + encode(params);
		},
		/**
			检测当前网站登录的用户名
			@returns {String} 当前登录的用户名，如果未登录或检测不到，返回空白字符串
		 */
		detectUsername = function() {
			var m = /(; )?UserName=([^;]+)/.exec(doc.cookie);
			return m && m[2] || '';
		},
		/**
			使用CSS选择器检索对应的DOM元素
			@param {String} selector - CSS选择器
			@returns {Array[HTMLElement]} HTML元素集合，如果浏览器不支持使用CSS选择器查找将返回 undefined，如果找不到任何元素返回0长度的近似数组
		 */
		querySelectorAll = function(selector) {
			if(doc.querySelectorAll) {
				return doc.querySelectorAll(selector);
			} else if(window.jQuery) {
				return window.jQuery(selector).get();
			}
		},
		/**
			检测当前页面包含的Tags
			@param {Array} site - 当前站点需要检测的Url和检测策略数组, 每3个一组，分别是[要匹配的Url片段，检索页面中Tags使用的选择器，选择器无法执行时替代的检测函数]
			@returns {String} 逗号分隔的Tags，如果找不到将返回 undefined
		 */
		detectTags = function(site) {
			var i, len, tags, path, selector, altfunc;
			if(site) {
				len = site.length;
				for(i = 0; i < len; i += 3) {
					path = site[i];
					if(loc.pathname.indexOf(path) < 0) {
						continue;
					}
					selector = site[i + 1];
					tags = querySelectorAll(selector);
					if(tags === undefined) {
						altfunc = site[i + 2];
						if(altfunc) {
							tags = altfunc();
						}
					}
					tags = tags && tags[0] ? map(tags, elementTag).join(',') : undefined;
				}
			}
			return tags;
		},
		/**
			不同站点检测Tags的策略配置
		 */
		siteCfgs = {
			// TODO 替代检测函数
			'bbs.csdn.net': ['/topics/', 'div.tag span'],
			'blog.csdn.net': ['/article/details/', 'div.tag2box a'],
			'ask.csdn.net': ['/questions/', 'div.tag_data a.tag span'],
			'download.csdn.net': ['/detail/', 'div.info a[href^="/tag/"]'],
			'www.csdn.net': ['/article/', 'div.tag a'],
			'www.csto.com': ['/p/', 'span.tech a']
		};

	crossdomainGet('http://status.api.csdn.net/track', {
		user_agent: nav.userAgent,
		url: loc.href,
		referrer: doc.referrer,
		user_name: detectUsername() || '',
		tag: detectTags(siteCfgs[loc.host]) || '',
		'x-acl-token': 'status_js_dkuyqthzbajmncbsb_token'
	});
}(window));