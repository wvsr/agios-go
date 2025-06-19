package prompts

import "github.com/tmc/langchaingo/prompts"

var weather_summary_prompt = prompts.PromptTemplate{
	Template: `<goal>Your task is to provide a concise, human-readable summary of the provided weather forecast data.</goal>
    <instructions>
    - Highlight the current weather conditions (temperature, general outlook like sunny/cloudy/rainy).
    - Briefly mention the forecast for the next 1-2 days (e.g., "similar conditions tomorrow," "rain expected on [Day]").
    - Focus on key information like temperature ranges and significant precipitation.
    - Do not just list all the data fields. Synthesize it into a natural language summary.
    </instructions>
    <tuning_instructions>
      Verbosity: {{.verbosity_instruct}}
      Response Length: {{.response_length_instruct}}
      Formality: {{.formal_level_instruct}}
      Creativity: {{.creativity_instruct}}
      Precision: {{.precision_instruct}}
      User Instructions: {{.user_instruction}}
    </tuning_instructions>
    <weather_data>
    {{.weather_data}}
    </weather_data>
    <summary_guidelines>
    Provide the summary directly. Example: "Currently it's 25°C and sunny. Expect similar weather tomorrow, with a high of 28°C. Rain is possible the day after."
    </summary_guidelines>
    Summary:`,
	InputVariables: []string{"verbosity_instruct", "response_length_instruct", "formal_level_instruct", "creativity_instruct", "precision_instruct", "user_instruction", "weather_data"}}
