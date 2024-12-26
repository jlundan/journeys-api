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
    * parsing: partial OK
      * TODO: IMPLEMENTATION: route_color: Route color designation that matches public facing material. Defaults to white (FFFFFF) when omitted or left empty.
      * TODO: IMPLEMENTATION: route_text_color: Legible color to use for text drawn against a background of route_color. Defaults to black (000000) when omitted or left empty.
    * validation: partial
      * TODO: VALIDATION: route_color: The color difference between route_color and route_text_color should provide sufficient contrast when viewed on a black and white screen.
      * TODO: VALIDATION: route_text_color: The color difference between route_color and route_text_color should provide sufficient contrast when viewed on a black and white screen.
      * TODO: VALIDATION: route_url: URL of a web page about the particular route. Should be different from the agency.agency_url value.
      * TODO: VALIDATION: network_id: Forbidden it the route_networks.txt file is present. 
    * tests: OK
* shapes
    * parsing: OK
    * validation: partial OK
      * TODO: VALIDATION: shape_pt_sequence: Sequence in which the shape points connect to form the shape. Values must increase along the trip but do not need to be consecutive.
    * tests: NOK
* stops
    * parsing: OK
    * validation: partial OK
      * TODO: VALIDATION: stop_desc: Should not be a duplicate of stop_name
      * TODO: VALIDATION: zone_id: If this record represents a station or station entrance, the zone_id is ignored
      * TODO: VALIDATION: stop_url: This should be different from the agency.agency_url and the routes.route_url field values.
      * TODO: VALIDATION: parent_station: Required for locations which are entrances (location_type=2), generic nodes (location_type=3) or boarding areas (location_type=4).
      * TODO: VALIDATION: parent_station: Optional for stops/platforms (location_type=0)
      * TODO: VALIDATION: parent_station: Forbidden for stations (location_type=1) (this field must be empty)
      * TODO: VALIDATION: parent_station: Stop/platform (location_type=0): the parent_station field contains the ID of a station.
      * TODO: VALIDATION: parent_station: Entrance/exit (location_type=2) or generic node (location_type=3): the parent_station field contains the ID of a station (location_type=1)
      * TODO: VALIDATION: parent_station: Boarding Area (location_type=4): the parent_station field contains ID of a platform
      * TODO: VALIDATION: stop_timezone: If the location has a parent station, it inherits the parent station’s timezone instead of applying its own.
      * TODO: VALIDATION: stop_timezone: Stations and parentless stops with empty stop_timezone inherit the timezone specified by agency.agency_timezone.
      * TODO: VALIDATION: level_id: Foreign ID referencing levels.level_id (must exist)
      * TODO: VALIDATION: platform_code: Words like “platform” or "track" (or the feed’s language-specific equivalent) should not be included.
      * TODO: VALIDATION: stop_id: ID must be unique across all stops.stop_id, locations.geojson id, and location_groups.location_group_id values.
    * tests: NOK
* stop_times
    * parsing: OK
    * validation: NOK
    * tests: NOK
* trips
    * parsing: OK
    * validation: NOK
    * tests: NOK