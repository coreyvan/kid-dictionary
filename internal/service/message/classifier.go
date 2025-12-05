package message

import (
	"strings"

	"github.com/coreyvan/kid-dictionary/internal/domain"
)

// Classifier determines the content tier for a user's message.
type Classifier struct{}

// NewClassifier creates a new content classifier.
func NewClassifier() *Classifier {
	return &Classifier{}
}

// sensitiveKeywords are topics that require soft guidance prefix.
var sensitiveKeywords = []string{
	"death", "die", "dying", "dead", "passed away",
	"divorce", "separated", "separation",
	"illness", "sick", "disease", "cancer", "hospital",
	"baby", "babies", "pregnant", "pregnancy", "birth", "born",
	"sex", "reproduction", "where do babies come from",
	"abuse", "violence", "hurt",
	"war", "terrorism", "terrorist",
	"drugs", "alcohol", "addiction",
	"suicide", "self-harm",
	"grief", "mourning", "funeral",
}

// contextualKeywords are topics that may need framing preference.
var contextualKeywords = []string{
	"religion", "god", "church", "faith", "belief",
	"politics", "election", "president", "government",
	"race", "racism", "discrimination",
	"gender", "transgender", "sexuality",
	"immigration", "immigrant",
	"abortion",
	"gun", "guns",
	"climate change", "global warming",
}

// offPurposePatterns detect requests that are not about explaining concepts.
var offPurposePatterns = []string{
	"write me a",
	"write a",
	"help me write",
	"can you write",
	"create a",
	"make me a",
	"help with my",
	"do my",
	"solve this",
	"calculate",
	"translate",
	"code",
	"program",
	"recipe",
	"directions to",
	"weather",
	"news",
	"stock",
	"price of",
	"tell me a joke",
	"sing",
	"play a game",
}

// onPurposePatterns detect valid child explanation requests.
var onPurposePatterns = []string{
	"explain",
	"what is",
	"what are",
	"what does",
	"why is",
	"why do",
	"why are",
	"how do",
	"how does",
	"how to explain",
	"how to talk to",
	"how to tell",
	"how should i",
	"meaning of",
	"what happens when",
	"help me explain",
	"how can i explain",
}

// Classify determines the content tier for the given message content.
func (c *Classifier) Classify(content string) domain.ContentTier {
	lower := strings.ToLower(content)

	// Check for off-purpose requests first
	if c.isOffPurpose(lower) {
		return domain.ContentTierRedirect
	}

	// Check for contextual topics (religion, politics, etc.)
	if c.containsAny(lower, contextualKeywords) {
		return domain.ContentTierContextual
	}

	// Check for sensitive topics
	if c.containsAny(lower, sensitiveKeywords) {
		return domain.ContentTierSensitive
	}

	return domain.ContentTierNormal
}

// isOffPurpose checks if the content is not related to explaining concepts to children.
func (c *Classifier) isOffPurpose(lower string) bool {
	// If it matches on-purpose patterns, it's valid even if it has off-purpose keywords
	for _, pattern := range onPurposePatterns {
		if strings.Contains(lower, pattern) {
			return false
		}
	}

	// Check if it contains off-purpose patterns
	for _, pattern := range offPurposePatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

// containsAny checks if the content contains any of the keywords.
func (c *Classifier) containsAny(lower string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}
