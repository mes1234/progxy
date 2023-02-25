import http from 'k6/http';
import { sleep } from "k6";

export default function () {
    // http.get('http://127.0.0.1:9100/test/integration/index.html');
    // http.get('http://127.0.0.1:9200');
    http.get('http://127.0.0.1:9700/ping');
    sleep(0.01);
}