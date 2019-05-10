package mod

type ModReader interface {
	ReadMeta()
}

type ModWriter interface {
	WriteMeta(mod *Mod) error
}

type ModGenerater interface {
	GenerateMeta()
}

type ModReadWriterCreater interface {
	ModReader
	ModWriter
}
