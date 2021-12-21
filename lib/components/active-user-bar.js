import { LitElement, html, css } from "lit";
import {activeUserBarStyles} from './styles.js';

class ActiveUserBar extends LitElement { 
    static properties = {
        userList: {type: String},
    };
    
    static styles = activeUserBarStyles;
    constructor(){
        super();
        console.log("constructing user list");
        if(this.userList){
            console.log("have data in constructor")
            console.log(this.userList);
        }
        console.log(this.userList);
    }
    // async performUpdate() {
    //     await this.userList;
    //     return super.performUpdate();
    //  }

    render() {
        var userListJSON = JSON.parse(this.userList)
        if(userListJSON["ActiveUsers"]){
            return html`
            <div class="user-list">
                <div class="user">
                    ${userListJSON["ActiveUsers"].map( 
                            (item, index) => html`
                            <img title="${item["Name"]} <${item["Email"]}" src="${item["Picture"]}">
                            `
                    )}
                </div>
                
            </div>`;
        } else {
            return html`
            <div class="user-list">
                
                
            </div>`;
        }
       
    }
}
customElements.define('active-user-bar', ActiveUserBar);