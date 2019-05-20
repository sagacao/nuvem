function micro_request() {
    let forms= new FormData();
    forms.append('uname','test');
    forms.append('psd',123456);
    forms.append('age','22'); 

    let xhr = new XMLHttpRequest();    
   xhr.onreadystatechange = function(){
   if(xhr.readyState==4){
       if(xhr.status>=200&&xhr.status<=300||xhr.status==304){
             console.log(xhr.response)
       }
       }else{
           console.log(xhr.status)
       }
   }
   
   xhr.open('POST','https://qgame-test.dayukeji.com:12004/micro/bytedance/login/verify',true);
   xhr.setRequestHeader("Content-Type","application/x-www-form-urlencoded");  //formdata数据请求头需设置为application/x-www-form-urlencoded
console.log(forms)
xhr.send(forms)
}
