import "@typespec/http";

using Http;

@service(#{ title: "Statusy" })
@server("https://statusy.app", "A Status Page Aggregator")
namespace Statuspage;

model Statuspage {
  @doc("Unique identifier for the status page")
  id: integer;

  @doc("Display name of the status page")
  name: string;

  @doc("URL-friendly slug for the status page")
  slug: string;
}

model Incident {
  @doc("Unique identifier for the incident")
  id: integer;

  @doc("Title or summary of the incident")
  title: string;

  @doc("Current status of the incident (investigating, identified, monitoring, resolved)")
  status: string;

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
  page_size?: integer
): Incident[];

// TODO: Change response
@doc("Get detailed information about a specific incident with paginated updates")
@get
@route("/statuspages/{statuspageSlug}/incidents/{incidentID}")
op incidentInfo(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,

  @doc("Unique identifier of the incident")
  @path
  incidentID: string,
): Incident[];

@doc("Get all scheduled maintenance windows for a status page")
@get
@route("/statuspages/{statuspageSlug}/schedule-maintenances")
op scheduleMaintenanceByStatuspage(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,
): Statuspage[];

//TODO: Change response
@doc("Get detailed information about a specific scheduled maintenance")
@get
@route("/statuspages/{statuspageSlug}/schedule-maintenances/{scheduleMaintenanceID}")
op scheduleMaintenanceInfo(
  @doc("Unique identifier of the status page")
  @path
  statuspageSlug: string,

  @doc("Unique identifier of the scheduled maintenance")
  @path
  scheduleMaintenanceID: string,
): Statuspage[];

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
