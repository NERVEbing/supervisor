package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v6/plumbing/object"
)

func detectLanguage(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}

	languageMap := map[string]string{
		".go":    "Go",
		".py":    "Python",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".jsx":   "JavaScript",
		".tsx":   "TypeScript",
		".java":  "Java",
		".c":     "C",
		".cpp":   "C++",
		".cc":    "C++",
		".h":     "C",
		".hpp":   "C++",
		".rs":    "Rust",
		".rb":    "Ruby",
		".php":   "PHP",
		".swift": "Swift",
		".kt":    "Kotlin",
		".scala": "Scala",
		".sh":    "Shell",
		".bash":  "Shell",
		".zsh":   "Shell",
		".sql":   "SQL",
		".md":    "Markdown",
		".json":  "JSON",
		".yaml":  "YAML",
		".yml":   "YAML",
		".xml":   "XML",
		".html":  "HTML",
		".css":   "CSS",
		".scss":  "SCSS",
		".sass":  "Sass",
		".proto": "Protocol Buffers",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}
	return ""
}

func isGeneratedFile(path string) bool {
	patterns := []string{
		".pb.go",
		".pb.gw.go",
		"_generated.go",
		".gen.go",
		"vendor/",
		"node_modules/",
		"dist/",
		"build/",
		".min.js",
		".min.css",
	}

	for _, pattern := range patterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}

func isTestFile(path string) bool {
	patterns := []string{
		"_test.go",
		"_test.py",
		".test.js",
		".test.ts",
		".spec.js",
		".spec.ts",
		"test/",
		"tests/",
		"__tests__/",
	}

	for _, pattern := range patterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}

func isConfigFile(path string) bool {
	configFiles := []string{
		"Makefile",
		"Dockerfile",
		".gitignore",
		".dockerignore",
		"go.mod",
		"go.sum",
		"package.json",
		"package-lock.json",
		"yarn.lock",
		"Cargo.toml",
		"Cargo.lock",
		"pom.xml",
		"build.gradle",
		"CMakeLists.txt",
		".golangci.yml",
		".golangci.yaml",
	}

	baseName := filepath.Base(path)
	for _, cf := range configFiles {
		if baseName == cf {
			return true
		}
	}

	if strings.HasPrefix(path, ".github/") ||
		strings.HasPrefix(path, ".vscode/") ||
		strings.HasPrefix(path, ".idea/") ||
		strings.Contains(path, "/config/") {
		return true
	}

	ext := filepath.Ext(path)
	configExts := []string{".yaml", ".yml", ".toml", ".ini", ".conf", ".config"}
	for _, ce := range configExts {
		if ext == ce {
			return true
		}
	}

	return false
}

func isBinaryFile(change *object.Change, fromTree, toTree *object.Tree) (bool, error) {
	path := getChangePath(change)
	if isBinaryExtension(path) {
		return true, nil
	}

	if change.To.Name != "" {
		file, err := toTree.File(change.To.Name)
		if err != nil {
			return false, fmt.Errorf("read file %s: %w", change.To.Name, err)
		}
		isBin, err := file.IsBinary()
		if err != nil {
			return false, fmt.Errorf("detect binary for %s: %w", change.To.Name, err)
		}
		return isBin, nil
	}

	if change.From.Name != "" && change.To.Name == "" {
		file, err := fromTree.File(change.From.Name)
		if err != nil {
			return false, fmt.Errorf("read file %s: %w", change.From.Name, err)
		}
		isBin, err := file.IsBinary()
		if err != nil {
			return false, fmt.Errorf("detect binary for %s: %w", change.From.Name, err)
		}
		return isBin, nil
	}

	return false, nil
}

func isBinaryExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := []string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp", ".ico", ".svg",
		".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar",
		".exe", ".dll", ".so", ".dylib", ".a", ".o",
		".wasm", ".class", ".pyc", ".pyo",
		".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv",
		".wav", ".ogg", ".m4a",
		".ttf", ".otf", ".woff", ".woff2", ".eot",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".db", ".sqlite", ".sqlite3",
		".bin", ".dat", ".data",
	}

	for _, be := range binaryExts {
		if ext == be {
			return true
		}
	}

	return false
}

func getChangePath(change *object.Change) string {
	if change.To.Name != "" {
		return change.To.Name
	}
	return change.From.Name
}
