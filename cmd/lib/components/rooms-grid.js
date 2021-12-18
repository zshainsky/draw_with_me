import { LitElement, html, css } from "lit";
import {gridStyles} from './styles.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class RoomsGrid extends LitElement {
    static properties = {
        rooms: {type: Array},
        loading: { type: Boolean },
        isRoomSelected: {type: Boolean},
        selectedRoomId: {type: String},
    };
    static styles = gridStyles;

    connectedCallback() {
        super.connectedCallback();
        this.isRoomSelected = false;
        this.selectedRoomId = "";
        console.log(this.isRoomSelected);
        console.log(this.rooms, !this.rooms);
        // if rooms doesn't exist, create it
        if(!this.rooms) {
            this.fetchRooms();
        }
    }
    async fetchRooms() {
        this.loading = true;
        const response = await fetch('/get-rooms');
        const jsonResponse = await response.json();
        this.rooms = jsonResponse["RoomsList"];
        this.loading = false;
        console.log("rooms list: " + this.rooms);
    }
    async createRoom() {
        this.loading = true;
        const response = await fetch('/create-room');
        console.log(response);
        const jsonResponse = await response.json();
        console.log(jsonResponse);
        // add to rooms
        this.rooms = [jsonResponse];
        console.log('/room-' + jsonResponse["Id"]);
        this.loading = false;
    }

    render () {
        if (this.loading) {
            return html` <p>Loading...</p> `;
        }
        return html`<div id="rooms">
            ${this.rooms.length > 0 ? 
                html`
                <div class="grid-title">Active Rooms</div>
                    <div class="grid-wrapper">
                        ${this.rooms.map( 
                            (item, index) => html `
                                    
                                <room-element class="grid-room clickable"  @click="${this.handleClick}" .id=${item["Id"]} .name=${item["Id"]} .isRoomSelected=${this.isRoomSelected} .selectedRoomId=${this.selectedRoomId}></room-element>
                                `
                        )}
                    </div>
                `
                :
                html`
                <div class="grid-wrapper">
                    <div class="grid-room empty-grid clicakable" @click="${this.createRoom}"> 
                        <p> Create your first Room </p>
                        
                    </div>
                </div>
                `
            }
            </div>
        </div>
        `
    }
    handleClick(e){
        // Get all items in grid
        var gridItems = e.currentTarget.parentElement.getElementsByClassName("click");
        // Check if any other items are already clicked
        if (gridItems.length > 0) {
            // look through all items in grid
            for (let item of gridItems) {
                // handle all items except the one that was just clicked (e)
                if (item.id != e.currentTarget.id) {
                    // check if item has click as a class
                    if (item.classList.value.includes('click')) {
                        item.classList.toggle('click');
                    }
                } else { // (item in preivew) if user selected the same item that is already clicked and displayed by parent
                    e.currentTarget.parentElement.classList.toggle("click");

                    this.selectedRoomId = ""
                    console.log("hide preview: " + this.selectedRoomId);
                }
            }
        } else { // no other elements in the grid have been clicked
            e.currentTarget.parentElement.classList.toggle("click");
            
            this.selectedRoomId = e.currentTarget.id;
            console.log("show preview with id: " + this.selectedRoomId);
        }
        
        e.currentTarget.classList.toggle("click");

    }
    openRoom(e) {

    }
    handleCreateNewRoom(e) {
        this.createRoom();
    }
    // handleMouseOver(e) {
    //     console.log("mouse over: " +e.currentTarget.classList);
    //     e.currentTarget.classList.toggle("mouseover");

    // }
    // handleMouseLeave(e) {
    //     console.log("mouse leave: " + e.currentTarget.id);
    //     e.currentTarget.classList.toggle("mouseover");
    // }
    //@mouseover="${this.handleMouseOver}" @mouseleave="${this.handleMouseLeave}"
}
customElements.define('rooms-grid', RoomsGrid);