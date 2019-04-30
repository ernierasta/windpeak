// Package nexusapi implements NexusMods api in go. Nexusmods is used here as:
// Nexus, NexusMods, nexusmods.com. All means the same in context of this library.
//
// Official API documentation:
// SSO Registration description:
//
// Generate a random unique id (we suggest uuid v4)
//
// Create a websocket connection to wss://sso.nexusmods.com
//
// When the connection is established, send a JSON encoded message containing the id you just generated and the appid you got on registration. Example: { "id": "4c694264-1fdb-48c6-a5a0-8edd9e53c7a6", "appid": "your_fancy_app" }
//
// From now on, until the connection is closed, send a websocket ping once every 30 seconds.
//
// Have the user open https://www.nexusmods.com/sso?id=xyz (id being the random id you generated in step 1) in the default browser
// On the website users will be asked to log-in to Nexus Mods if they aren't already. Then they will be asked to authorize your application to use their account.
//
// Once the user confirms, a message will be sent to your websocket with the APIKEY (not encoded, just the plain key). This is the only non-error message you will ever receive from the server.
//
// Save away the key and close the connection.
//
// taken from: https://github.com/Nexus-Mods/node-nexus-api readme file.
//
//
// Swagger:
// https://app.swaggerhub.com/apis-docs/NexusMods/nexus-mods_public_api_params_in_form_data/1.0#/
//
//
// Example usage:
//
//   fmt.Println("check your browser! it waits for confirmation")
//	 uu, err := uuid.NewV4()
//	 n := New("My app test", "0.1", myuuid.String(), "")
//	 c := make(chan os.Signal, 1)
//	 signal.Notify(c, os.Interrupt)
//
//	 apikey, err := n.RegisterTest(c)
//	 if err != nil {
//		 fmt.Println(err)
//	 }
//   fmt.Println(apikey)
//
// For more examples look into test file in this dir.
package nexusapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
)

const (
	URLBASE = "https://api.nexusmods.com"
	APIVER  = "v1"
)

// Mod contains data neccessary to work with API.
type Mod struct {
	ID      int
	FileID  int
	Game    string
	Expires int
	Key     string
	UserID  int
}

// ModDownload struct represents download info. Only URI is used.
type ModDownload []struct {
	Name      string `json:"-"`
	ShortName string `json:"-"`
	URI       string `json:"URI"`
}

// ModInfo contains all data returned from mod info query.
type ModInfo struct {
	FileID               int       `json:"file_id"`
	Name                 string    `json:"name"`
	Version              string    `json:"version"`
	CategoryID           int       `json:"category_id"`
	CategoryName         string    `json:"category_name"`
	IsPrimary            bool      `json:"is_primary"`
	Size                 int       `json:"size"`
	FileName             string    `json:"file_name"`
	UploadedTimestamp    int       `json:"uploaded_timestamp"`
	UploadedTime         time.Time `json:"uploaded_time"`
	ModVersion           string    `json:"mod_version"`
	ExternalVirusScanURL string    `json:"external_virus_scan_url"`
	Description          string    `json:"description"`
	SizeKb               int       `json:"size_kb"`
	ChangelogHTML        string    `json:"changelog_html"`
}

// Nexus represents NexusMods connection.
type Nexus struct {
	appName, appVersion string
	appUUID, apikey     string
	chanDone            chan struct{}
	chanInterrupt       chan os.Signal
	reterr              error
}

// New returns Nexus type.
// AppUUID should be generated every time you want to make new registration. If reused
// NexusMods will start terminating websocket connection, error: websocket: close 1006 (abnormal closure): unexpected EOF
// Store apikey somewhere safe, reuse.
func New(appName, appVersion, appUUID, apikey string) *Nexus {
	done := make(chan struct{})
	return &Nexus{
		appName:    appName,
		appVersion: appVersion,
		appUUID:    appUUID,
		apikey:     apikey,
		chanDone:   done,
	}
}

// Register registers NexusMods SSO via websocket,
// it returns apikey or error.
// It will open NexusMods page in the browser,
// waiting for user to confirm.
// At the end websocet is disconected and never used again.
// exit channel is used to correctly end connection with server,
// use it like that:
//   interrupt := make(chan os.Signal, 1)
//   signal.Notify(interrupt, os.Interrupt)
// You can always send signal to channel when needed:
//   interrupt <- os.Interrupt
//
// If apikey param is not empty string, it will not perform any action,
// but return apikey back.
func (n *Nexus) Register(exit chan os.Signal) (string, error) {
	return n.register(true, exit)
}

