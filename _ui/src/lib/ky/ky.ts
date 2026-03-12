import ky from "ky";
const KyClient = ky.create({
    prefixUrl: "http://localhost:8081",
    json: true,
})

export default KyClient;