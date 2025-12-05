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
  const host = __ENV.HOST || "http://localhost:8000";
  http.get(`${host}/v1/optimizations/sync`);
  sleep(1);
}
