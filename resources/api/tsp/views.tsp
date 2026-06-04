import "@typespec/http";

using Http;

namespace Statuspage;

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
  @doc("Display name of the view")
  name: string;

  @doc("Unique public identifier (Snowflake ID) for the view")
  public_id: string;

  @doc("Description of the view")
  description: string;

  @doc("Whether this is the default view")
  is_default: boolean;
}


@doc("Get or create default view")
@post
@route("/views/default")
op createOrGetDefaultView(): View | StandardErrors;

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
@route("/views/{publicId}/unconfigured-services")
op getUnconfiguredServices(
  @doc("The view public ID")
  @path
  publicId: string,

  @doc("Search term to filter services by name or slug")
  @query
  search?: string,
): Statuspage[] | StandardErrors;

@doc("Get all components and groups for a service")
@get
@route("/services/{serviceSlug}/components")
op getServiceComponents(
  @doc("The service slug")
  @path
  serviceSlug: string,
): ServiceComponents | StandardErrors;

model AddViewServiceRequest {
  @doc("The ID of the service to add to the view")
  service_id: integer;

  @doc("Whether to include all components of the service")
  include_all_components: boolean;

  @doc("Whether to monitor incidents for the service")
  monitor_incidents: boolean;

  @doc("Whether to monitor scheduled maintenances for the service")
  monitor_scheduled_maintenances: boolean;

  @doc("List of specific component IDs to include (ignored if include_all_components is true)")
  component_ids?: integer[];

  @doc("List of specific component group IDs to include (ignored if include_all_components is true)")
  component_group_ids?: integer[];
}

model EditViewServiceRequest {
  @doc("Whether to include all components of the service")
  include_all_components: boolean;

  @doc("Whether to monitor incidents for the service")
  monitor_incidents: boolean;

  @doc("Whether to monitor scheduled maintenances for the service")
  monitor_scheduled_maintenances: boolean;

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

  @doc("Whether to monitor incidents for the service")
  monitor_incidents: boolean;

  @doc("Whether to monitor scheduled maintenances for the service")
  monitor_scheduled_maintenances: boolean;

  @doc("List of specific component IDs to include (ignored if include_all_components is true)")
  component_ids?: integer[];

  @doc("List of specific component group IDs to include (ignored if include_all_components is true)")
  component_group_ids?: integer[];
}

@doc("Add a service to a view with component configuration")
@post
@route("/views/{publicId}/services")
op addViewService(
  @doc("The view public ID")
  @path
  publicId: string,

  @body
  body: AddViewServiceRequest,
): ViewServiceResponse | StandardErrors;

@doc("Get a service's component configuration in a view")
@get
@route("/views/{publicId}/services/{serviceId}")
op getViewService(
  @doc("The view public ID")
  @path
  publicId: string,

  @doc("The service ID")
  @path
  serviceId: integer,
): ViewServiceResponse | StandardErrors;

@doc("Edit a service's component configuration in a view")
@put
@route("/views/{publicId}/services/{serviceId}")
op editViewService(
  @doc("The view public ID")
  @path
  publicId: string,

  @doc("The service ID")
  @path
  serviceId: integer,

  @body
  body: EditViewServiceRequest,
): ViewServiceResponse | StandardErrors;

@doc("Remove a service from a view")
@delete
@route("/views/{publicId}/services/{serviceId}")
op deleteViewService(
  @doc("The view public ID")
  @path
  publicId: string,

  @doc("The service ID")
  @path
  serviceId: integer,
): void | StandardErrors;

model EditViewRequest {
  @doc("Display name of the view")
  name: string;

  @doc("Description of the view")
  description: string;

  @doc("Whether this is the default view")
  is_default: boolean;
}

@doc("Edit a view's details")
@put
@route("/views/{publicId}")
op editView(
  @doc("The view public ID")
  @path
  publicId: string,

  @body
  body: EditViewRequest,
): View | StandardErrors;

@doc("Delete a view")
@delete
@route("/views/{publicId}")
op deleteView(
  @doc("The view public ID")
  @path
  publicId: string,
): void | StandardErrors;
