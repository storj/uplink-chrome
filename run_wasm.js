'use strict';

window.initApp = function initApp(goApp) {
 // const access = goApp.funcs.access('addr', 'api-key', 'pwd')
  //console.log(access)
  const s = {name: 'ivan'}
  s.access = goApp.funcs.access
  s.accessThis = goApp.funcs.accessThis

  let res = s.access('satAddr', 'apiKey', 'pwd')
  console.log(res)

  res = s.accessThis('satAddr', 'apiKey', 'pwd')
  console.log(res)

  res = goApp.funcs.different('a string', true, 34.4, 152512, -85)
  console.log(res)
 }


const go = new Go()
go.argv.push('initApp')

WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject).then((result) => {
  go.run(result.instance);
})
