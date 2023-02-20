import http from 'k6/http';
import { sleep } from "k6";

export default function () {
    // http.get('http://127.0.0.1:9000/test/integration/index.html');
    // http.get('http://127.0.0.1:80/');
    http.get('http://127.0.0.1:9001/');
    sleep(0.01);
}