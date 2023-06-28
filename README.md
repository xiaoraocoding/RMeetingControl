# RMeetingControl
### 1.项目介绍（当前版本为V0.1）

使用go语言开发的一个在线会议demo。可以用来理解在线会议的部分业务。当前主要功能包含：

- 创建会议室
- 用户加入会议室
- 用户在会议室发送群聊信息
- 用户退出会议室(在修改中)
- 会议管理者导出当前参会人员
- 会议管理者修改当前会议室能否开启摄像头
- 会议管理者修改当前会议室能否发言
- 会议管理者修改当前会议室能否修改名称



### 2.启动项目

电脑只需要安装redis以及go相关环境即可运行。

1. 使用docker安装redis命令：

````shell
docker run --name redis -p 6379:6379 -d redis
````

2. 进入文件，使用命令：

   ````shell
   go run main.go
   ````

   

### 3.API文档

**http文档：**

[apipost](https://console-docs.apipost.cn/preview/dad7e1a7278b281e/2a34280aaed22dac)

**websocket（需要使用浏览器进行验证）：**

````html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>WebSocket测试页面</title>
</head>
<body>
<div>
    <button onclick="connectWebSocket()">连接 WebSocket</button>
    <button onclick="sendMessage()">发送消息</button>
</div>
<div>
    <textarea id="messageInput" rows="4" cols="50"></textarea>
</div>
<div>
    <textarea id="messageOutput" rows="4" cols="50" readonly></textarea>
</div>

<script>
    var socket;

    function connectWebSocket() {
        socket = new WebSocket("ws://localhost:8080/AddMeeting?username=a&meetingUid=7b49770e-c8fb-4750-8518-7b8e742bf5bf");
        socket.onopen = function(event) {
            console.log("WebSocket连接已建立");
        };

        socket.onmessage = function(event) {
            var message = event.data;
            console.log("收到消息：" + message);
            appendMessageToOutput("收到消息：" + message);
        };

        socket.onclose = function(event) {
            console.log("WebSocket连接已关闭");
        };

        socket.onerror = function(error) {
            console.error("WebSocket发生错误：" + error);
        };
    }

    function sendMessage() {
        var messageInput = document.getElementById("messageInput");
        var message = messageInput.value;

        socket.send(message);
        appendMessageToOutput("发送消息：" + message);

        messageInput.value = "";
    }

    function appendMessageToOutput(message) {
        var messageOutput = document.getElementById("messageOutput");
        messageOutput.value += message + "\n";
    }
</script>
</body>
</html>

````



### 4.业务流程（重点一些的模块）



**4.1 用户加入会议室，具体逻辑为**：

1. 客户端向服务端发送会议唯一ID，客户端的名称
2. 服务端将http升级链接为websocket
3. 服务端将发送有新用户的消息给当前会议室的所有成员
4. 服务端将当前会议室的其他状态发送给新成员（能否开视频，能否发言……等）

![截屏2023-06-26 下午5.42.09](https://p.ipic.vip/4m9x9i.png)

⚠️**4.2 用户发送消息业务逻辑（⚠️正在开发，还没有测试完成，当前为V0.1）：**

1. 维护一个work pool，服务端会先开启这个工作池，并且维护一个消息队列
2. 服务端开启一个goroutine监听链接的客户端（每一个链接成功的客户端有一个goroutine）
3. 每当监听到客户端有信息传递过来时，直接将消息发送给消息队列，因为消息队列有多个，可以采用不同的逻辑（哈希，加权，随机）选定将消息发送给哪一个消息队列。
4. work pool会将消息队列中的信息发送给当前会议的所有其他成员



![截屏2023-06-26 下午2.03.56](https://p.ipic.vip/33ttu6.png)



<img src="https://p.ipic.vip/74kt91.png" alt="截屏2023-06-26 下午5.41.10" style="zoom:50%;" />



**4.3 管理者修改会议能否开启摄像头（其他几个很类似）**

1. 管理者发送会议的唯一id，自身的唯一id给服务端
2. 服务端验证是否为管理者（这些信息存放在redis中，因为redis基于内存，比mysql快一些，我们的会议数据在退出后即可删除，所以业务中我没有使用mysql）
3. 服务端信息发生改变，将改变的信息发送给当前会议室所有已经连接的客户端



### future

- 将日志框架由原生log修改为zap，当前日志只打印了Err级别，需要精确到info级别
- 完成工作池和消息队列模块的开发
- 当前代码格式不规范，有部分模块可进行封装抽离



### 当前的一些思考和疑惑

- ⚠️在管控模块是否可以使用etcd对（能否静音，是否开启摄像头等）进行监听，只要其中的数据一发生修改，直接对所有链接的客户端进行广播

- 当前的所有服务考虑的都是单机，但是单机维护大量的长链接需要消费巨量的内存，将消息进行推送的瞬间单机cpu可能无法承受。能使用这样的架构进行设计吗？
- ![5331687762029_.pic](https://p.ipic.vip/wmovv2.jpg)

