WebFontConfig = {
  google: { 
    families: [ 
      'Gentium+Basic:400,400italic,700,700italic:latin',
      'Ubuntu+Mono:400,400italic,700,700italic',
      'Yanone+Kaffeesatz:latin',
    ] 
  }
};

(function() {
  var wf = document.createElement('script');
  wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
      '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
  wf.type = 'text/javascript';
  wf.async = 'true';
  var s = document.getElementsByTagName('script')[0];
  s.parentNode.insertBefore(wf, s);
})();
