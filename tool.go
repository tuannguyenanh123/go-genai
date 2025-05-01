package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

var fileWriteSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"filename": {
			Type:        genai.TypeString,
			Description: "The name of the file to write to. Don't include extension, it will be automatically added (.txt)",
		},
		"content": {
			Type:        genai.TypeString,
			Description: "The text content to write to the file.",
		},
	},
	Required: []string{"filename", "content"},
}

var FileTools = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "file_write",
			Description: "Write a text file to user local file system with specified name and context.",
			Parameters:  fileWriteSchema,
		},
	},
}

func WriteDesktop(filename string, content string) error {
	filename = filename + ".txt"
	home, _ := os.UserHomeDir()
	fullPath := filepath.Join(home, "Desktop", filename)

	formattedContent := strings.ReplaceAll(content, "\\n", "\n")

	err := os.WriteFile(fullPath, []byte(formattedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
