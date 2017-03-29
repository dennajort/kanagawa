package metadata

type File struct {
	Length uint64   `benc:"length"`
	Path   []string `benc:"path"`
}

type Info struct {
	PieceLength uint64 `benc:"piece length"`
	Pieces      []byte `benc:"pieces"`
	Name        string `benc:"name"`
	Length      uint64 `benc:"length"`
	Files       []File `benc:"files"`
}

type Metadata struct {
	Announce string                 `benc:"announce"`
	Info     map[string]interface{} `benc:"info"`
}
