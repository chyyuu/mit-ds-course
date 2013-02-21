 

$(document).ready(function() {
// setting the tabs in the sidebar hide and show, setting the current tab
	$('div.tabbed div').hide();
	$('div.t1').show();
	$('div.tabbed ul.tabs li.t1 a').addClass('tab-current');

// SIDEBAR TABS
$('div.tabbed ul li a').click(function(){
	var thisClass = this.className.slice(0,2);
	$('div.tabbed div').hide();
	$('div.' + thisClass).show();
	$('div.tabbed ul.tabs li a').removeClass('tab-current');
	$(this).addClass('tab-current');
	});
});