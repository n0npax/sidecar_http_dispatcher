import http from "k6/http";
import { check, group, sleep } from "k6";
import { Counter } from "k6/metrics";

// A simple counter for http requests

export const requests = new Counter("http_reqs");

// you can specify stages of your test (ramp up/down patterns) through the options object
// target is the number of VUs you are aiming for

export const options = {
    stages: [
        { target: 25, duration: "30s" },
        { target: 25, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 50, duration: "30s" },
        { target: 50, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 75, duration: "30s" },
        { target: 75, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 100, duration: "30s" },
        { target: 100, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 125, duration: "30s" },
        { target: 125, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 150, duration: "30s" },
        { target: 150, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 175, duration: "30s" },
        { target: 175, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 200, duration: "30s" },
        { target: 200, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 225, duration: "30s" },
        { target: 225, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 250, duration: "30s" },
        { target: 250, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 275, duration: "30s" },
        { target: 275, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 300, duration: "30s" },
        { target: 300, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 325, duration: "30s" },
        { target: 325, duration: "300s" },
        { target: 5, duration: "5s" },
        { target: 5, duration: "60s" },
        { target: 350, duration: "30s" },
        { target: 350, duration: "300s" },
        { target: 5, duration: "600s" },
        { target: 5, duration: "60s" },
        { target: 500, duration: "60s" },
        { target: 5, duration: "1s" },
        { target: 5, duration: "60s" },
        { target: 500, duration: "1s" },
        { target: 500, duration: "10s" },
        { target: 5, duration: "1s" },
        { target: 5, duration: "1s" },

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
