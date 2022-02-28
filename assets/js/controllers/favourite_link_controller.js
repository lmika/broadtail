import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    static classes = [ "loading", "active" ];

    static values = {
        isFavourite: Boolean
    };

    connect() {
        
    }

    toggleActive() {
        this.element.classList.add("loading");
        window.setTimeout(() => {
            this.element.classList.toggle(this.activeClass);
            this.element.classList.remove("loading");
        }, 1500);
    }
}