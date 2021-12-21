import {css} from 'lit';
export const globalStyles = css `
    .clickable {
        cursor: pointer;
    }
    
`;
export const gridStyles = css`
        
        @media only screen and (max-width: 899px) {
            .grid-wrapper {
                display: grid;
                grid-template-areas: 
                    "a"
                    "b";
                grid-template-rows: 200px 200px;
            }
        }
        @media only screen and (min-width: 900px) {
            .grid-wrapper {
                display: grid;
                grid-template-areas: 
                    "a a"
                    "b b";
                grid-template-rows: 250px 250px;
            }
        }
        #rooms {
            /* margin: 1rem 75px 1rem 75px; */
            min-height: 100%;
        }
        .grid-wrapper {
                display: grid;
                grid-gap: 1.5rem;
                color: #444;
                transition: 1s;
                margin: 4%;
        }
        .grid-room {
            background-color: #444;
            color: #fff;
            opacity: .85;
            border-radius: 5mm;
            padding: 20px;
            font-size: 150%;
            text-align: center;
            transition: transform 500ms;
            cursor: pointer;
        }
            
        .grid-wrapper.click {
            grid-template-rows: 50vmin;
            transition: 2s;
        }
        .grid-room.click {
            grid-area: a;
            transition: 2s;
            opacity: 1;
        }
        .grid-room:hover.click{
            background-color: #444;
            transform: scale(1);


        }
        .grid-room:hover {
            background-color: #009c9c;
            /* transform: translateY(-5x); */
            transform: scale(1.025);
            opacity: 1;
        }

        .empty-grid {
            grid-area: a;
            background-color: #03463a;
        }

        .grid-title {
            color: rgb(0, 0, 0);
            border-radius: 5mm;
            padding: 20px;
            font-size: 150%;
            text-align: center;
            margin-bottom: inherit;
            display: flex;
            align-items: center;
            place-content: stretch center;
            flex-flow: row nowrap;
        }
        .grid-title .title{
            padding-right: 25px;
        }
        img.create {
            width: 40px;
            border-radius: 75%;
            transition: transform 500ms ease 0s;
        }
        img.create:hover {
            background-color: #009c9c;
            transform: scale(1.025);
            opacity: 1;

        }
    
        
`;

export const roomElementStyles = css `
    .open-room-btn {

    }
    .open-room-btn:hover {
        color: #009c9c;
        /* transform: translateY(-5x); */
        transform: scale(1.025);
        opacity: 1;
    }
`;

export const navigationStyles = css`
    .header {
        /* background-color: #444; */
        color: #000;
        padding: 20px;
        font-size: 250%;
        text-align: center;
        margin-bottom: 10px;
    }
    .header .home {
        float: left;
    }
    .header .home:hover {
        color: #009c9c;
        
    }
    .header .profile:hover {
        color: #009c9c;
    }
    .header .profile {
        float: right;
    }
    .profile img {
        border-radius: 50%;
        height:45px;
    }

`;

export const footerStyles = css`
    .footer-bar {
        color: #000;
        padding: 20px;
        font-size: 25%;
        text-align: center;
        margin-top: 10px;
    }
    .footer-bar .author {
        float: left;
    }
    .footer-bar .attribution {
        float: right;
    }

`;

export const roomStyles = css `
    canvas {
        background-color: #fff;
        box-shadow: 0px 0px 10px 1.5px #4040407a;
        border-radius: 7px;
    }
    .canvas-parent {
        padding: 20px;
    }
`;

export const toolPaletteStyles = css `
        #palette-parent {
            top: 50%;
            float: right;
            vertical-align: top;
        }
        input {
            vertical-align: top;
            float: right;
            display: inline-block;
        }
        label {
            margin: 10px;
        }
    #color {
        -webkit-appearance: none;
        padding: 0;
        border: none;
        border-radius: 20px;
        width: 40px;    
        height: 40px;
    }
    #color::-webkit-color-swatch {
        border: none;
        border-radius: 20px;
        padding: 0;
    }
    #color::-webkit-color-swatch-wrapper {
        border: none;
        border-radius: 20px;
        padding: 0;
    }
`;

export const activeUserBarStyles = css`
    img {
        width: 50px;
        border-radius: 50%;
        padding: 0px 2px 0px 2px;

    }
    img:hover{
        transform: scale(1.05);
        opacity: 1;
    }
`;