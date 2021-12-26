/* 
 * This file is not used as part of the compiled application. 
 * This is intended to be used as a reference to the schema used in the application.
 */
CREATE DATABASE draw_with_me;

CREATE TABLE IF NOT EXISTS users (
    id varchar(36) PRIMARY KEY,
    authId varchar(21) UNIQUE NOT NULL, -- Map to 'sub' in JWT
    authType varchar(10) NOT NULL, -- draw.AuthType is a string. Using 'google' for the google type
    name text,
    email text,
    picture text
);
INSERT INTO users (id, authId, authType, name, email, picture) VALUES ('62769698-ca64-472f-6da7-20becadb522a', '113062537928538714745', 'google', 'test_username', 'test_email@test.com', 'https://test.pictures.com') RETURNING id;

CREATE TABLE IF NOT EXISTS rooms (
    id varchar(36) PRIMARY KEY NOT NULL,
    name text
);
INSERT INTO rooms (id, name) VALUES ('62769698-ca64-472f-6da7-20becadb522b', 'test_room_name') RETURNING id;

CREATE TABLE IF NOT EXISTS user_room (
    user_id varchar(36) NOT NULL,
    room_id varchar(36) NOT NULL,
    PRIMARY KEY (user_id, room_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (room_id) REFERENCES rooms (id)
);
INSERT INTO user_room (user_id, room_id) VALUES ('62769698-ca64-472f-6da7-20becadb522a', '62769698-ca64-472f-6da7-20becadb522b') RETURNING (user_id, room_id);

CREATE TABLE IF NOT EXISTS paint_event (
    user_id varchar(36) NOT NULL,
    room_id varchar(36) NOT NULL,
    lastX integer NOT NULL,
    lastY integer NOT NULL,
    curX integer NOT NULL,
    curY integer NOT NULL,
    color varchar(7) NOT NULL, -- should be hex color including '#' (ex: #0000FF)
    
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (room_id) REFERENCES rooms (id)
);
INSERT INTO paint_event (user_id, room_id, lastX, lastY, curX, curY, color) VALUES ('62769698-ca64-472f-6da7-20becadb522a', '62769698-ca64-472f-6da7-20becadb522b', 100, 100, 200, 200, '#000fff') RETURNING *;


