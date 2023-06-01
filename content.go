package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	chrome "github.com/bjornpagen/goplay"
)

func parseCustomFields(ctx context.Context, c *chrome.Browser) ([]CustomField, error) {
	js := `JSON.stringify({
		fields: Array.from(document.querySelectorAll("div#custom_fields div.field")).map((el) => {
			let inputNode = el.querySelector("select, input:not([type='hidden'])")
			let type = inputNode.tagName.toLowerCase()

			if (type === "input") {
				type = inputNode.type.toLowerCase()
			}

			const selector = '#' + inputNode.id

			// looks like 'LinkedIn Profile *\n'
			const rawLabel = el.innerText
			let label = rawLabel.slice(0, rawLabel.indexOf("*")).trim()
			const required = rawLabel !== label
			label = label.trim()

			let options
			if (type === "select") {
				labels = Array.from(el.querySelectorAll("select option")).map((el) => el.textContent)
				options = labels.map((label, _) => ({ label }))
			} else if (type === "radio") {
				let labels = Array.from(el.querySelectorAll("label:has(input[type='radio'])")).slice(1).map((el) => el.innerText.trim())
				let selectors = Array.from(el.querySelectorAll("input[type='radio']")).map((el) => el.id)
				options = labels.map((label, i) => ({ label, selector: "input[type='radio']#" + selectors[i] }))
			} else if (type === "checkbox") {
				let labels = Array.from(el.querySelectorAll("label:has(input[type='checkbox'])")).slice(1).map((el) => el.innerText.trim())
				let selectors = Array.from(el.querySelectorAll("input[type='checkbox']")).map((el) => el.id)
				options = labels.map((label, i) => ({ label, selector: "input[type='checkbox']#" + selectors[i] }))
			}

			return { type, label, required, selector, options }
		})
	})`

	s, err := c.Evaluate(ctx, js)
	if err != nil {
		return nil, err
	}

	var m struct {
		Fields []CustomField `json:"fields"`
	}

	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	return m.Fields, nil
}

type PageInfo struct {
	Title    string `json:"title,omitempty"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	Content  string `json:"content,omitempty"`
}

func parsePageInfo(ctx context.Context, c *chrome.Browser) (m PageInfo, err error) {
	js := `JSON.stringify({
		title: document.querySelector(".app-title").innerText,
		company: document.querySelector(".company-name").innerText,
		location: document.querySelector(".location").innerText,
		content: document.querySelector("#content").innerHTML
	})`

	s, err := c.Evaluate(ctx, js)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return m, err
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(m.Content)
	if err != nil {
		return m, err
	}

	m.Content = strings.TrimSpace(markdown)
	m.Company = strings.TrimSpace(m.Company)
	m.Location = strings.TrimSpace(m.Location)
	m.Content = strings.TrimSpace(m.Content)

	return m, err
}

type FieldType = string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeSelect   FieldType = "select"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeRadio    FieldType = "radio"
)

type CustomField struct {
	Type     FieldType     `json:"type"`
	Label    string        `json:"label"`
	Required bool          `json:"required"`
	Selector string        `json:"selector"`
	Options  []FieldOption `json:"options"`
}

type FieldOption struct {
	Selector string `json:"selector"`
	Label    string `json:"label"`
}

func fillField(ctx context.Context, c *chrome.Browser, sel, text string) error {
	// Get the DomRect for the selector.
	node, err := c.Select(ctx, sel)
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}

	// Click it.
	err = c.Click(ctx, node)
	if err != nil {
		return fmt.Errorf("click: %w", err)
	}

	err = c.Text(ctx, text)
	if err != nil {
		return fmt.Errorf("text: %w", err)
	}

	return nil
}

func fillFile(ctx context.Context, c *chrome.Browser, sel, filePath string) error {
	node, err := c.Select(ctx, sel)
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}

	err = c.File(ctx, node, []string{filePath})
	if err != nil {
		return fmt.Errorf("file: %w", err)
	}

	return nil
}
