package manifest

import (
	"encoding/json"

	"github.com/imdario/mergo"
	_ "github.com/piprate/json-gold/ld"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

var Context = map[string]interface{}{
	"@context": map[string]interface{}{
		"mm":           "https://github.com/SUNET/metadata-manifests",
		"dct":          "http://purl.org/dc/terms/",
		"identifier":   map[string]interface{}{"@id": "dct:identifier"},
		"manifest":     map[string]interface{}{"@id": "mm:manifest", "@type": "@id"},
		"schema":       map[string]interface{}{"@id": "mm:schema", "@type": "@id"},
		"publisher":    map[string]interface{}{"@id": "dct:publisher"},
		"creator":      map[string]interface{}{"@id": "dct:creator"},
		"rightsHolder": map[string]interface{}{"@id": "dct:rightsHolder"},
	},
}

// ManifestEntry example
type ManifestEntry struct {
	ID     string `json:"@id" example:"@base/.metadata/dc.xml"`
	Schema string `json:"mm:schema" example:"http://purl.org/dc/terms/"`
}

// Manifest example
type ManifestInfo struct {
	Context      map[string]interface{} `json:"@context" swaggerignore:"true"`
	Identifier   string                 `json:"mm:identifier" example:"https://s3.example.com/mybucket"`
	Manifest     []ManifestEntry        `json:"mm:manifest"`
	Publisher    string                 `json:"mm:publisher" example:"example.net"`
	Creator      string                 `json:"mm:creator" example:"joe@example.net"`
	RightsHolder string                 `json:"mm:rightsHolder" example:"example.com"`
} //@name Manifest

func (m *ManifestInfo) AddContext() {
	m.Context = Context
}

func (to *ManifestInfo) MergeManifests(froms ...ManifestInfo) error {
	var err error

	for _, from := range froms {
		err = mergo.Merge(to, &from)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewManifest() (*ManifestInfo, error) {
	m := new(ManifestInfo)
	m.AddContext()
	return m, nil
}

func NewSimpleManifest(creator string, publisher string) (*ManifestInfo, error) {
	m, err := NewManifest()
	if err != nil {
		return nil, err
	}
	m.Creator = creator
	m.Publisher = publisher
	return m, nil
}

func (m *ManifestInfo) ToString() string {
	data, err := json.Marshal(m)
	if err != nil {
		Log.Fatal(err)
	}
	return string(data)
}
