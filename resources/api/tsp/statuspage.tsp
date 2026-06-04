import "@typespec/http";

using Http;

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

@doc("List all available status pages with optional search filter")
@get
@route("/statuspages")
op listStatuspages(
  @doc("Search term to filter status pages by name or slug")
  @query
  search?: string,
): Statuspage[] | StandardErrors;
 
 @doc("Get statuspage by slug")
 @get
 @route("/statuspages/{statuspageSlug}")
 op statuspageBySlug(
   @doc("Unique identifier of the status page")
   @path
   statuspageSlug: string,
 ): Statuspage | StandardErrors;
 
 @doc("Get RSS feed for a status page")
 @get
 @route("/statuspages/{statuspageSlug}/feed.rss")
 op getRssFeed(
   @doc("Unique identifier of the status page")
   @path
   statuspageSlug: string,
 ): RssResponse | StandardErrors;
 
 @doc("Get Atom feed for a status page")
 @get
 @route("/statuspages/{statuspageSlug}/feed.atom")
 op getAtomFeed(
   @doc("Unique identifier of the status page")
   @path
   statuspageSlug: string,
 ): AtomResponse | StandardErrors;
