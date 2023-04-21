package ux

import (
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
)

func PromptSingle(question string, choices []string) (string, error) {
	sp := selection.New(question, choices)
	sp.PageSize = 4
	sp.Filter = nil

	choice, err := sp.RunPrompt()
	if err != nil {
		return "", err
	}

	// do something with the final choice
	return choice, nil
}
func Input(question, initialValue, placeholder string) (string, error) {
	input := textinput.New(question)
	input.InitialValue = initialValue
	input.Placeholder = placeholder

	val, err := input.RunPrompt()
	if err != nil {
		return "", err
	}

	// do something with the result
	return val, nil
}

func Confirm(question string) (bool, error) {
	input := confirmation.New(question, confirmation.Undecided)

	ready, err := input.RunPrompt()
	if err != nil {
		return false, err
	}

	// do something with the result
	return ready, nil
}
