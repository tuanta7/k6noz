import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
    vus: 10, // 10 virtual users
    duration: '10s',
};

export default function () {
    const url = 'ws://localhost:8080/ws';
    const params = {
        headers: { 'X-Source': 'k6-script' },
    };

    const res = ws.connect(url, params, function (socket) {})
    check(res, { 'status is 101': (r) => r && r.status === 101 });
}