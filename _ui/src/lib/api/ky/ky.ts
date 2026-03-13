import { API_SERVER_ROUTE } from "$env/static/private";
import ky from "ky";

const KyClient = ky.create({
    prefixUrl: API_SERVER_ROUTE,
    json: true,
})

export default KyClient;