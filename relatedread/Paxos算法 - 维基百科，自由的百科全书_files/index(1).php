// from http://en.wikipedia.org/w/index.php?title=MediaWiki:Gadget-edittop.js&oldid=488945990, modified

// **********************************************************************
// **                 ***WARNING GLOBAL GADGET FILE***                 **
// **             changes to this file affect many users.              **
// **           please discuss on the talk page before editing         **
// **                                                                  **
// **********************************************************************
// Imported from [[User:Alex Smotrov/edittop.js]], version as of: 2007-06-19T04:28:52
// Updated from [[User:TheDJ/Gadget-edittop.js]], version as of: 2009-04-28T11:54:22

if ($.inArray( mw.config.get('wgAction'), [ 'view', 'purge' ]) !== -1 && mw.config.get( 'wgNamespaceNumber' ) >=0) {
  $(function edittop_hook () {
    var localtitles = {
      cs: 'Editovat úvodní sekci',
      en: 'Edit lead section',
      fr: 'Modifier le résumé introductif',
      it: 'Modifica della sezione iniziale',
      ja: '導入部を編集',
      ko: '도입부를 편집',
      pt: 'Editar a seção superior',
      'pt-br': 'Editar a seção superior'
    };

    localtitles['zh'] = localtitles['zh-hans'] = localtitles['zh-cn'] = localtitles['zh-sg'] = localtitles['zh-my'] = '编辑首段';
    localtitles['zh-hant'] = localtitles['zh-hk'] = localtitles['zh-mo'] = localtitles['zh-tw'] = '編輯首段';
    var our_content = document.getElementById ("content") || document.getElementById ("mw_content") || document.body;
    var editspans = getElementsByClassName (our_content, "span", "editsection");
    var span1;

    for (var i = 0; editspans && i < editspans.length; i++) {
      if (editspans[i].className.indexOf ("plainlinks") == -1) {
        span1 = editspans[i];
        break;
      }
    }
    if (!span1) {
      return;
    }
    var span0 = span1.cloneNode (true);
    var editwidth = span1.offsetWidth;
    if (mw.config.get("skin") == "monobook") {
      mw.util.addCSS ("h1.firstHeading span.editsection {float: right;}");
    }
    if (mw.config.get("skin") == "modern") {
      mw.util.addCSS ("h1#firstHeading span.editsection {float: right;}");
    }
    if (mw.config.get("skin") == "vector") {
      mw.util.addCSS ("h1.firstHeading span.editsection {font-size: 50%;}");
    } else {
      editwidth += 10;
    }
    var topicons = getElementsByClassName (our_content, "div", "topicon");
    for (var el = 0; topicons && el < topicons.length; el++) {
      topicons[el].style.marginRight = editwidth + "px";
    }
    var mwfrtag = document.getElementById ("mw-fr-revisiontag");
    if (mwfrtag) {
      mwfrtag.style.marginRight = editwidth + "px";
    }

    our_content = document.getElementById ("mw_header") || document.getElementById ("content") || document.body;
    var parent = our_content.getElementsByTagName ("H1")[0];
    parent.insertBefore (span0, parent.firstChild);
    var a = span0.getElementsByTagName ("A")[0];
    if (a.href.indexOf ("&section=T") == -1) {
      a.title = a.title.replace (/(: |：).*$/, "$1" + "0");
      a.setAttribute ("href", a.getAttribute ("href", 2).replace (/&section=\d+/, "&section=0"));
    }
    else { //transcluded
      a.title = localtitles[mw.config.get( 'wgUserLanguage' )] || localtitles.en;
      a.setAttribute ("href", mw.util.wikiGetlink( mw.config.get( 'wgPageName' ) ) + "?action=edit&section=0");
    }
  });
}