import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    initialize() {
        this._socket = new WebSocket(`ws://${document.location.host}/ws/status`);
        this._socket.addEventListener('message', this._onMessage.bind(this));
    }

    _onMessage(event) {
        let msgJson = JSON.parse(event.data);

        window.dispatchEvent(new CustomEvent('jobUpdateMessage', {
            detail: {
                message: msgJson,
            },
            bubbles: true
        }))
    }
}