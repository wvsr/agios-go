package prompts

import "github.com/tmc/langchaingo/prompts"

var search_term_prompt = prompts.PromptTemplate{
	Template: `<prompt>
    <task>Extract the main specific search term from the given input string.</task>
    <instructions>
      <step>Read the input string associated with the given ID.</step>
      <step>Identify the most specific and central search term that best represents the user's query.</step>
      <step>Ignore general words, modifiers, or intent phrases (e.g., "how to", "best way to", "examples of").</step>
      <step>Return only the core search term, as a single value, inside the JSON key "search_terms".</step>
    </instructions>
    <input_string>{{.text}}</input_string>
    <output_format>
      <json>
        {{
          "search_term": str[]
        }}
      </json>
    </output_format>
  </prompt>`,
	InputVariables: []string{"text"}}
