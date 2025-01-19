package tests

import (
	"reflect"
	"testing"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
)

func TestExtractLeaguesFromResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		expectedOutput map[string]interface{}
		expectError    bool
	}{
		{
			name: "Valid Response with Leagues",
			input: map[string]interface{}{
				"fantasy_content": map[string]interface{}{
					"users": map[string]interface{}{
						"user": map[string]interface{}{
							"games": map[string]interface{}{
								"game": []interface{}{
									map[string]interface{}{
										"leagues": map[string]interface{}{
											"league": []interface{}{
												map[string]interface{}{
													"league_key": "403.l.84093",
													"name":       "Sad Degens",
												},
												map[string]interface{}{
													"league_key": "411.l.53877",
													"name":       "Another League",
												},
											},
										},
									},
									map[string]interface{}{
										"leagues": map[string]interface{}{
											"league": map[string]interface{}{
												"league_key": "453.l.29317",
												"name":       "Tip-Top League",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedOutput: map[string]interface{}{
				"leagues": []interface{}{
					map[string]interface{}{
						"league_key": "403.l.84093",
						"name":       "Sad Degens",
					},
					map[string]interface{}{
						"league_key": "411.l.53877",
						"name":       "Another League",
					},
					map[string]interface{}{
						"league_key": "453.l.29317",
						"name":       "Tip-Top League",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Missing Fantasy Content",
			input: map[string]interface{}{
				"invalid_key": map[string]interface{}{},
			},
			expectedOutput: nil,
			expectError:    true,
		},
		{
			name: "Missing Leagues in Game",
			input: map[string]interface{}{
				"fantasy_content": map[string]interface{}{
					"users": map[string]interface{}{
						"user": map[string]interface{}{
							"games": map[string]interface{}{
								"game": []interface{}{
									map[string]interface{}{},
								},
							},
						},
					},
				},
			},
			expectedOutput: map[string]interface{}{
				"leagues": []interface{}{},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := services.ExtractLeaguesFromResponse(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(result, tc.expectedOutput) {
					t.Errorf("Expected:\n%v\nGot:\n%v", tc.expectedOutput, result)
				}
			}
		})
	}
}
