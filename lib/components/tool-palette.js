import { LitElement, html, css } from "lit";
import {toolPaletteStyles, globalStyles} from './styles.js';

class ToolPalette extends LitElement { 
    static properties = {
        initColor: {},
        roomName: {},
        currentUser: {type: Object},
    };

    static styles = [toolPaletteStyles, globalStyles];
    connectedCallback() {
        super.connectedCallback();
    
    }
    constructor () {
        super();
        console.log();
        this.roomName = "";
        this.currentUser = {};
    }

    render () {
        return html`
            <div id="palette-parent">
                <div id="palette-inner">
                    <span class="" id="color_front" @click="${this.selectNewColor}"></span>
                    <input type="color" id="color" class="clickable color palette-element" value="${this.initColor}" @change="${this.dispatchChangeColor}" >
                    <a target="_blank" href="mailto:?subject=Come%20Draw%20with%20Me!&body=${this.currentUser.Name}%20has%20invited%20you%20to%20Draw%20with%20Me%20in%20this%20new%20room%20(${this.roomName}).%0A%0AClick%20this%20link%20to%20join%3A%20${document.URL}%0A%0A%2D%20The%20Draw%20with%20Me%20team">
                        <img class="clickable" src="static/img/invite.png" @click="${this.handleShare}"/>
                    </a>
                  </div>
            </div>
        `;
    }
    selectNewColor(e){
        var color = this.shadowRoot.querySelector("#color");
        color.click();
    }
    dispatchChangeColor(e) {
        var newColor = e.currentTarget.value;
        console.log(newColor + "Attempting to dispatchChangeCOlor");
        if (newColor){
            this._dispatchChangedColorEvent(newColor);
        }
    }
    _dispatchChangedColorEvent(newColor) {
        const options = {
            detail: { color: newColor },
            bubbles: true,
            composed: true
        };
        this.dispatchEvent(new CustomEvent('changed-color', options));
    }
    handleShare(e) {
        console.log("share clicked:", e);
    }


}
customElements.define('tool-palette', ToolPalette);