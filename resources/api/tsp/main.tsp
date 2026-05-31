import "@typespec/http";
import "@typespec/versioning";

import "./models.tsp";
import "./statuspage.tsp";
import "./incident.tsp";
import "./views.tsp";

using Http;

@service(#{ title: "Statusy" })
@server("https://statusy.app", "A Status Page Aggregator")
@route("/api/")
namespace Statuspage;
