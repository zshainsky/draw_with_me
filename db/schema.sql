/* Heroku cheatsheet: 
 * > heroku pg:psql -a draw-with-me-io
 * > heroku pg:psql -f ./schema.sql -a draw-with-me-io
 * > draw-with-me-io::DATABASE=> select room_id, jsonb_array_length(canvas_json -> 'CanvasState') as len from canvas_state;
 * 
 * This file is not used as part of the compiled application. 
 * This is intended to be used as a reference to the schema used in the application.
 */
-- CREATE DATABASE draw_with_me;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS user_room CASCADE;
DROP TABLE IF EXISTS canvas_state CASCADE;
DROP TABLE IF EXISTS paint_event CASCADE;

CREATE TABLE IF NOT EXISTS users (
    id varchar(36) PRIMARY KEY,
    authId varchar(21) UNIQUE NOT NULL, -- Map to 'sub' in JWT
    authType varchar(10) NOT NULL, -- draw.AuthType is a string. Using 'google' for the google type
    name text,
    email text,
    picture text
);
-- INSERT INTO users (id, authId, authType, name, email, picture) VALUES ('62769698-ca64-472f-6da7-20becadb522a', '113062537928538714745', 'google', 'test_username', 'test_email@test.com', 'https://test.pictures.com') RETURNING id;

CREATE TABLE IF NOT EXISTS rooms (
    id varchar(36) PRIMARY KEY NOT NULL,
    name text
);
-- INSERT INTO rooms (id, name) VALUES ('62769698-ca64-472f-6da7-20becadb522b', 'test_room_name') RETURNING id;

CREATE TABLE IF NOT EXISTS user_room (
    user_id varchar(36) NOT NULL,
    room_id varchar(36) NOT NULL,
    PRIMARY KEY (user_id, room_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (room_id) REFERENCES rooms (id)
);
-- INSERT INTO user_room (user_id, room_id) VALUES ('62769698-ca64-472f-6da7-20becadb522a', '62769698-ca64-472f-6da7-20becadb522b') RETURNING (user_id, room_id);
-- INSERT INTO user_room (user_id, room_id) VALUES ('62769698-ca64-472f-6da7-20becadb522c', '62769698-ca64-472f-6da7-20becadb522b') RETURNING (user_id, room_id);

-- Use this to load instate canvas on starting the hub
CREATE TABLE IF NOT EXISTS canvas_state (
    room_id varchar(36) UNIQUE NOT NULL, -- This is unique bcause there should only be 1 record per room
    canvas_json JSONB,
    FOREIGN KEY (room_id) REFERENCES rooms (id)
);
-- INSERT INTO canvas_state (room_id, canvas_json) VALUES ('fd5fe37c-e64f-4315-6a78-c546e179cb3a','{
--     "CanvasState": [
--         {
--             "CurX": 181,
--             "CurY": 175,
--             "LastX": 180,
--             "LastY": 177,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 182,
--             "CurY": 170,
--             "LastX": 181,
--             "LastY": 175,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 186,
--             "CurY": 161,
--             "LastX": 182,
--             "LastY": 170,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 206,
--             "CurY": 143,
--             "LastX": 186,
--             "LastY": 161,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 226,
--             "CurY": 133,
--             "LastX": 206,
--             "LastY": 143,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 248,
--             "CurY": 128,
--             "LastX": 226,
--             "LastY": 133,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 255,
--             "CurY": 128,
--             "LastX": 248,
--             "LastY": 128,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 266,
--             "CurY": 132,
--             "LastX": 255,
--             "LastY": 128,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 271,
--             "CurY": 135,
--             "LastX": 266,
--             "LastY": 132,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 276,
--             "CurY": 144,
--             "LastX": 271,
--             "LastY": 135,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 278,
--             "CurY": 150,
--             "LastX": 276,
--             "LastY": 144,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 276,
--             "CurY": 159,
--             "LastX": 278,
--             "LastY": 150,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 272,
--             "CurY": 165,
--             "LastX": 276,
--             "LastY": 159,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 253,
--             "CurY": 177,
--             "LastX": 272,
--             "LastY": 165,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 238,
--             "CurY": 182,
--             "LastX": 253,
--             "LastY": 177,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 214,
--             "CurY": 188,
--             "LastX": 238,
--             "LastY": 182,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 204,
--             "CurY": 188,
--             "LastX": 214,
--             "LastY": 188,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 192,
--             "CurY": 187,
--             "LastX": 204,
--             "LastY": 188,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 189,
--             "CurY": 185,
--             "LastX": 192,
--             "LastY": 187,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         },
--         {
--             "CurX": 188,
--             "CurY": 184,
--             "LastX": 189,
--             "LastY": 185,
--             "Color": "#F2500F",
--             "UserId": "62769698-ca64-472f-6da7-20becadb522a",
--             "RoomId": "62769698-ca64-472f-6da7-20becadb522b",
--             "EvtTime": 1640646814012146
--         }
--     ]
-- }') RETURNING canvas_json;


/* 
 * Use canvas_state to room canvas events into memory
 * Defering to canvas_state to reduce number of Writes to this table and reduce large number of rows queried.
 * This table may only be used for recreating canvas playbacks or other functions.
 * ...This was created for furture work and because I was on a roll like butter
*/
CREATE TABLE IF NOT EXISTS paint_event (
    evt_time bigint NOT NULL, -- epoch time (ex: 1640646814.012146)
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
-- INSERT INTO paint_event (evt_time, user_id, room_id, lastX, lastY, curX, curY, color) VALUES (1640646814.012146,'62769698-ca64-472f-6da7-20becadb522a', '62769698-ca64-472f-6da7-20becadb522b', 100, 100, 200, 200, '#000fff') RETURNING *;
