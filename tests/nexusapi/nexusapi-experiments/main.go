package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	uuid "github.com/satori/go.uuid"
	flag "github.com/spf13/pflag"
)

var (
	download string
)

func init() {
	flag.StringVarP(&download, "download", "d", "", "used for NexusMod downloads")
}

func main() {
	flag.Parse()
	binaryPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	logPath := ""

	if binaryPath[len(binaryPath)-4:] == ".exe" {
		logPath = binaryPath[:len(binaryPath)-3] + "log"
	} else {
		logPath = binaryPath + ".log"
	}
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("starting!")

	if download != "" {
		log.Println("We have something!")
		log.Println(download)
		os.Exit(0)
	}

	fid := "uuid.txt"
	fkey := "apikey.ini"
	err = writeConfig(fid)
	if err != nil {
		panic(err)
	}

	u, err := readConfig(fid)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(u)

	apikey, err := ioutil.ReadFile(fkey)
	if err != nil {
		log.Println("no apikey file", err)
		ws(u)
	}

	fmt.Println(string(apikey))

	self, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	registerNxm(self)
}

func writeConfig(file string) error {
	uu, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u := uu.Bytes()

	// write file only if do not exist
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil
	}
	defer f.Close()
	_, err = f.Write(u)
	if err != nil {
		return err
	}
	return nil
}

func readConfig(file string) (string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	u, err := uuid.FromBytes(f)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// ------------- Websocket

func ws(windpeakID string) string {

	jRegister := fmt.Sprintf("{ \"id\": \"%s\", \"appid\": \"Windpeak\" }", windpeakID)
	// TODO: talk to Nexus, register Windpeak as app there, than add this to the end of url:
	// &application=Windpeak
	lRegister := fmt.Sprintf("https://www.nexusmods.com/sso?id=%s", windpeakID)

	apikey := ""

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	u := url.URL{Scheme: "wss", Host: "sso.nexusmods.com"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			apikey = string(message)
			interrupt <- os.Interrupt
			return
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(jRegister))
	log.Println("sent: ", jRegister)
	browser.OpenURL(lRegister)

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			log.Println("returnded after 'done'")
			return apikey
		case <-ticker.C:
			// send pings every 30s
			log.Println("sending ping")
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("error sending ping:", err)
				return ""
			}
			//err := c.WriteMessage(websocket.TextMessage, []byte("hi! this is my time:"+t.String()))
			//if err != nil {
			//	log.Println("write:", err)
			//	return
			//}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return apikey
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			log.Println("aaa return here")
			return apikey
		}
	}
	return apikey
}
