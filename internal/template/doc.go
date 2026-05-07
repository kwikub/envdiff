// Package template generates .env template files from existing environment
// configurations. It strips real values and optionally replaces secret fields
// with placeholders, making it safe to commit templates to version control.
//
// Usage:
//
//	err := template.Generate(".env.production", ".env.template", template.Options{
//		MaskSecrets: true,
//		AddComments: true,
//	})
//
// The generated template can be shared with new team members or used as a
// reference for required environment variables in a project.
package template
