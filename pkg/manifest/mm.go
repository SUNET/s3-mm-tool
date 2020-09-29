package manifest

import (
	"encoding/json"

	"github.com/imdario/mergo"
	_ "github.com/piprate/json-gold/ld"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

var Context = Manifest{
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

type Manifest map[string]interface{}

func (m *Manifest) AddContext() error {
	return m.MergeManifests(Context)
}

func (to *Manifest) MergeManifests(froms ...Manifest) error {
	var err error

	for _, from := range froms {
		err = mergo.Merge(to, &from)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewManifest() (*Manifest, error) {
	m := new(Manifest)
	err := m.AddContext()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manifest) SetString(key string, value string) {
	(*m)[key] = value
}

func (m *Manifest) Get(key string) interface{} {
	return (*m)[key]
}

func NewSimpleManifest(creator string, publisher string) (*Manifest, error) {
	m, err := NewManifest()
	if err != nil {
		return nil, err
	}
	m.SetString("creator", creator)
	m.SetString("publisher", publisher)
	return m, nil
}

func (m *Manifest) ToString() string {
	data, err := json.Marshal(m)
	if err != nil {
		Log.Fatal(err)
	}
	return string(data)
}
