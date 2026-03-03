import "@typespec/http";

using Http;
@service(#{ title: "Statusy" })
namespace Service;

model Service {
  id: int;
  name: string;
  slug: string;
};

model Incident {
  id: int;
  title: string;
  
}