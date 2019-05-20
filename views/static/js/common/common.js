if (!String.prototype.format) {
    String.prototype.format = function () {
        var args = arguments;
        return this.replace(/{(\d+)}/g, function (match, number) {
            return typeof args[number] != 'undefined'
              ? args[number]
              : match
            ;
        });
    };
}

//将form中的值转换为键值对。
function getFormJson(frm) {
	
    var o = {};
    
    var a = $(frm).serializeArray();
    
    $.each(a, function () {
    	
        if (o[this.name] !== undefined) {
        	
            if (!o[this.name].push) {
            	
                o[this.name] = [o[this.name]];
                
            };
            
            o[this.name].push(this.value || '');
            
        } else {
        	
            o[this.name] = this.value || '';
            
        };
    });

    return o;
};

function addDays(date,n){
    return new Date(date.getTime()+n*24*60*60*1000); 
};

Date.prototype.format = function (fmt) {            //时间格式化
    var o = {
        "M+": this.getMonth() + 1,                 //月份 
        "d+": this.getDate(),                    //日 
        "h+": this.getHours(),                   //小时 
        "m+": this.getMinutes(),                 //分 
        "s+": this.getSeconds(),                 //秒 
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度 
        "S": this.getMilliseconds()             //毫秒 
    };
    if (/(y+)/.test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    }
    for (var k in o) {
        if (new RegExp("(" + k + ")").test(fmt)) {
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
        }
    }
    return fmt;
}

function Trim(str,is_global){
    var result;
    result = str.replace(/(^\s+)|(\s+$)/g,"");
    if(is_global.toLowerCase()=="g")
    {
        result = result.replace(/\s/g,"");
     }
    return result;
}

Array.prototype.indexOf = function( searchvalue ) {//给数组添加indexOf函数
	if(this.length == 0) return -1;
	for (var i=0;i<this.length;i++) {
		if(this[i] == searchvalue) return i;
	}
	return -1;
}

/* Array.prototype.containd = function( children ) {//判断一个数组里是否包含另一个数组
	if( !(children instanceof Array) ) return false;
	
	for(var i = 0, len = children.length; i < len; i++){

       if(this.indexOf(children[i]) == -1) return false;
    }
    return true;
	
}

Array.prototype.(arr){

} = function() {//数组去重
	var res = [];
	var json = {};
	for(var i = 0; i < this.length; i++) {
		if(!json[this[i]]) {
			res.push(this[i]);
			json[this[i]] = 1;
		}
	}
	return res;
}; */

function setCookie(name,value,expiredays){//设置cookie
	var exdate = new Date();
	exdate.setDate(exdate.getDate()+expiredays);
	document.cookie = name+ "=" +escape(value)+ ((expiredays==null) ? "" : ";expires="+exdate.toGMTString())
}

function getCookie(name){//获取cookie
	if (document.cookie.length > 0){
		var start = document.cookie.indexOf( name + "=" )
		if ( start != -1){
		    start = start + name.length + 1; 
		    end = document.cookie.indexOf(";", start);
		    if ( end== -1 ) end = document.cookie.length;
		    return unescape(document.cookie.substring(start , end))
	    } 
    }
	return "";
}

//数组去重
function unique(arr){
    return Array.from(new Set(arr))
}