console.log("in sockets.js")

var apikey = "13YqeL9Pd4FUpSMedjGo7BhxfQod8wE1S1KjJoUwYvRqBEHvsJtSDgF4QuxNRPmqfGLQV7AEZq6TDrBxZck3GHxQvujfomvpQuTgzqz";
var satellite = "1wKcviARJh1xUDKnAVT9RQY7UoFDVxtyj3mYSU7jyznbPmZX7P@127.0.0.1:10000";
var passphrase = "testpass";
var Reads = {}

function socketConnect(ip, port, callback) {
    console.log("socketConnect 1 (JS): " + ip + " : " + port)
    chrome.sockets.tcp.create({}, function (createInfo) {
        console.log("socketConnect 2 (JS): " + ip + " : " + port)
        chrome.sockets.tcp.connect(createInfo.socketId, ip, parseInt(port), function (result) {
            if (result >= 0) {
                Reads[createInfo.socketId] = {
                    data: new Uint8Array(0),
                    EOF: false
                }
            }
            console.log("socketConnect 3 (JS): #" + createInfo.socketId + " : " + result)
            callback(createInfo.socketId, result)
        });
    });
}

function socketRead(id, bytesRequested, callback) {
    if (bytesRequested == 0) {
        callback(new Uint8Array(0), false);
        return
    }
    console.log('read (JS) #' + id + ' : ' + Reads[id].data.length + " bytes")
    var toSend = Reads[id].data;
    if (Reads[id].data.length > bytesRequested) {
        toSend = Reads[id].data.slice(0, bytesRequested)
        Reads[id].data = Reads[id].data.slice(bytesRequested)
    } else {
        toSend = Reads[id].data
        Reads[id].data = new Uint8Array(0)
    }
    //force other functions to be allowed to run
    setTimeout(function () { callback(toSend, Reads[id].EOF); }, 1);
}

function socketWrite(id, buf, callback) {
    let idInt = parseInt(id);
    console.log("socket write (JS) #" + id + ' : ' + buf.length + " bytes")
    chrome.sockets.tcp.send(idInt, buf, function () {
        callback(0, buf.length)
    });
}

function socketDisconnect(id, callback) {
    console.log('disconnect socket #' + id)
    // chrome.sockets.tcp.getInfo(id, function (socketInfo) {
    //     console.log('disconnect socket #' + id + " : " + (socketInfo.connected ? "connected" : "disconnected"))
    //     if (socketInfo.connected && !Reads[id].EOF) {
    chrome.sockets.tcp.disconnect(id, callback)
    //     }
    // })
}

chrome.sockets.tcp.onReceive.addListener(onReceive)
function onReceive(info) {
    let id = info.socketId;
    let item = Reads[id].data;
    if (info.data.byteLength == 0) {
        console.log("zero length in #" + id)
        Reads[id].EOF = true;
    }
    console.log("onReceive on " + info)
    var newData = new Uint8Array(info.data)
    var tmp = new Uint8Array(item.byteLength + newData.byteLength);
    tmp.set(new Uint8Array(item), 0);
    tmp.set(new Uint8Array(newData), item.byteLength);
    Reads[id].data = tmp
}

chrome.sockets.tcp.onReceiveError.addListener(onReceiveError)
function onReceiveError(info) {
    Reads[info.socketId].EOF = true;
    console.log("onReceiveError in #" + info.socketId + ' : ' + netError(info.resultCode));
}

function onSentCallback() {
}
function onConnectedCallback() {
}

