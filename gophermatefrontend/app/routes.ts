import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("/login", "../src/pages/LoginPage.js"), // Corrected the path to LoginPage.js
] satisfies RouteConfig;
