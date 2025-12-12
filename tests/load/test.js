import papa from "https://jslib.k6.io/papaparse/5.1.1/index.js";
import http from "k6/http";
import { sleep } from "k6";

export function setup() {
  const vus = parseInt(__ENV.VUS || "1");
  const runTime = __ENV.RUN_TIME || "1m";
  const spawnRate = parseInt(__ENV.SPAWN_RATE || "1");

  // tempo de ramp-up calculado
  const rampSeconds = Math.max(1, Math.floor(vus / spawnRate));

  return { vus, runTime, rampSeconds };
}

export const options = {
  scenarios: {
    load_test: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: `${__ENV.RAMP}s`, target: parseInt(__ENV.VUS) },
        { duration: __ENV.RUN_TIME, target: parseInt(__ENV.VUS) },
      ],
    },
  },
};

export default function () {
  const host = __ENV.HOST;
  http.get(`${host}/v1/optimizations/sync`);
  sleep(1);
}

export function handleSummary(data) {
  const httpReqDuration = data.metrics.http_req_duration;
  const httpReqs = data.metrics.http_reqs;
  const httpFailed = data.metrics.http_req_failed;

  const rows = [];

  // Access the values directly from the metric
  if (httpReqDuration && httpReqDuration.values) {
    const stats = httpReqDuration.values;

    rows.push({
      "Request Count": httpReqs?.values?.count ?? 0,
      "Failure Count": httpFailed?.values?.passes ?? 0,
      "Median Response Time": stats["p(50)"],
      "Average Response Time": stats.avg,
      "Min Response Time": stats.min,
      "Max Response Time": stats.max,
      "Average Content Size":
        data.metrics.http_req_receiving?.values?.avg ?? "",
      "Requests/s": httpReqs?.values?.rate ?? 0,
      "90%": stats["p(90)"],
      "95%": stats["p(95)"],
      "99%": stats["p(99)"] ?? "",
      "100%": stats.max,
    });
  }

  const csv = papa.unparse(rows);

  return {
    [__ENV.FILE_NAME]: csv,
  };
}
