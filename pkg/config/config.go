package config

import (
	"os"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

type Configurer interface {
	GetSynonymCategory(category string) string
}

type Synonym struct {
	Categories map[string]string
}

func New(fileName string) (*Synonym, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't open file: %s", fileName)
	}
	defer file.Close()
	var synonym Synonym
	err = yaml.NewDecoder(file).Decode(&synonym)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't decode yaml")
	}
	return &synonym, nil
}

func (s *Synonym) GetSynonymCategory(category string) string {
	return s.Categories[category]
}
