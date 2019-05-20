$(function(){

    indexModel.systemChannelHide();

    var shareconfigViewModel = function () {
        this.dataValue               = ko.observable();
        this.gameid = ko.observableArray([ { name: "方块部落", id: 20201 },  
                                            { name: "国王", id: 20301 }, 
                                            { name: "春节福袋", id: 20401 }, 
                                            { name: "休闲射击", id: 20501 },
                                            { name: "猫咪", id: 20601 }
                                        ]);
        this.selectedGame = ko.observable();
        this.searchshareconfig = function () {
            shareconfigsearch();
        };

        this.mycontent        = ko.observable("");
        this.addshareconfig = function () {
            showaddshareconfig();
        };


        this.ensureshareconfig = function () {
            ensureshareconfig();
        };
    }

    var modelobj = new shareconfigViewModel();
    ko.applyBindings(modelobj, $("#shareconfigdiv").get(0));

    function shareconfigsearch(){
        var jsonData = {
            gameId: modelobj.selectedGame()
        }
        $.ajax({
            url:"config/search",
            async:false,
            dataType:"json",
            data: jsonData,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                console.log(info.data);
                if (info.data) {
                    codestr = uncompileStr(info.data)
                    modelobj.dataValue(codestr);
                } else {
                    console.log("notice info.rows is no data.")
                }
            }
        });
    }

    function showaddshareconfig() {
        $("#shareconfigpage").toggle();
    }

    function compileStr(code){ 
        //对字符串进行加密       
        var c=String.fromCharCode(code.charCodeAt(0)+code.length);
        for(var i=1;i<code.length;i++){      
            c+=String.fromCharCode(code.charCodeAt(i)+code.charCodeAt(i-1));
        }   
        return escape(c);   
    }

    function uncompileStr(code){
        code=unescape(code);      
        var c=String.fromCharCode(code.charCodeAt(0)-code.length);      
        for(var i=1;i<code.length;i++){      
            c+=String.fromCharCode(code.charCodeAt(i)-c.charCodeAt(i-1));      
        }      
        return c;   
    }
    

    function ensureshareconfig() {
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        valuedata = compileStr(modelobj.mycontent())
        formData.append("value", valuedata);
        console.log(modelobj.selectedGame(), modelobj.mycontent())
        $.ajax({
            url:"config/update",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                $("#shareconfigpage").toggle();
                shareconfigsearch();
            }
        });
    }
});