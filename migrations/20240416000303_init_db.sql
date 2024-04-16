-- +goose Up
-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    configurations TEXT, -- Assuming JSON or serialized string for configurations
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TABLE email_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_email VARCHAR(255) NOT NULL,
    recipient_emails TEXT NOT NULL, -- For JSON or array of emails
    cc_emails TEXT, -- Optional for JSON or array of CC emails
    bcc_emails TEXT, -- Optional for JSON or array of BCC emails
    subject VARCHAR(255),
    body TEXT NOT NULL,
    body_format VARCHAR(50) DEFAULT 'plain', -- 'plain' or 'html'
    attachments JSON, -- JSON array of attachment details
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TABLE email_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL,
    email_message_id UUID NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    status VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    FOREIGN KEY (provider_id) REFERENCES providers(id),
    FOREIGN KEY (email_message_id) REFERENCES email_messages(id)
);
-- +goose Down
DROP TABLE IF EXISTS email_logs;
DROP TABLE IF EXISTS email_messages;
DROP TABLE IF EXISTS providers;
