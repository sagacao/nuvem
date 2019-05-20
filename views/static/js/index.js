var indexViewModel = function(){
    var self                = this;
    this.typeOfSystem       = ko.observable(true); 
    this.filterByDate       = ko.observable(false); 
    this.filterByChannel    = ko.observable(false); 
    this.onlyTime           = ko.observable(true);
    
    //标签栏的各种状态
    this.allShow            = function () {
        self.typeOfSystem(true);
        self.filterByDate(true);
        self.filterByChannel(true);
    };
    this.allHide            = function(){
        self.typeOfSystem(false);
        self.filterByDate(false);
        self.filterByChannel(false);
    }
    this.channelZeuidHide   = function () {
        self.typeOfSystem(true);
        self.filterByDate(true);
        self.filterByChannel(false);
    };
    this.zeuidHide          = function(){
        self.typeOfSystem(true);
        self.filterByDate(true);
        self.filterByChannel(true);
    };
    this.channelHide = function () {
        self.typeOfSystem(true);
        self.filterByDate(true);
        self.filterByChannel(false);
    };
    this.dataChannelHide    = function(){
        self.typeOfSystem(true);
        self.filterByDate(false);
        self.filterByChannel(false);
        self.filterByZeuid(true);
    };
    this.systemChannelHide  = function(){
        self.typeOfSystem(false);
        self.filterByDate(true);
        self.filterByChannel(false);
    };
    this.onlyZeuidShow      = function(){
        self.typeOfSystem(false);
        self.filterByDate(false);
        self.filterByChannel(false);
        self.filterByZeuid(true);
    };
    this.onlySystemShow     = function(){
        self.typeOfSystem(true);
        self.filterByDate(false);
        self.filterByChannel(false);
    }
    this.onlyTimeShow     = function () {
        self.typeOfSystem(false);
        self.filterByDate(true);
        self.filterByChannel(false);
    }
    this.SystemHide         = function(){
        self.typeOfSystem(false);
        self.filterByDate(true);
        self.filterByChannel(true);
    }
    this.onlyTimeFn         = function () {
        this.onlyTime(false);
    }.bind(this);
    
    //configUrl.实时数据
    this.configUrl = {
        "实时在线数据": {
            url: "/views/templates/online.html",
        },
        "分享配置": {
            url: "/views/static/pages/share-config.html",
        },
        "网关": {
            url: "/views/static/pages/gates.html",
        },
        "激活码": {
            url: "/views/static/pages/magic.html",
        },
        "功能开关": {
            url: "/views/static/pages/func-switch.html",
        },
        "公告": {
            url: "/views/static/pages/notice.html",
        }
    };
    // //渠道值
    // this.channelIdChange    = ko.observable("");   
    // this.channelidList      = ko.observableArray();
    // this.channelIdChangeFn  = ko.computed(function () {
    //     if (this.channelIdChange() == "") {
    //         getChannelId = '';                  //渠道默认值清除
    //         $("#areaClothing").removeAttr("data-id");
    //     }
    // }, this);
    // this.getChannelValue    = function (data, event) {     
    //     var idCurrent = $(event.currentTarget).attr("id");
    //     $(event.currentTarget).parents("ul").siblings("input").attr("data-id", idCurrent);
    //     indexModel.channelIdChange($(event.currentTarget).children("a").text());
    //     $("#channelGo").trigger("click");
    // };
    // this.channelGo          = function (data, event) {
    //     if (indexModel.channelIdChange() == "") {
    //         getChannelId = '';                  //渠道默认值清除
    //         reload();   
    //     } else {
    //         getChannelId = $(event.currentTarget).prev("input").attr("data-id");
    //         reload();
    //     };
    // };
    
    //ztree树配置
    this.userSettingTree    = {          
        view: {							//可视界面相关配置
            dblClickExpand: false,		//双击节点时，是否自动展开父节点的标识
            showLine: true,				 //设置是否显示节点与节点之间的连线
            selectedMulti: false		//设置是否能够同时选中多个节点
        },
        data: {							//数据相关配置
            simpleData: {
                enable: true,			//设置是否启用简单数据格式
                idKey: "id",			//设置启用简单数据格式时id对应的属性名称
                pIdKey: "pId",			//设置启用简单数据格式时parentId对应的属性名称
                rootPId: ""                    
            }
        },
        callback: {
            //beforeClick: 用于捕获单击节点之前的事件回调函数，并且根据返回值确定是否允许单击操作
            beforeClick: function (treeId, treeNode) {
                //treeId: 对应 zTree 的 treeId
                //treeNode: 被单击的节点 JSON 数据对象
                var zTree = $.fn.zTree.getZTreeObj("tree");
                //获取 id 为 tree 的 zTree 对象
                if (treeNode.isParent) {    //查看当前被选中的节点是否是父节点
                    zTree.expandNode(treeNode, null, null, null, true);
                    //expandNode: 展开 / 折叠 指定的节点
                    return false;
                } else {
                    var name = treeNode.name;
                    if (indexModel.configUrl[name]) {
                        $(".pageContainer").attr("data-url",indexModel.configUrl[name].url);
                        $(".pageContainer").load(indexModel.configUrl[name].url+"?number="+Math.random(),function(){
                        });
                        
                        if (indexModel.configUrl[name].callback) {
                            (indexModel.configUrl[name].callback)();    //self.onlyTime(false)
                        } else{
                            self.onlyTime(true);
                        }

                        // $("#deStartTime").data("datetimepicker").setDate(self.oneDayTime());
                        // $(".time-start").text($("#deStartTime").val());
                    }
                    return true;
                }
            },
            
            beforeExpand: function (treeId, treeNode) {     //用于捕获父节点展开之前的事件回调函数，并且根据返回值确定是否允许展开操作

                var pNode = curExpandNode ? curExpandNode.getParentNode() : null;
                //getParentNode: 获取 treeNode 节点的父节点
                var treeNodeP = treeNode.parentTId ? treeNode.getParentNode() : null;
                
                var zTree = $.fn.zTree.getZTreeObj("tree");
                for (var i = 0, l = !treeNodeP ? 0 : treeNodeP.children.length; i < l; i++) {
                    if (treeNode !== treeNodeP.children[i]) {
                        zTree.expandNode(treeNodeP.children[i], false);
                    }
                }
                while (pNode) {
                    if (pNode === treeNode) {
                        break;
                    }
                    pNode = pNode.getParentNode();
                }
                if (!pNode) {
                    singlePath(treeNode);
                }

            },
            onExpand: function onExpand(event, treeId, treeNode) {
                curExpandNode = treeNode;
            }

        }
    };
    //logout & userName
    this.userName = ko.observable("shaonian");
    this.logout             = function(){
        $.ajax({
            type:"GET",
            url:"",
            async:true,
            cache:false,
            dataType:"json",
            error:function(){
                layer.open({
                    skin: 'layui-layer-lan',
                    closeBtn:1,
                    shadeClose:true,
                    content:'请求失败,请稍后再试！'
                });
            },
            success:function(){
                if (data.code == 200) {
                    window.location.herf = 'login.html';
                }else{
                    layer.open({
                        skin: 'layui-layer-lan',
                        closeBtn: 1,
                        shadeClose:true,
                        content:'请求失败！'
                    });
                }
            }
        });
    }
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
    this.timePicker         = ko.observable();
    
    self.deStartTime        = ko.observable("");
    self.deEndTime          = ko.observable("");

    self.oneDayTime         = ko.observable(new Date());
    self.twoDayStarTime     = ko.observable(addDays(new Date(), -6));
    self.twoDayEndTime      = ko.observable(new Date());      
    
    self.fourteenStarTime   = ko.observable(addDays(new Date(), -13));
    self.fourteenEndTime    = ko.observable(new Date());

    this.timeError          = ko.computed(function () {
        if (this.deStartTime() != "" && this.deEndTime() != "" && this.onlyTime()) {
            var startTime = this.deStartTime(),
                endTime   = this.deEndTime(),
                startNum  = parseInt(startTime.replace(/-/g, ''), 10),
                endNum    = parseInt(endTime.replace(/-/g, ''), 10);
            return startNum > endNum ? this.timeErrorShow(true) : this.timeErrorShow(false);
        }
    }, this);
    this.timeErrorShow      = ko.observable(false);
    this.today              = function () {
        $("#deStartTime").data("datetimepicker").setDate(new Date());
        $("#deEndTime").data("datetimepicker").setDate(new Date());
        $('#ensure').trigger("click");
    };
    this.yesterday          = function () {
        $("#deStartTime").data("datetimepicker").setDate(addDays(new Date(), -1));
        $("#deEndTime").data("datetimepicker").setDate(addDays(new Date(), -1));
        $('#ensure').trigger("click");
    };
    this.lastSeven          = function () {
        $("#deStartTime").data("datetimepicker").setDate(addDays(new Date(), -6));
        $("#deEndTime").data("datetimepicker").setDate(new Date());
        $('#ensure').trigger("click");
    };
    this.lastThirty         = function () {
        $("#deStartTime").data("datetimepicker").setDate(addDays(new Date(), -29));
        $("#deEndTime").data("datetimepicker").setDate(new Date());
        $('#ensure').trigger("click");
    };
    this.timeAll            = function () {
        $("#deStartTime").data("datetimepicker").setDate(new Date(0));
        $("#deEndTime").data("datetimepicker").setDate(new Date());
        $('#ensure').trigger("click");
    };
    this.resetTime          = function () {
        this.deStartTime("");
        this.deEndTime("");
    };
    this.correctTime        = function () {
        self.oneDayTime(new Date(self.deStartTime()));
        
        $(".time-start").text($("#deStartTime").val());
        $(".time-end").text($("#deEndTime").val());
        reload();
        $('#dropdown_time').trigger("click");
    };
}

