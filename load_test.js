import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { check } from 'k6';
import http from 'k6/http';

export const options = {
  scenarios: {
    create_prs: {
      executor: 'constant-vus',
      vus: 20,
      duration: '10s',
    },
  },
};


const BASE_URL = 'http://host.docker.internal:8080'; 

export function setup() {
  const teamName = `k6_team_${randomString(5)}`;
  const user1 = `k6_u1_${randomString(5)}`;
  const user2 = `k6_u2_${randomString(5)}`;
  const user3 = `k6_u3_${randomString(5)}`;

  const payload = JSON.stringify({
    team_name: teamName,
    members: [
        { user_id: user1, username: "K6_User1", is_active: true },
        { user_id: user2, username: "K6_User2", is_active: true },
        { user_id: user3, username: "K6_User3", is_active: true }
    ],
  });

  const params = { headers: { 'Content-Type': 'application/json' } };
  http.post(`${BASE_URL}/team/add`, payload, params);

  return { author_id: user1 };
}

export default function (data) {
  const prID = `pr_${randomString(10)}`;

  const payload = JSON.stringify({
    pull_request_id: prID,
    pull_request_name: `Load Test PR ${prID}`,
    author_id: data.author_id,
  });

  const params = { headers: { 'Content-Type': 'application/json' } };

  const res = http.post(`${BASE_URL}/pullRequest/create`, payload, params);

  check(res, {
    'is status 201': (r) => r.status === 201,
  });
}