package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"storj.io/uplink"
)

func main() {
	apikey := js.Global().Get("apikey").String()
	satellite := js.Global().Get("satellite").String()
	passphrase := js.Global().Get("passphrase").String()
	bucket := "test-bucket"
	path := "from/the/web"
	dataToUpload := []byte("one fish two fish red fish blue fish")
	ctx := context.Background()
	fmt.Println("starting")
	//make sure the connection parameters are all good
	project, err := PrepareBucket(ctx, satellite, apikey, passphrase, bucket)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer project.Close()
	//create JavaScript functions for upload and download
	js.Global().Set("Upload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Upload")
		err := UploadData(ctx, project, bucket, path, dataToUpload)
		if err != nil {
			fmt.Println("error:", err)
		}
		return nil
	}))
	js.Global().Set("Download", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Download")
		_, err := DownloadData(ctx, project, bucket, path, dataToUpload)
		if err != nil {
			fmt.Println("error:", err)
		}
		return nil
	}))
	//run indefinitely
	select {}
}

// PrepareBucket check for a bucket, using the specified Satellite, API key, and passphrase.
func PrepareBucket(ctx context.Context, satelliteAddress, apiKey, passphrase, bucketName string) (*uplink.Project, error) {
	// Request access grant to the satellite with the API key and passphrase.
	myConfig := uplink.Config{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			fmt.Println("dial context")
			addressParts := strings.Split(address, ":")
			port, _ := strconv.Atoi(addressParts[1])
			return NewJsConn(addressParts[0], port)
		},
	}
	access, err := myConfig.RequestAccessWithPassphrase(ctx, satelliteAddress, apiKey, passphrase)
	if err != nil {
		return nil, fmt.Errorf("could not request access grant: %v", err)
	}
	fmt.Println("\n\n>>> access grant requested successfully <<<\n")
	// Open up the Project we will be working with.
	project, err := myConfig.OpenProject(ctx, access)
	if err != nil {
		return nil, fmt.Errorf("could not open project: %v", err)
	}
	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("could not ensure bucket: %v", err)
	}
	return project, nil
}

// UploadData uploads the specified data to the specified key in the
// specified bucket, using the specified Satellite, API key, and passphrase.
func UploadData(ctx context.Context, project *uplink.Project, bucketName string, uploadKey string, dataToUpload []byte) error {
	// Intitiate the upload of our Object to the specified bucket and key.
	upload, err := project.UploadObject(ctx, bucketName, uploadKey, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}
	// Copy the data to the upload.
	buf := bytes.NewBuffer(dataToUpload)
	_, err = io.Copy(upload, buf)
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}
	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}
	return nil
}

// DownloadData uploads the specified data to the specified key in the
// specified bucket, using the specified Satellite, API key, and passphrase.
func DownloadData(ctx context.Context, project *uplink.Project, bucketName string, uploadKey string, dataToUpload []byte) ([]byte, error) {
	// Initiate a download of the same object again
	download, err := project.DownloadObject(ctx, bucketName, uploadKey, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open object: %v", err)
	}
	defer download.Close()
	// Read everything from the download stream
	receivedContents, err := ioutil.ReadAll(download)
	if err != nil {
		return nil, fmt.Errorf("could not read data: %v", err)
	}
	// Check that the downloaded data is the same as the uploaded data.
	if !bytes.Equal(receivedContents, dataToUpload) {
		return nil, fmt.Errorf("got different object back: %q != %q", dataToUpload, receivedContents)
	}
	fmt.Printf("**** got back \"%s\" ****\n", string(receivedContents))
	return receivedContents, err
}

//JsConn is a javascript thing
type JsConn struct {
	ip   string
	port int
	id   int
}

var uint8Array = js.Global().Get("Uint8Array")

//NewJsConn returns a new JsConn
func NewJsConn(ip string, port int) (*JsConn, error) {
	creating := make(chan struct{})
	fmt.Println("connect start (Go)")
	var socketID, resultCode js.Value
	var err error
	js.Global().Call("dialContext", ip, port,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			socketID = args[0]
			resultCode = args[1]
			fmt.Printf("connection created (Go) #%d\n", socketID.Int())
			close(creating)
			return nil
		}),
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("failed to create socket connection")
			err = fmt.Errorf("could not create socket connection")
			close(creating)
			return nil
		}))
	<-creating
	if err != nil {
		return nil, err
	}
	fmt.Printf("connect (Go) socket id = %s\n", socketID.String())
	return &JsConn{ip: ip, port: port, id: socketID.Int()}, nil
}

func (c *JsConn) Read(b []byte) (n int, err error) {
	reading := make(chan struct{})
	fmt.Println("read start (Go)")
	var retVal, eof js.Value
	js.Global().Call("readSocket", c.id, len(b),
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			retVal = args[0]
			eof = args[1]
			fmt.Printf("read close (Go) #%d\n", c.id)
			close(reading)
			return nil
		}),
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("failed to read from socket connection")
			err = fmt.Errorf("could not read from socket connection")
			close(reading)
			return nil
		}))
	<-reading
	if err != nil {
		return 0, err
	}
	fmt.Println("received read bytes (Go):")
	js.CopyBytesToGo(b, retVal)
	if eof.Bool() {
		fmt.Printf("EOF (Go) #%d\n", c.id)
		return retVal.Length(), io.EOF
	}
	return retVal.Length(), nil
}

func (c *JsConn) Write(b []byte) (n int, err error) {
	fmt.Println("write start (Go):")
	fmt.Println(len(b))
	buf := uint8Array.New(len(b))
	js.CopyBytesToJS(buf, b)

	writing := make(chan struct{})
	js.Global().Call("writeSocket", c.id, buf,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			close(writing)
			return nil
		}),
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("failed to write to socket connection")
			err = fmt.Errorf("could not write to socket connection")
			close(writing)
			return nil
		}))
	<-writing
	if err != nil {
		return 0, err
	}
	fmt.Println("write end (Go)")
	return len(b), nil
}
func (c *JsConn) Close() error {
	fmt.Println("closing socket (Go)")

	closing := make(chan struct{})
	js.Global().Call("closeSocket", c.id, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(closing)
		return nil
	}))
	<-closing
	return nil
}
func (c *JsConn) LocalAddr() net.Addr {
	fmt.Println("local addr")
	return &addr{}
}
func (c *JsConn) RemoteAddr() net.Addr {
	fmt.Println("remote addr")
	return &addr{}
}
func (c *JsConn) SetDeadline(t time.Time) error {
	fmt.Println("set deadline")
	return nil
}
func (c *JsConn) SetReadDeadline(t time.Time) error {
	fmt.Println("set read deadline")
	return nil
}
func (c *JsConn) SetWriteDeadline(t time.Time) error {
	fmt.Println("set write deadline")
	return nil
}

type addr struct {
}

func (*addr) Network() string {
	return ""
}
func (*addr) String() string {
	return ""
}
