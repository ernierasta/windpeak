package nexusapi

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
)

const apikeyFileName = "apikey.txt"

// tests in here are integration tests, not unit tests. Do NOT run them periodically!

func getApikey(t *testing.T) string {
	k, err := ioutil.ReadFile(apikeyFileName)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(k))
}

func TestRegisterTest(t *testing.T) {

	if _, err := os.Stat(apikeyFileName); err == nil {
		t.Skip("apikey.txt exists, skipping")
	}
	t.Log("check your browser! it waits for confirmation")
	uu, err := uuid.NewV4()
	n := New("My app test", "0.1", uu.String(), "")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	apikey, err := n.RegisterTest(c)
	if err != nil {
		t.Fatal(err)
	}

	ioutil.WriteFile(apikeyFileName, []byte(apikey), 0600)
	t.Log(apikey)
}

// you have to provide actual data from nxm link to make this test work
func TestGetDownloadLink(t *testing.T) {
	//"nxm://Oblivion/mods/48577/files/1000021855?key=1&expires=1556829151&user_id=1"
	m := &Mod{
		Game:    "Oblivion",
		ID:      48577,
		FileID:  1000021855,
		Key:     "1",
		Expires: 1556829151,
		UserID:  1,
	}
	apikey := getApikey(t)
	n := New("My app test", "0.1", "not used", apikey)
	s, err := n.GetModDownloadLink(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
	//https://cf-files.nexusmods.com/cdn/101/48577/NorthernUI-48577-1-2-1-1555628981.zip?md5=DmRyKCA4PcwlzolmW18nmQ&expires=1556673526&user_id=1
	// if there is error:
	//  json: cannot unmarshal object into Go value of type nexusapi.ModDownload
	// {"code":410,"message":"This link has expired - please visit the mod page again to get a new link"}
}

// you have to provide actual data from nxm link to make this test work
func TestGetModInfo(t *testing.T) {
	m := &Mod{
		Game:    "Oblivion",
		ID:      48577,
		FileID:  1000021855,
		Key:     "1",
		Expires: 1556829151,
		UserID:  1,
	}
	apikey := getApikey(t)
	n := New("My app test", "0.1", "not used", apikey)
	mi, err := n.GetModFileInfo(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", mi)
	//&{FileID:1000021855 Name:NorthernUI Version:1.2.1 CategoryID:1 CategoryName:MAIN IsPrimary:false Size:26359 FileName:NorthernUI-48577-1-2-1-1555628981.zip UploadedTimestamp:1555628981 UploadedTime:2019-04-18 23:09:41 +0000 +0000 ModVersion:1.2.1 ExternalVirusScanURL:https://www.virustotal.com/file/3a6c1acd9847a32c25680156d85b2f8c06278fe097b52568deaafc7f130e6811/analysis/1555629151/ Description:Main file. SizeKb:26359 ChangelogHTML:The "Potions Cooked" misc stat now increments properly when cooking potions using NorthernUI's enhanced Alchemy menu.
	//        Fixed a crash that could occur when using NorthernUI's enhanced Alchemy menu, if any ingredients in your inventory had fewer magic effects than you were capable of using (e.g. the Poisoned Apple, which has only one effect).
	//        The training menu now shows a message explaining why training is disabled, if you have outleveled the trainer or if you have trained the maximum number of times for your current level already.
	//        Improved compatibility with Dynamic Training Cost: We now hide the number of times you have trained for the current level if “unlimited training” is enabled.
	//        Improved compatibility with Dynamic Training Cost: We now properly show raw skill values in place of skill mastery levels (e.g. “Novice”) if you’ve enabled that setting.
	//        Improved compatibility with Dynamic Training Cost: We now properly display training adjust stats when the training menu is opened via a mouse click. (Getting that working with a gamepad will require more work; the DTC script is using a MenuQue-provide API.)
	//        Improved compatibility with Dynamic Training Cost: When the “training takes time” feature is enabled, the timer shown will no longer protrude out of the window; it will instead be cleanly positioned with the gold cost.
	//        Restored the extended console commands (“!xxn commandName”) patch. It was disabled during a previous version and development got hectic enough that I forgot to enable it again.
	//        Removed some unneeded debug logging code from a patch to the dialogue menu.}
	//
}

func TestGetModByMD5(t *testing.T) {

	apikey := getApikey(t)
	n := New("My app test", "0.1", "not used", apikey)
	md5i, err := n.GetModByMD5("99ac0362b889f8563505932028e14fbf", "skyrim")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", md5i)

}

func TestCreateHeaders(t *testing.T) {
	n := New("My app test", "0.1", "not used", "")
	h := n.createHeaders(http.Header{})
	t.Log(h)
}
