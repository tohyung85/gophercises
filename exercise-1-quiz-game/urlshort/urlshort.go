package urlshort

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if dest, inMap := pathsToUrls[r.RequestURI]; inMap {
			fmt.Printf("redirecting from %s to: %s", r.RequestURI, dest)
			http.Redirect(w, r, dest, http.StatusSeeOther)
			return
		}

		fallback.ServeHTTP(w, r)
	}

	return handler
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type URLStruct struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	urlMap := buildMap(parsedYaml)

	return MapHandler(urlMap, fallback), err
}

func parseYaml(yml []byte) ([]URLStruct, error) {
	var ymlSlice []URLStruct
	err := yaml.Unmarshal(yml, &ymlSlice)
	if err != nil {
		return nil, err
	}

	return ymlSlice, nil
}

func buildMap(parsedUrls []URLStruct) map[string]string {
	urlMap := make(map[string]string)
	for _, urlStruct := range parsedUrls {
		urlMap[urlStruct.Path] = urlStruct.Url
	}

	return urlMap
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}

	urlMap := buildMap(parsedJSON)

	return MapHandler(urlMap, fallback), err
}

func parseJSON(jsn []byte) ([]URLStruct, error) {
	var jsnSlice []URLStruct
	err := json.Unmarshal(jsn, &jsnSlice)
	if err != nil {
		return nil, err
	}

	return jsnSlice, nil
}

func DBHandler(db *sql.DB, fallback http.Handler) (http.HandlerFunc, error) {
	rows, err := db.Query(`SELECT * FROM "UrlMaps"`)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()

	structMap, err := parseRows(rows)
	urlMap := buildMap(structMap)

	return MapHandler(urlMap, fallback), err
}

func parseRows(rows *sql.Rows) ([]URLStruct, error) {
	var urlStructs []URLStruct

	for rows.Next() {
		var id int
		var path string
		var url string
		rows.Scan(&id, &path, &url)
		urlStructs = append(urlStructs, URLStruct{path, url})
	}

	return urlStructs, nil
}
