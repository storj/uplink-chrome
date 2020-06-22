# Storj-webassembly

Just tinkering!

This repo shows how WebAssembly may be used along with a custom JS network adapter, to allow the Go Uplink code to be used from a Chrome App.

The Storj Uplink client currently cannot run in a typical browser (though tools like Gateway facilitate this). 
This is because browsers do not allow opening TCP / UDP sockets, as this is generally considered to be a security vulnerability. 
A Chrome 'App', however, comes with a manifest that allows it to request permission to use sockets.

# How to build Go into WebAssembly

Just run `GOOS=js GOARCH=wasm go build -o main.wasm main.go` or `make build`

# Configure

Update `apikey`, `satellite`, and `passphrase` in sockets.js.

# How to run / debug the Chrome App

Go to chrome://extensions/. Turn on Developer mode if it's not on.  Click "Load unpacked", and browse-to and select the desired App directory structure.  It should add to the Chrome Apps section.  Then go to chrome://apps/.  You should see it there as well.  It will likely open a window.  Right click and choose "Inspect".  That should show the DevTools window where you can view the Console tab.  You can alternatively see errors and reload the App from chrome://extensions/.

# License

TBD


