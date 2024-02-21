CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE chat_rooms (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255), -- Nullable
    type VARCHAR(50) CHECK (type IN ('direct', 'group'))
);

-- Create the chat_room_participants table with foreign key constraints
CREATE TABLE chat_room_participants (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES chat_rooms(id),
    user_id INT REFERENCES users(id),
    UNIQUE (room_id, user_id) -- Ensure a user can only be in a room once
);

-- Create the messages table with foreign key constraints
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES chat_rooms(id),
    sender_id INT REFERENCES users(id),
    text TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE
);
