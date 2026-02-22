package postgres

import (
	"testing"
)

type mockParams struct {
	ID    int
	Name  string
	Email string
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name        string
		templString string
		values      mockParams
		transform   transformFunc[mockParams]
		expected    string
		wantErr     bool
	}{
		{
			name:        "simple template with single field",
			templString: "SELECT * FROM users WHERE id = {{.ID}}",
			values:      mockParams{ID: 1, Name: "John", Email: "john@example.com"},
			transform: func(params mockParams) map[string]any {
				return map[string]any{
					"ID": params.ID,
				}
			},
			expected: "SELECT * FROM users WHERE id = 1",
			wantErr:  false,
		},
		{
			name:        "template with multiple fields",
			templString: "INSERT INTO users (name, email) VALUES ('{{.Name}}', '{{.Email}}')",
			values:      mockParams{ID: 1, Name: "John", Email: "john@example.com"},
			transform: func(params mockParams) map[string]any {
				return map[string]any{
					"Name":  params.Name,
					"Email": params.Email,
				}
			},
			expected: "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
			wantErr:  false,
		},
		{
			name:        "template with no fields",
			templString: "SELECT * FROM users",
			values:      mockParams{},
			transform: func(params mockParams) map[string]any {
				return map[string]any{}
			},
			expected: "SELECT * FROM users",
			wantErr:  false,
		},
		{
			name:        "template with special characters",
			templString: "SELECT * FROM users WHERE name = '{{.Name}}'",
			values:      mockParams{Name: "O'Brien", Email: "obrien@example.com"},
			transform: func(params mockParams) map[string]any {
				return map[string]any{
					"Name": params.Name,
				}
			},
			expected: "SELECT * FROM users WHERE name = 'O'Brien'",
			wantErr:  false,
		},
		{
			name:        "invalid template syntax",
			templString: "SELECT * FROM users WHERE id = {{.ID",
			values:      mockParams{ID: 1},
			transform: func(params mockParams) map[string]any {
				return map[string]any{"ID": params.ID}
			},
			expected: "",
			wantErr:  true,
		},
		{
			name:        "template with numeric values",
			templString: "SELECT * FROM users WHERE id = {{.ID}} AND name = '{{.Name}}'",
			values:      mockParams{ID: 42, Name: "Alice"},
			transform: func(params mockParams) map[string]any {
				return map[string]any{
					"ID":   params.ID,
					"Name": params.Name,
				}
			},
			expected: "SELECT * FROM users WHERE id = 42 AND name = 'Alice'",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQueryWithValues(tt.templString, tt.values, tt.transform)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if got != tt.expected {
				t.Errorf("ParseQuery() = %q, want %q", got, tt.expected)
			}
		})
	}
}
