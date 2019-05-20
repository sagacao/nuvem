$(function(){

    indexModel.systemChannelHide();

    var funcswitchViewModel = function () {
        this.dataList               = ko.observableArray();
        this.gameid = ko.observableArray([ { name: "方块部落", id: 20201 }, 
                                            { name: "国王", id: 20301 }, 
                                            { name: "春节福袋", id: 20401 }, 
                                            { name: "休闲射击", id: 20501 },
                                            { name: "猫咪", id: 20601 }
                                        ]);
        this.selectedGame = ko.observable();
        this.searchfuncswitch = function () {
            funcswitchsearch();
        };

        this.funcname          = ko.observableArray([ { name: "invite", id: "invite" }, 
                                                    { name: "comic", id: "comic" }, 
                                                    { name: "unlockLevelMode", id: "unlockLevelMode" },
                                                    { name: "mute", id: "mute" }, 
                                                    { name: "push", id: "push" }, 
                                                    { name: "clickTips", id: "clickTips" },
                                                    { name: "forcar", id: "forcar" }, 
                                                    { name: "banner", id: "banner" }
                                                ]); //
        this.funcswitch        = ko.observableArray([ { name: "开", id: 0 }, { name: "关", id: 1 }]);
        this.selectedName = ko.observable();
        this.selectedSwitch = ko.observable();

        this.addfuncswitch = function () {
            showaddfuncswitch();
        };

        this.ensurefuncswitch = function () {
            ensurefuncswitch();
        };

        this.removefuncswitch = function () {
            removefuncswitch(this);
        };
    }

    var modelobj = new funcswitchViewModel();
    ko.applyBindings(modelobj, $("#funcswitchdiv").get(0));

    function funcswitchsearch(){
        var jsonData = {
            gameId: modelobj.selectedGame()
        }
        $.ajax({
            url:"funcswitch/search",
            async:false,
            dataType:"json",
            data: jsonData,
            error:function(){
                console.log("funcswitch ajax error.")
            },
            success:function (info) {
                console.log(info);
                if (info.data) {
                    modelobj.dataList(info.data);
                } else {
                    console.log("funcswitch info.rows is no data.")
                }
            }
        });
    }

    function showaddfuncswitch() {
        $("#funcswitchpage").toggle();
    }

    function ensurefuncswitch() {
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("funcname", modelobj.selectedName());
        formData.append("funcswitch", modelobj.selectedSwitch());
        console.log(modelobj.selectedGame(), modelobj.selectedName(), modelobj.selectedSwitch())
        $.ajax({
            url:"funcswitch/add",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("funcswitch ajax error.")
            },
            success:function (info) {
                $("#funcswitchpage").toggle();
                funcswitchsearch();
            }
        });
    }

    function removefuncswitch(data) {
        console.log(data)
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("funcname", data.funcname);
        $.ajax({
            url:"funcswitch/del",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("funcswitch ajax error.")
            },
            success:function (info) {
                funcswitchsearch();
            }
        });
    }
});