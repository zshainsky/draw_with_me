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
        paintJSON: {},

        userList: {type: String},

        wsConn: {type: Object},
    };

    static styles = roomStyles;

    connectedCallback() {
        super.connectedCallback();
        console.log("connectedCallback()");

    }
    constructor() {
        super();
        // default values
        this.width = 1500;
        this.height = 700;

        this.wsConn = this.connectToWS();
        console.log(this.wsConn);

        this.curX=0;
        this.curY=0;
        this.lastX=0;
        this.lastY=0;
        this.color = "#F2500F";

        // Empty user list
        this.userList = JSON.stringify({ActiveUsers:[]});
    }

    render() {
        return html`
            <div id="canvas-parent" class="canvas-parent" @changed-color="${this.handleChangedColor}">
                <!-- <p id="canvas-details">Meta data:</p> -->
                <active-user-bar .userList="${this.userList}"></active-user-bar>
                
                <tool-palette .initColor="${this.color}"></tool-palette>
                <canvas id="canvas" width="${this.width}" height="${this.height}" @mousedown="${this.handleMouseDown}" @mouseup="${this.handleMouseUp}" @mousemove="${this.handleMouseMove}" ></canvas>
            </div>  
        `
    }

    handleMouseUp(e) {
        if (this.isMouseDown) {
            this.isMouseDown = false;
        }
        console.log("mosueup (isMouseDown):" + this.isMouseDown);
    }

    handleMouseDown(e) {
            // React to the mouse down event
        this.isMouseDown = true;
        var ctx = e.target.getContext('2d');
        var canvas = e.target;
        
        console.log("mousedown: " + ctx);
        console.log(e.pageX - canvas.offsetLeft);
        this.curX = e.pageX - canvas.offsetLeft;
        this.curY = e.pageY - canvas.offsetTop;
        this.lastX = this.curX;
        this.lastY = this.curY;
        console.log(this.curX, ", ", this.curY, ", " ,this.lastX, ", ", this.lastY );
    }

    handleMouseMove(e) {
        var canvas = e.target;
        var ctx = e.target.getContext('2d');

        if (this.isMouseDown && ctx) {

            this.lastX = this.curX;
            this.lastY = this.curY;
            this.curX = e.pageX - canvas.offsetLeft;
            this.curY = e.pageY - canvas.offsetTop;
            // paint
            var paintJSON = this.paint(ctx, this.curX, this.curY, this.lastX, this.lastY, this.color);            // format paint event
            
            // this.renderRoot.querySelector('#canvas-details').innerText = "This Client's Canvas: " + JSON.stringify(paintJSON);

            this.wsConn.send(JSON.stringify(paintJSON));

        }
    }
    paint (ctx, pageX, pageY, lastX, lastY, color) {
        // Set line width
        ctx.lineWidth = 5;
        ctx.lineJoin = 'round';
        ctx.lineCap = 'round';
        ctx.strokeStyle = color;
    
        // Paint
        ctx.beginPath();
        ctx.moveTo(lastX, lastY);
        ctx.lineTo(pageX, pageY);
        ctx.closePath();
        ctx.stroke();

        // // Return values to send to ws
        return {"PaintEvent":{curX: this.curX, curY: this.curY, lastX: lastX, lastY: lastY, color: this.color}};
    }
    
    handleChangedColor (e) {
        console.log("updateColor inROom Canvas: " + e.detail.color, e.currentTarget);
        this.color = e.detail.color;
    }

    connectToWS(ctx) {
        if (window['WebSocket']) {
            // wsStatus = this.renderRoot.querySelector("#canvas-details");
            const conn = new WebSocket('ws://' + document.location.host + document.location.pathname + '/ws');
            console.log('ws://' + document.location.host + document.location.pathname + '/ws');
            
            conn.onopen = function () {
                // conn.send("WS Open")
                // wsStatus.innerHTML = "Selected Color: " + "WS Open";
                console.log("WS Open");
                // this.wsConn = conn;
                console.log(this.wsConn);
            }

            conn.onclose = evt => { 
                // wsStatus.innerHTML = "WS Closed";
                console.log("WS Closed");
            }

            conn.onerror = function (e) {
                console.log("ws error: " + e);
            }
            
            conn.onmessage = evt => {
                // this.renderRoot.querySelector('#canvas-details').innerText = "Selected Color: " + evt.data;
                var canvas = this.renderRoot.querySelector('#canvas');
                if (!canvas) {
                    console.log("could not find canvas");
                    return
                }
                var ctx = canvas.getContext('2d');
                if (!ctx) {
                    console.log("could not find canvas");
                    return
                }
                // console.log("onmessage: " + evt.data);
                // This only works if evt.data is receicing a Paint JSON event
                if (evt.data) {
                    try {
                        var jsonEvent = JSON.parse(evt.data);
                        console.log("jsonEvent");
                        console.log(jsonEvent);
                        if (!Object.keys(jsonEvent).length > 0) {
                            return
                        }

                        var key = Object.keys(jsonEvent)[0];
                        // console.log(key)
                        switch (key) {
                            case "PaintEvent":
                                console.log("PaintEvent: " + key);
                                jsonEvent = jsonEvent[key]
                                this.paint(ctx, jsonEvent.curX, jsonEvent.curY, jsonEvent.lastX, jsonEvent.lastY, jsonEvent.color);
                                break;
                            case "CanvasState":
                                console.log("CanvasState: " + key);
                                jsonEvent = jsonEvent[key]

                                for (let i = 0; i < jsonEvent.length; i++) {
                                    // draw from senders canvas
                                    this.paint(ctx, jsonEvent[i]["CurX"], jsonEvent[i]["CurY"], jsonEvent[i]["LastX"], jsonEvent[i]["LastY"], jsonEvent[i]["Color"]);
                                }

                                break;
                            case "ActiveUsers":
                                console.log("ActiveUsers");
                                console.log(jsonEvent[key])
                                this.userList = JSON.stringify(jsonEvent);
                                // this.renderRoot.querySelector('#active-users').innerText = "Active Users: " + JSON.stringify(jsonEvent);

                                break;
                            default:
                                console.log("Other event: "+ key);
                                break;
                        }
                        
                    } catch(e) {
                        console.log("error parsing json event from websocket: "+ e); // error in the above string (in this case, yes)!
                    }
                }
            }
            return conn;
        }
        return
        
    }
}
customElements.define('room-canvas', RoomCanvas)