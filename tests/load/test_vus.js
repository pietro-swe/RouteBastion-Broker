import papa from "https://jslib.k6.io/papaparse/5.1.1/index.js";
import http from "k6/http";
import { Trend, Counter } from "k6/metrics";

const durationSeconds = 5 * 60;

const metrics = {
  50: {
    duration: new Trend("http_req_duration_50"),
    requests: new Counter("http_reqs_50"),
  },
  75: {
    duration: new Trend("http_req_duration_75"),
    requests: new Counter("http_reqs_75"),
  },
  100: {
    duration: new Trend("http_req_duration_100"),
    requests: new Counter("http_reqs_100"),
  },
  200: {
    duration: new Trend("http_req_duration_200"),
    requests: new Counter("http_reqs_200"),
  },
  300: {
    duration: new Trend("http_req_duration_300"),
    requests: new Counter("http_reqs_300"),
  },
};


const userLevels = [
  50,
  75,
  100,
  200,
  300,
];
const stageDuration = "5m";

export const options = {
  scenarios: userLevels.reduce((acc, users, index) => {
    acc[`load_${users}_users`] = {
      executor: "constant-vus",
      vus: users,
      duration: stageDuration,
      startTime: `${index * 5}m`,
      tags: { concurrent_users: String(users) },
    };
    return acc;
  }, {}),
};

export default function () {
  const res = http.get(`${__ENV.HOST}/v1/optimizations/sync`);

  const vus = __VU;

  let level;
  if (vus <= 50) level = 50;
  else if (vus <= 75) level = 75;
  else if (vus <= 100) level = 100;
  else if (vus <= 200) level = 200;
  else level = 300;

  metrics[level].duration.add(res.timings.duration);
  metrics[level].requests.add(1);
}


export function handleSummary(data) {
  const rows = [];

  for (const level of userLevels) {
    const durationMetric = data.metrics[`http_req_duration_${level}`];
    const reqMetric = data.metrics[`http_reqs_${level}`];

    if (durationMetric?.values && reqMetric?.values) {
      rows.push({
        "Concurrent Users": level,
        "Request Count": reqMetric.values.count,
        "Requests/s": reqMetric.values.count / durationSeconds,
        "Avg Response Time (ms)": durationMetric.values.avg,
        "p90 (ms)": durationMetric.values["p(90)"],
        "p95 (ms)": durationMetric.values["p(95)"],
        "Max (ms)": durationMetric.values.max,
      });
    }
  }

  return {
    [__ENV.FILE_NAME]: papa.unparse(rows),
  };
}
