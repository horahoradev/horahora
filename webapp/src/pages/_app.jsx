import "antd/dist/antd.css";
import "../styles/index.css";

import Head from "next/head";

import { Footer } from "../components/footer";

function MyApp({ Component, pageProps }) {
  return (
    <>
      <Head>
        <meta charSet="utf-8" />
        <title>Horahora</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="theme-color" content="#000000" />
        <meta name="description" content="horahora" />
        <script
          id="theme-setup"
          async
          dangerouslySetInnerHTML={{
            __html: `
              const theme = document.cookie
                .split("; ")
                .find((row) => row.startsWith("theme="));
              const value = theme ? theme.split("=")[1] : "dark";
              document.documentElement.dataset.theme = value;`,
          }}
        ></script>
        <link rel="icon" href="/favicon.ico" />
        <link rel="apple-touch-icon" href="/logo192.png" />
        <link rel="manifest" href="%PUBLIC_URL%/manifest.json" />
      </Head>
      <div className="bg-yellow-50 dark:bg-gray-900 min-h-screen font-sans-serif">
        <Component {...pageProps} />
      </div>

      <Footer />
    </>
  );
}

export default MyApp;
