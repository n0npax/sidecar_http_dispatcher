import http from "k6/http";
import { check, group, sleep } from "k6";
import { Counter } from "k6/metrics";

// A simple counter for http requests

export const requests = new Counter("http_reqs");

// you can specify stages of your test (ramp up/down patterns) through the options object
// target is the number of VUs you are aiming for

export const options = {
    stages: [
        { target: 25, duration: "10s" },
        { target: 25, duration: "50s" },
        { target: 50, duration: "10s" },
        { target: 50, duration: "50s" },
        { target: 75, duration: "10s" },
        { target: 75, duration: "50s" },
        { target: 100, duration: "10s" },
        { target: 100, duration: "50s" },
        { target: 125, duration: "10s" },
        { target: 125, duration: "50s" },
        { target: 150, duration: "10s" },
        { target: 150, duration: "50s" },
        { target: 175, duration: "10s" },
        { target: 175, duration: "50s" },
        { target: 200, duration: "10s" },
        { target: 200, duration: "50s" },
        { target: 225, duration: "10s" },
        { target: 225, duration: "50s" },
        { target: 250, duration: "10s" },
        { target: 250, duration: "50s" },
        { target: 275, duration: "10s" },
        { target: 275, duration: "50s" },
        { target: 300, duration: "10s" },
        { target: 300, duration: "50s" },
        { target: 325, duration: "10s" },
        { target: 325, duration: "50s" },
        { target: 350, duration: "10s" },
        { target: 350, duration: "50s" },
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
