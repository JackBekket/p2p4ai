package main

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
)




func StartRelay(tcp string, udp string) (host.Host,error){
	
	/*
	relay1, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9876"))
	if err != nil {
		log.Printf("Failed to create relay1: %v", err)
		return
	}
	*/
	//sourceMultiAddrTCP, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	//sourceMultiAddrUDP, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/udp/4001/quic-v1")
	sourceMultiAddrTCP, _ := multiaddr.NewMultiaddr(tcp)
	sourceMultiAddrUDP, _ := multiaddr.NewMultiaddr(udp)



	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host,err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddrTCP, sourceMultiAddrUDP),

		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		libp2p.EnableHolePunching(),
		libp2p.EnableNATService(),

		libp2p.EnableRelayService(),
	)
	if err != nil {
		panic(err)
	}
	
	_, err = relay.New(host, relay.WithInfiniteLimits())
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
		return nil,err
	}
	_, err = dht.New(context.Background(), host, dht.Mode(dht.ModeServer))
	if err != nil {
		log.Printf("Failed to create DHT: %v", err)
		return nil,err
	}
	
	_, err = autonat.New(host)
	if err != nil {
		log.Printf("Failed to create AutoNAT: %v", err)
		return nil,err
	}
	
	log.Println("Relay started")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", "0.0.0.0", 4001, host.ID().String())
	return host,nil
	
	//select {}


}


/*
func RunRelayHost() (host.Host){

	// Now, normally you do not just want a simple host, you want
	// that is fully configured to best support your p2p application.
	// Let's create a second host setting some more options.

	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	var idht *dht.IpfsDHT

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}
	h2, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/9000",         // regular tcp connections
			"/ip4/0.0.0.0/udp/9000/quic-v1", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(context.Background(), h,dht.Mode(dht.ModeServer))
			return idht, err
		}),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),

		libp2p.EnableHolePunching(),

		libp2p.EnableRelayService(),
	)
	if err != nil {
		panic(err)
	}
	defer h2.Close()

	// The last step to get fully up and running would be to connect to
	// bootstrap peers (or any other peers). We leave this commented as
	// this is an example and the peer will die as soon as it finishes, so
	// it is unnecessary to put strain on the network.

	
		// This connects to public bootstrappers
		for _, addr := range dht.DefaultBootstrapPeers {
			pi, _ := peer.AddrInfoFromP2pAddr(addr)
			// We ignore errors as some bootstrap peers may be down
			// and that is fine.
			h2.Connect(ctx, *pi)
		}
	
	log.Printf("Your relay hosts ID is %s\n", h2.ID())
	log.Println("Relay host address: ", h2.Addrs())
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", "0.0.0.0", 9000, h2.ID().String())

	return h2
}
*/