package urlshort

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dest, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(yml []byte) ([]pathUrl, error) {
	var results []pathUrl
	err := yaml.Unmarshal(yml, &results)
	return results, err
}

func buildMap(yml []pathUrl) map[string]string {
	results := make(map[string]string)
	for _, element := range yml {
		results[element.Path] = element.Url
	}
	return results
}

func DBHandler(fallback http.Handler) (http.HandlerFunc, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dataSouece := fmt.Sprintf("postgres://%s:%s@localhost:5430/test_db", dbUser, dbPassword)
	db, err := sql.Open("pgx", dataSouece)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	pathMap, err := selectRows(db)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathMap, fallback), nil
}

func selectRows(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT * FROM path_url")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]string)
	for rows.Next() {
		var path string
		var url string
		if err := rows.Scan(&path, &url); err != nil {
			return nil, err
		}
		fmt.Printf("Selected path %s to %s", path, url)
		results[path] = url
	}
	return results, nil
}
