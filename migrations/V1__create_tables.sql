-- Create the Applications table
CREATE TABLE Applications (
    -- default index on id
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    -- default index on token
    token VARCHAR(255) NOT NULL UNIQUE,
    chats_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create the Chats table
CREATE TABLE Chats (
    -- default index on id
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT NOT NULL,
    subject TEXT NOT NULL,
    number INT NOT NULL,
    messages_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES Applications(id) ON DELETE CASCADE,
    -- default index on (application_id, number)
    UNIQUE (application_id, number)
);

-- Create the Messages table
CREATE TABLE Messages (
    -- default index on id
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    number INT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES Chats(id) ON DELETE CASCADE,
    -- default index on (chat_id, number)
    UNIQUE (chat_id, number)
);
