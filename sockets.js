console.log("in sockets.js")

var socketId = ""

function socketRead(x) {
    console.log('socket read ' + x)
    return new Uint8Array([21, 31]);
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
            });
        } else {
            chrome.sockets.tcp.send(socketId, buf, function() {
                console.log("tcp message sent")
                resolve()
            });
        }
    });

}
function onSentCallback() {
    console.log("onsentcallback")
}
function onConnectedCallback() {
    console.log("onconnectedcallback")
}
