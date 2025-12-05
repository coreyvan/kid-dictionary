package message_test

import (
	"testing"

	"github.com/coreyvan/kid-dictionary/internal/domain"
	"github.com/coreyvan/kid-dictionary/internal/service/message"
	"github.com/stretchr/testify/assert"
)

func TestClassifier_Classify(t *testing.T) {
	classifier := message.NewClassifier()

	testCases := []struct {
		name     string
		content  string
		expected domain.ContentTier
	}{
		// Normal topics (Tier 1)
		{
			name:     "normal question about sky",
			content:  "Why is the sky blue?",
			expected: domain.ContentTierNormal,
		},
		{
			name:     "normal question about animals",
			content:  "How do birds fly?",
			expected: domain.ContentTierNormal,
		},
		{
			name:     "normal question about science",
			content:  "What is gravity?",
			expected: domain.ContentTierNormal,
		},

		// Sensitive topics (Tier 2)
		{
			name:     "sensitive topic - death",
			content:  "How do I explain death to my child?",
			expected: domain.ContentTierSensitive,
		},
		{
			name:     "sensitive topic - divorce",
			content:  "What is divorce?",
			expected: domain.ContentTierSensitive,
		},
		{
			name:     "sensitive topic - where babies come from",
			content:  "Where do babies come from?",
			expected: domain.ContentTierSensitive,
		},
		{
			name:     "sensitive topic - illness",
			content:  "Why is grandma sick?",
			expected: domain.ContentTierSensitive,
		},

		// Contextual topics (Tier 3)
		{
			name:     "contextual topic - religion",
			content:  "What is religion?",
			expected: domain.ContentTierContextual,
		},
		{
			name:     "contextual topic - god",
			content:  "Does god exist?",
			expected: domain.ContentTierContextual,
		},
		{
			name:     "contextual topic - politics",
			content:  "What is politics?",
			expected: domain.ContentTierContextual,
		},
		{
			name:     "contextual topic - climate change",
			content:  "What is climate change?",
			expected: domain.ContentTierContextual,
		},

		// Off-purpose requests (Tier 4 - Redirect)
		{
			name:     "off-purpose - write a poem",
			content:  "Write me a poem about cats",
			expected: domain.ContentTierRedirect,
		},
		{
			name:     "off-purpose - help with taxes",
			content:  "Help with my taxes",
			expected: domain.ContentTierRedirect,
		},
		{
			name:     "off-purpose - recipe",
			content:  "Give me a recipe for cookies",
			expected: domain.ContentTierRedirect,
		},
		{
			name:     "off-purpose - code",
			content:  "Write some code for me",
			expected: domain.ContentTierRedirect,
		},

		// On-purpose patterns should NOT be redirected
		{
			name:     "on-purpose - explain concept",
			content:  "Can you explain photosynthesis?",
			expected: domain.ContentTierNormal,
		},
		{
			name:     "on-purpose - how to explain",
			content:  "How to explain to my child why the moon changes shape?",
			expected: domain.ContentTierNormal,
		},
		{
			name:     "on-purpose - what is",
			content:  "What is a volcano?",
			expected: domain.ContentTierNormal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := classifier.Classify(tc.content)
			assert.Equal(t, tc.expected, result, "content: %s", tc.content)
		})
	}
}
