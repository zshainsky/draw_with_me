import { LitElement, html, css } from "lit";
import {roomElementStyles} from './styles.js';
// import {styleMap} from 'lit/directives/style-map.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class RoomElement extends LitElement {
    static get properties() {
        return {
            id: {type: String},
            name: {type: String},
            isRoomSelected: {type: Boolean},
            selectedRoomId: {type: String, 
                hasChanged(newVal, oldVal) { 
                    
                    console.log("has changed: ", newVal, " : ", oldVal);
                }
            }
        };
        
    }

    static styles = roomElementStyles;

    constructor() {
        super();
        //Consider adding click event here: this.addEventListener('click', (e) => console.log(e.type, e.target.localName)));
        //Or shdow root example: createRenderRoot()
        //or bubble in the rooms-grid

    }

    render() {
        console.log(this.selectedRoomId + ", " + this.id, this.selectedRoomId != "");
        return html`
                    <div class="clickable open-room-btn center" @click="${this.openRoom}"> Open Room </div>
                    <div class="title">${this.name}</div>
                    <div class="canvas-preview icon">

                    </div>
                    `;
    }
    openRoom (e) {
        console.log("open room" + this.isRoomSelected);
        console.log(e.currentTarget);
        // open room in preview
        window.open("/room-"+this.id);
        
    }
}
customElements.define('room-element', RoomElement);