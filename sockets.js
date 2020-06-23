console.log("in sockets.js")

var apikey = "13YqgBsPw9cwbLNHtw5U52FTjPE8xaMU5PUEE5Vi2t6dhfThA5AqJGvY1rXqHTUM6bWATMWx1rywRhfatEmXLUSe6qSE4x7bmvEvZvF";
var satellite = "12ndYNoiwjW8zygsvLKEq9GEckV81W6mL6vQ2QgTeSpuimFWfLf@127.0.0.1:10000";
var passphrase = "testpass";

class Socket {
    constructor(info) {
        this.id = info.socketId;
        this.info = info;
        this.readData = {
            buffer: new Uint8Array(),
            EOF: false
        }
        this.resultCode = 0;

        this.receivedData = null;
        this.receivedFailure = null;
    }

    readBytes(count, read, failed) {
        if (this.readData.buffer.byteLength == 0) {
            if (this.resultCode != 0) {
                failed(this.resultCode);
                return;
            }
            let self = this;
            this.receivedData = function () {
                self.receivedData = null;
                self.receivedFailure = null;

                self.writeFromBuffer(count, read, failed);
            };
            this.receivedFailure = function () {
                self.receivedData = null;
                self.receivedFailure = null;

                failed()
            };
            return;
        }

        this.writeFromBuffer(count, read, failed);
    }

    writeFromBuffer(count, read, failed) {
        var n = Math.min(this.readData.buffer.byteLength, count);
        var data = this.readData.buffer.slice(0, n);
        this.readData.buffer = this.readData.buffer.slice(n);

        read(data, this.readData.EOF);
    }

    appendToReadBuffer(data) {
        if (this.resultCode != 0) return;

        if (data.byteLength == 0) {
            this.readData.EOF = true;
        }

        var tmp = new Uint8Array(this.readData.buffer.byteLength + data.byteLength);
        tmp.set(new Uint8Array(this.readData.buffer), 0);
        tmp.set(new Uint8Array(data), this.readData.buffer.byteLength);

        this.readData.buffer = tmp;
        if (this.receivedData !== null) this.receivedData();
    }

    appendError(resultCode) {
        if (this.resultCode != 0) return;

        this.readData.EOF = true;

        this.resultCode = resultCode;
        if (this.receivedFailure !== null) this.receivedFailure();
    }
}

let alive = {};

function dialContext(ip, port, connected, failed) {
    console.log('creating connection js')
    chrome.sockets.tcp.create({}, function (createInfo) {
        let sock = new Socket(createInfo);
        alive[sock.id] = sock;
        chrome.sockets.tcp.connect(sock.id, ip, port, function (result) {
            if (result < 0) {
                failed(result);
                delete (alive[sock.id]);
                return;
            }

            console.log('connected in js')
            connected(sock.id, result);
        })
    });
}

function writeSocket(id, data, wrote, failed) {
    let sock = alive[id];
    chrome.sockets.tcp.send(sock.id, data, function (sendInfo) {
        if (sendInfo.result < 0) {
            failed(sendInfo.bytesSent);
            return;
        }
        wrote(sendInfo.bytesSent);
    });
}

function readSocket(id, count, read, failed) {
    let sock = alive[id];
    sock.readBytes(count, read, failed);
}

function closeSocket(id, closed, failed) {
    chrome.sockets.tcp.disconnect(id, closed)
}

chrome.sockets.tcp.onReceive.addListener(function (info) {
    let sock = alive[info.socketId];
    sock.appendToReadBuffer(info.data);
});

chrome.sockets.tcp.onReceiveError.addListener(function (info) {
    let sock = alive[info.socketId];
    sock.appendError(sock.resultCode);
});
