import http from "k6/http";
import { check, group, sleep } from "k6";
import { Counter } from "k6/metrics";

// A simple counter for http requests

export const requests = new Counter("http_reqs");

// you can specify stages of your test (ramp up/down patterns) through the options object
// target is the number of VUs you are aiming for

export const options = {
    stages: [
        { target: 5, duration: "60s" },
        { target: 5, duration: "60s" },
        { target: 20, duration: "60s" },
        { target: 20, duration: "60s" },
        { target: 35, duration: "60s" },
        { target: 35, duration: "60s" },
        { target: 50, duration: "60s" },
        { target: 50, duration: "60s" },
        { target: 65, duration: "60s" },
        { target: 65, duration: "60s" },
        { target: 80, duration: "60s" },
        { target: 80, duration: "60s" },
        { target: 95, duration: "60s" },
        { target: 95, duration: "60s" },
        { target: 110, duration: "60s" },
        { target: 110, duration: "60s" },
        { target: 125, duration: "60s" },
        { target: 125, duration: "60s" },
        { target: 140, duration: "60s" },
        { target: 140, duration: "60s" },
        { target: 155, duration: "60s" },
        { target: 155, duration: "60s" },
    ],
    thresholds: {
        requests: ["count < 100"],
    }
};

export default function () {
    // our HTTP request, note that we are saving the response to res, which can be accessed later
    let params = {
        cookies: { testCookie: "value" },
        headers: { "environment": "qa" },
        redirects: 15,
        tags: { k6test: "yes" },
    };

    group("GET", function () {
        let res = http.get(`${__ENV.ENDPOINT}`, params);
        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body": (r) => r.body.indexOf("destination-service-app"),
            "header": (r) => r.headers["Destination_app"] === "ok",
        });

    });
    sleep(1);
}