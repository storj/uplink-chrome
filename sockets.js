console.log("in sockets.js")

function socketRead(x) {
    console.log(socket read)
    return new Uint8Array([21, 31]);
}
function socketWrite(ip, port, buf) {
    console.log("socket write " + buf)

    chrome.sockets.tcp.create({}, function(createInfo) {
        chrome.sockets.tcp.connect(createInfo.socketId,
            ip, port, onConnectedCallback);
    });

    chrome.sockets.tcp.send(socketId, buf, onSentCallback);
}
function onSentCallback() {
    console.log("onsentcallback")
}
function onConnectedCallback() {
    console.log("onconnectedcallback")
}
