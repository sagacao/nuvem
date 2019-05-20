$(document).ready(function () {

    var loginViewModel = function () {
        this.userName       = ko.observable("");
        this.userNameGet    = ko.computed(function () {
            if (getCookie("userName") != "") {
                return this.userName(getCookie("userName"));
            }
        }, this);
        this.password       = ko.observable("");
		/*this.passwordGet = function(){
			if( getCookie("userPass") != ""  ){
				this.password( getCookie("userPass") );
			}
		};*/
        this.passwordGet    = ko.computed(function () {
            if (getCookie("userPass") != "") {
                return this.password(getCookie("userPass"));
            }
        }, this);
        this.rememberPass   = ko.observable(true);
        this.passkey        = ko.observable();

        this.login = function () {
            dologin();
        };
    }
    var modelobj = new loginViewModel();
    ko.applyBindings(modelobj, document.getElementById("loginContainer"));

    function checkTime(i) {
        if (i < 10) { i = "0" + i }
        return i;
    }

    function getNowFormatDate() {
        var date = new Date();
        var month = date.getMonth() + 1;
        var strDate = date.getDate();
        var currentdate = date.getFullYear() + checkTime(month) + checkTime(strDate) + checkTime(date.getHours()) + checkTime(date.getMinutes());
        return currentdate;
    }

    // var loginForm = $("#loginForm").Validform({
    //     btnSubmit: "#loginSubmit",
    //     tiptype: 3,
    //     showAllError: true,
    //     callback: function (form) {
    //         if (modelobj.rememberPass() == true) {
    //             setCookie("userName", modelobj.userName(), 1095);
    //             setCookie("userPass", modelobj.password(), 1095);
    //         };
    //         var name = $.trim(modelobj.userName());
    //         console.log(form)
    //         // dologin()
    //         // var pssd = $.trim(modelobj.password());
    //         $.ajax({
    //             type: "get",
    //             url: "",
    //             async: false,
    //             cache: false,
    //             data: { name: name },
    //             dataType: "json",
    //             error: function (e) {
    //                 layer.alert("获取密钥失败，请稍后再试！", {
    //                     skin: 'layui-layer-lan',//样式类名
    //                     closeBtn: 0
    //                 });
    //             },
    //             success: function (data) {
    //                 console.log(data)
    //                 if (data.code == 200) {
    //                     modelobj.passkey(data.info);
    //                     dologin();
    //                 } else {
    //                     layer.alert(data.message + '!', {
    //                         skin: 'layui-layer-lan',//样式类名
    //                         closeBtn: 0
    //                     });
    //                 }
    //             }
    //         });

    //         return false; //阻止表单提交
    //     }
    // });

    function dologin() {
        var user = $.trim(modelobj.userName());
        var pssd = $.trim(modelobj.password());
        // var formData = new FormData();
        // formData.append("user", user);
        // formData.append("password", hex_md5(hex_md5(modelobj.passkey() + hex_md5(pssd)) + getNowFormatDate()));
        var jsonData = {
            user: user,
            password : hex_md5(hex_md5(modelobj.passkey() + hex_md5(pssd)) + getNowFormatDate()),
        }
        console.log(jsonData)
        $.ajax({
            url: "login",
            type: "POST",
            data: jsonData,
            async: true,
            cache: false,
            error: function () {
                layer.alert("登录失败，请稍后再试！", {
                    skin: 'layui-layer-lan',//样式类名
                    closeBtn: 0
                });
            },
            success: function (result) {
                console.log(result)
                if (result.code == 200) {
                    window.location.href = 'index';
                } else {
                    layer.alert(result.message + '!', {
                        skin: 'layui-layer-lan',//样式类名
                        closeBtn: 0
                    });
                }
            }
        });
    }

    // loginForm.addRule([
    //     {
    //         ele: "#username",
    //         datatype: "*1-15",
    //         nullmsg: "请输入用户名！",
    //         errormsg: "请输入1-15字符长度，支持汉字、字母、数字及_ !",
    //         sucmsg: ""
    //     },
    //     {
    //         ele: "#password",
    //         datatype: "*5-30",
    //         nullmsg: "请输入密码！",
    //         errormsg: "请输入5-30位密码，支持字母、数字及_ !",
    //         sucmsg: ""
    //     }
    // ]);
})






