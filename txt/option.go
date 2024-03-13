package txt

import "strings"

type Option func(*Analyzer)

func WithStopWords(words []string) Option {
	return func(ana *Analyzer) {
		ana.SetStopWords(words)
	}
}

func WithDefaultStopWords(ana *Analyzer) {
	ana.SetStopWords(stopWords)
}

func ToLower(ana *Analyzer) {
	ana.AddNormalizer(strings.ToLower)
}

func WithoutPunct(ana *Analyzer) {
	ana.AddNormalizer(StripPunct)
}

func WithStemmer(ana *Analyzer) {
	ana.AddNormalizer(Stem)
}

func WithNormalizers(n ...Normalizer) Option {
	return func(ana *Analyzer) {
		ana.AddNormalizer(n...)
	}
}
