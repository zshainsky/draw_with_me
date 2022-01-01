import { LitElement, html, css } from "lit";
import {roomStyles} from './styles.js';

class RoomCanvas extends LitElement {
    
    static properties = {
        width: {},
        height: {},
        isMouseDown: {type: Boolean},
        curX: {type: Number},
        curY: {type: Number},
        lastX: {type: Number},
        lastY: {type: Number},
        color: {type: String},
        
        canvasState: {},
        canvasEl: {},

        userList: {type: String},
        roomId: {type: String},
        roomName: {type: String},

        wsConn: {type: Object},
    };

    static styles = roomStyles;

    connectedCallback() {
        super.connectedCallback();
        //console.log("connectedCallback()");
        //console.log(this.shadowRoot.querySelector("#canvas"));
    }
    async firstUpdated() {
        // Give the browser a chance to paint
        await new Promise((r) => setTimeout(r, 0));
        this.canvasEl = this.shadowRoot.querySelector("#canvas");
        //console.log("Selected Canvas: ", this.canvasEl);
        window.addEventListener("resize", evt => this.handleResize(evt) );
        this.handleResize();
    }
    constructor() {
        super();
        // default values
        this.width = 1500;
        this.height = 700;

        this.wsConn = this.connectToWS();
        //console.log(this.wsConn);

        this.curX=0;
        this.curY=0;
        this.lastX=0;
        this.lastY=0;
        this.color = "#F2500F";

        // Empty user list
        this.userList = JSON.stringify({ActiveUsers:[]});
        // Empty CanvasState JSON object
        this.initCanvasState();

        // Set room ID
        this.userAuthId = "";
        this.roomId = "";
        this.roomName = "";
        
    }

    render() {
        return html`
            <div id="canvas-parent" class="canvas-parent" @changed-color="${this.handleChangedColor}">
                <active-user-bar .userList="${this.userList}"></active-user-bar>
                
                <tool-palette .initColor="${this.color}"></tool-palette>
                <canvas id="canvas" width="${this.width}" height="${this.height}" @mousedown="${this.handleMouseDown}" @mouseup="${this.handleMouseUp}" @mousemove="${this.handleMouseMove}"  @touchstart="${this.handleMouseDown}" @touchend="${this.handleMouseUp}" @touchmove="${this.handleMouseMove}"></canvas>
            </div>  
        `
    }
    handleResize(e){
        var rect = this.canvasEl.getBoundingClientRect();
        
        this.canvasEl.width = rect.width;
        this.canvasEl.height = rect.height;

        this.paintAllEvents(this.canvasEl, this.canvasState["CanvasState"]);
    }

    handleChangedColor (e) {
        this.color = e.detail.color;
    }

    handleMouseUp(e) {
        if (this.isMouseDown) {
            this.isMouseDown = false;
        }
    }

    handleMouseDown(e) {
        this.isMouseDown = true;
        var canvas = e.target;
        var firstTouch;

        if (e.type == "touchstart") { 
            // disable page scroll
            e.preventDefault();
            // only handle first touch
            firstTouch = e.touches[0]; 
        }
        
        // set mouse position variables
        this.curX = (e.pageX || parseInt(firstTouch.pageX)) - canvas.offsetLeft;
        this.curY = (e.pageY || parseInt(firstTouch.pageY)) - canvas.offsetTop;
        this.lastX = this.curX;
        this.lastY = this.curY;
    }

    handleMouseMove(e) {
        var canvas = e.target;
        var firstTouch;

        if (e.type == "touchmove") { 
            // disable page scroll
            e.preventDefault();
            firstTouch = e.touches[0]; 
        }

        if (this.isMouseDown) {
            // set mouse position variables
            this.lastX = this.curX;
            this.lastY = this.curY;
            this.curX = (e.pageX || parseInt(firstTouch.pageX)) - canvas.offsetLeft;
            this.curY = (e.pageY || parseInt(firstTouch.pageY)) - canvas.offsetTop;

            // paint
            // var paintJSON = this.paint(this.canvasEl, this.curX, this.curY, this.lastX, this.lastY, this.color, this.userAuthId, this.roomId, this.canvasEl.width, this.canvasEl.height);            // format paint event
            var paintJSON = this.paint(canvas, this.curX, this.curY, this.lastX, this.lastY, this.color, this.userAuthId, this.roomId, canvas.width, canvas.height);            // format paint event

            // make sure the 
            this.addPaintEventToCanvasState(paintJSON);
            this.wsConn.send(JSON.stringify(paintJSON));

        }
    }
    
    paint (canvasToDrawOn, curX, curY, lastX, lastY, color, userId, roomId, canvasToDrawFromWidth, canvasToDrawFromHeight) {
        // get the passed in canvas context variable
        var ctx = canvasToDrawOn.getContext('2d');
        if (!ctx) {
            console.log("issue getting canvasToDrawOn 2d context variable");
            return {"PaintEvent":{}}
        }
        // this is the default value for canvas size ... may never get called in production env
        if (canvasToDrawFromWidth == 0 && canvasToDrawFromHeight == 0) {
            canvasToDrawFromWidth = 1500;
            canvasToDrawFromHeight = 700;
        }  

        // create scale variables to make sure the current canvas displays the paint event in the correct plays based on the canvas it was sent from
        var scaleX = canvasToDrawOn.width / canvasToDrawFromWidth;
        var scaleY = canvasToDrawOn.height / canvasToDrawFromHeight;

        // set line width
        ctx.lineWidth = 5;
        ctx.lineJoin = 'round';
        ctx.lineCap = 'round';
        ctx.strokeStyle = color;
    
        // paint
        ctx.beginPath();
        ctx.moveTo(lastX*scaleX, lastY*scaleY);
        ctx.lineTo(curX*scaleX, curY*scaleY);
        ctx.closePath();
        ctx.stroke();

        // return values to send to ws
        return {"PaintEvent":{EvtTime: Date.now(), UserId: userId, RoomId: roomId, CurX: curX, CurY: curY, LastX: lastX, LastY: lastY, Color: color, CanvasWidth: canvasToDrawOn.width, CanvasHeight: canvasToDrawOn.height}};
    }

