package certificate

// Valid values for Entry.UsernameField.
const (
	UsernameFieldSubject    = "subject"
	UsernameFieldSubjectAlt = "subject-alt"
)

// Valid values for Entry.UsernameFieldValue when
// `UsernameField="subject"`.
const (
	UsernameFieldValueCommonName = "common-name"
)

// Valid values for Entry.UsernameFieldValue when
// `UsernameField="subject-alt"`.
const (
	UsernameFieldValueEmail         = "email"
	UsernameFieldValuePrincipalName = "principal-name"
)

const (
	singular = "certificate profile"
	plural   = "certificate profiles"
)
