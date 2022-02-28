import { Application } from "@hotwired/stimulus"
import $ from 'jquery'
import jqueryUjsInit from 'jquery-ujs'

import PagingController from "./controllers/paging_controller"
import FavouriteLinkController from "./controllers/favourite_link_controller"

window.Stimulus = Application.start();
Stimulus.register("paging", PagingController);
Stimulus.register("favourite-link", FavouriteLinkController);

jqueryUjsInit($, undefined);