import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

export const errorRate = new Rate('non_200_requests');

export let options = {
    stages: [
        // Ramp-up from 1 to 10 VUs in 10s.
        { duration: "10s", target: 10 },

        // Stay at rest on 10 VUs for 5s.
        { duration: "5s", target: 10 },

        // Linearly ramp down from 10 to 0 VUs over the last 15s.
        { duration: "15s", target: 0 }
    ],
    thresholds: {
        // We want the 95th percentile of all HTTP request durations to be less than 500ms
        "http_req_duration": ["p(95)<500"],
        // Thresholds based on the custom metric `non_200_requests`.
        // "non_200_requests": [
        //     // Global failure rate should be less than 1%.
        //     "rate<0.01",
        //     // Abort the test early if it climbs over 5%.
        //     { threshold: "rate<=0.05", abortOnFail: true },
        // ],
    },
};

export default function () {
    const url = 'http://localhost:8081/user/notification';
  
    const data = {
      type: "sms",
      description: "load testing",
    };
    check(http.post(url,JSON.stringify(data), {
      headers: { 'Content-Type': 'application/json' },
    }), {
      'status is 200': (r) => r.status == 200,
    }) || errorRate.add(1);
  
    sleep(0.5);
  }