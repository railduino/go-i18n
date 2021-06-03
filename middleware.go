package i18n

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type langWithWeight struct {
	Lang   string
	Weight float64
}

type langByWeight []langWithWeight

func Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if len(Locales) == 0 {
			h.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.URL.Path, Config.StaticPrefix) {
			h.ServeHTTP(w, r)
			return
		}

		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = r.Header.Get("Accept-Language")
		}
		accept := langAccept(lang)

		ctx := r.Context()
		ctx = context.WithValue(ctx, "language", accept.Lang)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func langAccept(header string) langWithWeight {
	var entry langWithWeight
	var candidates langByWeight

	for _, part := range strings.Split(header, ",") {
		code := strings.Split(strings.TrimSpace(part), ";")
		lang := strings.Replace(code[0], "-", "_", -1)

		if len(code) == 1 {
			entry = langWithWeight{Lang: lang, Weight: 1.0}
		} else {
			weight := strings.Split(code[1], "=")
			factor, err := strconv.ParseFloat(weight[1], 64)
			if err != nil {
				log.Printf("ERROR invalid language weight '%s' in '%s'", weight, header)
				continue
			}
			entry = langWithWeight{Lang: lang, Weight: factor}
		}

		// TODO if no locales exist, create the "most wanted" locale

		for _, locale := range Locales {
			if lang == locale.Lang {
				candidates = append(candidates, entry)
			}
		}
	}

	if len(candidates) == 0 {
		entry = langWithWeight{Lang: Locales[0].Lang, Weight: 0.0}
		candidates = append(candidates, entry)
	} else if len(candidates) > 1 {
		sort.Sort(langByWeight(candidates))
	}

	return candidates[0]
}

func (a langByWeight) Len() int {
	return len(a)
}

func (a langByWeight) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a langByWeight) Less(i, j int) bool {
	// N.B. descending order
	return a[j].Weight < a[i].Weight
}
