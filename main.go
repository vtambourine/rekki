package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ValidateRequest struct {
	Email string `json:"email"`
}

type Validator struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

type ValidatorResponse struct {
	Valid      bool                 `json:"valid"`
	Validators map[string]Validator `json:"validators"`
}

type Validators map[string]func(email string) (bool, string)

const (
	ReasonRegexpMismatch     = "REGEXP_MISMATCH"
	ReasonInvalidTLD         = "INVALID_TLD"
	ReasonInvalidHostname    = "INVALID_HOSTNAME"
	ReasonUntrustedDomain    = "UNTRUSTED_DOMAIN"
	ReasonUnableToConnect    = "UNABLE_TO_CONNECT"
	ReasonTimeout            = "CONNECTION_TIMEOUT"
	ReasonMailServerError    = "MAIL_SERVER_ERROR"
	ReasonUnavailableMailbox = "UNAVAILABLE_MAILBOX"
)

// TODO: Specify validators priority.
var validators = Validators{
	"regexp":    ValidateRegexp,
	"hostname":  ValidateDomain,
	"blacklist": ValidateBlacklist,
	"smtp":      ValidateSMTP,
}

var response ValidatorResponse

func validate(w http.ResponseWriter, r *http.Request) {
	// Accept only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	// Read request body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse body
	var v ValidateRequest
	err = json.Unmarshal(b, &v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Validate passed email
	validateEmail(v.Email)

	// Render result
	result, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Return validation result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(result)
}

func validateEmail(email string) {
	response.Valid = true
	response.Validators = make(map[string]Validator)

	// TODO: Do not runs all validation checks independently. Possibly return on first fail.
	// TODO: Run independent validators in parallel, for instance, with `sync.WaitGroup`.
	for name, validator := range validators {
		v, r := validator(email)
		response.Validators[name] = Validator{
			Valid:  v,
			Reason: r,
		}
		response.Valid = response.Valid && v
	}
}

func main() {
	http.HandleFunc("/email/validate", validate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
