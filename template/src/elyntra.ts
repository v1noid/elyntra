import Elysia from "elysia";
import { ElysiaWS } from "elysia/dist/ws";

type RxMessage = {
  type: "connection:close" | "request:response";
  data: any;
};
type RequestData = {
  headers: Record<string, string | undefined>;
  path: string;
  body: unknown;
  request: Request;
  method: string;
};
type ResponseData = {
  headers: any;
  status: string;
  body: any;
  id: string;
};

let conn: ElysiaWS | null = null;
const requests = new Map<string, RequestData>();
const requestPromises = new Map<
  string,
  (value: ResponseData | PromiseLike<ResponseData>) => void
>();
const TIMEOUT = 7000;

export function init(
  {
    hostname,
    port,
  }: {
    hostname?: string;
    port?: number;
  } = {
    hostname: process.env.HOST || "0.0.0.0",
    port: +process.env.PORT! || 3000,
  }
) {
  new Elysia()
    .all("*", async (c) => {
      try {
        const id = crypto.randomUUID();
        const splitUrl = c.request.url.split("/");
        const payload = {
          headers: c.headers,
          request: c.request,
          path: "/" + splitUrl.slice(3).join("/"),
          host: splitUrl[2],
          body: c.body,
          method: c.request.method,
        };

        requests.set(id, payload);

        if (conn) {
          conn.raw.send(
            JSON.stringify({ type: "request:handle", data: { id, ...payload } })
          );
        }

        const timeOutRace = new Promise<never>((res, rej) => {
          setTimeout(() => {
            requestPromises.delete(id);
            rej(new Error("TIMEOUT"));
          }, TIMEOUT);
        });

        const reqPromise = new Promise<ResponseData>((resolve) => {
          requestPromises.set(id, resolve);
        });

        const requestData = await Promise.race([timeOutRace, reqPromise]);
        const headers: Record<string, any> = {};
        for (const key in requestData.headers) {
          headers[key] = Array.isArray(requestData.headers[key])
            ? requestData.headers[key][0]
            : requestData.headers[key];
        }

        requestPromises.delete(requestData.id);

        return new Response(Buffer.from(requestData.body, "base64"), {
          headers: headers,
          status: +requestData.status,
        });
      } catch (error: any) {
        return new Response(error.message || "Unknown error", { status: 500 });
      }
    })
    .ws("/ws", {
      open(ws) {
        conn = ws;
        console.log(`New connection`);
      },
      close(ws) {
        conn = null as any;
      },
      message(ws, msg: RxMessage) {
        switch (msg.type) {
          case "connection:close": {
            let data = msg.data as { id: string };
            conn = null;
            data = null as any;
            break;
          }
          case "request:response": {
            let data = msg.data as ResponseData;
            const resolve = requestPromises.get(data.id);

            if (resolve) {
              resolve(data);
            }
            data = null as any;
            break;
          }
        }
      },
    })
    .listen(
      {
        port,
        hostname,
      },
      (server) => {
        console.log(`Server started on ${server.hostname}:${server.port}`);
      }
    );
}
