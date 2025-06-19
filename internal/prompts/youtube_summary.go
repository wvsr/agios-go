package prompts

import "github.com/tmc/langchaingo/prompts"

var youtube_summary_prompt = prompts.PromptTemplate{
	Template: `<goal>Your task is to provide a concise summary of the following YouTube video transcript.</goal>
    <instructions>
    - Focus on the main topics and key takeaways.
    - Keep the summary relatively short, around 3-5 sentences, unless the transcript is very long.
    - Do not include any personal opinions or interpretations not present in the transcript.
    - Write in a clear and neutral tone.
    </instructions>
    <tuning_instructions>
      Verbosity: {{.verbosity_instruct}}
      Response Length: {{.response_length_instruct}}
      Formality: {{.formal_level_instruct}}
      Creativity: {{.creativity_instruct}}
      Precision: {{.precision_instruct}}
      User Instructions: {{.user_instruction}}
    </tuning_instructions>
    <transcript>
    {{.transcript}}
    </transcript>
    <summary_guidelines>
    Provide the summary directly, without any introductory phrases like "Here is the summary:".
    </summary_guidelines>
    Summary:`, InputVariables: []string{"verbosity_instruct", "response_length_instruct", "formal_level_instruct", "creativity_instruct", "precision_instruct", "user_instruction", "transcript"}}
