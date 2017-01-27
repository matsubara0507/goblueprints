package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

func (b *BigHuge) Synonyms(term string) ([]string, error) {
	var syns []string
	response, err := http.Get(fmt.Sprint(
		"http://words.bighugelabs.com/api/2/", b.APIKey, "/", term, "/json"))
	if err != nil {
		return syns, fmt.Errorf("bighuge: %q の類似語検索に失敗しました: %v", term, err)
	}
	var data synonyms
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return syns, err
	}
	syns = append(syns, data.Noun.Syn...)
	syns = append(syns, data.Verb.Syn...)
	return syns, nil
}
