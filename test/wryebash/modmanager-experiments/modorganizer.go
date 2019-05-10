package modmanager

import (
	"fmt"
	"regexp"

	ini "gopkg.in/ini.v1"
)

// Here are functions related to reading Mod Organizer 2 files.

type ModOrganizer2 struct {
}

// ReadModMeta reads original MO meta files.
// Params:
// metaFile - path to file from downloads directory *.meta
// iniFile  - path to file from mods directory meta.ini
func (m *Mod) ReadModMeta(metaFile, iniFile string) error {
	cfg, err := ini.Load(metaFile)
	if err != nil {
		return err
	}
	d := cfg.Section("General")
	m.ID, err = d.Key("modID").Int64()
	m.FileID, err = d.Key("fileID").Int64()
	m.Url = d.Key("url").String()
	m.Description = d.Key("description").String()
	m.Name = d.Key("modName").String()
	m.Repository = d.Key("repository").In("Nexus", []string{"Nexus"})

	cfg2, _ := ini.Load(iniFile)
	if err != nil {
		return err
	}
	d2 := cfg2.Section("General")
	m.Version = d2.Key("version").String()
	m.FileName = d2.Key("installationFile").String()

	m.getMD5fromUrl

	return nil
}

//parseUrl
func (m *Mod) getMD5fromUrl() error {
	re, err := regexp.Compile("(?:md5=)(.*?)&")
	if err != nil {
		return err
	}
	s := re.FindStringSubmatch(m.Url)
	if len(s) >= 2 {
		m.MD5 = s[1]
	} else {
		return fmt.Errorf("no md5 signature in url: %s", m.Url)
	}

	return nil
}
