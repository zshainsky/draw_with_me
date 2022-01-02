import { LitElement, html, css } from "lit";
import {toolPaletteStyles, globalStyles} from './styles.js';

class ToolPalette extends LitElement { 
    static properties = {
        initColor: {},
    };

    static styles = [toolPaletteStyles, globalStyles];
    connectedCallback() {
        super.connectedCallback();
    
    }
    constructor () {
        super();
        console.log()
    }

    render () {
        return html`
            <div id="palette-parent">
                <div id="palette-inner">
                    <span class="" id="color_front" @click="${this.selectNewColor}"></span>
                    <input type="color" id="color" class="clickable color palette-element" value="${this.initColor}" @change="${this.dispatchChangeColor}" >
                    <img class="clickable" src="https://www.freeiconspng.com/uploads/share-sharing-icon-29.png" @click="${this.handleShare}"/>
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