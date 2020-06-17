console.log("in sockets.js")

var socketConns = {}
var socketId = ""
var readData = new Uint8Array([])

chrome.sockets.tcp.onReceive.addListener(onReceiveCallback)

function socketRead(ip, port, bytesRequested, done) {
    console.log('socket read (JS)')
    addr = ip + ":" + port
    bufferInfo = socketConns[addr]
    if (!bufferInfo) {
        return
    }
    setTimeout(function() {
        toSend = bufferInfo.readData
        if (bufferInfo.readData.length > bytesRequested) {
            toSend = bufferInfo.readData.slice(0, bytesRequested)
            bufferInfo.readData = bufferInfo.readData.slice(bytesRequested)
        }
        console.log('socket read end (JS); sending bytes: ' + toSend.length)
        done(toSend)
    }, 0)
    return;
}

function socketWrite(ip, port, buf, done) {
    console.log("socket write (JS): " + buf.length)
    addr = ip + ":" + port
    bufferInfo = socketConns[addr]
    if (!bufferInfo) {
        chrome.sockets.tcp.create({}, function(createInfo) {
            socketConns[addr] = {
                id: createInfo.socketId,
                readData: new Uint8Array([])
            }
            console.log('set socketCoons ' + addr)
            socketId = createInfo.socketId
            chrome.sockets.tcp.connect(socketId, ip, parseInt(port), function(result) {
                chrome.sockets.tcp.send(socketId, buf, function() {
                    console.log("tcp message sent (JS)")
                    done()
                });
            });
        });

    } else {
        chrome.sockets.tcp.send(bufferInfo.id, buf, function() {
            console.log("tcp message sent (JS)")
            done()
        });
    }
}
function onReceiveCallback(info) {
    for (let [key, item] of Object.entries(socketConns)) {
        if (item.id == info.socketId) {
            var newData = new Uint8Array(info.data)
            var tmp = new Uint8Array(item.readData.byteLength + newData.byteLength);
            tmp.set(new Uint8Array(item.readData), 0);
            tmp.set(new Uint8Array(newData), item.readData.byteLength);
            item.readData = tmp
        }
    }
}
function onSentCallback() {
}
function onConnectedCallback() {
}
function socketDisconnect(ip, port, done) {
    console.log('disconnect socket')
    var addr = ip + ":" + port
    if (socketConns[addr]) {
        chrome.sockets.tcp.disconnect(socketConns[addr].id, done)
    } else {
        done()
    }
}
