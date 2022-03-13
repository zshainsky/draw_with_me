# Draw with me
Draw with me is a real-time collaborative drawing application! Sign up for an account, create a room, and share the room URL with a friend to start drawing together!

Sign up for free: https://draw-with-me-io.herokuapp.com/

![](./static/img/readme/draw_with_me.gif)

# Technologies and Design
## Backend Organizational Units
Draw with me me is written in golang and can be understood through several architectural components:

### Server
The server runs the web application, establishes all of the public http routes, and loads existing users and rooms into memory

### Rooms
Rooms represent a place for multiple users to interact on a shared canvas. A creates the http end point for a specific room and creates a new Hub upon startup. Room:Hub has a 1:1 relationship. 

### Users
Users represent a single account in the application. Users can create and share rooms. 

## Backend Architecture
### Hub
A Hub handles all event communication (inbound & outbound) between clients for a specific room. The A Hub is composed of 1 goroutine that listens on incoming channels for specific events and broadcasts the events to all active clients for that room. The Hub handles the following operations:
1. Maintain a list of active clients by listening for activation and deactivation events from the client
2. Broadcast PaintEvents from a given client to all other clients
3. Save canvas state to the persistent layer (DB)

### Client
A Client handles connection to the FE via websocket connection, communication of FE events to the Hub, and receiving broadcasted events from the Hub to update the FE through the websocket connection. The Client runs 2 goroutines upon a new client creation. One goroutine handles listening for events from the FE to send to the Hub (`sendToHub()`) and the other handles listening for broadcasted events received from the Hub (`writeToWS()`).

The Client handles the following operations:
1. Establish Websocket connections upon user logging into a room
2. Send activation message to Hub
3. Send paint events to the Hub
4. Recieve paint events from the Hub and update FE with received events

__Network flow example__:

Hub <-`sendToHub()`-> Client <-`writeToWS()`-> FE (websocket)

### Database
Draw with me uses a simple PostgreSQL database to maintain application state. The Database schema can be found in the, [schema.sql](./db/schema.sql) file. The database is intended to maintain, Users, Rooms, and all PaintEvents that have been drawn in a specific room. The tables are as follows:
1. `users`: Simple table to maintain all users. PK is assigned to the `id` column
2. `rooms`: Simple table to maintain all rooms. PK is assigned to the `id` column
3. `user_room`: Table to maintain which users are found in which rooms. FK relationship between `user_id` and `id` in the users table and `room_id` and `id` in the rooms table
4. `canvas_state`: Table stores all paint events for each room in JSON format (JSONB data type). FK relationship between `room_id` and `id` in the rooms table. Paint event 

Example Canvas State JSON Format:
```
"CanvasState": [
    {
        "CurX": 181,
        "CurY": 175,
        "LastX": 180,
        "LastY": 177,
        "Color": "#F2500F",
        "UserId": "12345667-472f-472f-472f-20becadb522a",
        "RoomId": "12345667-472f-472f-472f-20becadb522b",
        "EvtTime": 1640646814012146
    }
]
```
5. `paint_event`: (_not used_) Table to store individual paint events for each user per room. FK relationship between `user_id` and `id` in the users table and `room_id` and `id` in the rooms table

## Frontend Architecture
The Frontend architecture was built from scratch using WebComponents on top of [lit element](https://lit.dev/). 

Elements are built in the [lib/components/](./lib/components/) folder. These elements are packaged using npm [rollup](https://www.npmjs.com/package/rollup) and stored to be served in the [lib/static/js/] folder. Details of the rollup config can be found here [lib/rollup.config.js](./lib/rollup.config.js). All HTML pages reference packaged js from the static folder. Please view [lib/package.json](./lib/package.json) for details on other dependencies. All CSS is defined in the [lib/components/styles.js](./lib/components/styles.js) file.


## Deployment
This app is deployed using the free tier of Heroku (Hobby Dev). It automatically deploys new versions when changes are committed to the heroku-main branch of this git repo. This app uses the `heroku/go` build pack and is deployed using a `web` Dyno. See the definition in the [Procfile](./Procfile). I have also added the PostgreSQL Add-on and configured the database manually using the [[schema.sql](./db/schema.sql) file.

## Resources
1. Architecture/Design: [https://outcrawl.com/realtime-collaborative-drawing-go](https://outcrawl.com/realtime-collaborative-drawing-go)
2. Go Testing: [https://quii.gitbook.io/learn-go-with-tests/](https://quii.gitbook.io/learn-go-with-tests/)
3. Go Concurrency: [https://www.udemy.com/course/up-n-running-with-concurrency-in-golang/](https://www.udemy.com/course/up-n-running-with-concurrency-in-golang/)
4. Lit Element: [https://lit.dev/](https://lit.dev/)