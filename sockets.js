console.log("in sockets.js")

var socketId = ""
var readData = new Uint8Array([])

function socketRead(bytesRequested, done) {
    console.log('socket read ')
    setTimeout(function() {
        console.log("after timeout")
        console.log("length of read data")
        console.log(readData.length)
        toSend = readData
        if (readData.length > bytesRequested) {
            toSend = readData.slice(0, bytesRequested)
            readData = readData.slice(bytesRequested)
        }
        done(toSend)
    }, 500)
    return;
}

function socketWrite(ip, port, buf, done) {
    console.log("socket write " + buf)
    if (socketId == "") {
        chrome.sockets.tcp.create({}, function(createInfo) {
            socketId = createInfo.socketId
            chrome.sockets.tcp.connect(socketId, ip, 10000, function(result) {
                chrome.sockets.tcp.send(socketId, buf, function() {
                    console.log("tcp message sent")
                    done()
                });
            });
            chrome.sockets.tcp.onReceive.addListener(onReceiveCallback)
        });
    } else {
        chrome.sockets.tcp.send(socketId, buf, function() {
            console.log("tcp message sent")
            done()
        });
    }
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
