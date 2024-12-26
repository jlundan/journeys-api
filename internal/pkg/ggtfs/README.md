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
    * tests: OK
* stop_times
    * parsing: OK
    * validation: partial OK
      * TODO: VALIDATION: stop_id: Required if stop_times.location_group_id AND stop_times.location_id are NOT defined.
      * TODO: VALIDATION: stop_id: Forbidden if stop_times.location_group_id or stop_times.location_id are defined.
      * TODO: VALIDATION: departure_time: Required for timepoint=1. Forbidden when start_pickup_drop_off_window or end_pickup_drop_off_window are defined. Otherwise optional.
      * TODO: VALIDATION: arrival_time: Required for the first and last stop in a trip (defined by stop_times.stop_sequence). Required for timepoint=1. Forbidden when start_pickup_drop_off_window or end_pickup_drop_off_window are defined. Otherwise optional.
      * TODO: VALIDATION: location_group_id: Forbidden if stop_times.stop_id or stop_times.location_id are defined
      * TODO: VALIDATION: location_id: Forbidden if stop_times.stop_id or stop_times.location_group_id are defined.
      * TODO: VALIDATION: stop_sequence: The values must increase along the trip but do not need to be consecutive.
      * TODO: VALIDATION: stop_sequence: Travel within the same location group or GeoJSON location requires two records in stop_times.txt with the same location_group_id or location_id.
      * TODO: VALIDATION: start_pickup_drop_off_window: Required if stop_times.location_group_id or stop_times.location_id is defined.  Required if end_pickup_drop_off_window is defined. Forbidden if arrival_time or departure_time is defined. Optional otherwise.
      * TODO: VALIDATION: pickup_type: pickup_type=0 forbidden if start_pickup_drop_off_window or end_pickup_drop_off_window are defined. pickup_type=3 forbidden if start_pickup_drop_off_window or end_pickup_drop_off_window are defined. Otherwise optional.
      * TODO: VALIDATION: drop_off_type: drop_off_type=0 forbidden if start_pickup_drop_off_window or end_pickup_drop_off_window are defined. Otherwise optional.
      * TODO: VALIDATION: drop_off_type: Forbidden if start_pickup_drop_off_window or end_pickup_drop_off_window are defined. Otherwise optional.
      * TODO: VALIDATION: continuous_drop_off: Forbidden if start_pickup_drop_off_window or end_pickup_drop_off_window are defined. Otherwise optional.
      * TODO: VALIDATION: Values used for shape_dist_traveled must increase along with stop_sequence; they must not be used to show reverse travel along a route.
      * TODO: VALIDATION: recommended for routes that have looping or inlining (the vehicle crosses or travels over the same portion of alignment in one trip). See shape_dist_traveled at https://github.com/google/transit/blob/master/gtfs/spec/en/reference.md#shapestxt
      * TODO: VALIDATION: pickup_booking_rule_id: Foreign ID referencing booking_rules.booking_rule_id (should reference a valid booking rule).
      * TODO: VALIDATION: drop_off_booking_rule_id: Foreign ID referencing booking_rules.booking_rule_id (should reference a valid booking rule).
      * TODO: VALIDATION: trip_id: Foreign ID referencing trips.trip_id (should reference a valid trip).
      * TODO: VALIDATION: stop_id: Referenced locations must be stops/platforms, i.e. their stops.location_type value must be 0 or empty.
      * TODO: VALIDATION: location_group_id: Foreign ID referencing location_groups.location_group_id (should reference a valid location group).
      * TODO: VALIDATION: location_id: Foreign ID referencing id from locations.geojson (should reference a valid location).
      * TODO: VALIDATION: stop_headsign: This field overrides the default trips.trip_headsign when the headsign changes between stops. If the headsign is displayed for an entire trip, trips.trip_headsign should be used instead.
      * TODO: VALIDATION: stop_headsign: A stop_headsign value specified for one stop_time does not apply to subsequent stop_times in the same trip. If you want to override the trip_headsign for multiple stop_times in the same trip, the stop_headsign value must be repeated in each stop_time row.
      * TODO: VALIDATION: drop_off_type: If this field is populated, it overrides any continuous pickup behavior defined in routes.txt.
      * TODO: VALIDATION: drop_off_type: If this field is empty, the stop_time inherits any continuous pickup behavior defined in routes.txt.
      * TODO: VALIDATION: continuous_drop_off: If this field is populated, it overrides any continuous drop-off behavior defined in routes.txt.
      * TODO: VALIDATION: continuous_drop_off: If this field is empty, the stop_time inherits any continuous drop-off behavior defined in routes.txt.
    * tests: OK
* trips
    * parsing: OK
    * validation: partial OK
      * TODO: VALIDATION: route_id: Foreign ID referencing routes.route_id (must refer to an existing route)
      * TODO: VALIDATION: service_id: Foreign ID referencing calendar.service_id or calendar_dates.service_id. (calendar_dates service_id relation is not checked)
      * TODO: VALIDATION: shape_id: Foreign ID referencing shapes.shape_id (must refer to an existing shape)
      * TODO: VALIDATION: shape_id: Required if the trip has a continuous pickup or drop-off behavior defined either in routes.txt or in stop_times.txt. Otherwise, optional.
    * tests: OK