import http from "k6/http";
import { check, group, sleep } from "k6";
import { Counter } from "k6/metrics";

// A simple counter for http requests

export const requests = new Counter("http_reqs");

// you can specify stages of your test (ramp up/down patterns) through the options object
// target is the number of VUs you are aiming for

export const options = {
    stages: [
        { target: 5, duration: "15s" },
        { target: 25, duration: "15s" },
        { target: 50, duration: "15s" },

    ],
    thresholds: {
        requests: ["count < 100"],
    },

    ext: {
        loadimpact: {
            projectID: 3492532,
            name: "CI github"
        }
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
    let res = http.get(`${__ENV.ENDPOINT}/`, params);

    sleep(1);

    group("GET", function () {
        let res = http.get("http://httpbin.org/get?verb=get");
        check(res, {
            "status is 200": (r) => r.status === 200,
            "response body": (r) => r.body.indexOf("destination-service-app"),
            "is verb correct": (r) => r.json().args.verb === "get",
        });
    });

    group("POST", function () {
        let res = http.post("http://httpbin.org/post", { verb: "post" });
        check(res, {
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => r.json().form.verb === "post",
        });
    });
}
