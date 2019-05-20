$(function(){

    indexModel.systemChannelHide();

    var magicViewModel = function () {
        this.dataList               = ko.observableArray();
        this.gameid = ko.observableArray([ { name: "方块部落", id: 20201 },  
                                { name: "国王", id: 20301 }, 
                                { name: "春节福袋", id: 20401 }, 
                                { name: "休闲射击", id: 20501 },
                                { name: "猫咪", id: 20601 }
                            ]);
        this.selectedGame = ko.observable();
        this.searchmagic = function () {
            magicsearch();
        };

        this.myetime        = ko.observable("");
        this.magicnumber          = ko.observable("uztj6cksgcSqlY");
        this.times        = ko.observable("1");
        this.coin          = ko.observable("100");
        this.addmagic = function () {
            showaddmagic();
        };

        //日历时间
        this.datetimepickerObj  = { 
            language: 'zh-CN',
            weekStart: 1,
            todayBtn: 1,
            autoclose: 1,
            todayHighlight: 1,
            startView: 2,
            minView: 2,
            forceParse: 0
        };

        this.ensuremagic = function () {
            ensuremagic();
        };

        this.removemagic = function () {
            removemagic(this);
        };
    }

    var modelobj = new magicViewModel();
    ko.applyBindings(modelobj, $("#magicdiv").get(0));

    function magicsearch(){
        var jsonData = {
            gameId: modelobj.selectedGame()
        }
        $.ajax({
            url:"magic/search",
            async:false,
            dataType:"json",
            data: jsonData,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                console.log(info.data);
                if (info.data) {
                    modelobj.dataList(info.data);
                } else {
                    console.log("notice info.rows is no data.")
                }
            }
        });
    }

    function showaddmagic() {
        $("#magicpage").toggle();
        modelobj.datetimepickerObj.format = 'yyyy-mm-dd';
        $('.form_datetime').datetimepicker(modelobj.datetimepickerObj);
        $("#myetime").data("datetimepicker").setDate(addDays(new Date(), 1));
        modelobj.myetime = ko.observable( $("#myetime").val());
        $('#ensure').trigger("click");
    }

    function ensuremagic() {
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("deadline", modelobj.myetime());
        formData.append("magicnumber", modelobj.magicnumber());
        formData.append("times", modelobj.times());
        formData.append("coin", modelobj.coin());
        console.log(modelobj.selectedGame(), modelobj.myetime(), modelobj.magicnumber(), modelobj.times())
        $.ajax({
            url:"magic/add",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                $("#magicpage").toggle();
                magicsearch();
            }
        });
    }

    function removemagic(data) {
        console.log(data)
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("magicnumber", data.keystr);
        $.ajax({
            url:"magic/del",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                magicsearch();
            }
        });
    }
});