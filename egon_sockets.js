class Socket {
    constructor(info) {
        this.id = info.socketId;
        this.info = info;
        this.readBuffer = new Uint8Array();
        this.resultCode = 0;

        this.receivedData = null;
        this.receivedFailure = null;
    }

    readBytes(count, read, failed) {
        if(this.readBuffer.byteLength == 0) {
            if(this.resultCode != 0) {
                failed(this.resultCode);
                return;
            }
            
            let self = this;
            this.receivedData = function() {
                self.receivedData = null;
                self.receivedFailure = null;

                self.writeFromBuffer(count, read, failed);
            };
            this.receivedFailure = function() {
                self.receivedData = null;
                self.receivedFailure = null;

                failed()
            };
            return;
        }

        this.writeFromBuffer(count, read, failed);
    }

    writeFromBuffer(count, read, failed) {
        var n = Math.min(this.readBuffer.byteLength, count);
        var data = this.readBuffer.slice(0, n);
        this.readBuffer = this.readBuffer.slice(n);

        read(data);
    }

    appendToReadBuffer(data) {
        if(this.resultCode != 0) return;

        var tmp = new Uint8Array(this.readBuffer.byteLength + data.byteLength);
        tmp.set(new Uint8Array(this.readBuffer), 0);
        tmp.set(new Uint8Array(data), this.readBuffer.byteLength);

        this.readBuffer = tmp;
        if(this.receivedData !== null) this.receivedData();
    }

    appendError(resultCode) {
        if(this.resultCode != 0) return;

        this.resultCode = resultCode;
        if(this.receivedFailure !== null) this.receivedFailure();
    }
}

let alive = {};

function dialContext(ip, port, connected, failed) {
    chrome.sockets.tcp.create({}, function(createInfo) {
        let sock = new Socket(createInfo.socketId);
        alive[sock.id] = sock;
        chrome.sockets.tcp.connect(sock.id, ip, port, function(result) {
            if(result < 0) {
                failed(result);
                delete(alive[sock.id]);
                return;
            }

            connected(sock.id);
        })
    });
}

function writeSocket(id, data, wrote, failed) {
    let sock = alive[id];
    chrome.sockets.tcp.send(sock.id, data, function(sendInfo) {
        if(sendInfo.result < 0) {
            failed(sendInfo.bytesSent);
            return;
        }
        wrote(sendInfo.bytesSent);
    });
}

function readSocket(id, countg, read, failed) {
    let sock = alive[id];
    sock.readBytes(count, read, failed);
}

chrome.sockets.tcp.onReceive.addListener(function(info){
    let sock = alive[info.socketId];
    sock.appendToReadBuffer(sock.data);
});

chrome.sockets.tcp.onReceiveError.addListener(function(info){
    let sock = alive[info.socketId];
    sock.appendError(sock.resultCode);
});
