import (
	"github.com/piprate/json-gold"
	"github.com/imdario/mergo"
	
)

Context := map[string]interface{}{
	"@context": {
        "mm": "https://github.com/SUNET/metadata-manifests",
        "dct": "http://purl.org/dc/terms/",
        "identifier": { "@id": "dct:identifier" },
        "manifest": { "@id": "mm:manifest", "@type": "@id" },
        "schema": { "@id": "mm:schema", "@type": "@id" },
        "publisher": { "@id": "dct:publisher" },
        "creator": { "@id": "dct:creator" },
        "rightsHolder": { "@id": "dct:rightsHolder" },
    }
} 


func NewManifest(creator string, publisher string) (map[string]interface{}, error) {
	m := new(map[string]manifest)
	err := mergo.Merge(m, Context)
	if err != nil {
		return nil, err
	}
	return m
}