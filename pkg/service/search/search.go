package search

import (
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/char/asciifolding"
	"github.com/blevesearch/bleve/analysis/lang/pt"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/registry"

	"rinha/pkg/entity"
)

type Searcher interface {
	Search(query string) ([]string, error)
	Save(user entity.User) error
}

type searcher struct {
	index bleve.Index
}

func (s *searcher) Search(searchTerm string) ([]string, error) {
	query := bleve.NewMatchQuery(searchTerm)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		log.Println(err)
	}

	ids := make([]string, 0, len(searchResult.Hits))
	for _, result := range searchResult.Hits {
		ids = append(ids, result.Fields["id"].(string))
	}

	return ids, nil
}

func (s *searcher) Save(user entity.User) error {
	return s.index.Index(user.ID, user)
}

const AnalyzerName = "pt-br"

func AnalyzerConstructor(_ map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed(unicode.Name)
	if err != nil {
		return nil, err
	}

	toLowerFilter, err := cache.TokenFilterNamed(lowercase.Name)
	if err != nil {
		return nil, err
	}

	stopPtFilter, err := cache.TokenFilterNamed(pt.StopName)
	if err != nil {
		return nil, err
	}

	stemmerPtFilter, err := cache.TokenFilterNamed(pt.LightStemmerName)
	if err != nil {
		return nil, err
	}

	asciiFilter, err := cache.CharFilterNamed(asciifolding.Name)
	if err != nil {
		return nil, err
	}

	rv := analysis.Analyzer{
		Tokenizer: tokenizer,
		TokenFilters: []analysis.TokenFilter{
			toLowerFilter,
			stopPtFilter,
			stemmerPtFilter,
		},
		CharFilters: []analysis.CharFilter{
			asciiFilter,
		},
	}

	return &rv, nil
}

func NewSearcher() (Searcher, error) {
	registry.RegisterAnalyzer(AnalyzerName, AnalyzerConstructor)
	textFiledMapping := bleve.NewTextFieldMapping()
	textFiledMapping.Analyzer = AnalyzerName

	userMapping := bleve.NewDocumentMapping()
	userMapping.AddFieldMappingsAt("apelido", textFiledMapping)
	userMapping.AddFieldMappingsAt("nome", textFiledMapping)
	userMapping.AddFieldMappingsAt("stack", textFiledMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("phrase", userMapping)
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = AnalyzerName

	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	return &searcher{index: index}, nil
}
