import "@typespec/http";

using Http;

namespace Statuspage;

model PaginationQuery {
  @doc("Page number for paginated results (default: 1)")
  @query
  page_number?: integer;

  @doc("Number of items per page (default: 20)")
  @query
  page_size?: integer;
}

@error
model ErrorResponse {
  @doc("Error code identifier")
  code: string;

  @doc("Human-readable error message")
  message: string;
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
