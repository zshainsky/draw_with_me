import { LitElement, html, css } from "lit";
import {roomElementStyles, roomStyles} from './styles.js';
// import {styleMap} from 'lit/directives/style-map.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class RoomElement extends LitElement {
    static get properties() {
        return {
            // room vars
            id: {type: String},
            name: {type: String},
            isRoomSelected: {type: Boolean},
            selectedRoomId: {type: String, 
                hasChanged(newVal, oldVal) { 

                    console.log("has changed: ", newVal, " : ", oldVal);
                }
            },
            // canvas vars
            canvasEL: {},
            canvasState: {},

            //animation vars
            stopAnimation: {type: Boolean},
            counter: {type: Number},
            maxCounter: {type: Number},
            
        };
        
    }

    static styles = [roomElementStyles, roomStyles];

    connectedCallback() {
        super.connectedCallback();
        
    }
    async firstUpdated() {
        // Give the browser a chance to paint
        await new Promise((r) => setTimeout(r, 0));
        // get canvas in shadow dom
        this.canvasEl = this.shadowRoot.querySelector("#canvas");
        
        // init paint for all canvases
        if(this.canvasState){
            this.paintAllEvents(this.canvasEl, this.canvasState["CanvasState"]);
            this.maxCounter = this.canvasState["CanvasState"].length;
        }


    }
    constructor() {
        super();
        this.stopAnimation = false;
        // console.log("constructor:",this.stopAnimation);
        this.counter = 0;
        this.maxCounter = 0;
    }

    render() {
        // console.log(this.selectedRoomId + ", " + this.id, this.selectedRoomId != "");
        return html`
                <div class="clickable open-room-btn center">
                    <div class="name">${this.name}</div>
                    <div id="${this.id}" class="canvas-preview icon">
                        <canvas id="canvas" @mouseenter="${this.handleMouseEnter}" @mouseleave="${this.handleMouseLeave}" @touchstart="${this.handleMouseEnter}" @touchend="${this.handleMouseLeave}"></canvas>
                    </div> 
                </div>
                `;
    }
    // paint on specified canvas (canvasToDrawOn)
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
        // console.log("scale: ",scaleX,scaleY)
        // set line width  
        ctx.lineWidth = 2;
        ctx.lineJoin = 'round';
        ctx.lineCap = 'round';
        ctx.strokeStyle = color;
    
        // paint
        ctx.beginPath();
        ctx.moveTo(lastX*scaleX, lastY*scaleY);
        ctx.lineTo(curX*scaleX, curY*scaleY);
        ctx.closePath();
        ctx.stroke();

    }
    // paint single event in the list (jsonPaintEventsList) at index (i) on specified canvas (canvasToDrawOn)
    paintEventAtIndex(i, canvasToDrawOn, jsonPaintEventsList) {
        i = parseInt(i);
        if (i >= jsonPaintEventsList.length){ return; }
        this.paint(canvasToDrawOn, jsonPaintEventsList[i]["CurX"], jsonPaintEventsList[i]["CurY"], jsonPaintEventsList[i]["LastX"], jsonPaintEventsList[i]["LastY"], jsonPaintEventsList[i]["Color"], jsonPaintEventsList[i]["UserId"], jsonPaintEventsList[i]["RoomId"], jsonPaintEventsList[i]["CanvasWidth"], jsonPaintEventsList[i]["CanvasHeight"]);     
    }

    // paint all events in the list (jsonPaintEventsList) on the specified canvas (canvasToDrawOn)
    paintAllEvents(canvasToDrawOn, jsonPaintEventsList) {
	    for (let i in jsonPaintEventsList) {
            this.paintEventAtIndex(i, canvasToDrawOn, jsonPaintEventsList);
        }
    }

    // animation functions
    startAnimating() {
        // must get local variables to use in the callback function animatePaint
        var ctx = this.canvasEl.getContext('2d');
        ctx.clearRect(0,0,this.canvasEl.width,this.canvasEl.height);
        
        var counter = 0;
        var maxCounter = this.maxCounter;
        var canvas = this.canvasEl;
        var ctx = canvas.getContext('2d');
        var canvasState = this.canvasState;

        var animatePaint = function() {
            // console.log(counter);
            // repeate until stopAnimation is set to true
            if (this.stopAnimation) {
                return; 
            }
            window.requestAnimationFrame(animatePaint);

            if (counter < maxCounter){
                // draw until all events paint events in list have been printed
                this.paintEventAtIndex(counter, canvas, canvasState["CanvasState"]);
                counter += 1;
                
            }
            // after all paint events have been displayed, clear canvas and loop through them again
            if (counter == maxCounter) {  
                // console.log("counter==maxCounter: ",counter,maxCounter);
                counter = 0; 
                ctx.clearRect(0,0,canvas.width,canvas.height);
            }
        }.bind(this);
        
        // invoke callback function to animate paint
        animatePaint();

    }
    resetCounters() {
        this.counter = 0;
        this.stopAnimation = false;
    }

    handleMouseEnter(e){
        // console.log("enter", e);
        // reset couters and start animating
        this.resetCounters();
        this.startAnimating();
    }
    handleMouseLeave(e){
        // console.log("leave");
        // Stop execution
        this.stopAnimation = true;
        // make sure canvasStae exists
        if (!this.canvasState) { return; }
        this.paintAllEvents(this.canvasEl, this.canvasState["CanvasState"]);

        if(e.type == "touchend"){
            _dispatchOpenRoom(this.id);
        }
    }
    _dispatchOpenRoom(id) {
        const options = {
            // detail: { roomId: id },
            bubbles: true,
            composed: true
        };
        this.dispatchEvent(new CustomEvent('open-room', options));
    }
}
customElements.define('room-element', RoomElement);