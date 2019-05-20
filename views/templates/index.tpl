<!doctype html>
<html lang="en">
<head>
    <title>Morning Go</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
{{.title}}

<button id="sendBtn">发送消息</button>

<button id="leaveBtn">离开</button>
<script type="text/javascript" src="./socket.io.js"></script>
<script type="text/javascript">
    var socket=io.connect('localhost:80'),//与服务器进行连接
        send=document.getElementById('sendBtn'),
        leave=document.getElementById('leaveBtn');

    send.onclick=function(){
        socket.emit('new player', {"nick" : "test"});
    }

    leave.onclick=function(){
        window.location.href="about:blank";
        window.close()
        socket.emit('disconnect', 'leave');
    }

    //接收来自服务端的信息事件c_hi
    socket.on('new player',function(msg){
        alert(msg)
    })

</script>
</body>
</html>
