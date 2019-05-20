$(function(){

    indexModel.systemChannelHide();

    var gatesViewModel = function () {
        this.dataList               = ko.observableArray();
        this.gameid = ko.observableArray([ { name: "方块部落-20201", id: 20201 },
                { name: "国王-20301", id: 20301 }, 
                { name: "春节福袋-20401", id: 20401 }, 
                { name: "休闲射击-20501", id: 20501 },
                { name: "猫咪-20601", id: 20601 }
            ]);
        this.selectedGame = ko.observable();
        this.searchgates = function () {
            gatessearch();
        };
    }

    var modelobj = new gatesViewModel();
    ko.applyBindings(modelobj, $("#gatesdiv").get(0));

    function gatessearch(){
        var jsonData = {
            gameId: modelobj.selectedGame()
        }
        $.ajax({
            url:"gates/search",
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

    function removegates(data) {
        console.log(data)
        var formData = new FormData();
        formData.append("gameId", modelobj.selectedGame());
        formData.append("name", data.name);
        $.ajax({
            url:"gates/del",
            type: 'POST',
            data: formData,
            processData:false,
            contentType:false,
            error:function(){
                console.log("notice ajax error.")
            },
            success:function (info) {
                gatessearch();
            }
        });
    }
});