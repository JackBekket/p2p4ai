package main

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
)




func StartRelay()  {
	
	/*
	relay1, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9876"))
	if err != nil {
		log.Printf("Failed to create relay1: %v", err)
		return
	}
	*/
	sourceMultiAddrTCP, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	sourceMultiAddrUDP, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/udp/4001/quic-v1")

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
		return
	}
	_, err = dht.New(context.Background(), host, dht.Mode(dht.ModeServer))
	if err != nil {
		log.Printf("Failed to create DHT: %v", err)
		return
	}
	
	_, err = autonat.New(host)
	if err != nil {
		log.Printf("Failed to create AutoNAT: %v", err)
		return
	}
	
	log.Println("Relay started")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", "0.0.0.0", 4001, host.ID().String())
	
	select {}


}