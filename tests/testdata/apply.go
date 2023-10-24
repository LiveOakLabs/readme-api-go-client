package testdata

import (
	"net/http"

	"github.com/liveoaklabs/readme-api-go-client/readme"
)

var ApplyOpenRoles = []readme.OpenRole{
	{
		Slug:        "front-end-engineer",
		Title:       "Front End Engineer",
		Description: "foo bar",
		Pullquote: "Collaborative: Your code reviews and pull requests are detailed, " +
			"you’re always happy to share knowledge, and you love a good pairing session.",
		Location:   "Remote, US",
		Department: "Engineering",
		URL:        "https://jobs.ashbyhq.com/readme/43ebc7c3-4653-4037-841e-075ad428a68d/application",
	},
}

var ApplyResponseSuccess = readme.ApplyResponse{
	Message:   "Thanks for applying, Gordon! We'll reach out to you soon!",
	Keyvalues: "https://www.keyvalues.com/readme",
	Careers:   "https://readme.com/careers",
	Questions: "greg@readme.io",
	Poem: []string{
		"Thanks for applying to work at ReadMe!",
		"Your application is lookin' spiffy",
		"We're going to review it ASAP",
		"And we'll get back to you in a jiffy!",
	},
}

var ApplyErrorResponseInvalidName = readme.APIErrorResponse{
	Error:      "APPLY_INVALID_NAME",
	Message:    "You need to provide a name.",
	Suggestion: "To apply for a job, you need to include your name as a body parameter.",
	Docs:       "https://docs.readme.com/main/logs/ddfb22d8-dfc6-43e8-b15f-4f4d9092f27d",
	Help: "https://docs.readme.com/main/reference and include the following link " +
		"to your API log: 'https://docs.readme.com/main/logs/ddfb22d8-dfc6-43e8-b15f-4f4d9092f27d'.",
	Poem: []string{
		"We normally support online anonymity",
		"Except when you're trying to apply,",
		"Because if we're gonna work together",
		"We need to know what you go by!",
	},
}

var ApplyCreateResponseInvalidName = &readme.APIResponse{
	APIErrorResponse: ApplyErrorResponseInvalidName,
	Body:             []byte(ToJSON(ApplyErrorResponseInvalidName)),
	HTTPResponse: &http.Response{
		StatusCode: 400,
	},
}

var ApplyApplication = readme.Application{
	Name:  "Gordon Ramsay",
	Email: "gordon@example.com",
	Job:   "Front End Engineer",
}

var ApplyGetAPIResponse = &readme.APIResponse{
	APIErrorResponse: readme.APIErrorResponse{},
	Body:             []byte(ToJSON(ApplyOpenRoles)),
	HTTPResponse: &http.Response{
		StatusCode: 200,
	},
}

var ApplyCreateAPIResponse = &readme.APIResponse{
	APIErrorResponse: readme.APIErrorResponse{},
	Body:             []byte(ToJSON(ApplyResponseSuccess)),
	HTTPResponse: &http.Response{
		StatusCode: 200,
	},
}
