package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bregydoc/gtranslate"

	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
)

// translate translates text from base language to target language
func translate(text string, src, dst string) (string, error) {
	if strings.TrimSpace(text) == "" || src == dst {
		return text, nil
	}

	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: src,
			To:   dst,
		},
	)
	if err != nil {
		return "", err
	}
	return translated, nil
}

// readJsonFile get map[string]string from json file
func readJsonFile(path string) (map[string]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// writeJsonMap writes map[string]string to a JSON file with sorted keys.
func writeJsonMap(path string, m map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	return os.WriteFile(path, data, 0o644)
}

// listLanguages returns all language directories under root, excluding the base language.
func listLanguages(root, baseLang string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var langs []string
	for _, e := range entries {
		if e.IsDir() {
			name := e.Name()
			if name != baseLang {
				langs = append(langs, name)
			}
		}
	}

	return langs, nil
}

// walkBaseLanguageFiles walks all files under the base language directory.
func walkBaseLanguageFiles(root, baseLang string, fn func(baseFilePath string, rel string) error) error {
	baseRoot := filepath.Join(root, baseLang)
	return filepath.WalkDir(baseRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		rel, err := filepath.Rel(baseRoot, path)
		if err != nil {
			return err
		}
		return fn(path, rel)
	})
}

// translateFileToLang translates a single base file to target language file.
func translateFileToLang(baseFilePath, rel, baseLang, targetLang string) error {
	// read base language file
	baseMap, err := readJsonFile(baseFilePath)
	if err != nil {
		return fmt.Errorf("read base json failed: %w", err)
	}

	// compute target file path
	root := filepath.Dir(filepath.Dir(baseFilePath)) // .../translations/<baseLang>
	translationsRoot := filepath.Dir(root)           // .../translations
	targetFile := filepath.Join(translationsRoot, targetLang, rel)

	// read target language file if exists
	targetMap := map[string]string{}
	if b, err := os.ReadFile(targetFile); err == nil {
		if len(strings.TrimSpace(string(b))) > 0 {
			if err := json.Unmarshal(b, &targetMap); err != nil {
				return fmt.Errorf("unmarshal target json (%s) failed: %w", targetFile, err)
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read target json (%s) failed: %w", targetFile, err)
	}

	changed := false
	for k, v := range baseMap {
		tv, exists := targetMap[k]
		// exists and not empty
		if exists && strings.TrimSpace(tv) != "" {
			continue
		}

		// use value or key as source
		source := strings.TrimSpace(v)
		if source == "" {
			source = strings.TrimSpace(k)
			if source == "" {
				continue
			}
		}

		tr, err := translate(source, baseLang, targetLang)
		if err != nil || strings.TrimSpace(tr) == "" {
			tr = source
		}

		targetMap[k] = tr
		changed = true
	}

	if changed {
		if err := writeJsonMap(targetFile, targetMap); err != nil {
			return fmt.Errorf("write target json %s failed: %w", targetFile, err)
		}
	}
	return nil
}

// AutoTranslateAll translate all files by base language files to all existing language directories (excluding base).
func AutoTranslateAll(path, baseLang string) error {
	langs, err := listLanguages(path, baseLang)
	if err != nil {
		return err
	}

	// translate to each existing target language (if any)
	return walkBaseLanguageFiles(path, baseLang, func(baseFilePath string, rel string) error {
		for _, lang := range langs {
			if err := translateFileToLang(baseFilePath, rel, baseLang, lang); err != nil {
				return err
			}
		}
		return nil
	})
}

// AutoTranslateToLang translates all base language files to the specified target language.
func AutoTranslateToLang(root, baseLang, targetLang string) error {
	// translate each base file to target language
	return walkBaseLanguageFiles(root, baseLang, func(baseFilePath string, rel string) error {
		return translateFileToLang(baseFilePath, rel, baseLang, targetLang)
	})
}

// AddLanguageAndTranslate adds a new language directory and translates from base language.
func AddLanguageAndTranslate(root, baseLang, targetLang string) error {
	return AutoTranslateToLang(root, baseLang, targetLang)
}

func main() {
	var (
		path = flag.String("path", "", "translations root directory")
		lang = flag.String("lang", "", "extra language to create and translate (optional)")
	)
	flag.Parse()
	baseLang := string(i18n.DefaultLanguage)
	if strings.TrimSpace(*path) == "" {
		fmt.Println("error: -path is required")
		return
	}

	// translate all base language files to all existing target languages
	if err := AutoTranslateAll(*path, baseLang); err != nil {
		fmt.Println("error: ", err)
		return
	}

	// add new language and translate
	if strings.TrimSpace(*lang) != "" {
		if err := AddLanguageAndTranslate(*path, baseLang, *lang); err != nil {
			fmt.Println("error: ", err)
			return
		}
	}

	fmt.Println("success")
}