var indexModel = new indexViewModel();
var curExpandNode = null;
function singlePath(newNode) {
    if (newNode === curExpandNode) return;

    var zTree = $.fn.zTree.getZTreeObj("tree"),
        rootNodes, tmpRoot, tmpTId, i, j, n;

    if (!curExpandNode) {
        tmpRoot = newNode;
        while (tmpRoot) {
            tmpTId = tmpRoot.tId;
            tmpRoot = tmpRoot.getParentNode();
        }
        rootNodes = zTree.getNodes();
        for (i = 0, j = rootNodes.length; i < j; i++) {
            n = rootNodes[i];
            if (n.tId != tmpTId) {
                zTree.expandNode(n, false);
            }
        }
    } else if (curExpandNode && curExpandNode.open) {
        if (newNode.parentTId === curExpandNode.parentTId) {
            zTree.expandNode(curExpandNode, false);
        } else {
            var newParents = [];
            while (newNode) {
                newNode = newNode.getParentNode();
                if (newNode === curExpandNode) {
                    newParents = null;
                    break;
                } else if (newNode) {
                    newParents.push(newNode);
                }
            }
            if (newParents != null) {
                var oldNode = curExpandNode;
                var oldParents = [];
                while (oldNode) {
                    oldNode = oldNode.getParentNode();
                    if (oldNode) {
                        oldParents.push(oldNode);
                    }
                }
                if (newParents.length > 0) {
                    zTree.expandNode(oldParents[Math.abs(oldParents.length - newParents.length) - 1], false);
                } else {
                    zTree.expandNode(oldParents[oldParents.length - 1], false);
                }
            }
        }
    }
    curExpandNode = newNode;
}
function reload() {
    
    // 遍历 input标签， 获取value存入 sessionStorage
    $(".form-inline").find('input').each(function (i, v) {
        window.sessionStorage.setItem(`autosaveinput${i}`, $(v).val())
    })

    // 加载页面
    $(".pageContainer").load($(".pageContainer").attr("data-url")+"?number="+Math.random());

    setTimeout(() => {
        if (window.sessionStorage.getItem("autosaveinput0")) {
            // 从 sessionStorage 获取值，填入input
            $(".form-inline").find('input').each(function (i, v) {
                $(v).val(window.sessionStorage.getItem(`autosaveinput${i}`))
            })
        }
    }, 50);
}
var getChannelId, getServerId, zeusidsGlobal, platGlobal, platGlobalKey, platGlobalValue, moneyChangeReason, spinner, rechargeGlobal;

