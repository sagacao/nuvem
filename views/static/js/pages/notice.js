$(function(){

    indexModel.systemChannelHide();

    var noticeViewModel = function () {
        this.dataList               = ko.observableArray();
        this.gameid = ko.observableArray([ { name: "方块部落", id: 20201 },   
                                            { name: "国王", id: 20301 }, 
                                            { name: "春节福袋", id: 20401 }, 
                                            { name: "休闲射击", id: 20501 },
                                            { name: "猫咪", id: 20601 }
                                        ]);
        this.selectedGame = ko.observable();
        this.searchnotice = function () {
            noticesearch();
        };

        this.mystime        = ko.observable("");
        this.myetime        = ko.observable("");
        this.mytitle          = ko.observable("test");
        this.mycontent        = ko.observable("");
        this.myimage          = ko.observable("");
        this.mysharetitle     = ko.observable("");
        this.myanswer         = ko.observable("");
        this.addnotice = function () {
            showaddnotice();
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

        this.ensurenotice = function () {
            ensurenotice();
        };

        this.removenotice = function () {
            removenotice(this);
        };
    }

    var modelobj = new noticeViewModel();
    ko.applyBindings(modelobj, $("#noticediv").get(0));

    function noticesearch(){
        var jsonData = {
            gameId: modelobj.selectedGame()
        }
        $.ajax({
            url:"notice/search",
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

    function showaddnotice() {
        $("#noticepage").toggle();
        modelobj.datetimepickerObj.format = 'yyyy-mm-dd';
        $('.form_datetime').datetimepicker(modelobj.datetimepickerObj);
        $("#mystime").data("datetimepicker").setDate(new Date());;
        $("#myetime").data("datetimepicker").setDate(addDays(new Date(), 1));
        modelobj.mystime = ko.observable( $("#mystime").val());
        modelobj.myetime = ko.observable( $("#myetime").val());
        $('#ensure').trigger("click");
    }

    function ensurenotice() {
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("stime", modelobj.mystime());
        formData.append("etime", modelobj.myetime());
        formData.append("title", modelobj.mytitle());
        formData.append("content", modelobj.mycontent());
        formData.append("image", modelobj.myimage());
        formData.append("sharetitle", modelobj.mysharetitle());
        formData.append("answer", modelobj.myanswer());
        console.log(modelobj.selectedGame(), modelobj.mytitle(), modelobj.mystime(), modelobj.myetime())
        $.ajax({
            url:"notice/add",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                $("#noticepage").toggle();
                noticesearch();
            }
        });
    }

    function removenotice(data) {
        console.log(data)
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("id", data.id);
        $.ajax({
            url:"notice/del",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                noticesearch();
            }
        });
    }
});