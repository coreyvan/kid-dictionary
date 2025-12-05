-- +migrate Up
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,  -- nullable for anonymous MVP
    title TEXT NOT NULL DEFAULT 'New Conversation',
    age_bracket INTEGER NOT NULL,  -- enum: 1=LITTLE_ONES, 2=GROWING_MINDS, 3=PRE_TEENS
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_conversations_user_id ON conversations(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_conversations_created_at ON conversations(created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS conversations;
