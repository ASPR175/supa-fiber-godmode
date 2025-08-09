ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
CREATE TABLE tasks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
      user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE POLICY "Allow logged-in access" 
ON tasks FOR ALL USING (auth.uid() = user_id);