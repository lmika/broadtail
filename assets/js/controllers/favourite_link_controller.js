import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    static classes = [ "loading", "active" ];

    static values = {
        feedItemId: String,
        favouriteId: String
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

    favouriteIdValueChanged() {
        if (this.favouriteIdValue !== "") {
            this.element.classList.add(this.activeClass);
        } else {
            this.element.classList.remove(this.activeClass);
        }
    }

    async _doToggleActive() {
        if (this.favouriteIdValue === "") {
            try {
                let newFavourite = await this._addFavourite();
                this.favouriteIdValue = newFavourite.id;
            } catch (e) {
                console.error("cannot add new favourite", e);
                alert("Error adding new favourite");
            }
        } else {
            try {
                await this._deleteFavourite();
                this.favouriteIdValue = "";
            } catch (e) {
                console.error("cannot remove favourite", e);
                alert("Error removing favourite");
            }
        }
    }

    async _addFavourite() {
        let requestBody = JSON.stringify({
            "origin": {
                "type": "feed-item",
                "id": this.feedItemIdValue,
            }
        });

        let resp = await fetch(`/favourites/`, {
            method: "POST",
            body: requestBody,
            headers: {
                "Content-type": "application/json"
            }
        });
        return await resp.json();
    }

    async _deleteFavourite() {
        let resp = await fetch(`/favourites/${this.favouriteIdValue}`, {
            method: "DELETE",
        });
    }
}