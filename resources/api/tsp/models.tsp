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

model ErrorResponse {
  @doc("Error code identifier")
  code: string;

  @doc("Human-readable error message")
  message: string;
}

model BadRequestResponse {
  @statusCode statusCode: 400;
  @body body: ErrorResponse;
}

model NotFoundResponse {
  @statusCode statusCode: 404;
  @body body: ErrorResponse;
}

model ConflictResponse {
  @statusCode statusCode: 409;
  @body body: ErrorResponse;
}

model InternalServerErrorResponse {
  @statusCode statusCode: 500;
  @body body: ErrorResponse;
}

union StandardErrors {
  BadRequestResponse,
  NotFoundResponse,
  ConflictResponse,
  InternalServerErrorResponse,
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
