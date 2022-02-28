import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    static classes = [ "loading", "active" ];

    static values = {
        feedItemId: String,
        favourite: Boolean
    };

    async toggleActive(ev) {
        ev.preventDefault();
        try {
            this.element.classList.add("loading");

            await this._doToggleActive();
        } finally {
            this.element.classList.remove("loading");
        }
    }

    favouriteValueChanged() {
        if (this.favouriteValue) {
            this.element.classList.add(this.activeClass);
        } else {
            this.element.classList.remove(this.activeClass);
        }
    }

    async _doToggleActive() {
        let requestBody = JSON.stringify({ "favourite": !this.favouriteValue });

        let resp = await fetch(`/feeditems/${this.feedItemIdValue}`, {
            method: "PATCH",
            body: requestBody,
            headers: {
                "Content-type": "application/json"
            }
        });
        let respJson = await resp.json();

        this.favouriteValue = respJson.favourite;
    }
}