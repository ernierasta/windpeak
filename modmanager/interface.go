package modmanager

//ModManager interface allows to abstract mod manager (MO2, Wrye Bash,...).
type ModManager interface {
	SetMetaFiles(a ...string) error
	ReadModMeta() error
	GetProfile(name string) []string
}

//Creator interface represents methods needed while creating modpack
type Creator interface {
	WatchModInstall() error
}

//User interface needed while installing modpack
type User interface {
}
