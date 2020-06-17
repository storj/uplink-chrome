console.log("in sockets.js")

var socketId = ""
var readData = new Uint8Array([])

function socketRead(x) {
    console.log('socket read ')
    setTimeout(function() {
        console.log("after timeout")
    }, 2000)
    console.log(readData.length)
    readData = new Uint8Array([])
    return readData;
}

function socketWrite(ip, port, buf) {
    console.log("socket write " + buf)
    return new Promise((resolve, reject) => {
        if (socketId == "") {
            chrome.sockets.tcp.create({}, function(createInfo) {
                socketId = createInfo.socketId
                chrome.sockets.tcp.connect(socketId, ip, 10000, function(result) {
                    chrome.sockets.tcp.send(socketId, buf, function() {
                        console.log("tcp message sent")
                        resolve()
                    });
                });
                chrome.sockets.tcp.onReceive.addListener(onReceiveCallback)
            });
        } else {
            chrome.sockets.tcp.send(socketId, buf, function() {
                console.log("tcp message sent")
                resolve()
            });
        }
    });

}
function onReceiveCallback(info) {
    readData  = new Uint8Array(info.data)
    console.log(readData.length)
}
function onSentCallback() {
    console.log("onsentcallback")
}
function onConnectedCallback() {
    console.log("onconnectedcallback")
}
