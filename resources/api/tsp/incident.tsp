import "@typespec/http";

using Http;

namespace Statuspage;

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

  @doc("Comma-separated list of component IDs to filter by")
  @query
  component_ids?: string,

  @doc("Comma-separated list of component group IDs to filter by")
  @query
  component_group_ids?: string,
): StatuspageIncidents | StandardErrors;

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

  @doc("Comma-separated list of component IDs to filter by")
  @query
  component_ids?: string,

  @doc("Comma-separated list of component group IDs to filter by")
  @query
  component_group_ids?: string,
): StatuspageScheduledMaintenances | StandardErrors;
