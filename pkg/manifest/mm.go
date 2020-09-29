package manifest

import (
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

func AddContext(manifest *Manifest) error {
	return MergeManifests(manifest, &Context)
}

func MergeManifests(to *Manifest, froms ...*Manifest) error {
	var err error

	for from, _ := range froms {
		err = mergo.Merge(to, from)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewManifest() (*Manifest, error) {
	m := new(Manifest)
	err := AddContext(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func NewSimpleManifest(creator string, publisher string) (*Manifest, error) {
	m, err := NewManifest()
	if err != nil {
		return nil, err
	}
	(*m)["creator"] = creator
	(*m)["publisher"] = publisher
	return m, nil
}
