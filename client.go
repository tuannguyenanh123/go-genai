package main

import "github.com/google/generative-ai-go/genai"

func NewModel(client *genai.Client, model string) *genai.GenerativeModel {
	genaiModel := client.GenerativeModel(model)
	genaiModel.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	return genaiModel
}
