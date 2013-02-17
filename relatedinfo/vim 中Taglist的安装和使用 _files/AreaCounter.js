// JScript 文件
function LogClickCount(link,areaId)
{
    var url=window.location.href;
    var myImage=new Image(0,0);
    if(link!=null && link!='')
    {
        myImage.src="http://counter.csdn.net/a/Counter.aspx?area=" + areaId+"&u="+ escape(link);
    }
    else
    {
        myImage.src="http://counter.csdn.net/a/Counter.aspx?area=" + areaId+"&u="+ escape(link);
    }
}