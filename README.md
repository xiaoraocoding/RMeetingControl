# RMeetingControl（V0.2）

### 1.项目介绍（当前版本为V0.2）

使用go语言开发的一个在线会议demo。可以用来理解在线会议的部分业务。当前主要功能包含：

- 创建会议室
- 用户加入会议室
- 用户在会议室发送群聊信息
- 用户退出会议室
- 会议管理者导出当前参会人员
- 会议管理者修改当前会议室能否开启摄像头
- 会议管理者修改当前会议室能否修改名称
- 服务端发送心跳链接，实时性的维护链接



### 2.启动项目

电脑只需要安装redis以及go相关环境即可运行。

1. 使用docker安装redis命令：

```shell
docker run --name redis -p 6379:6379 -d redis
```

1. 进入文件，使用命令：

   ```shell
   go run main.go
   ```

**使用流程：**

1. 用户调用 /CreateMeeting  创建会议拿到会议id等信息
2.  使用websocket链接调用 /AddMeeting  加入会议
3. 会议中途：
   - 用户可以使用websocket长链接发送群聊消息给会议室当前其他链接的用户
   - 管理员可以使用websocket长链接修改会议室的当前权限（是否能开启视频，是否能发言等。。。）
   - 管理员可以导出当前参会成员（json形式导出）
4. 用户可以调用 /LeaveMeeting 离开当前会议
5. 期间服务端会不断发送心跳来实时进行检测链接的客户端 



### 3.API文档

**http文档：**

[apipost](https://console-docs.apipost.cn/preview/dad7e1a7278b281e/2a34280aaed22dac)

ps：（此链接部分暂时不可用，很多接口已经由http转为websocket）

**websocket（需要使用浏览器进行验证，前端了解的较少，可以自己根据业务编写哈）：**

```html
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
```



### 4.业务流程（重点一些的模块）

**4.1 用户加入会议室，具体逻辑为**：

1. 客户端向服务端发送会议唯一ID，客户端的名称

2. 服务端将http升级链接为websocket

3. 服务端将发送有新用户的消息给当前会议室的所有成员

4. 服务端将当前会议室的其他状态发送给新成员（能否开视频，能否发言……等）

   

⚠️**4.2 对发送消息的封装：**

在会议业务中，很多情况需要发送信息给会议的其他成员，此模块对消息部分进行了封装，所有需要发送的信息（如会议中用户想要发送的信息，服务端传递给用户当前会议室的权限等）都将进入此模块。

1. 维护一个work pool，服务端会先开启这个工作池，并且维护一个消息队列
2. 服务端开启一个goroutine监听链接的客户端（每一个链接成功的客户端有一个goroutine）
3. 每当监听到客户端有信息传递过来时，直接将消息发送给消息队列，因为消息队列有多个，可以采用不同的逻辑（哈希，加权，随机）选定将消息发送给哪一个消息队列。
4. work pool会将消息队列中的信息发送给当前会议的所有其他成员

![20230626 20356](https://p.ipic.vip/j9vcf2.png)



**4.3 管理者修改会议能否开启摄像头（其他几个很类似）**

1. 管理者发送会议的唯一id（meetingUid），自身的唯一id（userUid），是否开启摄像头等消息传递给服务端（因为客户端与服务端之间的链接都是长链接，此处直接在长链接的基础上进行）

2. 服务端验证是否为管理者（这些信息存放在redis中，因为redis基于内存，比mysql快一些，我们的会议数据在退出后即可删除，所以业务中我没有使用mysql）

3. 服务端信息发生改变，将改变的信息发送给当前会议室所有已经连接的客户端（使用消息封装模块）

   

**4.4 服务端进行心跳链接**

1. 每当一个客户端与服务端建立链接之后，开启一个goroutine
2. 开启的goroutine每隔指定的时间发送一个信息给客户端
3. 只要当发送信息失败或者出现链接断开，即删除此链接的相关信息，并且进行广播（使用消息封装模块）



**4.5 客户端消息的封装**

1. 当客户端发送普通群聊消息的时候，无需注意，直接对websocket长链接发送消息即可。
2. 当客户端为管理员，且修改对会议室的权限时，一个类型为[]byte的message中,前四个byte为特殊消息，第5个为1的话发送的权限问题和发言相关，当发送为2的时候发送的权限问题和能否开启视频有关。（后续还在拓展中）



### future

- 将日志框架由原生log修改为zap，当前日志只打印了Err级别，需要精确到info级别
- 使用配置文件库，将写死的配置文件动态化

