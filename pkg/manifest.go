package manifest
import (
	_ "github.com/piprate/json-gold/ld"
	"github.com/imdario/mergo"
)

var Context = map[string]interface{}{
	"@context": map[string]interface{}{
        "mm": "https://github.com/SUNET/metadata-manifests",
        "dct": "http://purl.org/dc/terms/",
        "identifier": map[string]interface{}{ "@id": "dct:identifier", },
        "manifest": map[string]interface{}{ "@id": "mm:manifest", "@type": "@id", },
        "schema": map[string]interface{}{ "@id": "mm:schema", "@type": "@id", },
        "publisher": map[string]interface{}{ "@id": "dct:publisher", },
        "creator": map[string]interface{}{ "@id": "dct:creator", },
        "rightsHolder": map[string]interface{}{ "@id": "dct:rightsHolder", },
    },
}


func NewManifest(creator string, publisher string) (*map[string]interface{}, error) {
	m := new(map[string]interface{})
	err := mergo.Merge(m, Context)
	if err != nil {
		return nil, err
	}
	return m, nil
}
