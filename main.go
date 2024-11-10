package main

import (
	"bufio"
	"context"
	"flag"

	//"flag"
	//"flags"
	"fmt"
	"time"

	//"log"
	"os"
	"sync"

	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
	//"github.com/jackbekket/p2p4ai/flags"
)

/*
var (
  topicNameFlag = flag.String("topicName", "skynet", "name of topic to join")
)
*/

var logger = log.Logger("rendezvous")


//var topicNameFlag string


func main() {
	//flag.Parse()
	ctx := context.Background()


	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel("rendezvous", "info")
	help := flag.Bool("h", false, "Display Help")
	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	if *help {
		fmt.Println("This program demonstrates a simple p2p chat application using libp2p")
		fmt.Println()
		fmt.Println("Usage: Run './chat in two different terminals. Let them connect to the bootstrap nodes, announce themselves and connect to the peers")
		flag.PrintDefaults()
		return
	}

	/*
	_ = godotenv.Load()
	topicEnv := os.Getenv("TOPIC") 

	topicNameFlag := flag.String("topicName", topicEnv, "name of topic to join")
	*/

	//log.("topicName: ", topicNameFlag)
	logger.Infoln("topicName: ", config.RendezvousString)


		// libp2p.New constructs a new libp2p Host. Other options can be added
	// here.
	host, err := libp2p.New(libp2p.ListenAddrs([]multiaddr.Multiaddr(config.ListenAddresses)...))
	if err != nil {
		panic(err)
	}
	logger.Info("Host created. We are:", host.ID())
	logger.Info(host.Addrs())

	/*
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
	  panic(err)
	}
	*/
	go discoverPeers(ctx, host,config)
  
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
	  panic(err)
	}
	topic, err := ps.Join(*&config.RendezvousString)
	if err != nil {
	  panic(err)
	}

	// send test message
	//msg := 
	//SendMessageToTopic(ctx, topic, []byte("127.0.0.1:50051"))
	SendMessageToTopic(ctx, topic, []byte("this is a test. a chair is against the wall."))

	// Docker does not have a console
	//go streamConsoleTo(ctx, topic)
  
	sub, err := topic.Subscribe()
	if err != nil {
	  panic(err)
	}
	go handleIncomingMessages(ctx,sub)
	printMessagesFrom(ctx, sub)
  }

  

  func initDHT(ctx context.Context, h host.Host,cfg Config) *dht.IpfsDHT {

		// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	//ctx := context.Background()
	config := cfg
	bootstrapPeers := make([]peer.AddrInfo, len(config.BootstrapPeers))
	for i, addr := range config.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		bootstrapPeers[i] = *peerinfo
	}





	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h,dht.BootstrapPeers(bootstrapPeers...))
	if err != nil {
	  panic(err)
	}
	logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
	  panic(err)
	}
		// Wait a bit to let bootstrapping finish (really bootstrap should block until it's ready, but that isn't the case yet.)
	time.Sleep(1 * time.Second)
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
	  peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
	  wg.Add(1)
	  go func() {
		defer wg.Done()
		if err := h.Connect(ctx, *peerinfo); err != nil {
		  fmt.Println("Bootstrap warning:", err)
		}
	  }()
	}
	wg.Wait()
  
	return kademliaDHT
  }

  
  func discoverPeers(ctx context.Context, h host.Host, cfg Config) {
	// bootstrap DHT
	kademliaDHT := initDHT(ctx, h,cfg)

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	logger.Info("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, *&cfg.RendezvousString)
	logger.Debug("Successfully announced!")
  
	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
	  fmt.Println("Searching for peers...")
	  peerChan, err := routingDiscovery.FindPeers(ctx, *&cfg.RendezvousString)
	  if err != nil {
		panic(err)
	  }
	  for peer := range peerChan {
		if peer.ID == h.ID() {
		  continue // No self connection
		}
		err := h.Connect(ctx, peer)
		if err != nil {
		  fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
		} else {
		  fmt.Println("Connected to:", peer.ID)
		  anyConnected = true
		}
	  }
	}
	fmt.Println("Peer discovery complete")
  }

  

  func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
	  s, err := reader.ReadString('\n')
	  if err != nil {
		panic(err)
	  }
	  if err := topic.Publish(ctx, []byte(s)); err != nil {
		fmt.Println("### Publish error:", err)
	  }
	}
  }
  
  func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
	  m, err := sub.Next(ctx)
	  if err != nil {
		panic(err)
	  }
	  fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
	}
  }


  func SendMessageToTopic(ctx context.Context, topic *pubsub.Topic, message []byte) error {
	return topic.Publish(ctx, message)
  }

  func handleIncomingMessages(ctx context.Context, sub *pubsub.Subscription) {
	for {
	  m, err := sub.Next(ctx)
	  if err != nil {
		fmt.Println("Error receiving message:", err)
		continue
	  }
	  fmt.Println("Received message from:", m.ReceivedFrom, ":", string(m.Message.Data))
	  // Process the received message here, e.g., extract gRPC address:port info
	}
  }


  
  