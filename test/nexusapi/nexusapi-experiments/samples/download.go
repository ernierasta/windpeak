package main

// documentation: https://app.swaggerhub.com/apis-docs/NexusMods/nexus-mods_public_api_params_in_form_data/1.0#/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DownloadMod []struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	URI       string `json:"URI"`
}

func main() {

	url := "https://api.nexusmods.com/v1/games/Oblivion/mods/49266/files/1000021647/download_link.json?key=1PTuAnCrEN8GVDoQBW0j5w&expires=1556573130&user_id=3531134"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("apikey", "my-keyhere")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	dm := &DownloadMod{}
	err := json.Unmarshal(body, dm)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", dm)
	// returns
	// &[{Name:Nexus Global Content Delivery Network
	//	ShortName:Nexus CDN
	//	URI:https://supporter-files.nexus-cdn.com/101/49266/UHD Fonts for Darnified UI-49266-1-0-1553899266.7z?md5=U4gMX0SpBD1-Vp9ABnVwEw&expires=1556416613&user_id=3531134}]
}
