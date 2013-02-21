var button_dimension = {
   1 : {'w' : 60, 'h': 55 },
   2 : {'w' : 60, 'h': 69 },
   3 : {'w' : 80, 'h': 15 },
   4 : {'w' : 91, 'h': 17 },
   5 : {'w' : 95, 'h': 22 },
   6 : {'w' :115, 'h': 22 },
   7 : {'w' : 60, 'h': 66 },
   8 : {'w' :100, 'h': 34 },
   9 : {'w' :106, 'h': 55 },
  10 : {'w' : 60, 'h': 55 },
  11 : {'w' :100, 'h': 20 },
  12 : {'w' : 80, 'h': 15 },
  13 : {'w' : 80, 'h': 15 }
};

function funp_genButton(url,style) {
   if (!style) style = 1;
   if (!url) url='';
   else url="url="+encodeURIComponent(url)+'&';
   document.write('<iframe src="http://funp.com/tools/buttoniframe.php?'+url+'s='+style+
                  '" height="'+button_dimension[style]['h']
                  +'" width="'+button_dimension[style]['w']
                  +'" scrolling="no" frameborder="0" marginheight="0" marginwidth="0"></iframe>');
}
