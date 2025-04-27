-- Create generic_items table
CREATE TABLE IF NOT EXISTS generic_items (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_generic_items_status ON generic_items(status);
CREATE INDEX IF NOT EXISTS idx_generic_items_created_at ON generic_items(created_at);
CREATE INDEX IF NOT EXISTS idx_generic_items_name ON generic_items(name);

-- Add comments
COMMENT ON TABLE generic_items IS 'Generic items table for storing any type of entity';
COMMENT ON COLUMN generic_items.id IS 'Unique identifier for the item';
COMMENT ON COLUMN generic_items.name IS 'Name of the item';
COMMENT ON COLUMN generic_items.description IS 'Description of the item';
COMMENT ON COLUMN generic_items.status IS 'Status of the item (draft, active, inactive, archived)';
COMMENT ON COLUMN generic_items.created_at IS 'Timestamp when the item was created';
COMMENT ON COLUMN generic_items.updated_at IS 'Timestamp when the item was last updated';
