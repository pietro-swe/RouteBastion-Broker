import { sleep } from "k6";
import http from "k6/http";

export const options = {};

export default function () {
  const host = __ENV.HOST || "http://localhost:8000";
  http.get(`${host}/v1/optimizations/sync`);
  sleep(1);
}
