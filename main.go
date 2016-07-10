package main

import (
	"fmt"
	"flag"
	"log"
	"bufio"
	"os"
	"io"
	"github.com/sorcix/irc"
)

// TODO: connect to specified channels

var (
	name = flag.String("name", "human", "nickname")
	server = flag.String("server", "irc.freenode.net:6667", "server:port")
	channels = flag.String("chan", "#test #test1", "channels")
)

func init() {
	flag.Parse()
	fmt.Printf("name: %s\n", *name)
	fmt.Printf("server: %s\n", *server)
	fmt.Printf("channels: %s\n", *channels)
}

func defaultPrefix() *irc.Prefix {
	return &irc.Prefix{
		Name: *name,
		User: *name,
	}
}

func readMessages(conn *irc.Conn) {
	for {
		m, err := conn.Decode()
		if err != nil {
			log.Printf("could not decode %s", err.Error())
			if err == io.EOF {
				break
			}
		}
		fmt.Println(m)
	}
}

func setupConnection(conn *irc.Conn) error {
	messages := []*irc.Message{}
	messages = append(messages, &irc.Message{
		Command: irc.NICK,
		Params: []string{*name},
	})
	messages = append(messages, &irc.Message{
		Command: irc.USER,
		Params: []string{*name, "0", "*"},
		Trailing: *name,
	})
	for _, m := range messages {
		err := conn.Encode(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	c, err := irc.Dial(*server)
	if err != nil {
		log.Fatal(err.Error())	
	}

	// TODO: scan from network instead when connection working
	scanner := bufio.NewScanner(os.Stdin)

	go readMessages(c)
	err = setupConnection(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	for scanner.Scan() {
		m := irc.ParseMessage(scanner.Text())
		if m == nil {
			log.Println("could not parse")	
			continue
		}

		m.Prefix = defaultPrefix()	
		
		fmt.Println("sending: ", *m)	
		
		err = c.Encode(m)
		if err != nil {
			log.Printf("could not encode: %s", err.Error())
		}
		
		fmt.Println("receiving: ", *m)	
	}	
}
