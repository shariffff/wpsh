package prompt

import "github.com/AlecAivazis/survey/v2"

// PromptEditorSelection prompts the user to select their preferred editor
func PromptEditorSelection() (string, error) {
	editors := []string{
		"nano",
		"vim",
		"code",
		"Other",
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select your preferred editor:",
		Options: editors,
		Help:    "This will be saved for future use",
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", err
	}

	if selected == "Other" {
		var customEditor string
		inputPrompt := &survey.Input{
			Message: "Enter your editor command:",
			Help:    "e.g., subl, emacs, micro",
		}
		if err := survey.AskOne(inputPrompt, &customEditor, survey.WithValidator(survey.Required)); err != nil {
			return "", err
		}
		return customEditor, nil
	}

	return selected, nil
}
