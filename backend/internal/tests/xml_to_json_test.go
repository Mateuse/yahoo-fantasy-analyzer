package tests

import (
	"reflect"
	"testing"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func TestXMLtoJSON(t *testing.T) {
	tests := []struct {
		name           string
		inputXML       string
		expectedJSON   map[string]interface{}
		expectingError bool
	}{
		{
			name: "Valid XML to JSON",
			inputXML: `
			<root>
				<item>
					<name>Example</name>
					<value>42</value>
				</item>
				<item>
					<name>Test</name>
					<value>100</value>
				</item>
			</root>`,
			expectedJSON: map[string]interface{}{
				"root": map[string]interface{}{
					"item": []interface{}{
						map[string]interface{}{"name": "Example", "value": "42"},
						map[string]interface{}{"name": "Test", "value": "100"},
					},
				},
			},
			expectingError: false,
		},
		{
			name:           "Empty XML Input",
			inputXML:       ``,
			expectedJSON:   nil,
			expectingError: true,
		},
		{
			name:           "Malformed XML",
			inputXML:       `<root><item><name>Example</name><value>42</value></root>`,
			expectedJSON:   nil,
			expectingError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.XMLtoJSON([]byte(tc.inputXML))
			if tc.expectingError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				// Compare normalized JSON maps
				if !reflect.DeepEqual(tc.expectedJSON, result) {
					t.Errorf("Expected JSON:\n%v\nGot:\n%v", tc.expectedJSON, result)
				}
			}
		})
	}
}