// RegisterTest is the same as Register, but do not send application name,
// to NexusMods. This allow to register as Vortex and works without contacting
// NexusMods. Do not overuse this! Contact NexusMods to get your app registered.
func (n *Nexus) RegisterTest(exit chan os.Signal) (string, error) {
	return n.register(false, exit)
}

// GetModDownloadLink returns mod file download link.
func (n *Nexus) GetModDownloadLink(mod *Mod) (string, error) {
	jsonFile := "/download_link.json"

	body, err := n.getResponseBody(jsonFile, mod)
	if err != nil {
		return "", err
	}

	mdp := &ModDownload{}
	err = json.Unmarshal(body, mdp)
	if err != nil {
		return "", err
	}

	md := *mdp
	return md[0].URI, nil
}

// GetModInfo returns mod metadata.
func (n *Nexus) GetModInfo(mod *Mod) (*ModInfo, error) {
	jsonFile := ".json"

	body, err := n.getResponseBody(jsonFile, mod)
	if err != nil {
		return &ModInfo{}, err
	}

	mi := &ModInfo{}
	err = json.Unmarshal(body, mi)
	if err != nil {
		return mi, err
	}

	return mi, nil
}

// createHeaders creates correct headers, example output:
//   User-Agent:
func (n *Nexus) createHeaders(h http.Header) http.Header {
	h.Add("content-type", "application/json")
	h.Add("apikey", n.apikey)
	h.Add("User-Agent", fmt.Sprintf("%s/%s (%s; %s) %s", n.appName, n.appVersion, runtime.GOOS, runtime.GOARCH, runtime.Version()))
	return h
}

func (n *Nexus) getResponseBody(jsonFile string, m *Mod) ([]byte, error) {
	// non-premium link type
	url := URLBASE + "/" + APIVER + "/games/" + m.Game + "/mods/" + strconv.Itoa(m.ID) +
		"/files/" + strconv.Itoa(m.FileID) + jsonFile + "?key=" + m.Key + "&expires=" +
		strconv.Itoa(m.Expires) + "&user_id=" + strconv.Itoa(m.UserID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header = n.createHeaders(req.Header)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println(string(body))
	return body, nil
}

// register actually do the work, for description look at Register()
func (n *Nexus) register(sendAppName bool, interrupt chan os.Signal) (string, error) {
	if n.apikey != "" {
		return n.apikey, nil
	}
	n.chanInterrupt = interrupt
	jsonRegister := fmt.Sprintf("{ \"id\": \"%s\", \"appid\": \"%s\" }", n.appUUID, n.appName)
	linkRegister := ""
	if sendAppName {
		linkRegister = fmt.Sprintf("https://www.nexusmods.com/sso?id=%s&application=%s", n.appUUID, n.appName)
	} else {
		linkRegister = fmt.Sprintf("https://www.nexusmods.com/sso?id=%s", n.appUUID)
	}

	// connect to websocket
	ws, _, err := websocket.DefaultDialer.Dial("wss://sso.nexusmods.com", nil)
	if err != nil {
		return "", fmt.Errorf("error while connecting: %v", err)
	}
	defer ws.Close()

	// start goroutine listening for apikey
	go n.read(ws)

	// send SSO register request
	ws.WriteMessage(websocket.TextMessage, []byte(jsonRegister))
	browser.OpenURL(linkRegister) // open browser to allow user confirm

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	// ping every 30 seconds, exit cleanly if needed
	for {
		select {
		case <-n.chanDone:
			return n.apikey, n.reterr
		case <-ticker.C:
			// send pings every 30s
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return "", fmt.Errorf("error sending ping: %v", err)
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return n.apikey, n.reterr
			}
			select {
			case <-n.chanDone:
			case <-time.After(time.Second):
			}
			return n.apikey, n.reterr
		}
	}

	return "", nil
}

// read waits for apikey to appear in websocket, writes it to struct
//
// note: writing directly to n.reterr and n.apikey
// is considered as bug in concurrency world,
// but in this case - only one reader can exist anyway.
func (n *Nexus) read(ws *websocket.Conn) {
	defer close(n.chanDone)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			n.reterr = err
			return
		}
		n.apikey = string(message)
		n.chanInterrupt <- os.Interrupt
		return
	}

}
