(function (document, window) {
    var _ssp_ad = window._ssp_ad = {
        reqUrl: 'http://sax.sina.com.cn/impress',
        t: +new Date(),
        timeout: 5000,
        pdpslist: {},
        count: 0,
        getCount: function () {
            var rnd = Math.round(Math.random());
            if (_ssp_ad.getCookie("rotatecount")) {
                _ssp_ad.count = parseInt(_ssp_ad.getCookie("rotatecount"), 10) + 1;
            } else {
                _ssp_ad.count = rnd;
            }
            _ssp_ad.setCookie("rotatecount", _ssp_ad.count);
        },
        setCookie: function (key, value, expires) {
            var l = new Date();
            var z = new Date(l.getTime() + expires * 60000);
            document.cookie = key + "=" + escape(value) + ";path=/;expires=" + z.toGMTString() + ";domain=sina.com.cn";
        },
        getCookie: function (key) {
            var c = document.cookie.split("; ");
            for (var i = 0; i < c.length; i++) {
                var d = c[i].split("=");
                if (d[0] == key) {
                    return unescape(d[1]);
                }
            }
            return '';
        },
        _rmScr: function (scr) {
            if (scr && scr.parentNode) {
                scr.parentNode.removeChild(scr);
            }
            scr = null;
        },
        _creScr: function (scr, url, charset) {
            scr.setAttribute('type', 'text/javascript');
            charset && scr.setAttribute('charset', charset);
            scr.async = true;
            scr.setAttribute('src', url);
            var head = document.getElementsByTagName('head')[0];
            head.insertBefore(scr, head.childNodes[0]);
        },
        jsonp: function (url, opt_cb, opt_conf) {
            var scr = document.createElement('SCRIPT'),
                scrLoaded = 0,
                conf = opt_conf || {},
                cb = opt_cb ||
            function () {}, charset = conf['charset'] || 'utf-8', timeout = conf['timeout'] || 0, timer;
            if (scr.readyState) {
                scr.onreadystatechange = ready();
            } else {
                scr.onload = ready();
            }

            function ready() {
                return function () {
                    if (scrLoaded) {
                        return;
                    }
                    var readyState = scr.readyState;
                    if ('undefined' == typeof readyState || readyState == "loaded" || readyState == "complete") {
                        scrLoadeded = 1;
                        try {
                            cb();
                            clearTimeout(timer);
                        } finally {
                            scr.onerror = scr.onload = scr.onreadystatechange = null;
                            _ssp_ad._rmScr(scr);
                        }
                    }
                }
            };

            scr.onerror = function () {
                scr.onerror = scr.onload = scr.onreadystatechange = null;
                _ssp_ad._rmScr(scr);
                conf.onfailure && conf.onfailure();
                clearTimeout(timer);
            };

            if (timeout) {
                timer = setTimeout(function () {
                    scr.onerror = scr.onload = scr.onreadystatechange = null;
                    _ssp_ad._rmScr(scr);
                    conf.onfailure && conf.onfailure();
                }, timeout);
            }
            _ssp_ad._creScr(scr, url, charset);

        },
        callback: function (data) {
            if (typeof data == "object" && typeof data.ad == "object") {
                _ssp_ad.pdpslist[data.ad[0].id] = data;
            }
        },
        id2PDPS: function (id) {
            if (typeof id == "string") {
                return "PDPS0000000" + id.split("_")[1];
            } else {
                return id;
            }
        },
        isDZ: function (pdps) {
            return {
                'PDPS000000000000': 1
            }[pdps];
        },
        cookieMapping: function (mapping) {
            var map, i = 0,
                img;
            if (mapping instanceof Array && mapping.length > 0) {
                img = new Image();
                img.width = 1;
                img.height = 1;
                document.body.appendChild(img);
                while(map = mapping[i++]) {
                    img.src = map;
                }
                document.body.removeChild(img);
            }
        },
        dspCM: function (map) {
            var img = new Image();
            img.width = 1;
            img.height = 1;
            document.body.insertBefore(img, document.body.childNodes[0]);
            img.src = map;
            img.onload = function () {
                document.body.removeChild(img);
            }
        },
        showAE: function (src, el, w, h) {
            var rad = [];
            var nad = [];
            //eval(src);
            src = decodeURIComponent(src);
            (new Function('rad', 'nad', src))(rad, nad);
            if (rad.length > 0) {
                var formatSrc = rad[0][0];
                var formatUrl = rad[0][1];
            } else {
                var formatSrc = nad[0][0];
                var formatUrl = nad[0][1];
            }
            //el.innerHTML = formatSrc;
            var filetype = formatSrc.substring(formatSrc.length - 3).toLowerCase();
            switch(filetype) {
            case "swf":
                var of = new sinaFlash(formatSrc, el.id + '_swf', w, h, "7", "", false, "High");
                of.addParam("wmode", "opaque");
                of.addParam("allowScriptAccess", "always");
                of.write(el.id);
                if (formatUrl != "") {
                    of.addVariable("_did", formatUrl);
                }
                (function () {
                    var url = formatUrl;
                    if (url) {
                        var ell = document.createElement('a'),
                            flashEl = document.getElementById(el.id);
                        flashEl.style.position = "relative";
                        ell.setAttribute("href", url);
                        ell.setAttribute("target", "_blank");
                        ell.style.cssText += ";display:block;width:" + w + "px;height:" + h + "px;position:absolute;left:0px;top:0px;filter:alpha(opacity:0)";
                        if (ell.style.filter) {
                            ell.style.backgroundColor = "white";
                        }
                        flashEl.appendChild(ell);
                    }
                })();
                break;
            case "jpg":
            case "gif":
            case "png":
                el.innerHTML = '<a href="' + formatUrl + '" target="_blank"><img src="' + formatSrc + '" border="0" width="' + w + '" height="' + h + '"/></a>';
                break;
            case "htm":
            case "tml":
                el.innerHTML = '<iframe id="ifm_' + el.id + '" frameborder="0" scrolling="no" width="' + w + '" height="' + h + '" src="' + formatSrc + '"></iframe>';
                break;
            case ".js":
                break;
            default:

            }
        },
        showAMP: function (src, el, w, h) {
            src += (src.indexOf('?') > 0 ? '&' : '?') + 'i_ssp=1';
            el.innerHTML = '<iframe id="ifm_' + el.id + '" frameborder="0" scrolling="no" width="' + w + '" height="' + h + '" src="' + src + '"></iframe>';
        },

        showDSP: function (src, el, w, h) {
            el.innerHTML = src;
        },
        showNetwork: function (networkId, posId, el, w, h) {
            var src = "";
            switch(networkId) {
            case "1": // taobao
                src = '<iframe id="network_' + networkId + posId + '" frameborder="0" scrolling="no" width="' + w + '" height="' + h + '" src="http://d3.sina.com.cn/litong/zhitou/union/taobao.html?w=' + w + '&h=' + h + '&pid=' + posId + '"></iframe>'
                break;
            default:
            }
            el.innerHTML = src;
        },
        failCB: function (pdps, el, adid, cb, w, h) {
            if (_ssp_ad.isDZ(pdps)) {
                var srcList = {
                    "950*90": "http://d1.sina.com.cn/litong/zhitou/gongyi/gongyi-banner.html",
                    "300*250": "http://d1.sina.com.cn/litong/zhitou/gongyi/gongyi-pip.html",
                    "250*230": "http://d1.sina.com.cn/litong/zhitou/gongyi/gongyi-square.html"
                };
                var src = srcList[w + "*" + h];
                el.innerHTML = '<iframe id="ifm_' + adid + '" frameborder="0" scrolling="no" width="' + w + '" height="' + h + '" src="' + src + '"></iframe>';
            } else {
                cb();
            }
        },
        domReady: function (d, f) {
            var ie = !! (window.attachEvent && !window.opera);
            var wk = /webkit\/(\d+)/i.test(navigator.userAgent) && (RegExp.$1 < 525);
            var fn = [];
            var run = function () {
                    for (var i = 0; i < fn.length; i++) fn[i]();
                };
            if (!ie && !wk && d.addEventListener) {
                return d.addEventListener('DOMContentLoaded', f, false);
            }
            if (fn.push(f) > 1) {
                return;
            }
            if (ie) {
                (function () {
                    try {
                        d.documentElement.doScroll('left');
                        run();
                    } catch(err) {
                        setTimeout(arguments.callee, 0);
                    }
                })();
            } else if (wk) {
                var t = setInterval(function () {
                    if (/^(loaded|complete)$/.test(d.readyState)) {
                        clearInterval(t), run();
                    }
                }, 0);
            }
        },
        loadIdentityIframe: function () {
            if (!document.getElementById("identityFrame")) {
                var frameUrl = "http://d1.sina.com.cn/litong/zhitou/identity.html";
                var ifr = document.createElement('iframe');
                ifr.width = 0;
                ifr.height = 0;
                ifr.frameBorder = 0;
                ifr.src = frameUrl;
                ifr.id = "identityFrame";
                this.domReady(document, function () {
                    _ssp_ad.dspCM("http://ads.ad.sina.com.cn/cm?sina_nid=4"); //mediaV
                    document.body.insertBefore(ifr, document.body.childNodes[0]);
                });
            }
        },
        load: function (adid, cb, w, h, rotateId) {
            var url = _ssp_ad.reqUrl,
                pdps = this.id2PDPS(adid);
            var el = document.getElementById(adid);
            if (pdps && el) {
                el.setAttribute('data-asp', 1);
                this.isDZ(pdps) && el.setAttribute('data-dz', 1);

                _ssp_ad.jsonp(
                url + (url.indexOf('?') >= 0 ? '&' : '?rotate_count=' + (isNaN(rotateId) ? 0 : (rotateId + 1)) + '&adunitid=' + this.id2PDPS(adid) + '&TIMESTAMP=' + _ssp_ad.t + '&referral=' + encodeURIComponent(document.referrer || window.location.href)), function () {
                    var _ssp_ads = _ssp_ad.pdpslist[pdps];
                    if (_ssp_ads && _ssp_ads.ad instanceof Array && _ssp_ads.ad.length > 0 && _ssp_ads.ad[0].value instanceof Array && _ssp_ads.ad[0].value.length > 0 && _ssp_ads.ad[0].value[0].content) {
                        w = _ssp_ads.ad[0].size.split("*")[0];
                        h = _ssp_ads.ad[0].size.split("*")[1];
                        var src = _ssp_ads.ad[0].value[0].content;
                        switch(_ssp_ads.ad[0].engineType) {
                        case "sina":
                            if (_ssp_ads.ad[0].value[0].manageType === 'AE') {
                                _ssp_ad.showAE(src, el, w, h);
                            } else {
                                _ssp_ad.showAMP(src, el, w, h);
                            }
                            break;
                        case "dsp":
                            _ssp_ad.showDSP(src, el, w, h);
                            break;
                        case "network":
                            var networkId = _ssp_ads.ad[0].value[0].manageType;
                            _ssp_ad.showNetwork(networkId, src, el, w, h);
                            break;
                        default:
                            cb();
                        }
                        _ssp_ad.cookieMapping(_ssp_ads.mapUrl);
                        _ssp_ads = window["_ssp_ads"] = null;
                    } else {
                        _ssp_ad.failCB(pdps, el, adid, cb, w, h);
                    }
                    _ssp_ad.pdpslist[pdps] = false;
                }, {
                    timeout: _ssp_ad.timeout,
                    onfailure: function () {
                        _ssp_ad.failCB(pdps, el, adid, cb, w, h);
                    }
                });
            } else {
                cb();
            }
        },
        init: function (cb) {
            this.domReady(document, function () {
                var allDom = document.getElementsByTagName("*");
                for (var i = 0, il = allDom.length; i < il; i++) {
                    var pdps = allDom[i].getAttribute("pdps");
                    if (pdps) {
                        if (_ssp_ad.pdpslist[pdps] == undefined) {
                            _ssp_ad.pdpslist[pdps] = true;
                        }
                    }
                }
                for (var j in _ssp_ad.pdpslist) {
                    if (_ssp_ad.pdpslist[j]) {
                        var divid = "ad_" + j.substring(11, 16);
                        _ssp_ad.load(divid, cb);
                    }
                }
            });
        }
    };
})(document, window);
_ssp_ad.loadIdentityIframe();
_ssp_ad.getCount();