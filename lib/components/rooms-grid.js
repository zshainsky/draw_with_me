import { LitElement, html, css } from "lit";
import { classMap } from 'lit/directives/class-map.js';
import {globalStyles, gridStyles} from './styles.js';
// import { LitElement, html } from "https://unpkg.com/lit-element/lit-element.js?module";

class RoomsGrid extends LitElement {
    static properties = {
        rooms: {type: Array},
        loading: { type: Boolean },
        isRoomSelected: {type: Boolean},
        selectedRoomId: {type: String},
        createImg: {type: String},
        classes: {},
    };
    static styles = [globalStyles, gridStyles];
    
    connectedCallback() {
        super.connectedCallback();
        this.isRoomSelected = false;
        this.selectedRoomId = "";
        this.createImg = "static/img/create-new-room2.png"
        // console.log(this.isRoomSelected);
        // console.log(this.rooms, !this.rooms);
        this.classes = {"single-grid-element": true};
        // if rooms doesn't exist, create it
        if(!this.rooms) {
            this.fetchRooms();
        }
       
           
    }
    async fetchRooms() {
        this.loading = true;
        const response = await fetch('/get-rooms');
        const contentType = response.headers.get("content-type");

        // console.log(response.status, contentType);
        if (contentType == "application/json") {
            const jsonResponse = await response.json();
            // console.log(jsonResponse);  
            this.rooms = jsonResponse["RoomsList"];
            this.loading = false;
        } else {
            this.rooms = []
            console.log("/get-rooms response code: " + response.status);
        }
        this.loading = false;
        
    }
    async createRoom() {
        this.loading = true;
        const response = await fetch('/create-room');
        // console.log(response);
        const jsonResponse = await response.json();
        // console.log(jsonResponse);
        // add to rooms
        if(this.rooms){
            this.rooms.push(jsonResponse);
        } else {
            this.rooms = [jsonResponse];
        }
        // console.log('/room-' + jsonResponse["Id"]);
        this.loading = false;
    }

    render () {
        if (this.loading) {
            return html` <p>Loading...</p> `;
        }
        this.setGridColumns(this.rooms);
        return html`<div id="rooms">
            ${this.rooms.length > 0 ? 
                html`
                <div class="grid-title">
                    <div class="title">
                        Active Rooms
                    </div>
                    <img class="clickable create" src="${this.createImg}" @click="${this.createRoom}" title="Create New Room"/>
                </div>
                    <div class="grid-wrapper ${classMap(this.classes)}">
                        ${this.rooms.map( 
                            (item, index) => html `       
                                <room-element class="grid-room clickable"  @click="${this.handleOpenRoom}" @open-room="${this.handleOpenRoom}" .id=${item["Id"]} .name=${item["Name"]} .canvasState=${item["CanvasState"]} .isRoomSelected=${this.isRoomSelected} .selectedRoomId=${this.selectedRoomId}></room-element>
                                `
                        )}
                    </div>
                `
                :
                html`
                <div class="grid-wrapper ${classMap(this.classes)}">
                    <div class="grid-room empty-grid clicakable create" @click="${this.createRoom}"> 
                        <p> Create your first Room </p>
                        
                    </div>
                </div>
                `
            }
            </div>
        </div>
        `
    }
    
    // preview functionality ... depricated ...
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
                    // console.log("hide preview: " + this.selectedRoomId);
                }
            }
        } else { // no other elements in the grid have been clicked
            e.currentTarget.parentElement.classList.toggle("click");
            
            this.selectedRoomId = e.currentTarget.id;
            console.log("show preview with id: " + this.selectedRoomId);
        }
        
        e.currentTarget.classList.toggle("click");

    }
    
    handleOpenRoom (e) {
        var room_id = e.currentTarget.id;
        // console.log("open room" + this.isRoomSelected);
        // console.log(e.currentTarget.id);
        // open room in preview
        window.open("/room-"+room_id);
        
    }
    handleCreateNewRoom(e) {
        this.createRoom();
    }

    setGridColumns(roomsList){
        // console.log("rooms length: ", roomsList.length);
        if(roomsList) {
            if (roomsList.length > 1) { 
                // console.log("rooms length > 1");
                this.classes = {"single-grid-element": false}; 
            }
        }
    }
}
customElements.define('rooms-grid', RoomsGrid);