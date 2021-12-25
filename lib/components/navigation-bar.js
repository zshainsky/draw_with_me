import { LitElement, html, css } from "lit";
import {globalStyles, navigationStyles} from './styles.js';
// import {styleMap} from 'lit/directives/style-map.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class NavigationBar extends LitElement {
    static properties = {
        homeImg: {},
        profileImg: {},
        loading: { type: Boolean },
    };
    static styles =[globalStyles, navigationStyles];
    
    constructor() {
        super();
        this.homeImg = "static/img/favicon/favicon-32x32.png";
        if(!this.profileImg) {
            this.fetchUserInfo();
        }        
    }

    render () { 
        return html`
        <div class="header">Draw with me 
            <div class="home">
                <img class="clickable" src="${this.homeImg}" @click="${this.goHome}"><img>
            </div>
            
            <div class="profile">
                <img class="clickable" src="${this.profileImg}" @click="${this.goToSignin}"><img>
            </div>
        </div>`
    }
    async fetchUserInfo(){
        this.loading = true;
        const response = await fetch('/user-info');
        const contentType = response.headers.get("content-type");
        console.log(response.status, contentType);
        if (contentType == "application/json") {
            const jsonResponse = await response.json();
            this.profileImg = jsonResponse["Picture"] //The key value here might change if not Google auth type. TODO: Add check for authtype
            this.loading = false;
        } else {
            this.profileImg = "https://thumbs.dreamstime.com/b/palomino-shetland-pony-equus-caballus-17000908.jpg"
        }
        this.loading = false;
    }
    goToSignin() {
        location.href = "/signin";

    }
    goHome() {
        location.href = "/";
    }
}


customElements.define('navigation-bar', NavigationBar);