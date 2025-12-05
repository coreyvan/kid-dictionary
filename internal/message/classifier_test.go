package message_test

import (
	"testing"

	"github.com/coreyvan/kid-dictionary/internal/message"
	"github.com/stretchr/testify/assert"
)

func TestClassifier_Classify(t *testing.T) {
	classifier := message.NewClassifier()

	testCases := []struct {
		name     string
		content  string
		expected message.ContentTier
	}{
		// Normal topics (Tier 1)
		{
			name:     "normal question about sky",
			content:  "Why is the sky blue?",
			expected: message.ContentTierNormal,
		},
		{
			name:     "normal question about animals",
			content:  "How do birds fly?",
			expected: message.ContentTierNormal,
		},
		{
			name:     "normal question about science",
			content:  "What is gravity?",
			expected: message.ContentTierNormal,
		},

		// Sensitive topics (Tier 2)
		{
			name:     "sensitive topic - death",
			content:  "How do I explain death to my child?",
			expected: message.ContentTierSensitive,
		},
		{
			name:     "sensitive topic - divorce",
			content:  "What is divorce?",
			expected: message.ContentTierSensitive,
		},
		{
			name:     "sensitive topic - where babies come from",
			content:  "Where do babies come from?",
			expected: message.ContentTierSensitive,
		},
		{
			name:     "sensitive topic - illness",
			content:  "Why is grandma sick?",
			expected: message.ContentTierSensitive,
		},

		// Contextual topics (Tier 3)
		{
			name:     "contextual topic - religion",
			content:  "What is religion?",
			expected: message.ContentTierContextual,
		},
		{
			name:     "contextual topic - god",
			content:  "Does god exist?",
			expected: message.ContentTierContextual,
		},
		{
			name:     "contextual topic - politics",
			content:  "What is politics?",
			expected: message.ContentTierContextual,
		},
		{
			name:     "contextual topic - climate change",
			content:  "What is climate change?",
			expected: message.ContentTierContextual,
		},

		// Off-purpose requests (Tier 4 - Redirect)
		{
			name:     "off-purpose - write a poem",
			content:  "Write me a poem about cats",
			expected: message.ContentTierRedirect,
		},
		{
			name:     "off-purpose - help with taxes",
			content:  "Help with my taxes",
			expected: message.ContentTierRedirect,
		},
		{
			name:     "off-purpose - recipe",
			content:  "Give me a recipe for cookies",
			expected: message.ContentTierRedirect,
		},
		{
			name:     "off-purpose - code",
			content:  "Write some code for me",
			expected: message.ContentTierRedirect,
		},

		// On-purpose patterns should NOT be redirected
		{
			name:     "on-purpose - explain concept",
			content:  "Can you explain photosynthesis?",
			expected: message.ContentTierNormal,
		},
		{
			name:     "on-purpose - how to explain",
			content:  "How to explain to my child why the moon changes shape?",
			expected: message.ContentTierNormal,
		},
		{
			name:     "on-purpose - what is",
			content:  "What is a volcano?",
			expected: message.ContentTierNormal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := classifier.Classify(tc.content)
			assert.Equal(t, tc.expected, result, "content: %s", tc.content)
		})
	}
}
