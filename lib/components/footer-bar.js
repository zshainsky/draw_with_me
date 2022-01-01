import { LitElement, html, css } from "lit";
import {footerStyles} from './styles.js';
// import {styleMap} from 'lit/directives/style-map.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class FooterBar extends LitElement {
    static properties = {
    };
    static styles = footerStyles;
    
    constructor() {
        super();
    }

    render () { 
        return html`
            <div class="footer-bar"> 
                <div class="author clickable"><a target='_blank' href='https://www.linkedin.com/in/zack-shainsky/'>Zack Shainsky - Updated 2021/2022</a></div>
                <div class="github clickable"><a target='_blank' href='https://github.com/zshainsky/draw_with_me'>Github Project</a></div>
                <div class="attribution clickable"><a target='_blank' href='https://www.freepik.com/vectors/background'>attribution</a></div>
            </div>`
    }
}
customElements.define('footer-bar', FooterBar);