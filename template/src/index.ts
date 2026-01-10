import { init } from "./elyntra";
import fs from "fs";

// a simple configurable start,
init({
  hostname: process.env.HOST || "0.0.0.0", //default
  port: +process.env.PORT! || 3000, //default
});