    // display all evnets foudn in the jsonPaintEventsList list of PaintEvenets
    // jsonPaintEventsList: is an array of paint events: [{"PaintEvent":{evtTime: Date.now(), userId: userId, roomId: roomId, curX: curX, curY: curY, lastX: lastX, lastY: lastY, color: color, canvasWidth: canvasToDrawOn.width, canvasHeight: canvasToDrawOn.height}]
    paintAllEvents(canvasToDrawOn, jsonPaintEventsList) {
	    for (let i in jsonPaintEventsList) {
            this.paint(canvasToDrawOn, jsonPaintEventsList[i]["CurX"], jsonPaintEventsList[i]["CurY"], jsonPaintEventsList[i]["LastX"], jsonPaintEventsList[i]["LastY"], jsonPaintEventsList[i]["Color"], jsonPaintEventsList[i]["AuthId"], jsonPaintEventsList[i]["RoomId"], jsonPaintEventsList[i]["CanvasWidth"], jsonPaintEventsList[i]["CanvasHeight"]);        }
    }

    // add a paint event to the current CanvasState json boject
    // jsonPaintEvent: is a single json paint event: {"PaintEvent":{evtTime: Date.now(), userId: userId, roomId: roomId, curX: curX, curY: curY, lastX: lastX, lastY: lastY, color: color, canvasWidth: canvasToDrawOn.width, canvasHeight: canvasToDrawOn.height}
    addPaintEventToCanvasState(jsonPaintEvent) {
        this.canvasState["CanvasState"].push(jsonPaintEvent);
    }

    // initialize the canvas state variable which is a JSON representation of all events on the canvas
    initCanvasState() {
        this.canvasState = {"CanvasState": []};
    }
    
    connectToWS(ctx) {
        // check if window has a websocket
        if (window['WebSocket']) {
            // check websocket connection protocol (wss:// is a secure connection)
            var wsProtocol = 'ws://';
            if (location.protocol == "https:") {
                wsProtocol = 'wss://';
            }
            // create websocket connection
            const conn = new WebSocket(wsProtocol + document.location.host + document.location.pathname + '/ws');
            // console.log(wsProtocol + document.location.host + document.location.pathname + '/ws');
            
            conn.onopen = function () {
                console.log("WS Open");
            }

            conn.onclose = evt => { 
                // wsStatus.innerHTML = "WS Closed";
                console.log("WS Closed");
            }

            conn.onerror = function (e) {
                console.log("ws error: " + e);
            }
            
            conn.onmessage = evt => {
                // find canvas on page by tag
                var canvas = this.shadowRoot.querySelector('#canvas');
                if (!canvas) {
                    console.log("could not find canvas");
                    return
                }
                // get canvas 2d context
                var ctx = canvas.getContext('2d');
                if (!ctx) {
                    console.log("could not find canvas");
                    return
                }
                // This only works if evt.data is receicing a Paint JSON event
                if (evt.data) {
                    try {
                        // try to parse the websocket event (evt) into JSON. If error throw error
                        var jsonEvent = JSON.parse(evt.data);
                        if (!Object.keys(jsonEvent).length > 0) {
                            return
                        }

                        // get the "key" value of event. (i.e. if {"CanvasEvent":[]}) return "CanvasEvent")
                        var key = Object.keys(jsonEvent)[0];
                        switch (key) {
                            case "PaintEvent":
                                var jsonPaintEvent = jsonEvent[key]
                                // this.paint(this.canvasEl, jsonPaintEvent.CurX, jsonPaintEvent.CurY, jsonPaintEvent.LastX, jsonPaintEvent.LastY, jsonPaintEvent.Color, jsonPaintEvent.UserId, jsonPaintEvent.RoomId, jsonPaintEvent.CanvasWidth, jsonPaintEvent.CanvasHeight);
                                this.paint(canvas, jsonPaintEvent.CurX, jsonPaintEvent.CurY, jsonPaintEvent.LastX, jsonPaintEvent.LastY, jsonPaintEvent.Color, jsonPaintEvent.UserId, jsonPaintEvent.RoomId, jsonPaintEvent.CanvasWidth, jsonPaintEvent.CanvasHeight);
                                this.addPaintEventToCanvasState(jsonPaintEvent);
                                break;

                            case "CanvasState":
                                this.canvasState = jsonEvent;
                                var jsonPaintEventsList = jsonEvent[key];
                                this.paintAllEvents(canvas, jsonPaintEventsList)
                                break;

                            case "ActiveUsers":
                                this.userList = JSON.stringify(jsonEvent);
                                break;

                            case "CurrentUser":
                                this.userAuthId = jsonEvent[key]["AuthId"];
                                break;

                            case "CurrentRoom":
                                this.roomId = jsonEvent[key]["Id"];
                                this.roomName = jsonEvent[key]["Name"]
                                break;
                            default:
                                console.log("Other event: "+ key);
                                break;
                        }
                    } catch(e) {
                        console.log("error parsing json event from websocket: " + e); // error in the above string (in this case, yes)!
                    }
                }
            }
            return conn;
        }
        return
        
    }
}
customElements.define('room-canvas', RoomCanvas)