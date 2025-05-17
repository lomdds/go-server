CREATE TABLE exam_preparation_courses (
    id SERIAL PRIMARY KEY,
    subject VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    relevance INTEGER NOT NULL,
    number_of_classes INTEGER,
    contact_the_teacher BOOLEAN DEFAULT FALSE,
    individuality BOOLEAN DEFAULT FALSE,
    price INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);