import { LitElement, html, css } from "lit";
import {navigationStyles} from './styles.js';
// import {styleMap} from 'lit/directives/style-map.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class NavigationBar extends LitElement {
    static properties = {
    };
    static styles = navigationStyles;
    
    constructor() {
        super();
    }

    render () { 
        return html`
        <div class="header">Draw with me 
            <div class="home clickable">home</div>
            <div class="profile clickable">profile</div>
        </div>`
    }
}


customElements.define('navigation-bar', NavigationBar);