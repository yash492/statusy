import "@typespec/http";
import "@typespec/versioning";

using Http;

@service(#{ title: "Statusy" })
@server("https://statusy.app", "A Status Page Aggregator")
@route("/api/")
namespace Statuspage;

model Statuspage {
  @doc("Unique identifier for the status page")
  id: integer;

  @doc("Display name of the status page")
  name: string;

  @doc("URL-friendly slug for the status page")
  slug: string;

  @doc("Root domain URL of the status page")
  url: string;
}

model Incident {
  @doc("Unique identifier for the incident")
  id: integer;

  @doc("Title or summary of the incident")
  title: string;

  @doc("Current status of the incident (investigating, identified, monitoring, resolved)")
  status: string;

  @doc("The URL of the actual status page incident")
  incident_url: string;

  @doc("Timestamp when the incident was created by the provider")
  provider_created_at: utcDateTime;
}

model ScheduledMaintenance {
  @doc("Unique identifier for the incident")
  id: integer;

  @doc("Title or summary of the incident")
  title: string;

  @doc("Current status of the incident (investigating, identified, monitoring, resolved)")
  status: string;

  @doc("The URL of the actual scheduled maintenance event")
  scheduled_maintenance_url: string;

  @doc("The start time of the scheduled maintenance event in UTC.")
  starts_at: utcDateTime;

  @doc("The end time of the scheduled maintenance event in UTC")
  ends_at: utcDateTime;

  @doc("Timestamp when the incident was created by the provider")
  provider_created_at: utcDateTime;
}

model PaginationQuery {
  @doc("Page number for paginated results (default: 1)")
  @query
  page_number?: integer;

  @doc("Number of items per page (default: 20)")
  @query
  page_size?: integer;
}

model StatuspageIncidents {
  incidents: Incident[];
  statuspage: Statuspage;
  total_count: integer;
}

model StatuspageScheduledMaintenances {
  scheduled_maintenances: ScheduledMaintenance[];
  statuspage: Statuspage;
  total_count: integer;
}

model ErrorResponse {
  @doc("Error code identifier")
  code: string;

  @doc("Human-readable error message")
  message: string;
}

model NotFoundError is ErrorResponse {
  @statusCode statusCode: 404;
}

model RssResponse {
  @header("content-type")
  contentType: "application/rss+xml; charset=utf-8";

  @statusCode statusCode: 200;

  @body
  content: string;
}

model AtomResponse {
  @header("content-type")
  contentType: "application/atom+xml; charset=utf-8";

  @statusCode statusCode: 200;

  @body
  content: string;
}

@doc("List all available status pages with optional search filter")
@get
@route("/statuspages")
op listStatuspages(
  @doc("Search term to filter status pages by name or slug")
  @query
  search?: string,
): Statuspage[];

@doc("Get statuspage by slug")
@get
@route("/statuspages/{statuspageSlug}")
op statuspageBySlug(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,
): Statuspage;

@doc("Get all incidents for a specific status page")
@get
@route("/statuspages/{statuspageSlug}/incidents")
op incidentByStatuspage(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,

  @doc("Page number for incident updates")
  @query
  page_number?: integer,

  @doc("Number of updates per page")
  @query
  page_size?: integer,
): StatuspageIncidents;

// TODO: Change response
// @doc("Get detailed information about a specific incident with paginated updates")
// @get
// @route("/statuspages/{statuspageSlug}/incidents/{incidentID}")
// op incidentInfo(
//   @doc("Unique identifier of the status page")
//   @path
//   statuspageSlug: string,

//   @doc("Unique identifier of the incident")
//   @path
//   incidentID: string,
// ): Incident[];

@doc("Get all scheduled maintenance windows for a status page")
@get
@route("/statuspages/{statuspageSlug}/schedule-maintenances")
op scheduledMaintenanceByStatuspage(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,

  @doc("Page number for scheduled maintenance updates")
  @query
  page_number?: integer,

  @doc("Number of updates per page")
  @query
  page_size?: integer,
): StatuspageScheduledMaintenances;

// TODO: Change response
// @doc("Get detailed information about a specific scheduled maintenance")
// @get
// @route("/statuspages/{statuspageSlug}/schedule-maintenances/{scheduledMaintenanceID}")
// op scheduledMaintenanceInfo(
//   @doc("Unique identifier of the status page")
//   @path
//   statuspageSlug: string,

//   @doc("Unique identifier of the scheduled maintenance")
//   @path
//   scheduledMaintenanceID: string,
// ): Statuspage[];

@doc("Get RSS feed for a status page")
@get
@route("/statuspages/{statuspageSlug}/feed.rss")
op getRssFeed(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,
): RssResponse;

@doc("Get Atom feed for a status page")
@get
@route("/statuspages/{statuspageSlug}/feed.atom")
op getAtomFeed(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,
): AtomResponse;

model ViewServiceStatus {
  @doc("Unique identifier for the status page service")
  id: integer;

  @doc("Display name of the service")
  name: string;

  @doc("URL-friendly slug of the service")
  slug: string;

  @doc("Aggregated status of the service (up or down)")
  status: string;

  @doc("Active incident description (empty string if none)")
  last_incident: string;

  @doc("Whether all components of this service are included in the view")
  include_all_components: boolean;
}

model View {
  @doc("Unique identifier for the view")
  id: integer;

  @doc("Display name of the view")
  name: string;

  @doc("URL-friendly slug for the view")
  slug: string;

  @doc("Description of the view")
  description: string;

  @doc("Whether this is the default view")
  is_default: boolean;

  @doc("Configuration of services in this view")
  services: ViewServiceStatus[];
}

@doc("Get or create default view")
@post
@route("/views/default")
op getDefaultView(): View;

model Component {
  @doc("Unique identifier for the component")
  id: integer;

  @doc("Display name of the component")
  name: string;

  @doc("Provider-specific identifier for the component")
  provider_id: string;
}

model ComponentGroup {
  @doc("Unique identifier for the component group")
  id: integer;

  @doc("Display name of the component group")
  name: string;

  @doc("Provider-specific identifier for the component group")
  provider_id: string;

  @doc("List of components within this group")
  components: Component[];
}

model ServiceComponents {
  @doc("The service ID")
  service_id: integer;

  @doc("The service name")
  service_name: string;

  @doc("The service slug")
  service_slug: string;

  @doc("List of component groups for the service")
  grouped_components: ComponentGroup[];

  @doc("List of ungrouped components for the service")
  ungrouped_components: Component[];
}

@doc("Get all services not configured in a view")
@get
@route("/views/{viewSlug}/unconfigured-services")
op getUnconfiguredServices(
  @doc("The view slug")
  @path
  viewSlug: string,

  @doc("Search term to filter services by name or slug")
  @query
  search?: string,
): Statuspage[];

@doc("Get all components and groups for a service")
@get
@route("/services/{serviceSlug}/components")
op getServiceComponents(
  @doc("The service slug")
  @path
  serviceSlug: string,
): ServiceComponents;

model AddViewServiceRequest {
  @doc("The ID of the service to add to the view")
  service_id: integer;

  @doc("Whether to include all components of the service")
  include_all_components: boolean;

  @doc("List of specific component IDs to include (ignored if include_all_components is true)")
  component_ids?: integer[];

  @doc("List of specific component group IDs to include (ignored if include_all_components is true)")
  component_group_ids?: integer[];
}

model EditViewServiceRequest {
  @doc("Whether to include all components of the service")
  include_all_components: boolean;

  @doc("List of specific component IDs to include (ignored if include_all_components is true)")
  component_ids?: integer[];

  @doc("List of specific component group IDs to include (ignored if include_all_components is true)")
  component_group_ids?: integer[];
}

model ViewServiceResponse {
  @doc("Unique identifier for the view-service mapping")
  id: integer;

  @doc("The ID of the service")
  service_id: integer;

  @doc("Whether all components of this service are included")
  include_all_components: boolean;
}

@doc("Add a service to a view with component configuration")
@post
@route("/views/{viewSlug}/services")
op addViewService(
  @doc("The view slug")
  @path
  viewSlug: string,

  @body
  body: AddViewServiceRequest,
): ViewServiceResponse;

@doc("Edit a service's component configuration in a view")
@put
@route("/views/{viewSlug}/services/{serviceId}")
op editViewService(
  @doc("The view slug")
  @path
  viewSlug: string,

  @doc("The service ID")
  @path
  serviceId: integer,

  @body
  body: EditViewServiceRequest,
): ViewServiceResponse;

@doc("Remove a service from a view")
@delete
@route("/views/{viewSlug}/services/{serviceId}")
op deleteViewService(
  @doc("The view slug")
  @path
  viewSlug: string,

  @doc("The service ID")
  @path
  serviceId: integer,
): void;

model EditViewRequest {
  @doc("Display name of the view")
  name: string;

  @doc("URL-friendly slug for the view")
  slug: string;

  @doc("Description of the view")
  description: string;

  @doc("Whether this is the default view")
  is_default: boolean;
}

@doc("Edit a view's details")
@put
@route("/views/{viewSlug}")
op editView(
  @doc("The view slug")
  @path
  viewSlug: string,

  @body
  body: EditViewRequest,
): View;

@doc("Delete a view")
@delete
@route("/views/{viewSlug}")
op deleteView(
  @doc("The view slug")
  @path
  viewSlug: string,
): void;

