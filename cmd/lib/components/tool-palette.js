import { LitElement, html, css } from "lit";
import {toolPaletteStyles} from './styles.js';

class ToolPalette extends LitElement { 
    static properties = {
        initColor: {},
    };

    static styles = toolPaletteStyles;
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
                <div>
                    <label for="color"> Color</label>
                    <span id="color_front" @click="${this.selectNewColor}"></span>
                    <input type="color" id="color" class="color palette-element" value="${this.initColor}" @change="${this.dispatchChangeColor}" >
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

}
customElements.define('tool-palette', ToolPalette);