package prompts

import "github.com/tmc/langchaingo/prompts"

var business_summary_prompt = prompts.PromptTemplate{
	Template: `<goal>Your task is to provide a brief overview of the nearby places found.</goal>
    <instructions>
    - State the number of places found, perhaps by primary type if easily discernible (e.g., "Found 10 places, mostly cafes and restaurants.").
    - Mention any highly-rated or prominent places if evident from the data, but keep it brief.
    - Do not list all businesses. Provide a general summary.
    - If specific types were queried (e.g. "cafes"), focus the summary on that.
    </instructions>
    <tuning_instructions>
      Verbosity: {{.verbosity_instruct}}
      Response Length: {{.response_length_instruct}}
      Formality: {{.formal_level_instruct}}
      Creativity: {{.creativity_instruct}}
      Precision: {{.precision_instruct}}
      User Instructions: {{.user_instruction}}
    </tuning_instructions>
    <business_data>
    {{.business_data}}
    </business_data>
    <summary_guidelines>
    Provide the summary directly. Example: "Found 12 nearby places, including several cafes and a few highly-rated restaurants." or "Found 5 electronics stores in the area."
    </summary_guidelines>
    Summary:
	`,
	InputVariables: []string{"verbosity_instruct", "response_length_instruct", "formal_level_instruct", "creativity_instruct", "precision_instruct", "user_instruction", "business_data"},
}
