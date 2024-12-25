# Implementation status:
* agencies 
  * parsing: OK
  * validation: OK
  * tests: OK
* calendar 
  * parsing: OK
  * validation: OK
  * tests: OK
* calendar_dates 
  * parsing: OK
  * validation: OK
  * tests: OK
* routes
    * parsing: partial
      * IMPLEMENTATION: route_color: Route color designation that matches public facing material. Defaults to white (FFFFFF) when omitted or left empty.
      * IMPLEMENTATION: route_text_color: Legible color to use for text drawn against a background of route_color. Defaults to black (000000) when omitted or left empty.
    * validation: partial
      * VALIDATION: route_color: The color difference between route_color and route_text_color should provide sufficient contrast when viewed on a black and white screen.
      * VALIDATION: route_text_color: The color difference between route_color and route_text_color should provide sufficient contrast when viewed on a black and white screen.
      * VALIDATION: route_url: URL of a web page about the particular route. Should be different from the agency.agency_url value.
      * VALIDATION: network_id: Forbidden it the route_networks.txt file is present. 
    * tests: OK
* shapes
    * parsing: OK
    * validation: OK
      * TODO: VALIDATION: shape_pt_sequence: Sequence in which the shape points connect to form the shape. Values must increase along the trip but do not need to be consecutive.
    * tests: NOK
* stops
    * parsing: NOK
    * validation: NOK
    * tests: NOK
* stop_times
    * parsing: NOK
    * validation: NOK
    * tests: NOK
* trips
    * parsing: NOK
    * validation: NOK
    * tests: NOK