//go:build !tinygo

package ledgrid

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
type NetGridClient struct {
	conn      net.Conn
	rpcClient *rpc.Client
	sendWatch *Stopwatch
}

func NewNetGridClient(host string, network string, port, rpcPort uint) GridClient {
	var hostPortData, hostPortRPC string
	var err error

	p := &NetGridClient{}
	hostPortData = fmt.Sprintf("%s:%d", host, port)
	hostPortRPC = fmt.Sprintf("%s:%d", host, rpcPort)

	p.conn, err = net.Dial(network, hostPortData)
	if err != nil {
		log.Fatalf("Error in Dial(dataPort): %v", err)
	}

	p.rpcClient, err = rpc.DialHTTP("tcp", hostPortRPC)
	if err != nil {
		log.Fatalf("Error in Dial(rpcPort): %v", err)
	}
	p.sendWatch = NewStopwatch()

	return p
}

// Sendet die Bilddaten in der LedGrid-Struktur zum Controller.
func (p *NetGridClient) Send(buffer []byte) {
	var err error

	p.sendWatch.Start()
	_, err = p.conn.Write(buffer)
	if err != nil {
		log.Fatal(err)
	}
	p.sendWatch.Stop()
}

// Die folgenden Methoden verpacken die entsprechenden RPC-Calls zum
// Grid-Server.
func (p *NetGridClient) NumLeds() int {
	var reply NumLedsArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCNumLeds", 0, &reply)
	if err != nil {
		log.Fatal("NumLeds error:", err)
	}
	return int(reply)
}

func (p *NetGridClient) Gamma() (r, g, b float64) {
	var reply GammaArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCGamma", 0, &reply)
	if err != nil {
		log.Fatal("Gamma error:", err)
	}
	return reply.RedVal, reply.GreenVal, reply.BlueVal
}

func (p *NetGridClient) SetGamma(r, g, b float64) {
	var reply int
	var err error

	err = p.rpcClient.Call("GridServer.RPCSetGamma", GammaArg{r, g, b}, &reply)
	if err != nil {
		log.Fatal("SetGamma error:", err)
	}
}

func (p *NetGridClient) ModuleConfig() conf.ModuleConfig {
	var reply ModuleConfigArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCModuleConfig", 0, &reply)
	if err != nil {
		log.Fatal("ModuleConfig error:", err)
	}
	return reply.ModuleConfig
}

func (p *NetGridClient) Watch() *Stopwatch {
	return p.sendWatch
}

// Schliesst die Verbindung zum Controller.
func (p *NetGridClient) Close() {
	p.conn.Close()
}
