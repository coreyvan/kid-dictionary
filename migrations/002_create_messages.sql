-- +migrate Up
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role INTEGER NOT NULL,  -- enum: 1=USER, 2=ASSISTANT
    content TEXT NOT NULL,
    content_tier INTEGER,  -- enum: 1=NORMAL, 2=SENSITIVE, 3=CONTEXTUAL, 4=REDIRECT (nullable for user messages)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_conversation_created ON messages(conversation_id, created_at);

-- +migrate Down
DROP TABLE IF EXISTS messages;
