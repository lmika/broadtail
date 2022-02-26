import { Application } from "@hotwired/stimulus"
import $ from 'jquery'
import jqueryUjsInit from 'jquery-ujs'

import PagingController from "./controllers/paging_controller"

window.Stimulus = Application.start();
Stimulus.register("paging", PagingController);

jqueryUjsInit($, undefined);