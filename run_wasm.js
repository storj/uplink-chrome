'use strict'

window.initApp = function initApp(goApp) {
  const satAddr = '12rveQZBa7gDyYseU6tK5UVLFG3gjJcooNt1wDbumfMgwkLMD3@0.0.0.0:10000'
  const apiKey = '13YqcvnMrR2QV76cNLiGbyxwU8yTsFErq3h1FELyPu6uVx8f96uTEwL9Azfb3gaJxgCwuXibedn9C2yPDs6DM63p1K7F4Vud4cuJRVq'
  const passphrase = 'test'

  goApp.uplink.loadAccess(satAddr, apiKey, passphrase)
  let msg = goApp.up(satAddr, apiKey, passphrase, "mybucket", "myfile")
  console.log(msg)

}


const go = new Go()
go.argv.push('initApp')

WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject).then((result) => {
  go.run(result.instance);
})
