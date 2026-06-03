package report

type Options struct {
	Inputs    []string
	Addr      string
	Output    string
	WriteHTML bool
	Version   string
}

type Glyph struct {
	Char     string `json:"char"`
	Code     string `json:"code"`
	CodeInt  int    `json:"codeInt"`
	Category string `json:"category"`
	Name     string `json:"name"`
}

type Font struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Glyphs []Glyph `json:"glyphs"`
}
