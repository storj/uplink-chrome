package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"syscall/js"
	"time"

	"storj.io/uplink"
)

func main() {
	var satellite string
	var passphrase string
	var apikey string
	apikey = "13Yqcux9BQu4C1DyHjUbXispLV3qcTRws2NrGjBzi8MWZ4zLVBkZES3FPRD88y7ercGKKCDi7ud4aMEd2szmjL8HDYXXxEmXJs97CvQ"
	satellite = "12Wz6wJihX8yrnYht21kokZNiorNcLY5i5ai61sTLBR7qEhNqbi@127.0.0.1:10000"
	passphrase = "testpass"

	err := UploadAndDownloadData(context.Background(), satellite, apikey, passphrase,
		"my-first-bucket", "foo/bar/baz", []byte("one fish two fish red fish blue fish"))
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("success!")
}

// UploadAndDownloadData uploads the specified data to the specified key in the
// specified bucket, using the specified Satellite, API key, and passphrase.
func UploadAndDownloadData(ctx context.Context,
	satelliteAddress, apiKey, passphrase, bucketName, uploadKey string,
	dataToUpload []byte) error {

	// Request access grant to the satellite with the API key and passphrase.
	myConfig := uplink.Config{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			return &jsConn{
				ip:   "127.0.0.1",
				port: "10000",
			}, nil
		},
	}
	access, err := myConfig.RequestAccessWithPassphrase(ctx, satelliteAddress, apiKey, passphrase)
	if err != nil {
		return fmt.Errorf("could not request access grant: %v", err)
	}
	/*
		access, err := uplink.RequestAccessWithPassphrase(ctx, satelliteAddress, apiKey, passphrase)
		if err != nil {
			return fmt.Errorf("could not request access grant: %v", err)
		}
	*/

	// Open up the Project we will be working with.
	project, err := myConfig.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}

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

	// Initiate a download of the same object again
	download, err := project.DownloadObject(ctx, bucketName, uploadKey, nil)
	if err != nil {
		return fmt.Errorf("could not open object: %v", err)
	}
	defer download.Close()

	// Read everything from the download stream
	receivedContents, err := ioutil.ReadAll(download)
	if err != nil {
		return fmt.Errorf("could not read data: %v", err)
	}

	// Check that the downloaded data is the same as the uploaded data.
	if !bytes.Equal(receivedContents, dataToUpload) {
		return fmt.Errorf("got different object back: %q != %q", dataToUpload, receivedContents)
	}

	return nil
}

type jsConn struct {
	ip   string
	port string
}

var uint8Array = js.Global().Get("Uint8Array")

func (c *jsConn) Read(b []byte) (n int, err error) {
	reading := make(chan struct{})
	fmt.Println("read")
	var retVal js.Value
	js.Global().Call("socketRead", len(b), js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		retVal = args[0]
		close(reading)
		return nil
	}))
	<-reading
	fmt.Println("received bytes")
	js.CopyBytesToGo(b, retVal)
	fmt.Println(retVal.Length())
	return retVal.Length(), nil
}
func (c *jsConn) Write(b []byte) (n int, err error) {
	buf := uint8Array.New(len(b))
	js.CopyBytesToJS(buf, b)

	writing := make(chan struct{})
	js.Global().Call("socketWrite", c.ip, c.port, buf, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(writing)
		return nil
	}))
	<-writing
	return len(b), nil
}
func (c *jsConn) Close() error {
	fmt.Println("close")
	return nil
}
func (c *jsConn) LocalAddr() net.Addr {
	fmt.Println("local addr")
	return &addr{}
}
func (c *jsConn) RemoteAddr() net.Addr {
	fmt.Println("remote addr")
	return &addr{}
}
func (c *jsConn) SetDeadline(t time.Time) error {
	fmt.Println("set deadline")
	return nil
}
func (c *jsConn) SetReadDeadline(t time.Time) error {
	fmt.Println("set read deadline")
	return nil
}
func (c *jsConn) SetWriteDeadline(t time.Time) error {
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
