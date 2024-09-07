package transaction

type Data struct {
	Category       string       `json:"category"`
	IdentityNumber string       `json:"identityNumber"`
	Files          []*File      `json:"files"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Identifiers    []Identifier `json:"identifiers"`
}

type File struct {
	FileID     int    `json:"id_file"`
	Name       string `json:"name"`
	FileEncode string `json:"file_encode"`
	NameAws    string `json:"name_aws"`
}

type Identifier struct {
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
