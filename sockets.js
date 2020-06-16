console.log("in sockets.js")

function socketFunc(x) {
    console.log("SOCKETS FUNCTION" + x)
    return new Uint8Array([21, 31]);
}
function socketWrite(x) {
    console.log("socket write " + x)
}