function getStartDate(){
    return $(".time-start").text();
}
function getEndDate(){
    return $(".time-end").text();
}
function getTimeList(startTime, endTime,n) {          //返回特定时间格式的数组
    let startNum = new Date(startTime).getTime(),
        endNum = new Date(endTime).getTime();

    var timeStep = (endNum - startNum) / (1000 * 60 * 60 * 24) + 1;
    var ArrayTime = [];

    for (let i = 0; i < timeStep; i+=n) {
        ArrayTime.push(new Date(startNum).format("yyyy-MM-dd"));
        startNum += 1000 * 60 * 60 * 24 * n;
    }
    return ArrayTime;
}

$(document).ajaxError(function (event, jqxhr, settings) {
    layer.alert("Unexpected error, please try again later!", {
        skin: 'layui-layer-lan',
        closeBtn: 1
    });
});

$(function(){
    
    ko.applyBindings(indexModel);       //绑定视图模型
    var minHeight = window.innerHeight - $("header").height()-14;  
    $(".pageContainer").css("min-height",minHeight);

    (function () {                          //默认加载
        indexModel.datetimepickerObj.format = 'yyyy-mm-dd';
        $('.form_datetime').datetimepicker(indexModel.datetimepickerObj);
        
        // $("#deStartTime").data("datetimepicker").setDate(addDays(new Date(), -6));
        // $("#deEndTime").data("datetimepicker").setDate(new Date());
        // $(".time-start").text($("#deStartTime").val());
        // $(".time-end").text($("#deEndTime").val());
    })();

    $.getJSON("/views/static/index.json", function (result) {      //初始化zTree
        var zTreeObjUserIndex = $.fn.zTree.init($("#tree"), indexModel.userSettingTree, result.TreeList);
    });

    var opts = {                        //spinner 全局配置
        lines: 7, // The number of lines to draw
        length: 0, // The length of each line
        width: 10, // The line thickness
        radius: 18, // The radius of the inner circle
        scale: 0.85, // Scales overall size of the spinner
        corners: 1, // Corner roundness (0..1)
        color: '#3e4452', // CSS color or array of colors
        fadeColor: 'transparent', // CSS color or array of colors
        opacity: 0.25, // Opacity of the lines
        rotate: 0, // The rotation offset
        direction: 1, // 1: clockwise, -1: counterclockwise
        speed: 1, // Rounds per second
        trail: 60, // Afterglow percentage
        fps: 20, // Frames per second when using setTimeout() as a fallback in IE 9
        zIndex: 2e9, // The z-index (defaults to 2000000000)
        className: 'spinner', // The CSS class to assign to the spinner
        top: '26%', // Top position relative to parent
        left: '50%', // Left position relative to parent
        position: 'absolute' // Element positioning
    };
    spinner = new Spinner(opts);
    
});